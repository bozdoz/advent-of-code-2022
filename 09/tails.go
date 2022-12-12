package main

import (
	"container/list"
	"fmt"
	"image"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type Rope struct {
	nodes     list.List
	tailStops types.Set[image.Point]
}

func parseInput(data inType, nodes int) Rope {
	start := image.Point{0, 0}
	rope := Rope{
		tailStops: make(types.Set[image.Point]),
	}

	nodeList := list.New()

	// add nodes to the rope
	for ; nodes > 0; nodes-- {
		nodeList.PushFront(start)
	}

	rope.nodes = *nodeList

	// tail always starts with head
	rope.tailStops.Add(start)

	for _, line := range data {
		fields := strings.Fields(line)
		dir, steps := fields[0], utils.ParseInt(fields[1])

		switch dir {
		case "U":
			rope.move(image.Point{0, -steps})
		case "D":
			rope.move(image.Point{0, steps})
		case "L":
			rope.move(image.Point{-steps, 0})
		case "R":
			rope.move(image.Point{steps, 0})
		}
	}

	return rope
}

// convert vector into a sign vector
// ex: {4,0} -> {1,0} and {0,-4} -> {0,-1}
func getSign(point image.Point) image.Point {
	x := point.X
	y := point.Y

	if x != 0 {
		if x < 0 {
			x /= -x
		} else {
			x /= x
		}
	}

	if y != 0 {
		if y < 0 {
			y /= -y
		} else {
			y /= y
		}
	}

	return image.Point{x, y}
}

var zero = image.Point{0, 0}

// head moves, tail follows, in some way
func (rope *Rope) move(position image.Point) {
	step := getSign(position)

	// move in steps
	for !position.Eq(zero) {
		position = position.Sub(step)

		// head moves
		head := rope.nodes.Front()
		head.Value = head.Value.(image.Point).Add(step)

		// tails drag behind
		prev := head
		cur := prev.Next()

		for cur != nil {
			prevPoint := prev.Value.(image.Point)
			curPoint := cur.Value.(image.Point)
			diff := prevPoint.Sub(curPoint)

			if diff.Eq(zero) {
				// tail stopped moving
				break
			}

			if utils.Abs(diff.X) > 1 || utils.Abs(diff.Y) > 1 {
				// tail plays catch up
				cur.Value = curPoint.Add(getSign(diff))
			}

			prev = cur
			cur = prev.Next()
		}

		rope.tailStops.Add(rope.nodes.Back().Value.(image.Point))
	}
}

// print out a diagram of what the rope looks like
func Debug(l *list.List) {
	lowestX := l.Front().Value.(image.Point).X
	lowestY := l.Front().Value.(image.Point).Y

	for e := l.Front(); e != nil; e = e.Next() {
		x := e.Value.(image.Point).X
		y := e.Value.(image.Point).Y

		if x < lowestX {
			lowestX = x
		}

		if y < lowestY {
			lowestY = y
		}
	}

	bl := image.Point{lowestX, lowestY}
	length := l.Len()
	i := length - 1

	out := make([]string, length)

	for i := range out {
		out[i] = strings.Repeat(".", length)
	}

	for e := l.Back(); e != nil; e = e.Prev() {
		label := fmt.Sprint(i)
		if i == 0 {
			label = "H"
		}

		val := e.Value.(image.Point).Sub(bl)

		x := val.X
		y := val.Y

		out[y] = out[y][0:x] + label + out[y][x+1:]

		i--
	}

	for _, val := range out {
		fmt.Println(val)
	}
}
