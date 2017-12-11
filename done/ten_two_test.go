package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func encodeInput(s string) []byte {
	converted := []byte(s)
	return append(converted, []byte{17, 31, 73, 47, 23}...)
}

func TestEncodeInput(t *testing.T) {
	const input = "1,2,3"
	expected := []byte{49, 44, 50, 44, 51, 17, 31, 73, 47, 23}
	result := encodeInput(input)

	if bytes.Compare(result, expected) != 0 {
		t.Error(
			"For", input,
			"expected", expected,
			"got", result,
		)
	}
}

func makeDenseHash(sparseHash []byte) string {
	var result string
	var currentBlockValue byte

	for i, v := range sparseHash {
		currentBlockValue ^= v

		if (i+1)%16 == 0 {
			result += hex.EncodeToString([]byte{currentBlockValue})
			currentBlockValue = 0
		}
	}

	return result
}

func TestMakeDenseHash(t *testing.T) {
	input := []byte{65, 27, 9, 1, 4, 3, 40, 50, 91, 7, 6, 0, 2, 5, 68, 22, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expected := "4007ff"
	result := makeDenseHash(input)

	if result != expected {
		t.Error(
			"For", input,
			"expected", expected,
			"got", result,
		)
	}
}

func makeSparseHash(numberOfValues int, lengths []byte) []byte {
	adjustIndex := func(index int) int {
		if index >= numberOfValues {
			adjusted := index % numberOfValues
			// fmt.Printf("\tToo big! Adjusted %v to %v\n", index, adjusted)
			return adjusted
		}

		if index < 0 {
			// TODO: this is wrong, need the equivalent modulo
			adjusted := index + numberOfValues
			// fmt.Printf("\tToo small! Adjusted %v to %v\n", index, adjusted)
			return adjusted
		}

		// fmt.Printf("\tJust right! %v\n", index)
		return index
	}

	currentPosition := 0
	skipSize := 0
	// values := [numberOfValues]int
	values := make([]byte, numberOfValues)

	for i, _ := range values {
		values[i] = byte(i)
	}

	for round := 0; round < 64; round++ {
		// fmt.Printf(">>>>>>>>>>>>>>>>>>>>%v\n", round)
		// if round > 0 {
		// fmt.Printf("->->->->->->->->->->->->->->->->->->->->%v\n", round)
		// os.Exit(-1)
		// }
		for _, currentLengthAsByte := range lengths {
			// fmt.Printf("position: %v skip: %v length: %v\n",
			// 	currentPosition, skipSize, currentLength)

			// quick hack to avoid figuring outtypes
			currentLength := int(currentLengthAsByte)

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

				// fmt.Printf("<%v-%v>\n", reverseStartIndex, reverseEndIndex)
				// fmt.Println(values)
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
	}

	// fmt.Printf("%v * %v = %v\n", values[0], values[1], values[0]*values[1])
	// fmt.Println(values)
	return values
}

type cycletestpair struct {
	numberOfValues int
	input          []byte
	output         []byte
}

var testCycleInputs = []cycletestpair{
	{5, []byte{3, 4, 1, 5}, []byte{3, 4, 0, 1, 2}},
	// real input, real result
	// {256, []int{31, 2, 85, 1, 80, 109, 35, 63, 98, 255, 0, 13, 105, 254, 128, 33, 5547}, 6952},
}

func TestCycle(t *testing.T) {
	for _, pair := range testCycleInputs {
		v := makeSparseHash(pair.numberOfValues, pair.input)

		if bytes.Compare(v, pair.output) != 0 {
			t.Error(
				"For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}

func knotHash(input string) string {
	encodedInput := encodeInput(input)
	// TODO: taking the first argument of size is just for testing
	sparseHash := makeSparseHash(256, encodedInput)

	return makeDenseHash(sparseHash)
}

type knothashtestpair struct {
	input  string
	output string
}

var testKnotHashInputs = []knothashtestpair{
	{"", "a2582a3a0e66e6e86e3812dcb672a272"},
	{"AoC 2017", "33efeb34ea91902bb2f59c9920caa6cd"},
	{"1,2,3", "3efbe78a8d82f29979031a4aa0b16a9d"},
	{"1,2,4", "63960835bcdc130f0b66d7ff4f6a5a8e"},
	{"31,2,85,1,80,109,35,63,98,255,0,13,105,254,128,33", ""},
}

func TestKnotHash(t *testing.T) {
	for _, pair := range testKnotHashInputs {
		v := knotHash(pair.input)

		if v != pair.output {
			t.Error(
				"For \"", pair.input,
				"\"\nexpected \"", pair.output,
				"\"\ngot      \"", v,
				"\"",
			)
		}
	}
}
