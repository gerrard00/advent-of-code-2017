package main

import (
	"math"
	"testing"
)

func fn(input int) int {
	// TODO: could probably DRY this up

	// max number in ring = (1 + (2 * n))^2 from observation
	// we want the ring below our target, so we do the inverse
	rawRingStep := (math.Sqrt(float64(input)) - 1) / 2

	actualRingIndex := int(math.Ceil(rawRingStep))

	maxDistance := 2 * actualRingIndex

	// an exact match is easy, 2 * our ring index
	if rawRingStep == math.Trunc(rawRingStep) {
		return maxDistance
	}

	// TODO: can we avoid re-calculating this
	maxInRing := int(math.Pow(float64(1+(2*actualRingIndex)), 2))

	sideLength := int(math.Sqrt(float64(maxInRing)))
	sideOffset := sideLength - 1

	leftBottomCorner := int(maxInRing) - sideOffset
	leftTopCorner := leftBottomCorner - sideOffset
	rightTopCorner := leftTopCorner - sideOffset

	// if we are on a corner, return max distance
	if input == rightTopCorner ||
		input == leftTopCorner ||
		input == leftBottomCorner {
		return maxDistance
	}

	middleDistance := int(sideOffset / 2)

	// right column
	if input < rightTopCorner {
		rightMiddle := rightTopCorner - middleDistance

		if input >= rightMiddle {
			return maxDistance - (rightTopCorner - input)
		}

		return maxDistance - (rightMiddle - input)
	}

	// left column
	if input < leftBottomCorner &&
		input > leftTopCorner {
		leftMiddle := leftTopCorner - middleDistance

		if input >= leftMiddle {
			return maxDistance - (input - leftTopCorner)
		}

		return maxDistance - (input - leftMiddle)
	}

	// bottom row, not in column
	if input > leftBottomCorner {
		bottomCenter := maxInRing - middleDistance

		if input >= bottomCenter {
			return maxDistance - (maxInRing - input)
		}

		return maxDistance - (bottomCenter - input)
	}

	// top row, not in column
	topCenter := leftTopCorner - middleDistance

	if input >= topCenter {
		return maxDistance - (leftTopCorner - input)
	}

	return maxDistance - (input - rightTopCorner)
}

type testpair struct {
	input  int
	output int
}

var tests = []testpair{
	{1, 0},
	{12, 3},
	// gerrard test
	// left column
	{17, 4},
	{21, 4},
	// max ring
	{25, 4},
	// right column, not max
	{10, 3},
	{13, 4},
	// left column, not max
	{18, 3},
	{19, 2},
	// not column
	{22, 3},
	{23, 2},
	{14, 3},
	{16, 3},
	{32, 5},
	// gerrard test
	{23, 2},
	{1024, 31},
	// real input, real result
	{realInput, 419},
}

func Test(t *testing.T) {
	for _, pair := range tests {
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

const realInput = 289326
