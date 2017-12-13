package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type layer struct {
	layerRange       int
	scannerPosition  int
	directionReverse bool
}

func (l *layer) tick() int {
	var offset int

	if l.directionReverse {
		offset = -1
	} else {
		offset = 1
	}

	nextScannerPosition := l.scannerPosition + offset

	switch {
	case nextScannerPosition == l.layerRange:
		l.directionReverse = true
		l.scannerPosition = nextScannerPosition - 2
	case nextScannerPosition == -1:
		l.directionReverse = false
		l.scannerPosition = nextScannerPosition + 2
	default:
		l.scannerPosition = nextScannerPosition
	}

	return l.scannerPosition
}

func TestTick(t *testing.T) {
	const testRange = 3
	const numberOfTicksToTest = 2*testRange - 1
	l := layer{testRange, 0, false}
	expectedSequence := [5]int{1, 2, 1, 0, 1}

	for i := 0; i < numberOfTicksToTest; i++ {
		v := l.tick()
		if v != expectedSequence[i] {
			t.Error(
				"After step", i+1,
				"expected", expectedSequence[i],
				"got", l.scannerPosition,
			)
		}
	}
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
			fmt.Printf("padding line: %v\n", currentIndex)
			f.layers = append(f.layers, nil)
		}

		// fmt.Printf("current line: %v new layer index: %v : %v\n",
		// 	currentIndex, newLayerIndex, &newLayer.layerRange)

		f.layers = append(f.layers, &newLayer)
		if f.maxRange < newLayer.layerRange {
			f.maxRange = newLayer.layerRange
		}

		currentIndex = newLayerIndex + 1
	}

	return f
}

// just for func ;)
func printFirewall(f firewall, layerIndex int, packetPosition int) {
	for i := 0; i < len(f.layers); i++ {
		fmt.Printf(" %v  ", i)
	}
	fmt.Println()

	for i := 0; i < len(f.layers); i++ {
		fmt.Print("--- ")
	}
	fmt.Println()

	for rowIndex := 0; rowIndex < f.maxRange; rowIndex++ {
		for i := 0; i < len(f.layers); i++ {
			var scannerStatus string
			currentLayer := f.layers[i]

			if currentLayer == nil {
				if i == layerIndex && rowIndex == packetPosition {
					fmt.Print("( ) ")
				} else if rowIndex == 0 {
					fmt.Print("... ")
				} else {
					fmt.Print("    ")
				}
				continue
			}

			if rowIndex >= currentLayer.layerRange {
				fmt.Print("    ")
				continue
			}

			if f.layers[i].scannerPosition == rowIndex {
				scannerStatus = "S"
			} else {
				scannerStatus = " "
			}

			var cellFormatString string

			if i == layerIndex && rowIndex == packetPosition {
				cellFormatString = "(%v) "
			} else {
				cellFormatString = "[%v] "
			}

			fmt.Printf(cellFormatString, scannerStatus)
		}

		fmt.Println()
	}
	fmt.Println()
}

func fn(input string) int {
	f := convertToFirewall(input)
	spew.Dump(f)
	const packetPosition = 0

	result := 0

	for currentTick := 0; currentTick < len(f.layers); currentTick++ {
		printFirewall(f, currentTick, packetPosition)

		currentLayerToCheck := f.layers[currentTick]
		if currentLayerToCheck != nil && currentLayerToCheck.scannerPosition == packetPosition {
			fmt.Printf("\t** Caught on layer: %v at position %v **\n\n", currentTick, packetPosition)
			result += currentTick * currentLayerToCheck.layerRange
		}

		for _, currentLayerToUpdate := range f.layers {
			if currentLayerToUpdate != nil {
				currentLayerToUpdate.tick()
			}
		}
	}

	return result
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
		6: 4`, 24},
	// real input, real result
	{getRealInput(), 2264},
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

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("thirteen_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	return string(bytes)
}
