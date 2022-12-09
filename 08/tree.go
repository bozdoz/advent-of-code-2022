package main

import (
	"fmt"

	"github.com/bozdoz/advent-of-code-2022/types"
)

type TGrid [][]int

type Ttrees struct {
	grid          TGrid
	height, width int
}

// used for converting rune to int
const zero rune = '0'

// pretty standard grid creation
func parseInput(data inType) *Ttrees {
	height := len(data)
	width := len(data[0])

	trees := &Ttrees{
		height: height,
		width:  width,
	}

	grid := make(TGrid, height)

	for r, line := range data {
		grid[r] = make([]int, width)
		for c, char := range line {
			grid[r][c] = int(char - zero)
		}
	}

	trees.grid = grid

	return trees
}

// gets visible trees from all 4 directions
func (g *Ttrees) countVisible() int {
	// set of any visible trees with key "rowcol"
	visible := types.Set[string]{}

	check_height := func(r, c, max int) int {
		tree_height := g.grid[r][c]
		if tree_height > max {
			// tree is > max therefore visible
			visible.Add(fmt.Sprint(r, c))
			return tree_height
		}
		return max
	}

	// start now on inner trees HORIZONTAL
	for r := 1; r < g.height-1; r++ {
		// keep track of max so far
		// L->R (break if max == 9)
		max := g.grid[r][0]
		for c := 1; c < g.width-1 && max < 9; c++ {
			max = check_height(r, c, max)
		}

		// R->L (break if max == 9)
		max = g.grid[r][g.width-1]
		for c := g.width - 2; c > 0 && max < 9; c-- {
			max = check_height(r, c, max)
		}
	}

	// VERTICAL
	for c := 1; c < g.width-1; c++ {
		// keep track of max so far
		// T->B (break if max == 9)
		max := g.grid[0][c]
		for r := 1; r < g.height-1 && max < 9; r++ {
			max = check_height(r, c, max)
		}

		// B->T (break if max == 9)
		max = g.grid[g.height-1][c]
		for r := g.height - 2; r > 0 && max < 9; r-- {
			max = check_height(r, c, max)
		}
	}

	// include edges (top, bottom, left, right)
	edges := g.height*2 + (g.width-2)*2

	return len(visible) + edges
}

// scenic score is measured by how many trees can be seen by a given tree
// multiplied together
func (g *Ttrees) bestScenicScore() (best int) {
	// inner trees again
	for r := 1; r < g.height-1; r++ {
		for c := 1; c < g.width-1; c++ {
			score := g.score(r, c)
			if score > best {
				best = score
			}
		}
	}

	return
}

func (g *Ttrees) score(r, c int) (score int) {
	val := g.grid[r][c]

	// score starts as 1 to be a multiplier
	score = 1
	// count the trees you can see in each direction
	count := 0

	// trying desperately to cut down on duplication
	should_break := func(rr, cc int, last bool) bool {
		count++
		comp := g.grid[rr][cc]

		// if we're finished, multiply the score
		if comp >= val || last {
			score *= count
			// and reset count
			count = 0
			return true
		}

		return false
	}

	// looking down
	for i := r + 1; ; i++ {
		if should_break(i, c, i == g.height-1) {
			break
		}
	}

	// looking up
	for i := r - 1; ; i-- {
		if should_break(i, c, i == 0) {
			break
		}
	}

	// looking right
	for i := c + 1; ; i++ {
		if should_break(r, i, i == g.width-1) {
			break
		}
	}

	// looking left
	for i := c - 1; ; i-- {
		if should_break(r, i, i == 0) {
			break
		}
	}

	return
}
