package utils

import (
	"flag"
	"fmt"
	"time"
)

// A day is just a file reader and the functions to call
// with the input content
type Day[T any, O any] struct {
	FileReader func(string) T
	Fncs       []func(T) O
}

// abstract boilerplate to run solvers for each day
func RunSolvers[T any, O any](day Day[T, O]) {
	// depth 2 skips this function and the getInputFile function
	callerDepth := 2
	inputFile := getInputFile(callerDepth)
	data := day.FileReader(inputFile)

	var part int
	flag.IntVar(&part, "part", 0, "Which part to run in isolation")

	// parse flags from command line
	flag.Parse()

	var fun func(T) O
	var i int

	fncs := day.Fncs

	// output
	if part > 0 {
		fmt.Println("running part", part)
		fncs = fncs[part-1 : part]
		i = part - 1
	}

	for _, fun = range fncs {
		i++
		s := time.Now()
		val := fun(data)

		fmt.Printf("%v | %v (%v)\n", i, val, time.Since(s))
	}
}
