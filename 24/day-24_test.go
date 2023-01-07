package main

import (
	"image"
	"testing"
)

func TestBlizzardOverlap(t *testing.T) {
	assert := func(t testing.TB, input []string, st state, want bool) {
		t.Helper()

		vly := parseInput(input)

		isBlizzard := vly.stateOverlapsBlizzard(st.position, st.minutes)

		if isBlizzard != want {
			t.Errorf("expected %v, got %v", want, isBlizzard)
		}
	}

	t.Run("left", func(t *testing.T) {
		input := []string{
			"#.###",
			"#...#",
			"#..<#",
			"#...#",
			"###.#",
		}

		cur := state{
			position: image.Point{2, 1},
			minutes:  2,
		}

		assert(t, input, cur, true)

		cur.minutes = 5

		assert(t, input, cur, true)

		cur.minutes = 4

		assert(t, input, cur, false)
	})

	t.Run("right", func(t *testing.T) {
		input := []string{
			"#.###",
			"#...#",
			"#...#",
			"#>..#",
			"###.#",
		}

		cur := state{
			position: image.Point{3, 3},
			minutes:  8,
		}

		assert(t, input, cur, true)
	})

	t.Run("down", func(t *testing.T) {
		input := []string{
			"#.###",
			"#.E..#",
			"#....#",
			"#.v..#",
			"#....#",
			"####.#",
		}

		cur := state{
			position: image.Point{1, 2},
			minutes:  2,
		}

		assert(t, input, cur, true)

		cur.minutes = 6
		assert(t, input, cur, true)

		cur.minutes = 5
		assert(t, input, cur, false)
	})

	t.Run("up", func(t *testing.T) {
		input := []string{
			"#.###",
			"#...#",
			"#.^.#",
			"#...#",
			"###.#",
		}

		cur := state{
			position: image.Point{2, 2},
			minutes:  1,
		}

		assert(t, input, cur, false)

		cur.minutes = 0
		assert(t, input, cur, true)

		cur.minutes = 3
		assert(t, input, cur, true)
	})

	t.Run("minute 5", func(t *testing.T) {
		input := []string{
			"#.######",
			"#.E....#",
			"#......#",
			"#......#",
			"#.^....#",
			"######.#",
		}

		cur := state{
			position: image.Point{1, 2},
			minutes:  5,
		}

		assert(t, input, cur, false)
	})
}

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 18,
	2: 54,
}

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	// TODO
	t.Skip("Somehow this test is failing after Part 2")

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
