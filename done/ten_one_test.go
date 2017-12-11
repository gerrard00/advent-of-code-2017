package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func fn(numberOfValues int, lengths []int) int {
	adjustIndex := func(index int) int {
		if index >= numberOfValues {
			adjusted := index % numberOfValues
			// fmt.Printf("\tToo big! Adjusted %v to %v\n", index, adjusted)
			return adjusted
		}

		if index < 0 {
			adjusted := index + numberOfValues
			// fmt.Printf("\tToo small! Adjusted %v to %v\n", index, adjusted)
			return adjusted
		}

		return index
	}

	currentPosition := 0
	skipSize := 0
	// values := [numberOfValues]int
	values := make([]int, numberOfValues)

	for i, _ := range values {
		values[i] = i
	}

	for _, currentLength := range lengths {
		// fmt.Printf("position: %v skip: %v length: %v\n",
		// 	currentPosition, skipSize, currentLength)

		if currentLength > numberOfValues {
			fmt.Printf("Warning: invalid length %v ignored\n", currentLength)
			continue
		}

		reverseStartIndex := currentPosition
		reverseEndIndex := adjustIndex(currentPosition + currentLength - 1)

		// fmt.Printf("Starting reverse %v %v\n", reverseStartIndex, reverseEndIndex)

		reverseStepsLeft := int(currentLength / 2)
		for {
			if reverseStepsLeft == 0 {
				break
			}

			reverseStepsLeft--

			// fmt.Printf("<%v-%v>", reverseStartIndex, reverseEndIndex)
			tempHolder := values[reverseStartIndex]
			values[reverseStartIndex] = values[reverseEndIndex]
			values[reverseEndIndex] = tempHolder

			// fmt.Printf("\t\t%v => %v %v\n", reverseStartIndex, reverseEndIndex, values)

			reverseStartIndex = adjustIndex(reverseStartIndex + 1)
			reverseEndIndex = adjustIndex(reverseEndIndex - 1)
			// fmt.Printf("( %v-%v )", reverseStartIndex, reverseEndIndex)
		}
		// fmt.Println("finish reverse")

		// prepare for next cycle
		currentPosition = adjustIndex(currentPosition + currentLength + skipSize)
		skipSize++

		// fmt.Println(values)
		// fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	}

	fmt.Printf("%v * %v = %v\n", values[0], values[1], values[0]*values[1])
	return values[0] * values[1]
}

type testpair struct {
	numberOfValues int
	input          []int
	output         int
}

var testInputs = []testpair{
	{5, []int{3, 4, 1, 5}, 12},
	// real input, real result
	// 31506 is too high
	{256, []int{31, 2, 85, 1, 80, 109, 35, 63, 98, 255, 0, 13, 105, 254, 128, 33, 5547}, 6952},
}

func TestOne(t *testing.T) {
	for _, pair := range testInputs {
		v := fn(pair.numberOfValues, pair.input)

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
