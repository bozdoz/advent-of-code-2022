package main

import (
	"container/heap"
	"sort"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

func parseGroupedCalorieList(data dataType, num int) []int {
	out := make([]int, 0, len(data))

	for _, group := range data {
		sum := 0
		for _, num_str := range strings.Split(group, "\n") {
			sum += utils.ParseInt(num_str)
		}
		out = append(out, sum)
	}

	// DESC
	sort.Slice(out, func(i, j int) bool {
		return out[i] > out[j]
	})

	return out[:num]
}

func parseGroupedCalorieHeap(data dataType, num int) []int {
	intheap := make(types.IntHeap, 0, num+1)

	heap.Init(&intheap)

	for _, group := range data {
		sum := 0
		for _, num_str := range strings.Split(group, "\n") {
			sum += utils.ParseInt(num_str)
		}

		heap.Push(&intheap, sum)

		if len(intheap) > num {
			heap.Pop(&intheap)
		}
	}

	return intheap
}
