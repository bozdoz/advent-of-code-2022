package main

import (
	"fmt"
	"image"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

const (
	up    = '^'
	down  = 'v'
	left  = '<'
	right = '>'
)

type valley struct {
	height, width int
	start, end    image.Point
	walls         types.Set[image.Point] // walls don't move
	blizzards     map[rune]map[image.Point]struct{}
}

func parseInput(data inType) valley {
	height := len(data)
	width := len(data[0])
	valley := valley{
		start:     image.Point{0, 1},                  // hard-coding
		end:       image.Point{height - 1, width - 2}, // hard-coding
		height:    height,
		width:     width,
		walls:     types.Set[image.Point]{},
		blizzards: map[rune]map[image.Point]struct{}{},
	}

	valley.blizzards[up] = map[image.Point]struct{}{}
	valley.blizzards[down] = map[image.Point]struct{}{}
	valley.blizzards[left] = map[image.Point]struct{}{}
	valley.blizzards[right] = map[image.Point]struct{}{}

	// add wall above start to avoid moving out
	valley.walls.Add(image.Point{-1, 1})
	// add wall below end to avoid moving out
	valley.walls.Add(image.Point{height, width - 2})

	for r, row := range data {
		for c, char := range row {
			point := image.Point{r, c}
			switch char {
			case '#':
				valley.walls.Add(point)
			case up, down, left, right:
				valley.blizzards[char][point] = struct{}{}
			}
		}
	}

	return valley
}

type state struct {
	position image.Point
	minutes  int
}

func (vly valley) pathFinder(minutes int) int {
	// I think this may be A* or Dijkstra. I just don't know 🫣
	seen := types.Set[string]{}

	pq := types.PriorityQueue[state]{}
	pq.PushValue(&state{
		position: vly.start,
		minutes:  minutes,
	}, 0)

	for pq.Len() > 0 {
		state := pq.Get()

		if state.position == vly.end {
			// done
			// caches: 50727
			// fmt.Println("duplicate states", len(seen))
			// caches: 1148550
			// fmt.Println("blizzard overlap cache hits", len(overlapCache))
			return state.minutes
		}

		nextStates := vly.getNextStates(state)

		// push states to pq
		for i := range nextStates {
			next := nextStates[i]

			// remove seen
			hash := next.hash()
			if seen.Has(hash) {
				continue
			}
			seen.Add(hash)

			distance := vly.distanceToEnd(next.position)

			// prioritize states by distance and minutes?
			// this PQ is in ASC order 🤷‍♀️
			pq.PushValue(&next, distance+next.minutes)
		}
	}

	return -1
}

var moves = [...]image.Point{
	{-1, 0}, // up
	{1, 0},  // down
	{0, -1}, // left
	{0, 1},  // right
	{0, 0},  // still
}

func (vly valley) getNextStates(cur *state) []state {
	// moves up, down, left, right, or stays still
	states := make([]state, 0, 5)

	for i := range moves {
		move := cur.position.Add(moves[i])

		// if not a wall
		if !vly.walls.Has(move) {
			// if not a blizzard
			if !vly.stateOverlapsBlizzard(move, cur.minutes+1) {
				// add updated copy of state
				states = append(states, state{
					position: move,
					minutes:  cur.minutes + 1,
				})
			}
		}
	}

	return states
}

var overlapCache = map[string]bool{}

func updateCache(hash string, ret bool) {
	overlapCache[hash] = ret
}

// determine if a blizzard will be at a position at a given time
func (vly valley) stateOverlapsBlizzard(pos image.Point, min int) bool {
	hash := fmt.Sprint(pos, min)
	overlap, ok := overlapCache[hash]

	if ok {
		return overlap
	}

	width := vly.width - 2   // ignores walls
	height := vly.height - 2 // ignores walls

	// check if '<' is `minutes` to the right of `position`
	check := pos
	// Y is actually column; also +/- 1 to handle walls
	check.Y = (check.Y+min-1)%width + 1
	_, ok = vly.blizzards[left][check]

	if ok {
		updateCache(hash, true)
		return true
	}

	// check '>'
	check = pos
	// Y is actually column
	check.Y = (check.Y - min) % width

	// < 1 handles walls
	if check.Y < 1 {
		// wrap around
		check.Y += width
	}
	_, ok = vly.blizzards[right][check]

	if ok {
		updateCache(hash, true)
		return true
	}

	// check 'v'
	check = pos
	// X is actually row
	check.X = (check.X - min) % height

	// < 1 handles walls
	if check.X < 1 {
		// wrap around
		check.X += height
	}
	_, ok = vly.blizzards[down][check]

	if ok {
		updateCache(hash, true)
		return true
	}

	// check '^'
	check = pos
	// X is actually row
	check.X = (check.X+min-1)%height + 1

	_, ok = vly.blizzards[up][check]

	updateCache(hash, ok)

	return ok
}

func (vly valley) distanceToEnd(pos image.Point) int {
	return utils.ManhattanDistance(pos, vly.end)
}

func (st state) hash() (out string) {
	return fmt.Sprint(st)
}
