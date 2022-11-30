package utils

import (
	"path/filepath"
	"runtime"
)

// we always read from a relative "input.txt" file
const input string = "input.txt"

// get the relative "input.txt" file for a given day
func GetInputFile() string {
	// read input relative to the caller that called this file
	// "1" will skip this "getInputFile.go" file
	_, currentFile, _, _ := runtime.Caller(1)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, input)
}
