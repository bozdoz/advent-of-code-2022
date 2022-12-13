package main

import (
	"sort"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadEmptyLineGroups

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	monkeys := parseInput(data)

	for range [20]struct{}{} {
		monkeys.inspect(1)
	}

	most_active := make([]int, len(*monkeys))

	for _, monkey := range *monkeys {
		most_active = append(most_active, monkey.inspected)
	}

	sort.Slice(most_active, func(i, j int) bool {
		return most_active[i] > most_active[j]
	})

	return most_active[0] * most_active[1]
}

func partTwo(data inType) (ans outType) {
	monkeys := parseInput(data)

	for range [10000]struct{}{} {
		monkeys.inspect(2)
	}

	most_active := make([]int, len(*monkeys))

	for _, monkey := range *monkeys {
		most_active = append(most_active, monkey.inspected)
	}

	sort.Slice(most_active, func(i, j int) bool {
		return most_active[i] > most_active[j]
	})

	return most_active[0] * most_active[1]
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
