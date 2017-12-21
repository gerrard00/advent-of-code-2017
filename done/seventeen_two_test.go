package main

import (
	"fmt"
	"testing"
)

func fnNaive(stepsPer int, iterations int) int {
	if iterations == 0 {
		return 0
	}

	targetBufferIndex := 0
	valueAtSwitch := -1
	currentBufferLength := 1
	for iteration := 0; iteration < iterations; iteration++ {
		rawTargetBufferIndex := targetBufferIndex + stepsPer
		targetBufferIndex = (rawTargetBufferIndex % currentBufferLength) + 1
		if targetBufferIndex == 1 {
			valueAtSwitch = iteration + 1
			fmt.Println(valueAtSwitch)
		}
		nextBufferLength := currentBufferLength + 1
		currentBufferLength = nextBufferLength

	}

	return valueAtSwitch
}

type testpair struct {
	stepsPer   int
	iterations int
	output     int
}

var testsOne = []testpair{
	{377, 50000000, 39051595},
}

func runTest(t *testing.T, testInputs []testpair) {

	for _, pair := range testInputs {
		v := fnNaive(pair.stepsPer, pair.iterations)

		if v != pair.output {
			t.Error(
				"expected", pair.output,
				"got", v,
			)
		}
	}
}

func TestOne(t *testing.T) {
	runTest(t, testsOne)
}
