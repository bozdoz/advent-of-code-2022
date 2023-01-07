package utils

import "image"

// convert int into a signed number
// ex: 4 -> 1; -4 -> -1; 0 -> 0
func getSignInt(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	default:
		return 0
	}
}

// convert vector into a sign vector
// ex: {4,0} -> {1,0} and {0,-4} -> {0,-1}
func GetSignPoint(point image.Point) image.Point {
	x := getSignInt(point.X)
	y := getSignInt(point.Y)

	return image.Point{x, y}
}

func ManhattanDistance(a, b image.Point) int {
	return Abs(a.X-b.X) + Abs(a.Y-b.Y)
}
