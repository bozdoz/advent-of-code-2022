package types

import "fmt"

type Stack[T any] []*T

func (stack *Stack[T]) Push(item *T) {
	*stack = append(*stack, item)
}

func (stack *Stack[T]) Pop() (item *T) {
	n := len(*stack) - 1
	last := (*stack)[n]

	*stack = (*stack)[:n]

	return last
}

func (stack *Stack[T]) Peek() (item *T) {
	n := len(*stack) - 1
	if n < 0 {
		return nil
	}

	return (*stack)[n]
}

func (stack Stack[T]) String() (out string) {
	out += "[ "
	for _, val := range stack {
		out += fmt.Sprint(*val, " ")
	}
	out += "]"

	return
}
