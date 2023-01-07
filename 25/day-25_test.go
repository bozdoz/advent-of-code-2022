package main

import (
	"fmt"
	"testing"
)

var runs = map[string]int{
	"2=-01":         976,
	"1":             1,
	"2":             2,
	"1=":            3,
	"1-":            4,
	"10":            5,
	"11":            6,
	"12":            7,
	"2=":            8,
	"2-":            9,
	"20":            10,
	"1=0":           15,
	"1-0":           20,
	"1=11-2":        2022,
	"1-0---0":       12345,
	"1121-1110-1=0": 314159265,
}

func TestSnafuToDecimal(t *testing.T) {
	for input, want := range runs {
		t.Run(fmt.Sprintf("%v = %v", input, want), func(t *testing.T) {
			got := snafuToDecimal(input)

			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestDecimalToSnafu(t *testing.T) {
	for want, input := range runs {
		t.Run(fmt.Sprintf("%v = %v", input, want), func(t *testing.T) {
			got := decimalToSnafu(input)

			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: "2=-1=0",
	2: "",
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
