package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 24000,
	2: 45000,
}

var data, _ = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val, err := partOne(dataType(data))

	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val, err := partTwo(dataType(data))

	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}
