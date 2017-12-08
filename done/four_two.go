package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// https://stackoverflow.com/a/22698017/1011470
type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func SortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}

func makeKey(s string) string {
	return SortString(s)
}

func main() {
	reader, openErr := os.Open("four_one_input")

	if openErr != nil {
		panic(openErr)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	count := 0

	for scanner.Scan() {
		foundWords := make(map[string]bool)
		words := strings.Split(scanner.Text(), " ")
		valid := true

		for wordIndex := range words {
			key := makeKey(words[wordIndex])

			if foundWords[key] {
				valid = false
				break
			}

			foundWords[key] = true
		}

		// no dupes
		if valid {
			count++
		}
	}

	fmt.Printf("Number of valid phrases: %v\n", count)
}
