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

	for mixes > 0 {
		mixes--
		for i := 0; i < length; i++ {
			node := originalNodes[i]

			item := node.Value.(sorted)

			// this *feels* wrong; can't this just be a slice?
			oldI := indexOf(currentOrder, node)
			newI := (oldI + item.value) % (length - 1)

			if newI < 0 {
				newI = length - 1 + newI
			}

			// compare newI to oldI, and move in that direction
			diff := newI - oldI

			if diff < 0 {
				// move back
				for diff != 0 {
					diff++
					currentOrder.MoveBefore(node, node.Prev())
				}
			} else {
				// move forward
				for diff != 0 {
					diff--
					currentOrder.MoveAfter(node, node.Next())
				}
			}
		}
	}

	out := make([]int, 0, length)

	for e := currentOrder.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(sorted).value)
	}

	return out
}

func indexOf(l *list.List, x *list.Element) (i int) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e == x {
			return
		}
		i++
	}

	return -1
}

func debug(l *list.List) {
	out := make([]int, 0, l.Len())

	for e := l.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(sorted).value)
	}

	fmt.Println("[debug]", out)
}
