package main

import "github.com/bozdoz/advent-of-code-2022/utils"

// today's input data type
type inType = string

// how to read today's input
var fileReader = utils.ReadFile

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	return uniqueLettersIndex(data, 4)
}

func partTwo(data inType) (ans outType) {
	return uniqueLettersIndex(data, 14)
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
