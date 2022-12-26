package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
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

var shapes = []shape{
	{[]int{
		binaryStringToInt("1111") << 3,
	},
		4,
	},
	{[]int{
		binaryStringToInt("010") << 4,
		binaryStringToInt("111") << 4,
		binaryStringToInt("010") << 4,
	},
		3,
	},
	// these shapes are actually upside-down :S
	{[]int{
		binaryStringToInt("111") << 4,
		binaryStringToInt("001") << 4,
		binaryStringToInt("001") << 4,
	},
		3,
	},
	{[]int{
		1 << 6,
		1 << 6,
		1 << 6,
		1 << 6,
	},
		1,
	},
	{[]int{
		binaryStringToInt("11") << 5,
		binaryStringToInt("11") << 5,
	},
		2,
	},
}

func generator[T any](arr []T) func() T {
	i := -1
	l := len(arr)
	return func() T {
		i++
		return arr[i%l]
	}
}

func play(jetPattern string, iterations int) (height int) {
	// the play area (7 columns)
	// thought about starting with a floor: (1 << 7) - 1, // all mask of 7
	space := []int{}

	getJet := generator([]byte(jetPattern))
	getShape := generator(shapes)
	heights := []int{}

	// need to find cycle length, and height per cycle
	// use sequence of x positions, because it's better than checking
	// the pattern in `space` (there is a pattern but it doesn't appear
	// at the start or the end: only in the middle)
	sequence := []int{}
	var cycleLength, heightPerCycle int

	for i := 0; i < iterations; i++ {
		// loop shapes
		shape := getShape()

		// iterate air blowing and gravity pulling
		x := applyForces(&space, shape, getJet)

		if cycleLength == 0 {
			heights = append(heights, len(space))
			// find cycle pattern
			sequence = append(sequence, x)

			patternLength := findRepeatingPatternFromEnd(sequence)

			if patternLength > 0 {
				cycleLength = patternLength
				// current height - height at beginning of cycle
				heightPerCycle = len(space) - heights[len(heights)-cycleLength-1]
				// remaining iterations
				remaining := iterations - i
				// number of cycles remaining
				cycles := remaining / cycleLength
				// height of remaining cycles
				height += cycles * heightPerCycle
				// increment past patterns
				todo := remaining % cycleLength
				// skip to remaining iterations
				i = iterations - todo
			}
		}
	}

	return height + len(space)
}

func applyForces(space *[]int, shape shape, getJet func() byte) int {
	// 1. add shape to space
	top := len(*space)
	// add the shape to the top, 2 from left
	x := 2
	// 3 from the top
	y := top + 3

	isCollision := func() bool {
		for i := 0; i < len(shape.outline); i++ {
			j := y + i
			if j > top-1 {
				// nothing at this row yet
				return false
			}
			if j < 0 {
				// didn't make a bottom, so this is gone past
				return true
			}
			if (*space)[j]&(shape.outline[i]>>x) != 0 {
				// testing any other collision
				return true
			}
		}
		return false
	}

	for {
		// 2. air pushes left or right
		air := getJet()

		if air == airL {
			x = utils.Max(0, x-1)
		} else {
			x = utils.Min(7-shape.width, x+1)
		}

		if isCollision() {
			// reverse moving shape
			if air == airL {
				x++
			} else {
				x--
			}
		}

		// 3. shape falls 1 unit
		y--

		// check collision again
		if isCollision() {
			// rests at bottom+1
			y++
			// add shape to space
			for i := 0; i < len(shape.outline); i++ {
				j := y + i
				alignedPart := shape.outline[i] >> x
				if j > len(*space)-1 {
					*space = append(*space, alignedPart)
				} else {
					(*space)[j] ^= alignedPart
				}
			}
			break
		}
	}

	// we'll keep track of the sequence of x's to determine a cycle pattern
	return x
}

// debugging my insane bitmask concept
func PrintSpace(space []int) {
	fmt.Println()
	for i := len(space) - 1; i >= 0; i-- {
		row := space[i]
		str := fmt.Sprintf("|%07b|", row)
		str = strings.ReplaceAll(str, "0", ".")
		fmt.Println(str)
	}
	fmt.Println("+-------+")
}

func PrintSpaceWithShape(space []int, shape shape, x, y int) {
	// don't mutate space
	copy := append([]int{}, space...)

	for len(copy) < y+len(shape.outline) {
		copy = append(copy, 0)
	}

	// add shape to space
	for i := 0; i < len(shape.outline); i++ {
		j := y + i
		copy[j] ^= (shape.outline[i] >> x)
	}

	PrintSpace(copy)
}

// simple hash function from:
// https://golangprojectstructure.com/hash-functions-go-code/
func djb2(data []int) uint64 {
	hash := uint64(5381)

	for _, b := range data {
		hash += uint64(b) + hash + hash<<5
	}

	return hash
}

// try rolling hash?
// we believe `source` ENDS with a pattern
func findRepeatingPatternFromEnd(source []int) (length int) {
	end := len(source) - 1
	lowerLimit := 5
	upperLimit := len(source)/2 - 1

	for w := lowerLimit; w < upperLimit; w++ {
		a := source[end-w:]
		b := source[end-(w*2)-1 : end-w]

		if djb2(a) == djb2(b) {
			// hashes match; we have a pattern
			return w + 1
		}
	}

	return
}
