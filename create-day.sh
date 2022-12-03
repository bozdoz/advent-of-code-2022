#!/bin/bash

NEW_DAY=$1
NEW_DAY_NAME="day-$NEW_DAY"

usage() {
    cat >&2 <<END_USAGE

Create a new boilerplate directory from a template

USAGE:
    ./create-day.sh 01
END_USAGE
}

if [ -z $NEW_DAY ]; then
  echo "Provide ## for new day directory"
	usage
  exit 1
fi

mkdir $NEW_DAY

cd $NEW_DAY

# start touching things
touch README.md
touch input.txt
touch example.txt

# create main go file
cat > $NEW_DAY_NAME.go <<EOF
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bozdoz/advent-of-code-2022/utils"
)

// today's input data type
type dataType []string

// how to read today's inputs
var fileReader = utils.ReadLines

func partOne(data dataType) (ans int) {
	return
}

func partTwo(data dataType) (ans int) {
	return
}

// initialize the app by setting log flags
func init() {
	log.SetFlags(log.Llongfile)
}

// run the solvers
func main() {
	filename := utils.GetInputFile()
	data := fileReader(filename)

	fncs := map[string]func(dataType) int{
		"partOne": partOne,
		"partTwo": partTwo,
	}

	// run partOne and partTwo
	for k, fun := range fncs {
		s := time.Now()
		val := fun(dataType(data))

		fmt.Printf("%s: %v (%v)\n", k, val, time.Since(s))
	}
}

EOF

# create test file
cat > ${NEW_DAY_NAME}_test.go << EOF
package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 0,
	2: 0,
}

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(dataType(data))

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val := partTwo(dataType(data))

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

EOF