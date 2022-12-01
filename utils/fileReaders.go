package utils

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

// helper for parsing lines of text from a file
func ReadLinesFunc[T comparable](file string, parser func(s string) T) (lines []T, err error) {
	readFile, err := os.Open(file)

	if err != nil {
		return
	}

	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	for scanner.Scan() {
		lines = append(lines, parser(scanner.Text()))
	}

	return
}

// read lines into []string
func ReadLines(file string) (lines []string, err error) {
	return ReadLinesFunc(file, func(s string) string {
		return s
	})
}

// read entire file as a string
func ReadFile(file string) (content string, err error) {
	readFile, err := os.Open(file)

	if err != nil {
		return
	}

	defer readFile.Close()

	fileBytes, err := io.ReadAll(readFile)

	return string(fileBytes), err
}

// read lines into []int
func ReadInts(file string) (content []int, err error) {
	return ReadLinesFunc(file, func(s string) int {
		as_int, err := strconv.Atoi(s)

		if err != nil {
			panic(err)
		}

		return as_int
	})
}
