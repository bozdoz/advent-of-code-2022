package utils

import (
	"path/filepath"
	"runtime"
)

// we always read from a relative "input.txt" file
const input string = "input.txt"

// get the relative "input.txt" file for a given day
func GetInputFile(args ...any) string {
	// "1" will skip this "getInputFile.go" file
	depth := 1

	if len(args) > 0 {
		// get a different depth (used in runSolvers.go)
		depth = args[0].(int)
	}

	// read input relative to the caller that called this file
	_, caller, _, _ := runtime.Caller(depth)
	dir := filepath.Dir(caller)

	return filepath.Join(dir, input)
}
