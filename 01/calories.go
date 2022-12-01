package main

import (
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

func parseGroupedCalorieList(data dataType) []int {
	out := make([]int, 0, len(data))

	for _, group := range data {
		sum := 0
		for _, num_str := range strings.Split(group, "\n") {
			sum += utils.ParseInt(num_str)
		}
		out = append(out, sum)
	}

	return out
}
