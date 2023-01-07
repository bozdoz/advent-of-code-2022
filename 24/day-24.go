package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	vly := parseInput(data)
	return vly.pathFinder(0)
}

func partTwo(data inType) (ans outType) {
	vly := parseInput(data)

	// just swap the start and end to go back and forth
	swap := func() {
		start := vly.start
		end := vly.end
		vly.start = end
		vly.end = start
	}

	first := vly.pathFinder(0)

	swap()

	// start at the minute you finished the last trip
	andBack := vly.pathFinder(first)

	swap()

	finally := vly.pathFinder(andBack)

	return finally
}

//
// BOILERPLATE BELOW
//

func main() {
	// pass file reader and functions to call with input data
	utils.RunSolvers(utils.Day[inType, outType]{
		FileReader: fileReader,
		Fncs: []func(inType) outType{
			partOne,
			partTwo,
		},
	})
}
