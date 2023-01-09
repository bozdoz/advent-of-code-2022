package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type resource int

const (
	ore resource = 1 << iota
	clay
	obsidian
	geode
)

type blueprint struct {
	// robot_type->cost_type->cost: ore,clay,obsidian,geode
	robots map[resource]map[resource]int
	// bitmask saves time with comparing
	robotbitmask map[resource]int
	// figure out when we should stop buying robots
	maxRobots map[resource]int
}

// tried to do nested matches; impossible
var costRegex = regexp.MustCompile(`robot costs ((?:\d+ \w+(?: and )?)+)\.`)

func parseInput(data inType) []blueprint {
	blueprints := make([]blueprint, len(data))

	for i, line := range data {
		// first time using FindAllStringSubmatch
		matches := costRegex.FindAllStringSubmatch(line, -1)
		for j, match := range matches {
			robot := resource(1 << (j))

			costs := match[1]

			for _, cost := range strings.Split(costs, " and ") {
				fields := strings.Fields(cost)
				num, res := utils.ParseInt(fields[0]), fields[1]

				var resIndex resource

				switch res {
				case "ore":
					resIndex = ore
				case "clay":
					resIndex = clay
				case "obsidian":
					resIndex = obsidian
				case "geode":
					resIndex = geode
				}

				if blueprints[i].robots == nil {
					blueprints[i].robots = map[resource]map[resource]int{}
				}
				if blueprints[i].robots[robot] == nil {
					blueprints[i].robots[robot] = map[resource]int{}
				}

				blueprints[i].robots[robot][resIndex] = num

				if blueprints[i].robotbitmask == nil {
					blueprints[i].robotbitmask = map[resource]int{}
				}

				blueprints[i].robotbitmask[robot] |= int(resIndex)
			}
		}
	}

	for i, bp := range blueprints {
		maxRobots := map[resource]int{}
		for _, robots := range bp.robots {
			for res, num := range robots {
				maxRobots[res] = utils.Max(maxRobots[res], num)
			}
		}
		blueprints[i].maxRobots = maxRobots
	}

	return blueprints
}

type state struct {
	time              int
	robotbits         int // bitmask
	resources, robots map[resource]int
	parent            *state
}

func (bp blueprint) bestPath(timeLimit int) (best int) {
	cache := map[string]int{}
	cacheHits := 0

	pq := types.PriorityQueue[state]{}

	pq.PushValue(&state{
		// we start with 1 ore robot
		robots: map[resource]int{
			ore: 1,
		},
		// one ore robot
		robotbits: int(ore),
		resources: map[resource]int{},
	}, 0)

	i := 0

	// outer:
	for pq.Len() > 0 {
		i++

		// if debug && i == 20000 {
		// 	break
		// }

		state := pq.Get()

		if state.time == timeLimit {
			// done
			// get geodes cracked
			// update best
			if state.resources[geode] > best {
				best = state.resources[geode]
				// fmt.Println("best", best)

				// if debug {
				// 	cur := state
				// 	for cur.parent != nil {
				// 		fmt.Println(cur)
				// 		cur = cur.parent
				// 	}
				// }
			}
			continue
		}

		nextStates := state.getNextStates(bp, timeLimit)

		for i := range nextStates {
			next := nextStates[i]
			key := next.hash()

			// hash, cache
			geodes, ok := cache[key]
			isBetter := geodes < next.resources[geode]

			if !ok || isBetter {
				cache[key] = next.resources[geode]
			} else {
				// worse or equal state
				cacheHits++
				continue
			}

			// pruning
			timeLeft := timeLimit - next.time
			noGeodes := next.resources[geode] == 0

			// no time left, or no geode robot at one minute earlier
			if noGeodes && (timeLeft == 0 || timeLeft == 1 && next.robots[geode] == 0) {
				continue
			}

			// pq is in ASC order
			pq.PushValue(&next, -next.priority())
		}
	}

	// fmt.Println("cache hits", cacheHits)

	return
}

func (cur state) getNextStates(bp blueprint, end int) []state {
	nextStates := []state{}

	// possible next states:
	// buy a single robot, according to available resources
	// each minute, each robot that you had last minute collects a resource

	// BASICALLY: which robots could we buy next and when can we buy them?
	// and can it happen before the `end`?
	for robot, bitmask := range bp.robotbitmask {
		// if we have every robot that is needed
		if bitmask&cur.robotbits == bitmask {
			// we don't buy non-geode robots we don't need
			// changes test from 5sec to 0.3sec
			if robot != geode && cur.robots[robot] >= bp.maxRobots[robot] {
				continue
			}

			// we *might* be able to buy (if time permits), but when (and how)?
			clone := cur.copy()

			maxTime := 0
			for res, num := range bp.robots[robot] {
				// costs `num` of `res`
				time := int(math.Ceil(float64(num-cur.resources[res]) / float64(cur.robots[res])))

				if (cur.time + time) >= end-1 {
					// we can't actually buy this before the end of time
					// skip ahead to the end
					maxTime = end - clone.time - 1
					break
				}

				// we have to wait for the resources which take the longest
				if time > maxTime {
					maxTime = time
				}
				// reduce resources by cost
				// (could go negative, but we'll increase it soon)
				clone.resources[res] -= num
			}

			// increment time
			totalTime := maxTime + 1
			clone.time += totalTime

			// fmt.Println("we can buy", robot, "at", totalTime, "min")

			// add robot
			clone.robots[robot]++
			// add to bitmask
			clone.robotbits |= int(robot)

			// increment resources from current robots
			for res, num := range cur.robots {
				clone.resources[res] += num * totalTime
			}

			nextStates = append(nextStates, clone)
		}
	}

	// if debug {
	// 	fmt.Println("state", cur)
	// 	fmt.Println("next", nextStates)
	// 	fmt.Println("---")
	// }

	return nextStates
}

func (s *state) copy() state {
	clone := state{
		time:      s.time,
		robotbits: s.robotbits,
		resources: map[resource]int{},
		robots:    map[resource]int{},
		parent:    s,
	}

	for k, v := range s.resources {
		clone.resources[k] = v
	}

	for k, v := range s.robots {
		clone.robots[k] = v
	}

	return clone
}

func (s state) priority() (p int) {
	for robot, count := range s.robots {
		// values latter robots more
		p += count * int(robot)
	}
	return
}

func (s state) hash() string {
	// cache by time & robots (not geodes)
	return fmt.Sprint(s.time, s.resources[ore], s.resources[clay], s.resources[obsidian], s.robots[ore], s.robots[clay], s.robots[obsidian])
}
