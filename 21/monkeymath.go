package main

import (
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type monkey struct {
	value              int
	hasValue           bool
	needs              [2]string
	operator, neededBy string
}

type monkeys map[string]monkey

func parseInput(data inType) monkeys {
	monkeys := monkeys{}

	for _, line := range data {
		fields := strings.Fields(line)
		name := fields[0][:len(fields[0])-1]

		monkey := monkey{}

		if len(fields) == 2 {
			// monkey has number
			monkey.hasValue = true
			monkey.value = utils.ParseInt(fields[1])
		} else {
			// monkey depends on other monkeys
			a, op, b := fields[1], fields[2], fields[3]
			monkey.needs = [2]string{a, b}
			monkey.operator = op
		}

		monkeys[name] = monkey
	}

	// link monkeys via neededBy
	for name, monkey := range monkeys {
		if !monkey.hasValue {
			for _, needed := range monkey.needs {
				// cannot assign to struct field ... in map
				mo := monkeys[needed]
				mo.neededBy = name
				monkeys[needed] = mo
			}
		}
	}

	return monkeys
}

func (m *monkeys) getMonkey(a string) int {
	monkey := (*m)[a]

	if monkey.hasValue {
		return monkey.value
	}

	monkeyA := m.getMonkey(monkey.needs[0])
	monkeyB := m.getMonkey(monkey.needs[1])

	var value int

	switch monkey.operator {
	case "+":
		value = monkeyA + monkeyB
	case "-":
		value = monkeyA - monkeyB
	case "*":
		value = monkeyA * monkeyB
	case "/":
		value = monkeyA / monkeyB
	}

	monkey.hasValue = true
	monkey.value = value

	// actually update the map for part 2
	(*m)[a] = monkey

	return value
}

func (m *monkeys) whatToYell(a string) int {
	// move from humn to root, then back, to figure out what each monkey needs to yell
	stack := types.Stack[string]{}

	for a != "root" {
		// can't push &a, because a never changes address;
		// need to create a new local variable
		b := a
		stack.Push(&b)
		monkey := (*m)[a]
		a = monkey.neededBy
	}

	// figure out what number we need
	cur := (*m)["root"]
	// root operator is "=" for part 2
	cur.operator = "="

	// keep track of what number we need to yell for the previous monkey's equation
	var prevNeed int

	for len(stack) > 0 {
		// get next monkey, figure out what it needs to yell
		nextName := *stack.Pop()

		// index of monkey whose value determines what we need to yell
		var index int

		if cur.needs[0] == nextName {
			index = 1
		} else {
			index = 0
		}

		other := (*m)[cur.needs[index]].value

		// number we need the next monkey to yell
		var nextNeed int

		switch cur.operator {
		case "=":
			// we need to equal the other
			nextNeed = other
		case "/":
			// mult or divide depending on whether we're looking for
			// numerator or denominator
			if index == 1 {
				// x / other = prevNeed
				nextNeed = other * prevNeed
			} else {
				// other / x = prevNeed
				nextNeed = other / prevNeed
			}
		case "+":
			nextNeed = prevNeed - other
		case "*":
			nextNeed = prevNeed / other
		case "-":
			if index == 1 {
				// x - other = prevNeed
				nextNeed = prevNeed + other
			} else {
				// other - x = prevNeed
				nextNeed = other - prevNeed
			}
		}

		cur = (*m)[nextName]

		prevNeed = nextNeed
	}

	return prevNeed
}
