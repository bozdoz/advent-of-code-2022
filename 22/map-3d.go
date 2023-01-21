package main

import (
	"fmt"
	"image"
	"math"

	"github.com/bozdoz/advent-of-code-2022/types"
)

/**
example:
--#
###
--##

input:
-##
-#
##
#
*/

type neighbour struct {
	sq  *square
	dir direction
}

type square struct {
	n, e, w, s      *neighbour
	rect            image.Rectangle
	neighboursFound int // used to increment in cubeFold
	neighbourRects  types.Set[image.Rectangle]
}

func (a *square) hasNeighbourInDir(dir direction) bool {
	switch dir {
	case up:
		return a.n != nil
	case right:
		return a.e != nil
	case left:
		return a.w != nil
	case down:
		return a.s != nil
	}

	return false
}

func (a *square) setNeighbour(b *square, fromDir, toDir direction) {
	switch fromDir {
	case up:
		a.n = &neighbour{b, toDir}
	case right:
		a.e = &neighbour{b, toDir}
	case left:
		a.w = &neighbour{b, toDir}
	case down:
		a.s = &neighbour{b, toDir}
	}
}

func (dir direction) inverse() direction {
	return direction((dir + 2) % 4)
}

// add b as a neighbour to a in direction dir
// and vice-versa with the inverse direction
func (a *square) addNeighbour(b *square, fromDir, toDir direction) (added bool) {
	// a rect is already joined to b by some other direction
	if a.neighbourRects.Has(b.rect) {
		return false
	}

	inverseToDir := toDir.inverse()
	inverseFromDir := fromDir.inverse()

	// if either square already have a neighbour assigned, don't add
	if a.hasNeighbourInDir(fromDir) || b.hasNeighbourInDir(inverseToDir) {
		return false
	}

	switch fromDir {
	case up:
		a.n = &neighbour{b, toDir}
		b.setNeighbour(a, inverseToDir, inverseFromDir)
	case down:
		a.s = &neighbour{b, toDir}
		b.setNeighbour(a, inverseToDir, inverseFromDir)
	case left:
		a.w = &neighbour{b, toDir}
		b.setNeighbour(a, inverseToDir, inverseFromDir)
	case right:
		a.e = &neighbour{b, toDir}
		b.setNeighbour(a, inverseToDir, inverseFromDir)
	}

	a.neighbourRects.Add(b.rect)
	b.neighbourRects.Add(a.rect)

	a.neighboursFound++
	b.neighboursFound++

	return true
}

// overwritten in tests
var squareSize = 50

// turn the 2d map into 6 cube faces of {size}x{size}
func cubeFold(board board) map[image.Rectangle]*square {
	// iterate board, save squares
	// label neighbours
	squares := map[image.Rectangle]*square{}

	for r := 0; r <= board.height; r += squareSize {
		for c := 0; c <= board.width; c += squareSize {
			cell := board.get(r, c)
			if cell == open || cell == wall {
				min := image.Point{r, c}
				max := image.Point{r + squareSize, c + squareSize}
				rect := image.Rectangle{min, max}
				squares[rect] = &square{
					rect:           rect,
					neighbourRects: types.Set[image.Rectangle]{},
				}
			}
		}
	}

	// right, down, left, up (like direction type)
	neighbours := [4]image.Point{
		{0, squareSize},
		{squareSize, 0},
		{0, -squareSize},
		{-squareSize, 0},
	}

	// lazy naming for path-finding priority queue
	// first time internal type definition
	type state struct {
		square         *square
		fromDir, toDir direction
		currentRect    image.Rectangle
		priority       int
	}

	pq := make(types.PriorityQueue[state], len(squares)*len(neighbours))

	// populate priority queue
	index := 0
	for rect, sq := range squares {
		for i, vec := range neighbours {
			dir := direction(i)
			// this is the square that is adjacent to `rect`
			cur := rect.Add(vec)

			pq.NewItem(&state{
				square:      sq,
				fromDir:     dir,
				toDir:       dir,
				currentRect: cur,
				priority:    0,
			}, 0, index)
			index++
		}
	}

	pq.Init()

	neighboursFound := 0
	// 6 faces have 4 neighbours each
	// we can break after 6*4 neighbours found
	neighboursTotal := 6 * 4

	for neighboursFound != neighboursTotal {
		cur := pq.Get()

		if cur.square.neighboursFound == 4 {
			continue
		}

		sq, isSquare := squares[cur.currentRect]

		if sq == cur.square {
			// thou shalt not be neighbours with thyself
			continue
		}

		if isSquare {
			// found something!
			added := cur.square.addNeighbour(sq, cur.fromDir, cur.toDir)

			if added {
				neighboursFound += 2
			}
			continue
		}

		// else, continue path-finding
		for i, vec := range neighbours {
			dir := direction(i)
			nextRect := cur.currentRect.Add(vec)

			pq.PushValue(&state{
				square:      cur.square,
				fromDir:     cur.fromDir,
				toDir:       dir,
				currentRect: nextRect,
				priority:    cur.priority + 1,
			}, cur.priority+1)
		}
	}

	return squares
}

