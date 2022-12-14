package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Shape int

const (
	Rock     Shape = 1
	Paper    Shape = 2
	Scissors Shape = 3
)

type Outcome int

const (
	Win  Outcome = 6
	Draw Outcome = 3
	Lose Outcome = 0
)

func main() {
	part2 := flag.Bool("part2", false, "print part 2 solution")
	flag.Parse()

	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	totScore := 0

	for _, line := range lines {
		left, right, found := strings.Cut(line, " ")
		if !found {
			log.Fatalf("Couldn't split line '%s'", line)
		}

		op := mustReadShape(left)

		if !*part2 {
			// Part 1
			me := mustReadShape(right)
			totScore += score(me, round(op, me))
		} else {
			// Part 2
			out := mustReadOutcome(right)
			totScore += score(chooseShape(op, out), out)
		}
	}

	fmt.Printf("Score: %d\n", totScore)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func mustReadShape(s string) Shape {
	switch s {
	case "A", "X":
		return Rock
	case "B", "Y":
		return Paper
	case "C", "Z":
		return Scissors
	default:
		// Return error instead?
		panic("Invalid shape")
	}
}

func mustReadOutcome(s string) Outcome {
	switch s {
	case "X":
		return Lose
	case "Y":
		return Draw
	case "Z":
		return Win
	default:
		// Return error instead?
		panic("Invalid outcome")
	}
}

func chooseShape(op Shape, out Outcome) Shape {
	// Outcome:  0: Lose, 3: Draw, 6: Win
	// Div by 3: 0: Lose, 1: Draw, 2: Win
	// Sub 1:    0: Draw, 1: Win,  2: Lose
	// Add to shape
	// (Note that -1 ≡ 2 mod 3)
	delta := ((int(out) / 3) + 2) % 3

	me := (int(op) + delta) % 3

	if me == 0 {
		me = 3
	}

	return Shape(me)
}

func round(op, me Shape) Outcome {
	// The '3 +' is just to make the 'me - op mod 3' calculation always positive
	// 'me - op mod 3' => 0: Draw, 1: Win,  2: Lose
	// 1+ then makes it:  0: Lose, 1: Draw, 2: Win
	// 3* then makes it:  0: Lose, 3: Draw, 6: Win
	return Outcome(3 * ((3 + 1 + me - op) % 3))
}

func score(me Shape, out Outcome) int {
	return int(me) + int(out)
}
