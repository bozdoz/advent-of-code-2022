package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

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
		// *first time* using FindAllStringSubmatch
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
	resources, robots map[resource]int
}

// *first time* recurse inside a closure :D
var dfs func(st state)

func (bp blueprint) bestPath(timeLimit int) (best int) {
	cache := map[string]struct{}{}

	dfs = func(st state) {
		if st.time == timeLimit {
			// get geodes cracked
			// update best
			if st.resources[geode] > best {
				best = st.resources[geode]
			}
			return
		}

		nextStates := bp.getNextStates(st, timeLimit)

		for i := range nextStates {
			next := nextStates[i]

			// hash, cache
			key := next.hash()
			_, visited := cache[key]

			if !visited {
				cache[key] = struct{}{}
			} else {
				// worse or equal state
				continue
			}

			// pruning
			// could it possibly be better?
			if best > 0 {
				max := getMaxGeodesFromTimeLeft(next, timeLimit-next.time)
				if max <= best {
					// ! this saved perhaps the most time
					continue
				}
			}

			dfs(next)
		}
	}

	// we start with 1 ore robot
	dfs(state{
		robots: map[resource]int{
			ore: 1,
		},
	})

	return
}

var resources = [4]resource{geode, obsidian, clay, ore}

// which robots could we buy next and when can we buy them?
// and can it happen before the `end`?
func (bp blueprint) getNextStates(cur state, end int) []state {
	nextStates := []state{}

	for _, robot := range resources {
		// if we have every robot that is needed
		if robot != geode && cur.robots[robot] >= bp.maxRobots[robot] {
			// we don't buy non-geode robots we don't need
			// saves ~23s on Part 1
			continue
		}

		// can't buy geodes
		if robot == geode && cur.robots[obsidian] == 0 {
			continue
		}

		// can't buy obsidian
		if robot == obsidian && cur.robots[clay] == 0 {
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

		// add robot
		clone.robots[robot]++

		// increment resources from current robots
		for res, num := range cur.robots {
			clone.resources[res] += num * totalTime
		}

		nextStates = append(nextStates, clone)
	}

	return nextStates
}

func getMaxGeodesFromTimeLeft(st state, time int) int {
	// current geodes
	cur := st.resources[geode]
	// previous robots accumulate geodes
	rate := st.robots[geode] * time
	// new robots accumulate geodes
	// example: given diff == 5 -> 4+3+2+1 == 5*4/2
	maxBots := (time * (time - 1)) / 2

	return cur + rate + maxBots
}

func (s *state) copy() state {
	clone := state{
		time:      s.time,
		resources: map[resource]int{},
		robots:    map[resource]int{},
	}

	for k, v := range s.resources {
		clone.resources[k] = v
	}

	for k, v := range s.robots {
		clone.robots[k] = v
	}

	return clone
}

func (s *state) hash() string {
	// !incorrect to ignore geodes
	// cache by time & resources & robots

	return fmt.Sprint(s.time, s.resources[ore], s.resources[clay], s.resources[obsidian], s.resources[geode], s.robots[ore], s.robots[clay], s.robots[obsidian], s.robots[geode])
}
