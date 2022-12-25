package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 3068,
	2: 1514285714288,
}

var data = fileReader("example.txt")

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

func TestFindRepeating(t *testing.T) {
	repeating := []int{1, 2, 9, 21, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	got := findRepeatingPatternFromEnd(repeating)

	expect := 8

	if got != expect {
		t.Errorf("expected %v, got %v", expect, got)
	}
}
