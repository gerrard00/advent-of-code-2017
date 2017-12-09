package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func findTargetBank(banks []int) int {
	result := 0
	currentMaxSize := -1

	for i := range banks {
		if banks[i] > currentMaxSize {
			currentMaxSize = banks[i]
			result = i
		}
	}

	return result
}

func dumbHash(banks []int) string {
	asStrings := make([]string, len(banks))

	for _, bank := range banks {
		asStrings = append(asStrings, strconv.Itoa(bank))
	}

	return strings.Join(asStrings, "-")
}

func balance(banks *[]int) {
	var startingBankIndex = findTargetBank(*banks)
	// fmt.Printf("starting index: %v value: %v\n", startingBankIndex, (*banks)[startingBankIndex])

	blocksToMove := (*banks)[startingBankIndex]

	(*banks)[startingBankIndex] = 0

	currentIndex := startingBankIndex + 1

	for ; blocksToMove > 0; blocksToMove-- {
		if currentIndex == len((*banks)) {
			currentIndex = 0
		}

		(*banks)[currentIndex] += 1
		currentIndex++
	}
}

func fn(banks []int) int {
	previousConfigurations := make(map[string]bool)
	var currentConfiguration string
	cycles := 0

	for {
		cycles++
		// fmt.Println(banks)
		balance(&banks)

		currentConfiguration = dumbHash(banks)

		if previousConfigurations[currentConfiguration] {
			fmt.Println("Found duplicate configuration")
			return cycles
		}

		previousConfigurations[currentConfiguration] = true
	}
}

type testpair struct {
	input  []int
	output int
}

var testsOne = []testpair{
	{[]int{0, 2, 7, 0}, 5},
	// real input, real result
	{realInput, 7864},
}

func runTest(t *testing.T, testInputs []testpair) {
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

func TestOne(t *testing.T) {
	runTest(t, testsOne)
}

var realInput = []int{0, 5, 10, 0, 11, 14, 13, 4, 11, 8, 8, 7, 1, 4, 12, 11}