// move around a cube by an amount (dir can change)
func (b board) move3d(cur [2]int, dir direction, steps int) ([2]int, direction) {
	i := 0

	for i < steps {
		next, hitWall, newDir := b.getNext3d(cur[0], cur[1], dir)

		if hitWall {
			break
		}

		cur = next
		dir = newDir

		i++
	}

	return cur, dir
}

// move around a cube, somehow (dir can change)
func (b *board) getNext3d(r, c int, dir direction) (next [2]int, hitWall bool, newDir direction) {
	orig := [2]int{r, c}

	// assume we'll keep going this direction
	newDir = dir

	next = orig
	moveOne(&next, dir)

	// for loop because the `empty` case may hit a wall
	for {
		var tile tile
		// check if out of bounds
		if next[0] < 0 || next[1] < 0 || next[0] >= b.height || next[1] >= b.width {
			// next cell is empty
			tile = empty
		} else {
			// we can get next cell
			tile = b.get(next[0], next[1])
		}

		switch tile {
		case open:
			return
		case wall:
			// return original (didn't move)
			return orig, true, dir
		case empty:
			// move to new cube face in another direction
			// where were we, and what direction are we heading?
			// using orig, because we need to know what square this point is in
			next, newDir = b.rotateTile(orig, dir)
		default:
			return
		}
	}
}

// moving from one cube face to another requires changing directions
// and rotating tiles from one square to another
func (b board) rotateTile(pos [2]int, dir direction) (next [2]int, newDir direction) {
	curPoint := image.Point{pos[0], pos[1]}

	for rect, sq := range b.squares {
		if curPoint.In(rect) {
			var neighbour *neighbour
			switch dir {
			case up:
				neighbour = sq.n
			case down:
				neighbour = sq.s
			case left:
				neighbour = sq.w
			case right:
				neighbour = sq.e
			}
			newDir = neighbour.dir

			// figure out the tile we step on
			// rotate current position until current dir lines up with new dir
			// directions are clockwise, so + means clockwise and - means counter
			rotations := int(newDir - dir)

			size := rect.Size()

			// rotate a point around an origin algorithm
			px, py := float64(curPoint.X), float64(curPoint.Y)
			origin := rect.Min.Add(size.Div(2))
			ox, oy := float64(origin.X)-0.5, float64(origin.Y)-0.5
			// angle in radians 🤦‍♂️
			// inverse to go clockwiise
			theta := -(math.Pi / 2) * float64(rotations)
			rx := math.Cos(theta)*(px-ox) - math.Sin(theta)*(py-oy) + ox
			ry := math.Sin(theta)*(px-ox) + math.Cos(theta)*(py-oy) + oy

			// move in new dir again
			next = [2]int{int(rx), int(ry)}
			moveOne(&next, newDir)

			// annoying transitions here
			curPoint = image.Point{next[0], next[1]}

			// mod to next cube face
			curPoint = curPoint.Mod(neighbour.sq.rect)

			// set next for outer for loop to check for walls
			next = [2]int{curPoint.X, curPoint.Y}

			return
		}
	}

	return
}

//
// String representations
//

func (sq *square) String() (out string) {
	out += fmt.Sprintf(" rect: %v ", sq.rect)
	out += fmt.Sprintf("\tneighbours: %v", sq.neighboursFound)
	out += fmt.Sprintf("\n\t\tn: %v \n\t\te: %v \n\t\tw: %v \n\t\ts: %v \n", sq.n, sq.e, sq.w, sq.s)

	return
}

func (a neighbour) String() (out string) {
	if a.sq == nil {
		out += "[nil]"
	} else {
		out += fmt.Sprintf("{ %v ", a.sq.rect)
		out += fmt.Sprintf("dir: %v }", a.dir)
	}
	return
}
