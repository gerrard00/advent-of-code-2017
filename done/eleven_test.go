package bang

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"testing"
)

type step struct {
	description string
	x           int
	y           int
}

func getStep(s string) step {
	var result step

	// TODO: should really just build these once
	switch s {
	case "n":
		result = step{y: -10}
	case "ne":
		result = step{y: -5, x: 10}
	case "se":
		result = step{y: 5, x: 10}
	case "s":
		result = step{y: 10}
	case "sw":
		result = step{y: 5, x: -10}
	case "nw":
		result = step{y: -5, x: -10}
	default:
		panic(fmt.Sprintf("Unknown direction %v\n", s))
	}

	result.description = s
	return result
}

func convertCommandsToSteps(commands string) []step {
	result := make([]step, 0)

	for _, command := range strings.Split(commands, ",") {
		result = append(result, getStep(command))
	}

	return result
}

// wow I way overcomplicated this
// func calculateSteps(totalX int, totalY int) int {
// 	// I guess it would be more idiomatic to just write integer version of abs?
// 	// I'm lazy
// 	remainingX := math.Abs(float64(totalX))
// 	remainingY := math.Abs(float64(totalY))
// 	count := 0

// 	// TODO: I don't like that we are duplicating the knowledge of the directions
// 	for {
// 		switch {
// 		case totalX > 0 && totalY > 0:
// 			totalX--
// 			totalY--
// 		}
// 		case totalX > 0 && totalY > 0:
// 			totalX--
// 			totalY--
// 		}

// 		if remainingX+remainingY == 0 {
// 			return count
// 		}
// 	}
// panic("We should never get here")
// }

func abs(i int) uint {
	if i < 0 {
		return uint(0 - i)
	}

	return uint(i)
}

func getNumberOfSteps(x int, y int) uint {
	absoluteX := abs(x)
	absoluteY := abs(y)

	var rawCount uint

	if absoluteX > absoluteY {
		rawCount = absoluteX
	} else {
		rawCount = absoluteY
	}

	return uint(math.Ceil(float64(rawCount) / 10))
}

func fn(input string) (final uint, largest uint) {
	// fmt.Println(input)
	steps := convertCommandsToSteps(input)

	// fmt.Println(steps)

	totalX := 0
	totalY := 0
	var currentNumberOfSteps uint
	var largestNumberOfSteps uint

	for _, currentStep := range steps {
		totalX += currentStep.x
		totalY += currentStep.y
		currentNumberOfSteps = getNumberOfSteps(totalX, totalY)

		if largestNumberOfSteps < currentNumberOfSteps {
			largestNumberOfSteps = currentNumberOfSteps
		}
	}

	return currentNumberOfSteps, largestNumberOfSteps
}

type testpair struct {
	input    string
	total    uint
	furthest uint
}

var testFnInputs = []testpair{
	{"ne,ne,ne", 3, 3},
	{"ne,ne,sw,sw", 0, 2},
	{"ne,ne,s,s", 2, 2},
	{"se,sw,se,sw,sw", 3, 3},
	{getRealInput(), 773, 1560},
}

func TestFn(t *testing.T) {
	for _, pair := range testFnInputs {
		total, furthest := fn(pair.input)

		if total != pair.total || furthest != pair.furthest {
			t.Error(
				"For \"", pair.input,
				"\nexpected total: ", pair.total,
				"furthest: ", pair.furthest,
				"\ngot      total: ", total,
				"furthest", furthest,
			)
		}
	}
}

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("eleven_input")

	if readFileErr != nil {
		panic(readFileErr)
	}
	s := string(bytes)
	// get rid of the newline at the end of the file
	return strings.Replace(s, "\n", "", -1)
}
