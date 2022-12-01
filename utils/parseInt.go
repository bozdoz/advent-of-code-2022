package utils

import (
	"fmt"
	"strconv"
)

// panics if int cannot be parsed from string
func ParseInt(str string) int {
	num, err := strconv.Atoi(str)

	if err != nil {
		panic(fmt.Sprintf("Could not parse number: %s - %v", str, err))
	}

	return num
}
