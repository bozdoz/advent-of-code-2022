package main

import (
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	valves := parseInput(data)

	ans, _ = start(valves, 30)

	return
}

func partTwo(data inType) (ans outType) {
	valves := parseInput(data)

	// return all answers
	// get the best combo between
	_, bestDuoPaths := start(valves, 26)

	// 2422 with filter
	// fmt.Println("bestDuoPath Count:", len(bestDuoPaths))
	disjoint := 0

	for a, aVal := range bestDuoPaths {
	inner:
		for b, bVal := range bestDuoPaths {
			if a == b {
				continue
			}

			for _, valve := range strings.Split(a, " ") {
				if strings.Contains(b, valve) {
					// paths must be completely disjoint
					continue inner
				}
			}

			disjoint++

			ans = utils.Max(ans, aVal+bVal)
		}
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
