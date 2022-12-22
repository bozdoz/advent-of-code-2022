package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// helper for parsing lines of text from a file
func ReadLinesFunc[T comparable](file string, parser func(s string) T) (lines []T) {
	readFile, err := os.Open(file)

	if err != nil {
		panic(fmt.Sprint("could not open file", file, err))
	}

	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	for scanner.Scan() {
		lines = append(lines, parser(scanner.Text()))
	}

	return
}

// read lines into []string
func ReadLines(file string) (lines []string) {
	return ReadLinesFunc(file, func(s string) string {
		return s
	})
}

// read entire file as a string
func ReadFile(file string) (content string) {
	readFile, err := os.Open(file)

	if err != nil {
		panic(fmt.Sprint("could not open file", file, err))
	}

	defer readFile.Close()

	fileBytes, err := io.ReadAll(readFile)

	if err != nil {
		panic(fmt.Sprint("could not read file", file, err))
	}

	// only trim the end of the document
	return strings.TrimRightFunc(string(fileBytes), unicode.IsSpace)
}

// read lines into []int
func ReadInts(file string) (content []int) {
	return ReadLinesFunc(file, func(s string) int {
		return ParseInt(s)
	})
}

// splits and trims file to return empty-new-line-separated groups
func ReadEmptyLineGroups(file string) (content []string) {
	return strings.Split(ReadFile(file), "\n\n")
}
