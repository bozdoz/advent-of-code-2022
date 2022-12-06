package main

import (
	"container/list"
	"fmt"
)

// returns 1-based index of first unique sequential characters of length `count`
func uniqueLettersIndex(in inType, count int) int {
	l := list.New()
	// add first element
	l.PushBack(in[0])

outer:
	for i := 1; i < len(in); i++ {
		cur := in[i]

		// check if current is in list already, moving backwards
		for e, j := l.Back(), l.Len(); e != nil; e, j = e.Prev(), j-1 {
			if e.Value == cur {
				// match means we need to clean the list
				removeFront(l, j)

				// add element
				l.PushBack(cur)
				continue outer
			}
		}

		// check if we're done
		if l.Len() == count-1 {
			// index is 1-based
			return i + 1
		}

		// add element
		l.PushBack(cur)
	}

	return -1
}

// we remove items from the front of the list, depending on where we found the match
func removeFront(l *list.List, count int) {
	for i := count; i > 0; i-- {
		l.Remove(l.Front())
	}
}

// print list out for debugging
func ListToString(l *list.List) (out string) {
	e := l.Front()

	out += "["

	for e != nil {
		out += fmt.Sprintf("%q", e.Value)
		e = e.Next()
	}

	out += "]"

	return
}
