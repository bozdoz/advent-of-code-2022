package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

type grid map[int]map[int]struct{}

func parseInput(data inType) grid {
	grid := grid{}

	for r, line := range data {
		for c, char := range line {
			if char == '#' {
				grid.set(r, c)
			}
		}
	}

	return grid
}

func (g *grid) set(r, c int) {
	if (*g)[r] == nil {
		(*g)[r] = map[int]struct{}{}
	}
	(*g)[r][c] = struct{}{}
}

func (g *grid) get(r, c int) bool {
	if (*g)[r] == nil {
		(*g)[r] = map[int]struct{}{}
	}
	_, ok := (*g)[r][c]

	return ok
}

func (g *grid) delete(r, c int) {
	delete((*g)[r], c)
}

func (g *grid) bounds() [4]int {
	var miny, minx, maxy, maxx int
	for r, row := range *g {
		for c := range row {
			// can't check 'r' in the first loop,
			// because cells may have been removed,
			// and the row may be empty
			if r < miny {
				miny = r
			} else if r > maxy {
				maxy = r
			}
			if c < minx {
				minx = c
			} else if c > maxx {
				maxx = c
			}
		}
	}

	return [4]int{miny, maxx, maxy, minx}
}

func (g *grid) start(part int) int {
	// iterate elves
	// they don't move if they're isolated
	// they move by alternating checks for N, S, W, E
	// they don't move if multiple elves are moving the same place

	// part 1 checks at 10
	const checkAt int = 10

	rounds := 0
	neighbourIter := 0

	for {
		rounds++

		// map moveTo -> moveFrom
		planned := map[image.Point]image.Point{}
		contested := map[image.Point]struct{}{}

		elves := 0
		still := 0

		// TODO: maybe better to keep track of neighbours instead
		for r, row := range *g {
			for c := range row {
				elves++

				neighbours := make([]bool, 12)

				// top
				neighbours[0] = g.get(r-1, c-1)
				neighbours[1] = g.get(r-1, c)
				neighbours[2] = g.get(r-1, c+1)
				// bottom
				neighbours[3] = g.get(r+1, c-1)
				neighbours[4] = g.get(r+1, c)
				neighbours[5] = g.get(r+1, c+1)
				// left
				neighbours[6] = neighbours[3]
				neighbours[7] = g.get(r, c-1)
				neighbours[8] = neighbours[0]
				// right
				neighbours[9] = neighbours[5]
				neighbours[10] = g.get(r, c+1)
				neighbours[11] = neighbours[2]

				if utils.All(neighbours, func(x bool, _ int) bool {
					// check if all are empty
					return !x
				}) {
					// elf doesn't do anything
					still++
					continue
				}

				// alternate between which direction we check first
				j := neighbourIter % 4
				for i := j; i < j+4; i++ {
					k := i % 4
					arr := neighbours[k*3 : k*3+3]

					if utils.All(arr, func(x bool, i int) bool {
						// check if empty
						return !x
					}) {
						// elf can move here
						var moveTo image.Point
						switch k {
						case 0:
							moveTo = image.Point{r - 1, c}
						case 1:
							moveTo = image.Point{r + 1, c}
						case 2:
							moveTo = image.Point{r, c - 1}
						case 3:
							moveTo = image.Point{r, c + 1}
						}

						_, ok := planned[moveTo]

						if ok {
							// contested; can't move there
							contested[moveTo] = struct{}{}
						} else {
							// not contested; plan to move there
							moveFrom := image.Point{r, c}
							planned[moveTo] = moveFrom
						}

						break
					}
				}
			}
		}

		// update next diirectional search
		neighbourIter += 1

		// move to planned, minus any that are contested
		for moveTo, moveFrom := range planned {
			_, ok := contested[moveTo]
			if ok {
				continue
			}
			// move is a delete and a set
			g.delete(moveFrom.X, moveFrom.Y)
			g.set(moveTo.X, moveTo.Y)
		}

		// stop at 10 if part 1
		if part == 1 && rounds == checkAt {
			// get square from bounds
			bounds := g.bounds()

			height := bounds[2] - bounds[0] + 1
			width := bounds[1] - bounds[3] + 1
			size := height * width
			empty := size - elves

			return empty
		}

		// all the elves stopped moving
		if elves == still {
			return rounds
		}
	}
}

// sanity check
func PrintGrid(g *grid) {
	bounds := g.bounds()

	miny := bounds[0]
	minx := bounds[3]
	width := bounds[1] - minx + 1

	out := []string{}

	for i := bounds[0]; i <= bounds[2]; i++ {
		out = append(out, strings.Repeat(".", width))
	}

	for r, row := range *g {
		for c := range row {
			i := r - miny
			j := c - minx
			out[i] = out[i][0:j] + "#" + out[i][j+1:]
		}
	}

	fmt.Print(strings.Join(out, "\n") + "\n\n")
}
