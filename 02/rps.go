package main

import (
	"strings"
)

const (
	ROCK = iota + 1
	PAPER
	SCISSORS
)

const (
	LOSS = iota * 3
	DRAW
	WIN
)

var guide = map[string]int{
	"A": ROCK,
	"B": PAPER,
	"C": SCISSORS,
	"X": ROCK,
	"Y": PAPER,
	"Z": SCISSORS,
}

type tournament struct {
	yourScore int
}

func parseGuide(data dataType) tournament {
	yourScore := 0

	for _, round := range data {
		fields := strings.Fields(round)

		opponent, you := guide[fields[0]], guide[fields[1]]

		switch you - opponent {
		case 0:
			yourScore += DRAW + you
		case 1, -2:
			// win
			yourScore += WIN + you
		case -1, 2:
			// loss
			yourScore += LOSS + you
		}
	}

	return tournament{yourScore}
}

// X,Y,Z means you need to Lose, Draw, Win
func parseSuggestiveGuide(data dataType) tournament {
	yourScore := 0

	lose := "X"
	draw := "Y"
	win := "Z"

	for _, round := range data {
		fields := strings.Fields(round)

		opponent, you := guide[fields[0]], fields[1]

		switch you {
		case lose:
			choice := opponent - 1
			if choice == 0 {
				choice = 3
			}
			yourScore += LOSS + choice
		case draw:
			yourScore += DRAW + opponent
		case win:
			choice := opponent + 1
			if choice == 4 {
				choice = 1
			}
			yourScore += WIN + choice
		}
	}

	return tournament{yourScore}
}
