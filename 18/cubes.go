package main

import (
	"fmt"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

// cube, aka lava, is a type alias?
type cube = types.Vector3d

// can't have nested structures without explicit field names; so this function helps
var NewCube = types.NewVector3d

func parseInput(data inType) map[cube]struct{} {
	cubes := map[cube]struct{}{}

	for _, line := range data {
		var x, y, z int
		n, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if n != 3 || err != nil {
			panic(fmt.Sprint("Got:", n, "You can't format this line!", line, err))
		}
		cubes[NewCube(x, y, z)] = struct{}{}
	}

	return cubes
}

func touchingSides(cubes map[cube]struct{}) (touching int) {
	for cube := range cubes {
		for other := range cubes {
			if cube == other {
				continue
			}
			distance := cube.DistanceTo(other)

			if distance == 1 {
				touching++
			}
		}
	}

	return
}

// directions to iterate when checking neighbours
var cubeNeighbours = []cube{
	NewCube(1, 0, 0),
	NewCube(-1, 0, 0),
	NewCube(0, 1, 0),
	NewCube(0, -1, 0),
	NewCube(0, 0, 1),
	NewCube(0, 0, -1),
}

// more recursive closures
var visit func(cur cube)

// get the number of faces of lava cubes that are on the outside
func floodFill(cubes map[cube]struct{}) (faces int) {
	// need to create a cube that will envelope all cubes

	// funny way to get first from a map
	var first cube
	for cube := range cubes {
		first = cube
		break
	}

	// initialize min/max
	minX := first.X
	maxX := minX
	minY := first.Y
	maxY := minY
	minZ := first.Z
	maxZ := minZ

	// iterate to get size
	for cube := range cubes {
		minX = utils.Min(cube.X, minX)
		minY = utils.Min(cube.Y, minY)
		minZ = utils.Min(cube.Z, minZ)
		maxX = utils.Max(cube.X, maxX)
		maxY = utils.Max(cube.Y, maxY)
		maxZ = utils.Max(cube.Z, maxZ)
	}

	// make a cube around the perimeter, buffering by 1
	minX--
	minY--
	minZ--
	maxX++
	maxY++
	maxZ++

	visited := types.Set[cube]{}

	// start recursive function
	visit = func(cur cube) {
		// if lava cube, then we increment the face that was touched
		_, isLava := cubes[cur]
		if isLava {
			faces++
			return
		}

		// if visited, then exit
		if visited.Has(cur) {
			return
		}

		visited.Add(cur)

		// iterate 6 neighbours max
		for _, dir := range cubeNeighbours {
			next := cur.Add(dir)

			// check that it's not out of bounds
			if next.X < minX || next.X > maxX || next.Y < minY || next.Y > maxY || next.Z < minZ || next.Z > maxZ {
				continue
			}

			visit(next)
		}
	}

	// start at any corner, and visit neighbours
	visit(NewCube(minX, minY, minZ))

	return
}
