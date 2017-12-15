package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"testing"
)

const useCache = false

type layer struct {
	layerRange int
	period     int
	// can probably calculate this, but for now
	layerValues []int
}

type firewall struct {
	layers   []*layer
	maxRange int
}

func convertToFirewall(input string) firewall {
	// 0: 3
	// 1: 2
	// 4: 4
	// 6: 4
	scanner := bufio.NewScanner(strings.NewReader(input))
	f := firewall{maxRange: 0}
	currentIndex := 0
	// f.maxRange = 0

	for scanner.Scan() {
		var newLayerIndex int
		var newLayer layer

		currentLine := scanner.Text()

		_, scanErr := fmt.Sscanf(currentLine, "%d: %d", &newLayerIndex, &newLayer.layerRange)
		if scanErr != nil {
			panic(scanErr)
		}

		// handle empty layers
		for ; currentIndex < newLayerIndex; currentIndex++ {
			// fmt.Printf("padding line: %v\n", currentIndex)
			f.layers = append(f.layers, nil)
		}

		// fmt.Printf("current line: %v new layer index: %v : %v\n",
		// 	currentIndex, newLayerIndex, &newLayer.layerRange)

		newLayer.period = (2 * newLayer.layerRange) - 1
		newLayer.layerValues = make([]int, newLayer.period-1)

		halfWayIndex := int(math.Floor(float64(newLayer.period) / 2))
		// fmt.Printf("half %v\n", halfWayIndex)
		for valueIndex := 0; valueIndex < newLayer.period-1; valueIndex++ {
			var currentValue int

			if valueIndex <= halfWayIndex {
				currentValue = valueIndex
			} else {
				currentValue = newLayer.period - valueIndex - 1
				// fmt.Printf("%v - %v - 1 = %v\n", newLayer.period, valueIndex, currentValue)
			}

			newLayer.layerValues[valueIndex] = currentValue
		}

		// fmt.Println(newLayer)
		f.layers = append(f.layers, &newLayer)
		if f.maxRange < newLayer.layerRange {
			f.maxRange = newLayer.layerRange
		}

		currentIndex = newLayerIndex + 1
	}

	return f
}

// just a func for easy testing
func (l layer) getScannerLocation(targetTime int) int {
	// fmt.Printf("\trange %v period %v time: %v\n", l.layerRange, l.period, time)

	valueOffset := targetTime % (l.period - 1)

	// gross, but I just want to get done :)
	// fmt.Printf("\tvalue offset: %v\n", valueOffset)
	return l.layerValues[valueOffset]
}

func testGetScannerLocation() {
	l := layer{
		layerRange:  3,
		period:      5,
		layerValues: []int{0, 1, 2, 1},
	}
	for i := 0; i < 20; i++ {
		fmt.Printf("time: %v scanner: %v\n\n", i, l.getScannerLocation(i))
	}
}

func fn(input string) int {
	// testGetScannerLocation()
	// os.Exit(1)

	f := convertToFirewall(input)
	// spew.Dump(f)
	const packetRowIndex = 0

	// so we don't run forever, would be nice to calculate this
	// based on the maximum range and depth
	const maxTests = 10000000
	// const maxTests = 35

	var wasCaught bool

	for delay := 0; delay < maxTests; delay++ {
		if delay%1000 == 0 {
			fmt.Printf("\t***testing %v***\n", delay)
		}
		wasCaught = false

		for layerIndex := 0; layerIndex < len(f.layers); layerIndex++ {
			currentLayerToCheck := f.layers[layerIndex]
			if currentLayerToCheck != nil {
				// fmt.Printf("Layer: %v\n", layerIndex)
				scannerLocation := currentLayerToCheck.getScannerLocation(delay + layerIndex)
				// fmt.Printf("\tscanner location: %v\n", scannerLocation)
				if scannerLocation == 0 {
					wasCaught = true
					break
				}
			}
		}

		if !wasCaught {
			return delay
		}
	}

	return -1
}

type testpair struct {
	input  string
	output int
}

var testsOne = []testpair{
	{
		`0: 3
		1: 2
		4: 4
		6: 4`, 10},
	// real input, real result
	{getRealInput(), 3875838},
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

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("thirteen_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	return string(bytes)
}
