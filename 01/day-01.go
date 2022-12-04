package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type dataType = []string

// how to read today's inputs
var fileReader = utils.ReadNewLineGroups

func partOne(data dataType) (ans int) {
	groups := parseGroupedCalorieList(data, 1)

	ans = groups[0]

	return
}

func partTwo(data dataType) (ans int) {
	groups := parseGroupedCalorieList(data, 3)

	ans = groups[0] + groups[1] + groups[2]

	return
}

func main() {
	// pass file reader and functions to call with input data
	utils.RunSolvers(utils.Day[dataType]{
		FileReader: fileReader,
		Fncs: []func(dataType) int{
			partOne,
			partTwo,
		},
	})
}
