package main

import "github.com/bozdoz/advent-of-code-2022/utils"

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadEmptyLineGroups

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	board := parseInput(data)

	state := start(board, 1)

	return 1000*(state.row+1) + 4*(state.col+1) + int(state.face)
}

func partTwo(data inType) (ans outType) {
	board := parseInput(data)

	state := start(board, 2)

	return 1000*(state.row+1) + 4*(state.col+1) + int(state.face)
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
