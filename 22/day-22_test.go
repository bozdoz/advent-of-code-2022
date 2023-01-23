package main

import (
	"testing"
)

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 6032,
	2: 5031,
}

var data = fileReader("example.txt")

func init() {
	// example square size is 4
	squareSize = 4
}

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val := partTwo(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestCubeFold(t *testing.T) {
	board := parseInput(data, 4)

	board.squares = cubeFold(board)

	rotate := func(t testing.TB, cur [2]int, dir direction, wantPos [2]int, wantDir direction) {
		next, dir := board.rotateTile(cur, dir)

		if dir != wantDir {
			t.Errorf("want direction %v, but got %v", wantDir, dir)
		}

		if next != wantPos {
			t.Errorf("wanted %v, got: %v", wantPos, next)
		}
	}

	rotate(t, [2]int{10, 15}, right, [2]int{1, 11}, left)

	rotate(t, [2]int{4, 5}, up, [2]int{1, 8}, right)
}

// debug actual data
func TestCubeFoldFifty(t *testing.T) {
	t.Skip("Don't want to commit input files to repo")

	var actual = fileReader("input.txt")

	board := parseInput(actual, 50)

	board.squares = cubeFold(board)

	rotate := func(t testing.TB, cur [2]int, dir direction, wantPos [2]int, wantDir direction) {
		next, dir := board.rotateTile(cur, dir)

		if dir != wantDir {
			t.Errorf("want direction %v, but got %v", wantDir, dir)
		}

		if next != wantPos {
			t.Errorf("wanted %v, got: %v", wantPos, next)
		}
	}

	// getting a rounding error on these
	rotate(t, [2]int{0, 59}, up, [2]int{159, 0}, right)
	rotate(t, [2]int{159, 0}, left, [2]int{0, 59}, down)
}
