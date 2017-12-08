package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
			if foundWords[words[wordIndex]] {
				valid = false
				break
			}

			foundWords[words[wordIndex]] = true
		}

		// no dupes
		if valid {
			count++
		}
	}

	fmt.Printf("Number of valid phrases: %v\n", count)
}
