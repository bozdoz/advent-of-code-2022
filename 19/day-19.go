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
	blueprints := parseInput(data)

	for i, bp := range blueprints {
		best := bp.bestPath(24)
		ans += best * (i + 1)
	}

	return
}

func partTwo(data inType) (ans outType) {
	blueprints := parseInput(data)

	ans = 1

	for i, bp := range blueprints {
		if i == 3 {
			return
		}

		best := bp.bestPath(32)
		ans *= best
	}

	return
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
