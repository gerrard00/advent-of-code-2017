package main

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

type grid [128][128]int

func generateGrid(keyString string) grid {
	var newGrid [128][128]int

	for rowIndex := 0; rowIndex < len(newGrid); rowIndex++ {
		rowKey := fmt.Sprintf("%v-%v", keyString, rowIndex)
		// fmt.Println(rowKey)

		rowHash := KnotHash(rowKey)
		// fmt.Println(rowHash)

		columnIndex := 0
		for hexByteIndex := 0; hexByteIndex < len(rowHash); hexByteIndex++ {
			// TODO: could just use bits.OnesCount
			// fmt.Printf("bang: %v\n", rowHash)
			byteValue, conversionError := strconv.ParseInt(string(rowHash[hexByteIndex]), 16, 8)

			if conversionError != nil {
				panic(conversionError)
			}
			// fmt.Printf("zoom: %v\n", byteValue)

			for byteIndexInBock := 3; byteIndexInBock > -1; byteIndexInBock-- {
				currentTestValue := byte(math.Pow(2, float64(byteIndexInBock)))
				flag := (currentTestValue&byte(byteValue) == 0)
				if flag {
					newGrid[rowIndex][columnIndex] = 0
				} else {
					newGrid[rowIndex][columnIndex] = 1
				}
				// fmt.Printf("val: %v test: %v = %v/%v\n", byteValue, currentTestValue, currentTestValue&byte(byteValue), flag)
				columnIndex++
			}
		}
	}

	return newGrid
}

func getCellKey(y int, x int) int {
	result := (y * 10000) + x
	// fmt.Printf("\t\ty: %v x: %v key: %v\n", y, x, result)
	return result
}

func visitRegion(regionId int, y int, x int, g grid, visited *map[int]bool) {
	if g[y][x] == 0 {
		// fmt.Printf("\t-%v x %v\n", y, x)
		return
	}

	cellKey := getCellKey(y, x)
	if (*visited)[cellKey] {
		return
	}

	// fmt.Printf("\t+%v x %v\n", y, x)

	// TODO: check out our neighbors
	(*visited)[cellKey] = true

	if y > 0 {
		visitRegion(regionId, y-1, x, g, visited)
	}

	if y < (len(g) - 1) {
		visitRegion(regionId, y+1, x, g, visited)
	}

	if x > 0 {
		visitRegion(regionId, y, x-1, g, visited)
	}

	if x < (len(g[0]) - 1) {
		visitRegion(regionId, y, x+1, g, visited)
	}
}

func getNumberOfRegions(g grid) int {
	numberOfRegions := 0
	visited := make(map[int]bool, 128*128)
	for rowIndex := range g {
		for columnIndex := range g[rowIndex] {
			cellKey := getCellKey(rowIndex, columnIndex)
			// fmt.Println(visited)
			if !visited[cellKey] {
				if g[rowIndex][columnIndex] == 1 {
					numberOfRegions++
					// fmt.Printf("Examining region %v rooted at y: %v x: %v\n",
					// 	numberOfRegions, rowIndex, columnIndex)
					// fmt.Printf("Visited before: %v\n", len(visited))
					visitRegion(numberOfRegions, rowIndex, columnIndex, g, &visited)
					// fmt.Printf("Visited after: %v\n", len(visited))

					// if numberOfRegions > 3 {
					// 	return -1
					// }
				}

				// either way make sure we add this cell to out v
				visited[cellKey] = true
			}
		}
	}

	return numberOfRegions
}

func fn(keyString string) int {
	g := generateGrid(keyString)

	result := getNumberOfRegions(g)

	return result
}

type testpair struct {
	input  string
	output int
}

var testsOne = []testpair{
	{"flqrgnkx", 1242},
	// real input, real result
	{"xlqgujun", 1089},
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
