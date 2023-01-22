package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

type tile rune

const (
	empty tile = ' '
	open  tile = '.'
	wall  tile = '#'
)

type direction int

const (
	right direction = iota
	down
	left
	up
)

type cmd int

const (
	move cmd = iota
	rotate
)

type instruction struct {
	command cmd
	value   int // if rotate, then -1 for L and +1 for R
}

type board struct {
	grid          map[int]map[int]tile
	height, width int
	start         [2]int
	instructions  []instruction
	squareSize    int                         // example and input differ
	squares       map[image.Rectangle]*square // for 3d
	steps         map[[2]int]int              // for debug
}

func parseInput(data inType, squareSize int) board {
	tiles := strings.Split(data[0], "\n")
	width := 0
	board := board{
		grid:       map[int]map[int]tile{},
		height:     len(tiles),
		steps:      map[[2]int]int{},
		squareSize: squareSize,
	}

	var start *[2]int

	for r, row := range tiles {
		if len(row) > width {
			width = len(row)
		}
		for c, val := range row {
			t := tile(val)
			switch t {
			case empty:
				continue
			case wall, open:
				board.set(r, c, t)
				if start == nil && t == open {
					start = &[2]int{r, c}
				}
			}
		}
	}

	board.width = width
	board.start = *start

	instructions := []instruction{}

	// index of last value
	last := -1
	for i, field := range data[1] {
		if field == 'L' || field == 'R' {
			value := utils.ParseInt(data[1][last+1 : i])
			// the previous value was an int
			instructions = append(instructions, instruction{
				command: move,
				value:   value,
			})
			// this value is a rotate (clockwise)
			value = 1
			if field == 'L' {
				// (counter-clockwise)
				value = -1
			}
			instructions = append(instructions, instruction{
				command: rotate,
				value:   value,
			})
			last = i
		}
	}
	// get last move value
	instructions = append(instructions, instruction{
		command: move,
		value:   utils.ParseInt(data[1][last+1:]),
	})

	board.instructions = instructions

	return board
}

func (b *board) set(r, c int, val tile) {
	grid := b.grid
	if grid[r] == nil {
		grid[r] = map[int]tile{}
	}
	grid[r][c] = val
}

func (b *board) get(r, c int) (val tile) {
	tile, ok := b.grid[r][c]

	// believe I solved this differently in part 1
	// needs to check if empty space here, as empty
	// should be the default, and not `0`
	if !ok {
		return empty
	}

	return tile
}

// alters `pos` given a direction
func moveOne(pos *[2]int, dir direction) {
	switch dir {
	case right:
		pos[1]++
	case left:
		pos[1]--
	case up:
		pos[0]--
	case down:
		pos[0]++
	}
}

// rope around if there is empty space
func (b *board) getNext(r, c int, dir direction) (next [2]int, hitWall bool) {
	cur := [2]int{r, c}

	for {
		next = cur
		moveOne(&next, dir)

		// negative numbers
		if next[0] < 0 {
			next[0] = b.height + next[0]
		} else if next[1] < 0 {
			next[1] = b.width + next[1]
		}

		// wrap around (I wonder if it's better to check if greater first)
		next[0] %= b.height
		next[1] %= b.width

		cell := b.get(next[0], next[1])

		// we ignored empty here and just kept looping
		if cell == open {
			return
		}
		if cell == wall {
			// return original
			return [2]int{r, c}, true
		}

		cur = next
	}
}

type state struct {
	// TODO: maybe just [2]int?
	row, col int
	face     direction
}

func start(b board, part int) state {
	cur := b.start
	// start looking right
	dir := right
	st := state{cur[0], cur[1], dir}

	if part == 2 {
		// 3d needs the faces of the cube
		b.squares = cubeFold(b)
	}

	// get movement instructions
	for _, inst := range b.instructions {
		if inst.command == rotate {
			dir += direction(inst.value)
			if dir == -1 {
				dir = 3
			} else if dir == 4 {
				dir = 0
			}
		} else {
			if part == 1 {
				// move by an amount
				cur = b.move(cur, dir, inst.value)
			} else {
				// move by an amount
				cur, dir = b.move3d(cur, dir, inst.value)
			}
		}
		st = state{cur[0], cur[1], dir}
	}

	return st
}

// for debugging
var stepsNum = 0

func (b board) move(cur [2]int, dir direction, steps int) [2]int {
	i := 0

	for i < steps {
		next, hitWall := b.getNext(cur[0], cur[1], dir)

		if hitWall {
			break
		}

		cur = next

		// debugging
		b.steps[cur] = stepsNum
		stepsNum = (stepsNum % 9) + 1

		i++
	}

	return cur
}

//
// string representations
//

func (st state) String() (out string) {
	dir := map[direction]string{
		down:  "D",
		left:  "L",
		up:    "U",
		right: "R",
	}
	return fmt.Sprintf("[ %d %d ] %s", st.row, st.col, dir[st.face])
}

func (d direction) String() string {
	switch d {
	case right:
		return ">"
	case left:
		return "<"
	case up:
		return "^"
	case down:
		return "v"
	}
	return ""
}

func (i instruction) String() (out string) {
	if i.command == move {
		return fmt.Sprintf("[M %v]", i.value)
	} else {
		return fmt.Sprintf("[R %v]", i.value)
	}
}

// print the whole board similar to the examples
func (b board) Debug() {
	// print the board
	out := ""
	for r := 0; r < b.height; r++ {
		for c := 0; c < b.width; c++ {
			tile := b.get(r, c)
			if tile == empty {
				out += " "
			} else {
				pos := [2]int{r, c}
				n, ok := b.steps[pos]

				switch {
				case ok:
					out += fmt.Sprint(n)
				case tile == wall:
					out += "#"
				case tile == open:
					out += "."
				}
			}
		}
		out += "\n"
	}

	fmt.Println(out)
}
