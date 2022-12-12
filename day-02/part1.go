package main

import (
	"bufio"
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
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	score := 0

	for scanner.Scan() {
		line := scanner.Text()

		left, right, found := strings.Cut(line, " ")
		if !found {
			log.Fatalf("Couldn't split line '%s'", line)
		}

		// Use errors instead?
		op := mustReadShape(left)
		me := mustReadShape(right)

		score += int(me) + int(round(op, me))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scan error: %s", err)
	}

	fmt.Printf("Score: %d\n", score)
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

func round(op, me Shape) Outcome {
	// The '3 +' is just to make the calculation always positive
	switch (3 + me - op) % 3 {
	case 0:
		return Draw
	case 1:
		return Win
	case 2:
		return Lose
	default:
		// What's a better way of avoiding this?
		panic("Unexpected value")
	}
}
