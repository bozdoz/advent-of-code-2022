package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type dataType = []string

// how to read today's inputs
var fileReader = utils.ReadLines

func partOne(data dataType) (ans int) {
	tournament := parseGuide(data)

	return tournament.yourScore
}

func partTwo(data dataType) (ans int) {
	tournament := parseSuggestiveGuide(data)

	return tournament.yourScore
}

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
