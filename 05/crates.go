package main

import (
	"fmt"
	"strings"
	"unicode"
)

// keep []byte from parsing; no need to convert
// everything to string until the end
type crates []byte

func (c *crates) pop() byte {
	old := *c
	n := len(old)
	out := old[n-1]
	*c = old[0 : n-1]

	return out
}

func (c *crates) push(b byte) {
	*c = append(*c, b)
}

type stacks map[int]*crates

type supplies struct {
	stacks    *stacks
	procedure [][3]int // move x from y to z...
}

// check if we have an empty column in the input
var space byte = " "[0]

func parseSupplies(data inType) *supplies {
	stacks := &stacks{}

	drawing, instructions := data[0], data[1]

	// parse drawing

	lines := strings.Split(drawing, "\n")

	// -2 because we omit the stack numbers
	for i := len(lines) - 2; i >= 0; i-- {
		// we can ignore trailing space
		v := strings.TrimRightFunc(lines[i], unicode.IsSpace)

		for j := 0; j < len(v); j += 4 {
			// get letter inside of "[]"
			letter := v[j+1]

			if letter != space {
				stack := j / 4
				// map initialization continues to be annoying
				_, ok := (*stacks)[stack]

				if !ok {
					(*stacks)[stack] = &crates{}
				}

				(*stacks)[stack].push(letter)
			}
		}
	}

	// parse instructions
	lines = strings.Split(instructions, "\n")

	procedure := make([][3]int, len(lines))

	for i, instruction := range lines {
		var move, from, to int
		num, err := fmt.Sscanf(instruction, "move %d from %d to %d", &move, &from, &to)

		if num != 3 || err != nil {
			panic(fmt.Sprintf("could not parse (%d) ints from instruction: %s", num, err))
		}

		procedure[i] = [3]int{move, from, to}
	}

	return &supplies{
		stacks:    stacks,
		procedure: procedure,
	}
}

// moves crates one-by-one, gets top crates in order of stacks
func (s *supplies) runProcedure() string {
	for _, p := range s.procedure {
		move, from, to := p[0], p[1], p[2]

		s.stacks.move(move, from, to)
	}

	return s.stacks.getTop()
}

// moves multiple crates at a time,
// gets top crates in order of stacks
func (s *supplies) runProcedure9001() string {
	for _, p := range s.procedure {
		move, from, to := p[0], p[1], p[2]

		s.stacks.slice(move, from, to)
	}

	return s.stacks.getTop()
}

// crane9000 moves crates one at a time (for loop with .pop())
func (s *stacks) move(num, from, to int) {
	getFrom := (*s)[from-1]
	addTo := (*s)[to-1]

	for i := 0; i < num; i++ {
		addTo.push(getFrom.pop())
	}
}

// crane 9001 moves multiple crates in order
func (s *stacks) slice(num, from, to int) {
	getFrom := (*s)[from-1]
	addTo := (*s)[to-1]

	// move
	l := len(*getFrom)
	*addTo = append(*addTo, (*getFrom)[l-num:l]...)
	// remove
	*getFrom = (*getFrom)[:l-num]
}

// concatenates the strings associated with the top crates
// in the stacks, in order
func (s *stacks) getTop() (out string) {
	length := len(*s)

	for i := 0; i < length; i++ {
		v := *(*s)[i]
		out += string(v[len(v)-1])
	}

	return
}

// String representation (bytes to string)
func (c *crates) String() (out string) {
	out += "[ "
	for _, a := range *c {
		out += fmt.Sprint(string(a), " ")
	}
	out += "] "

	return
}

// String representation (bytes to string)
func (s stacks) String() (out string) {
	length := len(s)
	for i := 0; i < length; i++ {
		v := s[i]
		out += fmt.Sprintf("%d.", i+1)
		out += v.String()
	}

	return
}
