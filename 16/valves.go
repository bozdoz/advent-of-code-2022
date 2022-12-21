package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type valve struct {
	name    string
	flow    int
	leadsTo []string
}

func parseInput(data inType) map[string]valve {
	valves := map[string]valve{}

	// not sure wtf is up with the "s?" in the input/example
	var re = regexp.MustCompile(`^Valve ([A-Z]{2}) has flow rate=(\d+); tunnels? leads? to valves? (.*?)$`)

	for _, line := range data {
		parsed := re.FindStringSubmatch(line)

		name := parsed[1]
		flow := utils.ParseInt(parsed[2])
		leadsTo := strings.Split(parsed[3], ", ")

		valves[name] = valve{
			name:    name,
			flow:    flow,
			leadsTo: leadsTo,
		}
	}

	return valves
}

// copied the idea from: https://github.com/bozdoz/advent-of-code-2021/blob/354349e4943eba626edd877507c85f5df25d235b/23/shuffle.go

type state struct {
	time, pressure int
	position       string
	valvesOpen     types.Set[string]
}

func start(valves map[string]valve, time int) (max int, bestDuoPaths map[string]int) {
	// part 1 cache hits = 56K
	cached := map[string]int{}

	// keep track of all viable "solution" states, for part 2?
	bestDuoPaths = map[string]int{}

	var viableValveCount float64

	for _, valve := range valves {
		if valve.flow > 0 {
			viableValveCount++
		}
	}

	// then priority queue of next states
	pq := types.PriorityQueue[state]{}
	pq.PushValue(&state{
		// You start at valve AA
		position:   "AA",
		valvesOpen: types.Set[string]{},
	}, 0)

	for pq.Len() > 0 {
		state := pq.Get()

		if state.time == time-1 {
			// we're done

			// update `max` for part 1
			if state.pressure > max {
				max = state.pressure
			}

			// update allPaths for part 2
			// logic here is that if there are 2 of us, we should each open ~50%
			percentOpen := float64(len(state.valvesOpen)) / viableValveCount

			// test passes with 40% and 60%, but
			// actual puzzle does not pass after 8 seconds:
			// 		bestDuoPath Count: 146
			// trying 30% and 70% with actual puzzle worked:
			// 		bestDuoPath Count: 2422
			if percentOpen > 0.3 && percentOpen < 0.7 {
				// we can cache here again on valvesOpen, and update max(pressure)
				key := state.hashValvesOpen()
				pathPressure, ok := bestDuoPaths[key]

				if !ok || state.pressure > pathPressure {
					bestDuoPaths[key] = state.pressure
				}
			}

			continue
		}

		nextStates := getNextStates(*state, valves)

		// push states to pq
		for i := range nextStates {
			next := nextStates[i]
			key := next.hash()

			pressure, ok := cached[key]
			isBetter := pressure < next.pressure

			if !ok || isBetter {
				cached[key] = next.pressure
			} else {
				// cheaper or equal sequence
				continue
			}

			// not sure what priority to give here
			// this PQ is in ASC order 🤷‍♀️
			pq.PushValue(&next, -next.pressure)
		}
	}

	return
}

func getNextStates(cur state, valves map[string]valve) []state {
	// 1. where are you?
	valve := valves[cur.position]

	// capacity represents movements to any `leadsTo`, or opening
	nextStates := make([]state, 0, len(valve.leadsTo)+1)

	// 2. is the valve open and flow>0?
	if !cur.valvesOpen.Has(valve.name) && valve.flow > 0 {
		// we can open this valve
		nextState := cur.copy()
		// one minute to open a valve
		nextState.time++
		nextState.valvesOpen.Add(valve.name)
		nextState.addPressure(valves)
		nextStates = append(nextStates, nextState)
	}

	// 3. where can you go?
	for _, next := range valve.leadsTo {
		nextState := cur.copy()
		// one minute to move to a new valve
		nextState.time++
		nextState.position = next
		nextState.addPressure(valves)
		nextStates = append(nextStates, nextState)
	}

	return nextStates
}

func (cur state) copy() state {
	// would love to know if there's a better way to copy data
	copy := state{
		position:   cur.position,
		time:       cur.time,
		valvesOpen: types.Set[string]{},
		pressure:   cur.pressure,
	}

	for key := range cur.valvesOpen {
		copy.valvesOpen.Add(key)
	}

	return copy
}

// TODO: should valves be global, or a pointer?
func (cur *state) addPressure(valves map[string]valve) {
	for open := range cur.valvesOpen {
		cur.pressure += valves[open].flow
	}
}

// we probably need caching
func (cur state) hash() string {
	// we can probably cache based on:
	// position and valvesOpen (and time?)
	return fmt.Sprint(cur.position, "-", cur.hashValvesOpen(), "-", cur.time)
}

// useful for part 1 and part 2
func (cur state) hashValvesOpen() string {
	open := make([]string, 0, len(cur.valvesOpen))

	// tempted to do a bitmask
	for key := range cur.valvesOpen {
		open = append(open, key)
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i] < open[j]
	})

	asString := fmt.Sprint(open)

	// lazily omit the "[]" from the Sprint
	return asString[1 : len(asString)-1]
}
