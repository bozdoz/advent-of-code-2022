package main

import (
	"fmt"
	"math"
)

var snafuMap = map[rune]int{
	'2': 2,
	'1': 1,
	'0': 0,
	'-': -1,
	'=': -2,
}

func snafuToDecimal(snafu string) (dec int) {
	l := len(snafu)

	for _, char := range snafu {
		l--
		mapped := snafuMap[char]
		pow := int(math.Pow(5, float64(l)))

		dec += mapped * pow
	}

	return
}

// got this from https://meteorconverter.com/conversions/number-bases/10-to-5?input=600
func decimalToSnafu(dec int) (snafu string) {
	remainders := []int{}

	for dec > 0 {
		quotient := dec / 5
		remaining := dec % 5
		remainders = append(remainders, remaining)
		dec = quotient
	}

	out := []string{}
	carried := 0

	for _, val := range remainders {
		// 3 == 5-2 = "1="
		// 4 == 5-1 = "1-"
		withCarried := val + carried

		if withCarried > 2 {
			carried = 1
		} else {
			carried = 0
		}

		switch withCarried {
		case 3:
			out = append(out, "=")
		case 4:
			out = append(out, "-")
		default:
			// sometimes this could be 5, and should be 0
			out = append(out, fmt.Sprint(withCarried%5))
		}
	}

	if carried != 0 {
		out = append(out, "1")
	}

	// reverse
	for i := len(out) - 1; i >= 0; i-- {
		snafu += out[i]
	}

	return
}
