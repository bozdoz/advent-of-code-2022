package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadEmptyLineGroups

// today's output data type
type outType = string

func partOne(data inType) (ans outType) {
	supplies := parseSupplies(data)

	return supplies.runProcedure()
}

func partTwo(data inType) (ans outType) {
	supplies := parseSupplies(data)

	return supplies.runProcedure9001()
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
