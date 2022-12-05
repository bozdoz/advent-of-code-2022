package main

import "github.com/bozdoz/advent-of-code-2022/utils"

// today's input data type
type dataType = []string

// how to read today's inputs
var fileReader = utils.ReadLines

func partOne(data dataType) (ans int) {
	return parseSectionContains(data)
}

func partTwo(data dataType) (ans int) {
	return parseSectionOverlaps(data)
}

//
// BOILERPLATE BELOW
//

func main() {
	// pass file reader and functions to call with input data
	utils.RunSolvers(utils.Day[dataType, int]{
		FileReader: fileReader,
		Fncs: []func(dataType) int{
			partOne,
			partTwo,
		},
	})
}
