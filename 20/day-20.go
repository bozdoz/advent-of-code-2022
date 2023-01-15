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
	list := parseInput(data, 1)

	ordered, zeroIndex := reOrder(list, 1)

	first := ordered[(zeroIndex+1000)%len(ordered)]
	second := ordered[(zeroIndex+2000)%len(ordered)]
	third := ordered[(zeroIndex+3000)%len(ordered)]

	return *first + *second + *third
}

func partTwo(data inType) (ans outType) {
	list := parseInput(data, 811589153)

	ordered, zeroIndex := reOrder(list, 10)

	first := ordered[(zeroIndex+1000)%len(ordered)]
	second := ordered[(zeroIndex+2000)%len(ordered)]
	third := ordered[(zeroIndex+3000)%len(ordered)]

	return *first + *second + *third
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
