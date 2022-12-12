package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type CharSet map[rune]bool

func main() {
	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	sum := 0

	for _, rucksack := range lines {
		firstSet := charSet(rucksack[:len(rucksack)/2])
		secondSet := charSet(rucksack[len(rucksack)/2:])

		for c := range secondSet {
			if firstSet[c] {
				sum += priority(c)
			}
		}
	}

	fmt.Printf("Sum: %d\n", sum)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func charSet(s string) CharSet {
	set := make(CharSet)

	for _, c := range s {
		set[c] = true
	}

	return set
}

func priority(c rune) int {
	if c >= 'a' {
		return int(c-'a') + 1
	} else {
		return int(c-'A') + 27
	}
}
