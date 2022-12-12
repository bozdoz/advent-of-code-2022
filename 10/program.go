package main

import (
	"fmt"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

type Tprogram struct {
	cycle, x, signal_sum int
	crt                  []string
}

// interval of signal strength assessment
const signal_mod = 40

// when to check for strength
const signal_cycle = 20

func parseInput(data inType) *Tprogram {
	signal_sum := 0
	x := 1
	cycle := 0
	crt := make([]string, 240)

	inc_cycle := func() {
		pixel := cycle % signal_mod
		// draw crt if 3-width pixel overlaps
		if pixel >= x-1 && pixel <= x+1 {
			crt[cycle] = "█"
		} else {
			crt[cycle] = " "
		}

		cycle++
		if cycle%signal_mod == signal_cycle {
			signal_sum += cycle * x
		}
	}

	for _, line := range data {
		inc_cycle()
		if line == "noop" {
			continue
		}
		addx := utils.ParseInt(line[len("addx "):])
		// add waits an extra cycle
		inc_cycle()

		x += addx
	}

	return &Tprogram{
		x:          x,
		cycle:      cycle,
		signal_sum: signal_sum,
		crt:        crt,
	}
}

// prints crt
func (p *Tprogram) print() {
	for i := 40; i <= 240; i += 40 {
		fmt.Println(strings.Join(p.crt[i-40:i], ""))
	}
}
