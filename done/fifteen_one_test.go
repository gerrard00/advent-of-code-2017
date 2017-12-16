package main

import "testing"

func compareResults(a int, b int) bool {
	// n = 16
	// 2^(n - 1) == "1111111111111111"
	const mask = 65535

	// fmt.Printf("\t%016d\n", a)
	// fmt.Printf("\t%016d\n", b)
	// fmt.Println()
	// fmt.Printf("\t%016b\n", mask)
	// fmt.Printf("\t%016b\n", a)
	// fmt.Printf("\t%016b\n", b)
	// fmt.Println()

	aToCompare := a & mask
	bToCompare := b & mask

	// fmt.Printf("\t%v\n", aToCompare)
	// fmt.Printf("\t%v\n", bToCompare)

	return aToCompare == bToCompare
}

type testpairForCompareResults struct {
	inputA int
	inputB int
	output bool
}

var testsForCompare = []testpairForCompareResults{
	{65546, 10, true},
	{65546, 11, false},
}

func TestCompareResults(t *testing.T) {
	for _, pair := range testsForCompare {
		actual := compareResults(pair.inputA, pair.inputB)

		if actual != pair.output {
			t.Error("Nope")
		}
	}
}

const divisor = 2147483647

type generator struct {
	factor        int
	previousValue int
}

func (g *generator) generate() int {
	result := g.previousValue * g.factor % divisor
	g.previousValue = result
	return result
}

func judge(gA *generator, gB *generator) bool {
	aResult := gA.generate()
	bResult := gB.generate()

	// fmt.Printf("%10v\t%10v\n", aResult, bResult)

	// TODO: compare first 16 bits
	return compareResults(aResult, bResult)
}

func fn(numberOfTests int, startingValueForA int, startingValueForB int) int {
	result := 0
	generatorA := generator{16807, startingValueForA}
	generatorB := generator{48271, startingValueForB}

	for i := 0; i < numberOfTests; i++ {
		if judge(&generatorA, &generatorB) {
			result++
		}
	}

	return result
}

type testpair struct {
	numberOfTests     int
	startingValueForA int
	startingValueForB int
	output            int
}

var testsOne = []testpair{
	{5, 65, 8921, 1},
	// real input, real result
	{40000000, 634, 301, 573},
}

func runTest(t *testing.T, testInputs []testpair) {
	for _, pair := range testInputs {
		v := fn(pair.numberOfTests, pair.startingValueForA, pair.startingValueForB)

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
