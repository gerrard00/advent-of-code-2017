package main

import (
	// "fmt"
	"fmt"
	"testing"
)

func print(buff []int, currentIndex int) {
	for index := range buff {
		if index != currentIndex {
			fmt.Printf(" %d  ", buff[index])
		} else {
			fmt.Printf("(%d) ", buff[index])
		}
	}
	fmt.Println()
}

func fnNaive(stepsPer int, iterations int) int {
	if iterations == 0 {
		return 0
	}

	circularBuffer := []int{0}
	targetBufferIndex := 0
	resultIndex := -1
	for iteration := 0; iteration < iterations; iteration++ {
		// fmt.Printf("iteration %v targetBufferIndex %v\n", iteration, targetBufferIndex)
		// fmt.Println(circularBuffer)
		targetBufferIndex = ((targetBufferIndex + stepsPer) % len(circularBuffer)) + 1
		// fmt.Printf("\t add at %v\n", targetBufferIndex)
		newBuffer := make([]int, len(circularBuffer)+1)
		copy(newBuffer, circularBuffer[0:targetBufferIndex])
		// fmt.Printf("\t%v\n", newBuffer)
		newBuffer[targetBufferIndex] = iteration + 1
		// fmt.Printf("\t%v\n", newBuffer)

		for copyIndex := targetBufferIndex; copyIndex < len(circularBuffer); copyIndex++ {
			// fmt.Printf("\t\t%v->\n", copyIndex)
			newBuffer[copyIndex+1] = circularBuffer[copyIndex]
			// fmt.Printf("\t\t%v\n", newBuffer)
		}

		circularBuffer = newBuffer

		// fmt.Println()
		// print(circularBuffer, targetBufferIndex)
		resultIndex = (targetBufferIndex + 1) % len(circularBuffer)
		// fmt.Printf("Iteration: %v Result: %v Guess: %v\n",
		// 	iteration, circularBuffer[resultIndex], iteration+1)
		// fmt.Printf("Iteration: %v Result: %v\n",
		// 	iteration, circularBuffer[resultIndex])
	}

	return circularBuffer[resultIndex]
}

type testpair struct {
	stepsPer   int
	iterations int
	output     int
}

var testsOne = []testpair{
	// short sample
	// {3, 10, 5},
	// {3, 140, 5},
	// {3, 4, 5},
	// real sample
	{3, 2017, 638},
	// real input, real result
	{377, 2017, 596},
}

func runTest(t *testing.T, testInputs []testpair) {

	for _, pair := range testInputs {
		v := fnNaive(pair.stepsPer, pair.iterations)

		if v != pair.output {
			t.Error(
				// "For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}

func TestOne(t *testing.T) {
	runTest(t, testsOne)
}

// func getRealInput() string {
// 	bytes, readFileErr := ioutil.ReadFile("thirteen_input")

// 	if readFileErr != nil {
// 		panic(readFileErr)
// 	}

// 	return string(bytes)
// }
