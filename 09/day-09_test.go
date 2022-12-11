package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 13,
	2: 1,
	3: 36,
}

var data1 = fileReader("example1.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(data1)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val := partTwo(data1)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

var data2 = fileReader("example2.txt")

func TestExampleThree(t *testing.T) {
	expected := answers[3]

	val := partTwo(data2)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}
