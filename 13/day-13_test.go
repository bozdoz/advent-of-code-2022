package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Llongfile)
}

func TestParsing(t *testing.T) {
	runs := []string{
		"[[1],[2,3,4]]",
		"[[4,4],4,4]",
		"[1,1,3,1,1]",
		"[[[5,[1,8,1,0,9],[7,4,1,3]],[[6],[3],7],[[],0,[8],[10,6,8,10]]],[2,7,4],[4,[[],8,5,[8,3,3]],[],4,[[9,9,4],2,[],5,8]],[[],8],[[[1,3,10],[2,7,8]],[],[9],10,4]]",
		"[[[[8,5,0],[3,3],[4],8],10,[5,[9,8,4,6],0,7],[[5,7,10],[10,10,6,3],5],5],[[],9,[8]],[],[]]",
		"[[0,[[0,2,10,4],[7,3],6],9],[3,[0,[8,8,3,9,7]],[],[[4,9],6,[10,1,2,0],10,[]],[7,7,[9,8,1],1,2]]]",
		"[[9,[9],0,[6,[1,10,8,8,6],[],6,2]],[[[2,1,9],6,[]]],[[],[[2,3,5,8],[2,9,5,7,7],7],[],[8,[7,3,3],[8],3]],[7,[[5,9,5],4]]]",
		"[[[6,[10,7,6],[10,3,6,3],[]],[]]]",
		"[[[[7,5,8,9],6],[8,6,[10,7,9,5,2],[10,1,8]],3,10,8],[[[4,9,10,3,9],8],[[10,0,10,1],3,[8,0,4],[7,10,1,8],[9,9]],[3]],[[10,[10,2,6],1],[[8],[9,6,9,4],[2,9]],0]]",
		"[[],[7,[[3,0,3,7,3],[6,6,10,8,5],[],10],0,10,[[],3,[],5,[5,5,5,0]]],[[[0,8,6,7],[0,4,5,5,9]],[4,2],[[10,1,1,1,0],4,[],[3,6],[0,9]],0],[[8,2,2,6],[[1,9],9,[4]],6]]",
	}

	for _, data := range runs {
		t.Run(data, func(t *testing.T) {
			item := parseItem(data)

			expected := strings.ReplaceAll(data, ",", " ")
			got := fmt.Sprint(item)
			if expected != got {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func TestSorted(t *testing.T) {
	runs := []pair{{
		left:  []item{1, 1, 3, 1, 1},
		right: []item{1, 1, 5, 1, 1},
	}, {
		left:  []item{[]item{1}, []item{2, 3, 4}},
		right: []item{[]item{1}, 4},
	}, {
		left:  []item{[]item{4, 4}, 4, 4},
		right: []item{[]item{4, 4}, 4, 4, 4},
	}, {
		left:  []item{},
		right: []item{3},
	}, {
		left:  parseItem("[[],[[3]],[[[9,8,2]],0,6,2],[[[],0],6,9,8,5]]"),
		right: parseItem("[[3]]"),
	}}

	for _, pair := range runs {
		t.Run(fmt.Sprintf("%v is less than %v", pair.left, pair.right), func(t *testing.T) {
			got := pair.isOrdered()

			if !got {
				t.Errorf("expected %v, got %v", true, got)
			}
		})
	}
}

func TestNotSorted(t *testing.T) {
	runs := []pair{{
		left:  []item{9},
		right: []item{[]item{8, 7, 6}},
	}, {
		left:  []item{7, 7, 7, 7},
		right: []item{7, 7, 7},
	}, {
		left:  []item{[]item{[]item{}}},
		right: []item{[]item{}},
	}, {
		left:  parseItem("[1,[2,[3,[4,[5,6,7]]]],8,9]"),
		right: parseItem("[1,[2,[3,[4,[5,6,0]]]],8,9]"),
	}}

	for _, pair := range runs {
		t.Run(fmt.Sprintf("%v is gt %v", pair.left, pair.right), func(t *testing.T) {
			got := pair.isOrdered()

			if got {
				t.Errorf("expected %v, got %v", false, got)
			}
		})
	}
}

// fill in the answers for each part (as they come)
var answers = map[int]outType{
	1: 13,
	2: 140,
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
