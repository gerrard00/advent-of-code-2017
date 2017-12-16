package main

import "testing"

func fn(i int) int {
	return i + 1
}

type testpair struct {
	input  int
	output int
}

var testsOne = []testpair{
	{1, 2},
	// real input, real result
	// {getRealInput(), 3875838},
}

func runTest(t *testing.T, testInputs []testpair) {
	for _, pair := range testInputs {
		v := fn(pair.input)

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
