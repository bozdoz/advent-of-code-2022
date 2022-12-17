package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	area := parseInput(data)

	// prevent infinite loops
	done := make(chan bool)
	go func() {
		select {
		case <-time.After(100 * time.Millisecond):
			fmt.Println("--- timeout! ---")
			os.Exit(1)
		case <-done:
		}
	}()

	for area.dropSand() != nil {
		ans++
	}

	// stop timeout
	done <- true

	// TODO: should we track iterations instead?
	return
}

func partTwo(data inType) (ans outType) {
	area := parseInput(data)

	// floor is "two plus the highest y coordinate"
	area.bottom += 2

	// simulate the falling sand until
	// the source of the sand becomes blocked
	source := image.Point{500, 0}

	// lazy: add the floor as new rocks without refactoring
	area.addRocks([]image.Point{
		{source.X - area.bottom, area.bottom},
		{source.X + area.bottom, area.bottom},
	})

	// prevent infinite loops
	done := make(chan bool)
	go func() {
		timeout := time.After(1000 * time.Millisecond)

		for {
			select {
			case <-time.After(500 * time.Millisecond):
				// debugging sand count every so often
				fmt.Println("sand:", ans)
			case <-timeout:
				fmt.Println("--- timeout! ---")
				os.Exit(1)
			case <-done:
			}
		}
	}()

	// drop sand, one at a time
	for {
		lastSand := area.dropSand()

		ans++

		if lastSand.Eq(source) {
			break
		}
	}

	done <- true

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
