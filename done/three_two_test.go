package main

import (
	"fmt"
	"math"
	"testing"
)

func makeMemory(size int) [][]int {
	// make sure size is odd
	if size%2 == 0 {
		panic("Make it odd")
	}
	result := make([][]int, size)
	buffer := make([]int, size*size)

	for i := range result {
		result[i], buffer = buffer[:size], buffer[size:]
	}

	return result
}

type fillFunc func(memory [][]int, y int, x int) int

var incrementCounter = 1

func fillIncrement(memory [][]int, y int, x int) int {
	incrementCounter++
	return incrementCounter
}

func fillStress(memory [][]int, y int, x int) int {
	max := len(memory) - 1

	minimumY := y - 1
	if minimumY < 0 {
		minimumY = 0
	}

	minimumX := x - 1
	if minimumX < 0 {
		minimumX = 0
	}

	maximumY := y + 1
	if maximumY > max {
		maximumY = max
	}

	maximumX := x + 1
	if maximumX > max {
		maximumX = max
	}

	result := 0

	for examineY := minimumY; examineY <= maximumY; examineY++ {
		for examineX := minimumX; examineX <= maximumX; examineX++ {
			fmt.Printf("\ttesting -> %v %v\n", examineY, examineX)
			if examineX != x || examineY != y {
				fmt.Printf("\t\twill add -> %v\n", memory[examineY][examineX])
				result += memory[examineY][examineX]
			}
		}
	}

	return result
}

func printMemory(memory [][]int) {
	for y := range memory {
		for x := range memory[y] {
			fmt.Printf("%v\t", memory[y][x])
		}

		fmt.Println()
	}
}

func fn(input int, filler fillFunc) int {
	const memorySize = 11
	memory := makeMemory(memorySize)

	currentY := int(math.Ceil(memorySize / 2))
	currentX := currentY

	memory[currentY][currentX] = 1

	// TODO: declare these closer to where they are used?
	var yOffset, xOffset int

	for currentRingIndex := 0; currentRingIndex < 4; currentRingIndex++ {
		sideLength := 2 + (2 * currentRingIndex)
		fmt.Printf("current side length: %v\n", sideLength)
		// TODO: loop in ring here
		for numberOfTurnsTaken := 0; numberOfTurnsTaken < 5; numberOfTurnsTaken++ {
			fmt.Printf("Turn # %v\n", numberOfTurnsTaken)
			switch numberOfTurnsTaken {
			case 0:
				fallthrough
			case 4:
				yOffset = 0
				xOffset = 1
			case 1:
				yOffset = -1
				xOffset = 0
			case 2:
				yOffset = 0
				xOffset = -1
			case 3:
				yOffset = 1
				xOffset = 0
			}

			fmt.Printf("Turn %v offset y: %v x: %v\n",
				numberOfTurnsTaken, yOffset, xOffset)

			// the max length of a side is 1 for the
			// beginning of the first side and then
			// side length for the rest of them
			var maxLengthForThisSide int
			switch numberOfTurnsTaken {
			case 0:
				maxLengthForThisSide = 1
			case 1:
				maxLengthForThisSide = sideLength - 1
			default:
				maxLengthForThisSide = sideLength
			}

			// TODO: loop straight here
			for numberOfStraightStepsTaken := 0; numberOfStraightStepsTaken < maxLengthForThisSide; numberOfStraightStepsTaken++ {
				currentY += yOffset
				currentX += xOffset
				currentValue := filler(memory, currentY, currentX)
				memory[currentY][currentX] = currentValue

				if currentValue > input {
					printMemory(memory)
					return currentValue
				}

				fmt.Printf("debug: %v y: %v x:% v\nsteps: %v, length: %v\n",
					memory[currentY][currentX], currentY,
					currentX,
					numberOfStraightStepsTaken,
					maxLengthForThisSide)

				printMemory(memory)
			}
		}
	}

	return -1
}

type testpair struct {
	input  int
	output int
}

var tests = []testpair{
	// {1, 1},
	// {2, 1},
	// {3, 2},
	// {4, 4},
	// {5, 5},
	// real input, real result
	// {806, 4},
	{realInput, 419},
}

func Test(t *testing.T) {
	// fillFunc := fillIncrement
	fillFunc := fillStress

	for _, pair := range tests {
		v := fn(pair.input, fillFunc)

		if v != pair.output {
			t.Error(
				"For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}

const realInput = 289326
