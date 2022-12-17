package main

import (
	"image"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type obstacle = types.Set[image.Point]

type area struct {
	obstacle obstacle
	bottom   int
}

func parseInput(data inType) area {
	area := area{
		obstacle: obstacle{},
	}

	for _, line := range data {
		points := []image.Point{}

		for _, point := range strings.Split(line, " -> ") {
			parts := strings.Split(point, ",")
			x, y := utils.ParseInt(parts[0]), utils.ParseInt(parts[1])

			points = append(points, image.Point{x, y})

			if y > area.bottom {
				area.bottom = y
			}
		}

		area.addRocks(points)
	}

	return area
}

func (area *area) addRocks(points []image.Point) {
	r := area.obstacle

	for i := 0; i < len(points)-1; i++ {
		a := points[i]
		b := points[i+1]

		diff := b.Sub(a)
		step := utils.GetSignPoint(diff)

		for a != b {
			r.Add(a)
			a = a.Add(step)
		}
		r.Add(b)
	}
}

// y increases
var acc = image.Point{0, 1}
var diagLeft = image.Point{-1, 1}
var diagRight = image.Point{1, 1}

// drops from 500,0, one at a time,
// until sand falls beyond area.bottom
func (area *area) dropSand() *image.Point {
	// part 2 made me change 500,0 to 500,-1
	newSand := image.Point{500, -1}

	for {
		next := newSand.Add(acc)

		if next.Y > area.bottom {
			// area is full
			return nil
		}

		if area.obstacle.Has(next) {
			// collision
			// A unit of sand always falls down one step if possible
			// diagonally left, then diagonally right
			maybeLeft := newSand.Add(diagLeft)

			if area.obstacle.Has(maybeLeft) {
				// try diagonally right
				maybeRight := newSand.Add(diagRight)

				if area.obstacle.Has(maybeRight) {
					// left and right blocked
					// sand stand still
					area.obstacle.Add(newSand)
					break
				} else {
					// move right
					newSand = maybeRight
				}
			} else {
				// move left
				newSand = maybeLeft
			}
		} else {
			// keep dropping
			newSand = next
		}
	}

	// return last dropped sand to check if it's full
	return &newSand
}
