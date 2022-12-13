package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	tm := loadTreeMap(lines)

	fmt.Printf("Trees:\n%s\n\n", tm)

	// Part 1 //

	// Map of highest value seen when scanning each row from right to left
	leftRightMax := scanMap(tm, false, false, maxSeen())
	// Map of highest value seen when scanning each row from left to right
	rightLeftMax := scanMap(tm, false, true, maxSeen())
	// Map of highest value seen when scanning each column from top to bottom
	topBottomMax := scanMap(tm, true, false, maxSeen())
	// Map of highest value seen when scanning each column from bottom to top
	bottomTopMax := scanMap(tm, true, true, maxSeen())

	visible := 0

	for pos := range tm.Values {
		hidden := isHidden(pos, tm,
			leftRightMax,
			rightLeftMax,
			topBottomMax,
			bottomTopMax,
		)

		if !hidden {
			visible++
		}
	}

	fmt.Printf("Visible trees: %d\n", visible)

	// Part 2 //

	// Map of visibility to the left by scanning each row from right to left
	leftVisibility := scanMap(tm, false, false, visibility())
	// Map of visibility to the right by scanning each row from left to right
	rightVisibility := scanMap(tm, false, true, visibility())
	// Map of visibility upwards by scanning each column from top to bottom
	upVisibility := scanMap(tm, true, false, visibility())
	// Map of visibility downwards by scanning each column from bottom to top
	downVisibility := scanMap(tm, true, true, visibility())

	bestScore := 0

	for pos := range tm.Values {
		score := scenicScore(pos,
			leftVisibility,
			rightVisibility,
			upVisibility,
			downVisibility,
		)

		if score > bestScore {
			bestScore = score
		}
	}

	fmt.Printf("Best scenic score: %d\n", bestScore)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func isHidden(pos int, tm TreeMap, maxHeightMaps ...TreeMap) bool {
	for _, m := range maxHeightMaps {
		if tm.Values[pos] > m.Values[pos] {
			return false
		}
	}

	return true
}

func scenicScore(pos int, visibilityMaps ...TreeMap) int {
	score := 1
	for _, m := range visibilityMaps {
		score *= m.Values[pos]
	}

	return score
}
