package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 33,
	2: 56 * 62,
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

func TestGetNextStatesOne(t *testing.T) {
	blueprints := parseInput(data)

	// state &{18 map[1:13 2:1] map[1:1 2:1 4:1]}
	state := state{
		time: 18,
		resources: map[resource]int{
			1: 13,
			2: 1,
		},
		robots: map[resource]int{
			1: 1,
			2: 1,
			4: 1,
		},
	}

	next := blueprints[0].getNextStates(state, 24)

	for _, st := range next {
		if st.time <= 18 {
			t.Errorf("next states must increment time! Got %v from %v", st, state)
		}
	}
}

func TestGetNextStatesTwo(t *testing.T) {
	blueprints := parseInput(data)

	// state &{7 map[1:1 2:6] map[1:1 2:3]}
	state := state{
		time: 7,
		resources: map[resource]int{
			1: 1,
			2: 6,
		},
		robots: map[resource]int{
			1: 1,
			2: 3,
		},
	}

	next := blueprints[0].getNextStates(state, 24)

	// first obsidian should be minute 11, not 10
	for _, st := range next {
		if st.robots[obsidian] == 1 {
			if st.time != 11 {
				t.Errorf("Why on earth didn't we get an obsidian at minute 11? %v", st)
			}
			break
		}
	}
}
