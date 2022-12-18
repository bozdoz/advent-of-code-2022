package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 24000,
	2: 45000,
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

func TestExampleOneHeap(t *testing.T) {
	expected := answers[1]

	val := parseGroupedCalorieHeap(data, 1)

	if val[0] != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val[0])
	}
}

func TestExampleTwoHeap(t *testing.T) {
	expected := answers[2]

	group := parseGroupedCalorieHeap(data, 3)

	val := group[0] + group[1] + group[2]

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

var exampleTwo = fileReader("example-large.txt")

func BenchmarkSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGroupedCalorieList(exampleTwo, 3)
	}
}

func BenchmarkHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGroupedCalorieHeap(exampleTwo, 3)
	}
}
