package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

func parseInput(data inType, mul int) []int {
	list := make([]int, len(data))

	for i, line := range data {
		list[i] = utils.ParseInt(line) * mul
	}

	return list
}

func reOrder(orig []int, mixes int) []int {
	length := len(orig)

	// original order with pointers
	originalNodes := make([]*int, length)
	// current order with pointers to the original
	currentOrder := make([]*int, length)

	for i := range orig {
		originalNodes[i] = &orig[i]
		currentOrder[i] = originalNodes[i]
	}

	for mixes > 0 {
		mixes--
		for i := 0; i < length; i++ {
			node := originalNodes[i]

			// need to get current index; so have to loop over the slice, I believe
			oldI := indexOf(currentOrder, node)
			newI := (oldI + *node) % (length - 1)

			if newI < 0 {
				newI = length - 1 + newI
			}

			// remove
			currentOrder = append(currentOrder[:oldI], currentOrder[oldI+1:]...)
			// insert
			currentOrder = append(currentOrder[:newI], append([]*int{node}, currentOrder[newI:]...)...)
		}
	}

	out := make([]int, length)

	for i, v := range currentOrder {
		out[i] = *v
	}

	return out
}

func indexOf(l []*int, x *int) (i int) {
	for i := range l {
		if x == l[i] {
			return i
		}
	}

	return -1
}
