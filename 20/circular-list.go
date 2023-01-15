package main

import (
	"github.com/bozdoz/advent-of-code-2022/utils"
)

func parseInput(data inType, mul int) []*int {
	list := make([]*int, len(data))

	for i, line := range data {
		val := utils.ParseInt(line) * mul
		list[i] = &val
	}

	return list
}

// sorts list
func reOrder(orig []*int, mixes int) (sorted []*int, zeroIndex int) {
	length := len(orig)

	// current order with pointers to the original
	currentOrder := make([]*int, length)

	copy(currentOrder, orig)

	for lastIndex := length - 1; mixes > 0; mixes-- {
		for _, node := range orig {
			// need to get current index; so have to loop over the slice, I believe
			oldI := indexOf(currentOrder, node)
			newI := (oldI + *node) % lastIndex

			if newI < 0 {
				newI = lastIndex + newI
			}

			if oldI == newI {
				continue
			}

			// remove
			currentOrder = append(currentOrder[:oldI], currentOrder[oldI+1:]...)
			// insert
			currentOrder = append(currentOrder[:newI], append([]*int{node}, currentOrder[newI:]...)...)
		}
	}

	// find zero
	for i, v := range currentOrder {
		if *v == 0 {
			zeroIndex = i
			break
		}
	}

	return currentOrder, zeroIndex
}

// get index of pointer, not value
func indexOf(l []*int, x *int) (i int) {
	for i := range l {
		if x == l[i] {
			return i
		}
	}

	return -1
}
