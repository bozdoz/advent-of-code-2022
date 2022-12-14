package main

import (
	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type square struct {
	height, distance int
	visited          bool
	neighbours       []*square
}

type heightmap struct {
	grid          [][]*square
	height, width int
	end           *square
}

// used for converting rune to int
const letter_a rune = 'a'

// distance should be +inf
const max = int(^uint(0) >> 1)

func parseInput(data inType, part int) *heightmap {
	height := len(data)
	width := len(data[0])
	hm := &heightmap{
		height: height,
		width:  width,
	}
	grid := make([][]*square, height)

	for r, line := range data {
		grid[r] = make([]*square, width)
		for c, char := range line {
			square := &square{
				distance: max,
			}
			if char == 'S' {
				char = 'a'
				// starting square
				square.distance = 0
			} else if char == 'E' {
				char = 'z'
				hm.end = square
			}

			if part == 2 && char == 'a' {
				// all a's are starting spots
				square.distance = 0
			}

			square.height = int(char - letter_a)
			grid[r][c] = square
		}
	}

	hm.grid = grid
	hm.updateNeighbours()

	return hm
}

func (hm *heightmap) updateNeighbours() {
	maxRow := hm.height - 1
	maxCol := hm.width - 1

	for r, row := range hm.grid {
		for c, cell := range row {
			indices := [][]int{
				{r + 1, c},
				{r, c + 1},
				{r - 1, c},
				{r, c - 1},
			}

			for _, coords := range indices {
				r, c := coords[0], coords[1]

				if r < 0 || c < 0 || r > maxRow || c > maxCol {
					continue
				}

				cell.neighbours = append(cell.neighbours, hm.grid[r][c])
			}
		}
	}
}

func (hm *heightmap) pathFinder() {
	pq := make(types.PriorityQueue[square], hm.height*hm.width)

	for r, row := range hm.grid {
		for c, square := range row {
			index := r*hm.width + c
			pq.NewItem(
				square,
				square.distance,
				index,
			)
		}
	}

	pq.Init()

	for pq.Len() > 0 {
		square := pq.Get()

		for i := range square.neighbours {
			// using [i], otherwise for i, square is a copy
			neighbour := square.neighbours[i]

			if neighbour.visited || neighbour.height-square.height > 1 {
				// already visited, or
				// we can only walk up a square at most 1 higher
				continue
			}

			// update dijkstra's distance
			neighbour.distance = utils.Min(
				// we've walked one extra step (+1)
				square.distance+1,
				neighbour.distance,
			)

			pq.Update(neighbour, neighbour.distance)
		}

		square.visited = true
	}
}
