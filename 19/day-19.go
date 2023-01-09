package main

import "github.com/bozdoz/advent-of-code-2022/utils"

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	blueprints := parseInput(data)

	for i, bp := range blueprints {
		ans += bp.bestPath(24) * (i + 1)
	}

	return
}

func partTwo(data inType) (ans outType) {
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
