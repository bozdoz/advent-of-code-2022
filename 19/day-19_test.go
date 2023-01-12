package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 33,
	2: 0,
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

	// state &{18 7 map[1:13 2:1] map[1:1 2:1 4:1]}
	state := state{
		time:      18,
		robotbits: 7,
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

	next := state.getNextStates(blueprints[0], 24)

	for _, st := range next {
		if st.time <= 18 {
			t.Errorf("next states must increment time! Got %v from %v", st, state)
		}
	}
}

func TestGetNextStatesTwo(t *testing.T) {
	blueprints := parseInput(data)

	// state &{7 3 map[1:1 2:6] map[1:1 2:3]}
	state := state{
		time:      7,
		robotbits: 3,
		resources: map[resource]int{
			1: 1,
			2: 6,
		},
		robots: map[resource]int{
			1: 1,
			2: 3,
		},
	}

	next := state.getNextStates(blueprints[0], 24)

	// first obsidian should be minute 11, not 10
	for _, st := range next {
		if st.robotbits&int(obsidian) == int(obsidian) {
			if st.time != 11 {
				t.Errorf("Why on earth didn't we get an obsidian at minute 11? %v", st)
			}
			break
		}
	}
}

func TestTerminalVelociy(t *testing.T) {
	end := 10
	st := state{
		time: 8,
		resources: map[resource]int{
			geode: 0,
		},
		robots: map[resource]int{
			geode: 0,
		},
	}

	assert := func(t testing.TB, bots, geodes int) {
		t.Helper()

		terminal := terminalVelocity(st, end)
		got := terminal.robots[geode]

		if got != bots {
			t.Errorf("expected bots: %v got: %v", bots, got)
		}

		got = terminal.resources[geode]

		if got != geodes {
			t.Errorf("expected geodes: %v got: %v", geodes, got)
		}
	}

	assert(t, 2, 1)

	st.time = 7

	assert(t, 3, 3)

	st.resources[geode] = 5

	assert(t, 3, 3+5)

	st.robots[geode] = 5

	assert(t, 3+5, 3+5+(3*5))

	terminal := terminalVelocity(st, end)

	got := terminal.time
	want := end

	if got != want {
		t.Errorf("expected time: %v got: %v", want, got)
	}
}
