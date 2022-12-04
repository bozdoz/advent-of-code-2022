package utils

import (
	"fmt"
	"time"
)

// A day is just a file reader and the functions to call
// with the input content
type Day[T any] struct {
	FileReader func(string) T
	Fncs       []func(T) int
}

// abstract boilerplate to run solvers for each day
func RunSolvers[T any](day Day[T]) {
	callerDepth := 2
	inputFile := GetInputFile(callerDepth)
	data := day.FileReader(inputFile)

	for i, fun := range day.Fncs {
		s := time.Now()
		val := fun(data)

		fmt.Printf("%v | %v (%v)\n", i+1, val, time.Since(s))
	}
}
