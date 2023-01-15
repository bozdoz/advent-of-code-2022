package main

import (
	"container/list"
	"fmt"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

func parseInput(data inType, mul int) []int {
	list := make([]int, len(data))

	for i, line := range data {
		list[i] = utils.ParseInt(line) * mul
	}

	return list
}

type sorted struct {
	index, value int
}

func reOrder(orig []int, mixes int) []int {
	length := len(orig)

	// original order with pointers to current
	originalNodes := make([]*list.Element, length)
	// current order
	currentOrder := list.New()

	for i, item := range orig {
		originalNodes[i] = currentOrder.PushBack(sorted{i, item})
	}

	// TODO: do this X number of mixes
	for i := 0; i < length; i++ {
		node := originalNodes[i]

		item := node.Value.(sorted)

		oldI := item.index
		newI := (oldI + item.value) % (length - 1)

		// fmt.Println(item.value, oldI, newI)

		if newI < 0 {
			newI %= length
			newI = length - 1 + newI
		}
		if newI == 0 {
			newI = length - 1
		}

		// compare newI to oldI, and move in that direction
		diff := newI - oldI

		if diff < 0 {
			// move back
			for diff != 0 {
				diff++
				if node.Prev() != nil {
					currentOrder.MoveBefore(node, node.Prev())
				}
				// wrap around
				if currentOrder.Front() == node {
					currentOrder.MoveToBack(node)
				}
			}
		} else {
			// move forward
			for diff != 0 {
				diff--
				if node.Next() != nil {
					currentOrder.MoveAfter(node, node.Next())
				}
				// wrap around
				if currentOrder.Back() == node {
					currentOrder.MoveToFront(node)
				}
			}
		}

		// fmt.Print("i", i, " ")
		// debug(currentOrder)
	}

	out := make([]int, 0, length)

	for e := currentOrder.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(sorted).value)
	}

	return out
}

func debug(l *list.List) {
	out := make([]int, 0, l.Len())

	for e := l.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(sorted).value)
	}

	fmt.Println("[debug]", out)
}
