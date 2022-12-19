package main

import (
	"fmt"
	"image"
	"time"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

// override in tests for example input
var partOneRow = 2000000

func partOne(data inType) (ans outType) {
	space := parseInput(data)

	y := partOneRow

	for x := space.xmin; x <= space.xmax; x++ {
		// check if x,y is definitely not a beacon
		definitelyNot := !space.couldBeBeacon(x, y)
		if definitelyNot {
			ans++
		}
	}

	return
}

// override in tests
var partTwoMax = 4000000

func partTwo(data inType) (ans outType) {
	space := parseInput(data)

	space.xmin = 0
	space.ymin = 0
	space.xmax = partTwoMax
	space.ymax = partTwoMax

	done := make(chan image.Point)

	go func() {
		done <- space.findMissingBeacon()
	}()

	i := 0
	for {
		select {
		case <-time.After(5 * time.Second):
			i++
			fmt.Printf("%d... ", i)
		case point := <-done:
			return point.X*4000000 + point.Y
		}
	}
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
