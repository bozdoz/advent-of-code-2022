package main

import (
	"fmt"
	"strconv"
	"strings"
)

// airR = ">"[0] // but who cares?
var airL = "<"[0]

type shape struct {
	outline []int // 1 is rock, 0 is empty
	width   int
}

// sanity
func binaryStringToInt(bin string) int {
	val, err := strconv.ParseInt(bin, 2, 10)

	if err != nil {
		panic(err)
	}

	return int(val)
}

var shapes = [...]shape{
	{[]int{
		binaryStringToInt("1111"),
	},
		4,
	},
	{[]int{
		binaryStringToInt("010"),
		binaryStringToInt("111"),
		binaryStringToInt("010"),
	},
		3,
	},
	// these shapes are actually upside-down :S
	{[]int{
		binaryStringToInt("111"),
		binaryStringToInt("001"),
		binaryStringToInt("001"),
	},
		3,
	},
	{[]int{
		1,
		1,
		1,
		1,
	},
		1,
	},
	{[]int{
		binaryStringToInt("11"),
		binaryStringToInt("11"),
	},
		2,
	},
}

func play(jetPattern string, iterations int) (height int) {
	// 9 columns (starts with the floor and two walls)
	// tempted to do a bitmask again
	space := []int{
		(1 << 9) - 1, // all mask of 9
	}

	// closure generator for jet pattern
	getJet := func() func() byte {
		i := -1
		l := len(jetPattern)

		return func() byte {
			i++
			return jetPattern[i%l]
		}
	}()

	for i := 0; i < iterations; i++ {
		// loop shapes
		s := shapes[i%len(shapes)]

		// create a copy
		copy := shape{
			outline: append([]int{}, s.outline...),
			width:   s.width,
		}

		// shape begins falling 2 from left, and 3 from bottom
		padSpace(&space)

		// iterate air blowing and gravity pulling
		applyForces(&space, copy, getJet)
	}

	return getTowerHeight(space)
}

var emptyRow = binaryStringToInt("100000001")

func applyForces(space *[]int, shape shape, getJet func() byte) {
	// 1. add shape to space
	// add the shape to the top, 2 from left
	// width - shape.width - leftPadding + leftwall
	shift := 7 - shape.width - 2 + 1
	for i, v := range shape.outline {
		// TODO: should this be in the definition of the shape?
		shape.outline[i] = v << shift
		// add extra padding for space
		*space = append(*space, emptyRow)
	}

	// place the shape's bottom at a specific row
	bottom := len(*space) - len(shape.outline)

	// PrintSpaceWithShape(*space, shape, bottom)

	isCollision := func() bool {
		for i := 0; i < len(shape.outline); i++ {
			j := bottom + i
			if (*space)[j]&shape.outline[i] != 0 {
				return true
			}
		}
		return false
	}

	for {
		// 2. air pushes left or right
		air := getJet()

		// TODO: could copy here instead
		for i, b := range shape.outline {
			if air == airL {
				shape.outline[i] = b << 1
			} else {
				shape.outline[i] = b >> 1
			}
		}

		// check collision
		if isCollision() {
			// reverse moving shape
			for i, b := range shape.outline {
				if air == airL {
					shape.outline[i] = b >> 1
				} else {
					shape.outline[i] = b << 1
				}
			}
		}

		// 3. shape falls 1 unit
		bottom--

		// check collision again
		if isCollision() {
			// rests at bottom+1
			bottom++
			// add shape to space
			for i := 0; i < len(shape.outline); i++ {
				j := bottom + i
				(*space)[j] ^= shape.outline[i]
			}
			break
		}
	}
}

func getTowerHeight(space []int) int {
	rows := len(space)
	// find last row with a rock unit
	for r := rows - 1; r >= 0; r-- {
		row := space[r]
		if row != emptyRow {
			return r
		}
	}

	// can't get here, right?
	return -1
}

// prepare space for new shape by padding top at least 3 rows
func padSpace(space *[]int) {
	towerHeight := getTowerHeight(*space)

	// add shape 3 from bottom, including height of floor (1)
	newLen := towerHeight + 3 + 1
	// last shape may have left more padding than needed
	if len(*space) > newLen {
		*space = (*space)[:newLen]
	} else {
		// add empty rows
		for len(*space) < newLen {
			*space = append(*space, emptyRow)
		}
	}
}

// debugging my insane bitmask concept
func PrintSpace(space []int) {
	fmt.Println()
	for i := len(space) - 1; i >= 0; i-- {
		row := space[i]
		str := fmt.Sprintf("|%09b|", row)
		str = strings.ReplaceAll(str, "0", " ")
		fmt.Println(str)
	}
}

func PrintSpaceWithShape(space []int, shape shape, bottom int) {
	// don't mutate space
	copy := append([]int{}, space...)
	// add shape to space
	for i := 0; i < len(shape.outline); i++ {
		j := bottom + i
		copy[j] ^= shape.outline[i]
	}

	PrintSpace(copy)
}
