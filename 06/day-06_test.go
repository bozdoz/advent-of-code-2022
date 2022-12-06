package main

import (
	"fmt"
	"testing"
)

func TestExampleOne(t *testing.T) {
	runs := map[string]int{
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      5,
		"nppdvjthqldpwncqszvftbrmjlhg":      6,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 10,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  11,
	}

	for k, v := range runs {
		t.Run(fmt.Sprintf("%q should be %d", k, v), func(t *testing.T) {
			val := uniqueLettersIndex(k, 4)

			if val != v {
				t.Errorf("Answer should be %v, but got %v", v, val)
			}
		})
	}
}

func TestExampleTwo(t *testing.T) {
	runs := map[string]int{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb":    19,
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      23,
		"nppdvjthqldpwncqszvftbrmjlhg":      23,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 29,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  26,
	}

	for k, v := range runs {
		t.Run(fmt.Sprintf("%q should be %d", k, v), func(t *testing.T) {
			val := uniqueLettersIndex(k, 14)

			if val != v {
				t.Errorf("Answer should be %v, but got %v", v, val)
			}
		})
	}
}
