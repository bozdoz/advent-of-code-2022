package utils

import (
	"bufio"
	"io"
	"os"
)

func ReadLines(file string) (lines []string, err error) {
	readFile, err := os.Open(file)

	if err != nil {
		return
	}

	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}

func ReadFile(file string) (content string, err error) {
	readFile, err := os.Open(file)

	if err != nil {
		return
	}

	defer readFile.Close()

	fileBytes, err := io.ReadAll(readFile)

	return string(fileBytes), err
}
