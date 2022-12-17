package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bozdoz/advent-of-code-2022/types"
)

type item any

type pair struct {
	left, right []item
}

type pairs []pair

// part one checks pairs
func parsePairs(data inType) pairs {
	pairs := make(pairs, len(data))

	for i, group := range data {
		items := strings.Split(group, "\n")
		one, two := items[0], items[1]

		pairs[i] = pair{
			left:  parseItem(one),
			right: parseItem(two),
		}
	}

	return pairs
}

// parse a single line from the pair input data
func parseItem(data string) []item {
	// we're manually decoding the arrays
	dec := json.NewDecoder(strings.NewReader(data))

	// keep each new array in a stack
	stack := types.Stack[[]item]{}

	// current pointer to current array (initially nil)
	var cur *[]item

	for {
		t, err := dec.Token()

		if err != nil {
			panic(err)
		}

		switch t {
		case json.Delim('['):
			if cur == nil {
				cur = &[]item{}
			} else {
				// copy cur to local var, because cur's pointer can't change
				archive := append([]item{}, (*cur)...)
				stack.Push(&archive)
				// wipe cur (truncate slice)
				*cur = (*cur)[0:0]
			}
		case json.Delim(']'):
			if len(stack) == 0 {
				// done
				return *cur
			}

			prev := *stack.Pop()

			// updating/appending cur to prev is painful
			*cur = append([]item{}, append(prev, *cur)...)

		default:
			// json decode is always float64
			num, ok := t.(float64)

			if !ok {
				panic(fmt.Sprintf("What!, not a float64!? %T %v", t, t))
			}

			// append int
			*cur = append(*cur, int(num))
		}
	}
}

// compare two lists (left vs right);
// returns -1 if l < r, 0 if l == r, and 1 if l > r
func compareLists(l, r []item) int {
	for i := 0; i < len(l); i++ {
		if i == len(r) {
			// left has more than right
			return 1
		}

		switch compare(l[i], r[i]) {
		case -1: // l < r
			return -1
		case 0: // l == r
			continue
		case 1: // l > r
			return 1
		}
	}

	if len(l) < len(r) {
		return -1
	}

	return 0
}

// whether left is less than or equal to right-side
func (pair *pair) isOrdered() bool {
	return compareLists(pair.left, pair.right) < 1
}

// compare either []any or int, in any combination
// returns -1 if l < r, 0 if l == r, and 1 if l > r
func compare(a, b item) int {
	log.Printf("-- compare %v and %v\n", a, b)

	// catch panics just in case the type assertions are incorrect
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic: (a) %[1]T %[1]v (b) %[2]T %[2]v\n", a, b)
		}
	}()

	// figure out types
	aAsInt, aok := a.(int)
	bAsInt, bok := b.(int)

	// quality variable names here
	if aok && bok {
		// BOTH INTS
		log.Println("both ints")
		// we can actually compare
		switch {
		case aAsInt < bAsInt:
			return -1
		case aAsInt == bAsInt:
			return 0
		case aAsInt > bAsInt:
			return 1
		}
	}

	if !aok && !bok {
		// BOTH SLICES
		log.Println("both slices")
		aList := a.([]item)
		bList := b.([]item)

		return compareLists(aList, bList)
	}

	// ONLY ONE IS A SLICE
	if !aok {
		log.Println("a is slice")
		// a is a slice; b is an int
		// convert int into slice
		return compare(a, []item{b})
	}

	if !bok {
		log.Println("b is slice")
		// convert int into slice
		return compare([]item{a}, b)
	}

	log.Println("no way we get here right?")

	return 0
}
