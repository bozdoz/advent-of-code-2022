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
	position       string
	valvesOpen     types.Set[string]
	time, pressure int
}

func start(valves map[string]valve) (max int) {
	// part 1 cache hits = 56K
	cached := map[string]int{}

	// then priority queue of next states
	pq := types.PriorityQueue[state]{}
	pq.PushValue(&state{
		// You start at valve AA
		position:   "AA",
		valvesOpen: types.Set[string]{},
	}, 0)

	for pq.Len() > 0 {
		state := pq.Get()

		if state.time == 30-1 {
			// we're done

			// update `max`
			if state.pressure > max {
				max = state.pressure
			}

			continue
		}

		nextStates := getNextStates(state, valves)

		// push states to pq
		for i := range nextStates {
			next := nextStates[i]
			key := next.hash()

			pressure, ok := cached[key]
			isBetter := ok && (pressure < next.pressure)

			if !ok || isBetter {
				cached[key] = next.pressure
			} else {
				// cheaper or equal sequence
				continue
			}

			// not sure what priority to give here
			// this PQ is in ASC order 🤷‍♀️
			pq.PushValue(next, -next.pressure)
		}
	}

	return
}

func getNextStates(cur *state, valves map[string]valve) []*state {
	// 1. where are you?
	valve := valves[cur.position]

	// capacity represents movements to any `leadsTo`, or opening
	nextStates := make([]*state, 0, len(valve.leadsTo)+1)

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

func (cur state) copy() *state {
	// would love to know if there's a better way to copy data
	copy := &state{
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
func (cur *state) hash() string {
	// we can probably cache based on:
	// position and valvesOpen (and time?)
	// tempted to do a bitmask
	open := make([]string, 0, len(cur.valvesOpen))

	for key := range cur.valvesOpen {
		open = append(open, key)
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i] < open[j]
	})

	return fmt.Sprint(cur.position, "-", open, "-", cur.time)
}

func (state *state) String() (out string) {
	return fmt.Sprint(*state)
}
