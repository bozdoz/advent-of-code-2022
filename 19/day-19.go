package main

import (
	"fmt"

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

		// example should be 9 and 12
		fmt.Println("best of all:", best)
	}

	// acutal input should be 1144
	return
}

func partTwo(data inType) (ans outType) {
	blueprints := parseInput(data)

	for i, bp := range blueprints {
		if i > 1 {
			break
		}
		best := bp.bestPath(32) * (i + 1)

		// example should be 56 and 62
		fmt.Println("best for", i+1, ":", best)

		ans += best
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
