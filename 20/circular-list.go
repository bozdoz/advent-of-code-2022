package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

func parseInput(data inType) []int {
	list := make([]int, len(data))

	for i, line := range data {
		list[i] = utils.ParseInt(line)
	}

	return list
}

type sorting struct {
	value   int
	visited bool
}

func reOrder(list []int) []int {
	ordered := make([]sorting, len(list))

	for i, v := range list {
		ordered[i] = sorting{value: v}
	}

	for i := 0; i < len(ordered); {
		s := ordered[i]

		if s.visited {
			i++
			continue
		}
		newI := (i + s.value) % (len(ordered) - 1)
		if newI < 0 {
			newI %= len(ordered)
			newI = len(ordered) + newI - 1
		}
		if newI == 0 {
			newI = len(ordered) - 1
		}

		// remove
		ordered = append(ordered[:i], ordered[i+1:]...)
		// insert
		ordered = append(ordered[:newI], append([]sorting{{
			value:   s.value,
			visited: true,
		}}, ordered[newI:]...)...)

		// don't adjust i; we revisit this index
	}

	out := make([]int, len(list))

	for i, s := range ordered {
		out[i] = s.value
	}

	return out
}
