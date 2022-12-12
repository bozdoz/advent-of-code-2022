package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 13140,
}

// TODO: this passes, but not in parallel with TestExampleOne
// func TestProgram(t *testing.T) {
// 	data = []string{
// 		"noop",
// 		"addx 3",
// 		"addx -5",
// 	}

// 	program := parseInput(data)

// 	if program.x != -1 {
// 		t.Errorf("Expected %v, got %v", -1, program.x)
// 	}
// }

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}
