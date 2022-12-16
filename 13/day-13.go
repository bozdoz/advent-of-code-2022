package main

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadEmptyLineGroups

// today's output data type
type outType = int

func init() {
	// lots of logging today
	log.SetOutput(io.Discard)
}

func partOne(data inType) (ans outType) {
	pairs := parsePairs(data)

	for i := 0; i < len(pairs); i++ {
		pair := pairs[i]

		if pair.isOrdered() {
			// 1-based index
			ans += i + 1
		}
	}

	return
}

func partTwo(data inType) (ans outType) {
	packets := make([]item, 0, len(data)*2+2)

	for _, group := range data {
		lines := strings.Split(group, "\n")
		for _, line := range lines {
			packets = append(packets, parseItem(line))
		}
	}

	// add divider packets
	one, two := parseItem("[[2]]"), parseItem("[[6]]")
	packets = append(packets, one, two)

	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) == -1
	})

	// multiplier
	ans = 1

	pointerOne := fmt.Sprintf("%p", one)
	pointerTwo := fmt.Sprintf("%p", two)

outer:
	for i := range packets {
		switch fmt.Sprintf("%p", packets[i]) {
		case pointerOne:
			ans *= i + 1
		case pointerTwo:
			ans *= i + 1
			break outer
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
