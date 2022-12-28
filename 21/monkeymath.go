package main

import (
	"fmt"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

type monkey struct {
	value    int
	hasValue bool
	needs    [2]string
	operator string
	neededBy []string // doubly linked list?
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
				mo := monkeys[needed]
				mo.neededBy = append(mo.neededBy, name)
			}
		}
	}

	fmt.Println(monkeys)

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

	return value
}
