package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 26,
	2: 56000011,
}

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	partOneRow = 10

	expected := answers[1]

	val := partOne(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	partTwoMax = 20

	expected := answers[2]

	val := partTwo(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}
