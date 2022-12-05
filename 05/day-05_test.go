package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: "CMZ",
	2: "MCD",
}

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(data)

	if val != expected {
		t.Errorf("Answer should be %q, but got %q", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val := partTwo(data)

	if val != expected {
		t.Errorf("Answer should be %q, but got %q", expected, val)
	}
}
