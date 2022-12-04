package utils

import (
	"path/filepath"
	"runtime"
)

// we always read from a relative "input.txt" file
const input string = "input.txt"

// get the relative "input.txt" file for a given day
// depth "1" will skip this "getInputFile.go" file
func getInputFile(depth int) string {
	// read input relative to the caller that called this file
	_, caller, _, _ := runtime.Caller(depth)
	dir := filepath.Dir(caller)

	return filepath.Join(dir, input)
}
