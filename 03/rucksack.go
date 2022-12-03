package main

import (
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
)

// go characters 'a' is 97, but should be 1, and 'A' is 65, and should be 27
func getLetterScore(l rune) int {
	out := int(l) - 96

	if out < 0 {
		out += 58
	}

	return out
}

func parseRucksack(data dataType) int {
	sum := 0

	for _, items := range data {
		half := len(items) / 2

		first, last := items[:half], items[half:]

		for _, letter := range first {
			if strings.ContainsRune(last, letter) {
				score := getLetterScore(letter)
				sum += score
				break
			}
		}
	}

	return sum
}

const GROUP = 3

func parseCommonRucksackItem(data dataType) int {
	sum := 0

outer:
	for i := 0; i < len(data); i += GROUP {
		rucksacks := data[i : i+GROUP]
		first := types.Set[rune]{}
		second := types.Set[rune]{}

		for j, rucksack := range rucksacks {
			for _, letter := range rucksack {
				switch j {
				case 0:
					// add all
					first.Add(letter)
				case 1:
					// add all that match first set
					if first.Has(letter) {
						second.Add(letter)
					}
				case 2:
					// add all that match superset
					if second.Has(letter) {
						// found it
						sum += getLetterScore(letter)
						continue outer
					}
				}
			}
		}
	}

	return sum
}
