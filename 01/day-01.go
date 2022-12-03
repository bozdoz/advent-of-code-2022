package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type dataType []string

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

// initialize the app by setting log flags
func init() {
	log.SetFlags(log.Llongfile)
}

// run the solvers
func main() {
	filename := utils.GetInputFile()
	data := fileReader(filename)

	fncs := map[string]func(dataType) int{
		"partOne": partOne,
		"partTwo": partTwo,
	}

	// run partOne and partTwo
	for k, fun := range fncs {
		s := time.Now()
		val := fun(dataType(data))

		fmt.Printf("%s: %v (%v)\n", k, val, time.Since(s))
	}
}
