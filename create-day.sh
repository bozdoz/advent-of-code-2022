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

kill() {
	echo $1
	exit 1
}

mkdir $NEW_DAY || kill

cd $NEW_DAY

# start touching things
touch README.md
touch input.txt
touch example.txt

# create main go file
cat > $NEW_DAY_NAME.go <<EOF
package main

import "github.com/bozdoz/advent-of-code-2022/utils"

// today's input data type
type inType = []string

// how to read today's input
var fileReader = utils.ReadLines

// today's output data type
type outType = int

func partOne(data inType) (ans outType) {
	return
}

func partTwo(data inType) (ans outType) {
	return
}

//
// BOILERPLATE BELOW
//

func main() {
	// pass file reader and functions to call with input data
	utils.RunSolvers(utils.Day[inType, outType]{
		FileReader: fileReader,
		Fncs: []func(inType) outType{
			partOne,
			partTwo,
		},
	})
}

EOF

# create test file
cat > ${NEW_DAY_NAME}_test.go << EOF
package main

import "testing"

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 0,
	2: 0,
}

var data = fileReader("example.txt")

func TestExampleOne(t *testing.T) {
	expected := answers[1]

	val := partOne(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

func TestExampleTwo(t *testing.T) {
	expected := answers[2]

	val := partTwo(data)

	if val != expected {
		t.Errorf("Answer should be %v, but got %v", expected, val)
	}
}

EOF