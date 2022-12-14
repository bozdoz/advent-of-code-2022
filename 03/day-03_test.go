package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 157,
	2: 70,
}

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(dataType(data))

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val := partTwo(dataType(data))

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}
