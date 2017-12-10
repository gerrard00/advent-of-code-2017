package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

type stack []string

func (sl stack) push(s string) []string {
	return append(sl, s)
}

func (sl stack) pop() []string {
	return sl[:len(sl)-1]
}

func (sl stack) peek() string {
	return sl[len(sl)-1]
}

func fn(input string) int {
	var totalScore int
	var currentScore int
	stateStack := stack{"starting"}
	var currentState string
	// treat cancel "!" differently than normal states
	inCancelMode := false

	skippedGarbageChars := 0

	for index, char := range input {
		fmt.Printf("char(%v): %v %c", index, char, char)
		currentState = stateStack.peek()
		fmt.Printf("state: %v\n", currentState)

		// cancel mode overrides any other state
		// if we are in cancel mode just skip this rune
		if inCancelMode {
			inCancelMode = false
			continue
		}

		// if we aren't in cancel mode than a ! overrides the next char
		if char == '!' && !inCancelMode {
			inCancelMode = true
			continue
		}

		switch currentState {
		case "starting":
			switch char {
			case '{':
				currentScore = 1
				stateStack = stateStack.pop()
				stateStack = stateStack.push("ingroup")
			// dirty hack, but I'm sleepy...let it start with garbage
			case '<':
				stateStack = stateStack.push("garbage")
			default:
				panic("group not started")
			}
		case "ingroup":
			switch char {
			case '}':
				totalScore += currentScore
				currentScore--
				fmt.Printf("Before %v\n", stateStack)
				stateStack = stateStack.pop()
				fmt.Printf("After %v\n", stateStack)
			case '{':
				currentScore += 1
				stateStack = stateStack.push("ingroup")
			case ',':
				// no op
			case '<':
				stateStack = stateStack.push("garbage")
			default:
				panic("unknown state")
			}
		case "garbage":
			switch char {
			case '>':
				stateStack = stateStack.pop()
			default:
				skippedGarbageChars++
			}
		}
	}

	fmt.Println("##############################")
	return skippedGarbageChars
}

type testpair struct {
	input  string
	output int
}

var testInputs = []testpair{
	{"<>", 0},
	{"<random characters>", 17},
	{"<<<<>", 3},
	{"<{!>}>", 2},
	{"<!!>", 0},
	{"<!!!>>", 0},
	{"<{o\"i!a,<{i<a>", 10},
	// real input, real result
	{getRealInput(), 5547},
}

func TestOne(t *testing.T) {
	for _, pair := range testInputs {
		v := fn(pair.input)

		if v != pair.output {
			t.Error(
				"For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("nine_input")

	if readFileErr != nil {
		panic(readFileErr)
	}
	s := string(bytes)
	// get rid of the newline at the end of the file
	return strings.Replace(s, "\n", "", -1)
}
