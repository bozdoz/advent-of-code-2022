package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

type Toperator int

const (
	PLUS Toperator = iota
	MULTIPLY
	EXPONENT
)

type Tmonkey struct {
	items                            []int
	operator                         Toperator
	operationNum, testDiv, inspected int
	throwTo                          [2]int // [true, false]
}

type Tmonkeys []*Tmonkey

const (
	STARTING = iota
	OPERATION
	TEST
	IFTRUE
)

// get the common denominator between the divisors
var divisor int = 1

func parseInput(data inType) *Tmonkeys {
	// reset divisor
	divisor = 1

	monkeys := make(Tmonkeys, len(data))

	for i, desc := range data {
		monkey := Tmonkey{}

		lines := strings.Split(desc, "\n")

		// ignore "Monkey 0:" line, and last line
		for j, line := range lines[1 : len(lines)-1] {
			switch j {
			case STARTING:
				list := line[len("  starting items: "):]
				items := strings.Split(list, ", ")
				ints := make([]int, len(items))

				for k, item := range items {
					ints[k] = utils.ParseInt(item)
				}

				monkey.items = ints
			case OPERATION:
				op := line[len("  Operation: new = "):]

				if op == "old * old" {
					monkey.operator = EXPONENT
					monkey.operationNum = 2
				} else {
					fields := strings.Fields(op)

					operator, num := fields[1], utils.ParseInt(fields[2])

					switch operator {
					case "*":
						monkey.operator = MULTIPLY
					case "+":
						monkey.operator = PLUS
					}

					monkey.operationNum = num
				}
			case TEST:
				num := line[len("  Test: divisible by "):]
				monkey.testDiv = utils.ParseInt(num)
				// update divisor for modulus
				divisor *= monkey.testDiv
			case IFTRUE:
				// get both true and false lines here
				testMonkeys := [2]int{}
				// feeling lazy here...
				tests := []string{line, lines[len(lines)-1]}
				for k := 0; k < 2; k++ {
					fields := strings.Fields(tests[k])
					num := utils.ParseInt(fields[len(fields)-1])
					testMonkeys[k] = num
				}

				monkey.throwTo = testMonkeys
			}
		}

		monkeys[i] = &monkey
	}

	return &monkeys
}

// each monkey takes a round in inspecting
func (monkeys *Tmonkeys) inspect(part int) {
	for i := range *monkeys {
		// `for i, monkey` is a copy of `monkey`
		monkey := (*monkeys)[i]
		monkey.inspected += len(monkey.items)

		// inspect
		for _, val := range monkey.items {
			switch monkey.operator {
			case PLUS:
				val += monkey.operationNum
			case MULTIPLY:
				val *= monkey.operationNum
			case EXPONENT:
				val = int(math.Pow(float64(val), float64(monkey.operationNum)))
			}

			if part == 1 {
				val /= 3
			} else {
				val %= divisor
			}

			isDivisible := val%monkey.testDiv == 0

			var next int
			if isDivisible {
				next = monkey.throwTo[0]
			} else {
				next = monkey.throwTo[1]
			}
			// throw item to next monkey
			(*monkeys)[next].items = append((*monkeys)[next].items, val)
		}
		// wipe
		monkey.items = []int{}
	}
}

func (m *Tmonkeys) String() (out string) {
	for i, m := range *m {
		out += fmt.Sprintf("Monkey %d: %s\n", i, m.String())
	}
	return
}

func (m *Tmonkey) String() (out string) {
	out += fmt.Sprint(m.items, "\n")
	out += fmt.Sprint("  inspected: ", m.inspected, "\n")
	return
}
