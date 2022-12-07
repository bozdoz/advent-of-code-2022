package main

import "github.com/bozdoz/advent-of-code-2022/utils"

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	root := parseHistory(data)

	traverseDirectories(root, func(d *entry) {
		if d.size <= 100000 {
			ans += d.size
		}
	})

	return
}

func partTwo(data inType) (ans outType) {
	const (
		total int = 70000000
		need  int = 30000000
	)

	root := parseHistory(data)

	unused_space := total - root.size
	remaining_need := need - unused_space

	smallest := root

	traverseDirectories(root, func(d *entry) {
		if d.size > remaining_need && d.size < smallest.size {
			smallest = d
		}
	})

	return smallest.size
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
