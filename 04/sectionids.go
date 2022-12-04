package main

import (
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

type pair struct {
	min, max int
}

func parsePair(data string) pair {
	nums := strings.Split(data, "-")
	min, max := utils.ParseInt(nums[0]), utils.ParseInt(nums[1])

	return pair{min, max}
}

func parsePairs(data string) (pair, pair) {
	pairs := strings.Split(data, ",")

	return parsePair(pairs[0]), parsePair(pairs[1])
}

func (p *pair) contains(a *pair) bool {
	return p.min <= a.min && p.max >= a.max
}

func (p *pair) overlaps(a *pair) bool {
	return !(p.max < a.min || p.min > a.max)
}

func parseSectionContains(data dataType) (sum int) {
	for _, item := range data {
		first, second := parsePairs(item)

		if first.contains(&second) || second.contains(&first) {
			sum++
		}
	}

	return
}

func parseSectionOverlaps(data dataType) (sum int) {
	for _, item := range data {
		first, second := parsePairs(item)

		if first.overlaps(&second) || second.overlaps(&first) {
			sum++
		}
	}

	return
}
