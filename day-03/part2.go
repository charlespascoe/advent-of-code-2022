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

	groups, err := groupLines(lines, 3)
	if err != nil {
		log.Fatalf("Invalid input: %s", err)
	}

	sum := 0

	for i, group := range groups {
		var overlap CharSet

		for _, rucksack := range group {
			set := charSet(rucksack)

			if overlap == nil {
				overlap = set
			} else {
				overlap = charIntersection(overlap, set)
			}
		}

		if len(overlap) != 1 {
			log.Fatalf("Group %d did not have exactly one item in common", i+1)
		}

		// Remember, this loop will only run once
		for c := range overlap {
			sum += priority(c)
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

func groupLines(lines []string, count int) ([][]string, error) {
	var groups [][]string
	var group []string

	for _, line := range lines {
		group = append(group, line)

		if len(group) == count {
			groups = append(groups, group)
			group = nil
		}
	}

	if len(group) != 0 {
		return nil, fmt.Errorf(
			"couldn't group into groups of %[1]d: %[2]d lines not multiple of %[1]d",
			count,
			len(lines),
		)
	}

	return groups, nil
}

func charSet(s string) CharSet {
	set := make(CharSet)

	for _, c := range s {
		set[c] = true
	}

	return set
}

func charIntersection(set1, set2 CharSet) CharSet {
	res := make(CharSet)

	for c := range set1 {
		if set2[c] {
			res[c] = true
		}
	}

	return res
}

func priority(c rune) int {
	if c >= 'a' {
		return int(c-'a') + 1
	} else {
		return int(c-'A') + 27
	}
}
