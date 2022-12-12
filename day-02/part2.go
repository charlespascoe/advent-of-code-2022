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
		out := mustReadOutcome(right)

		score += int(chooseShape(op, out)) + int(out)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scan error: %s", err)
	}

	fmt.Printf("Score: %d\n", score)
}

func mustReadShape(s string) Shape {
	switch s {
	case "A":
		return Rock
	case "B":
		return Paper
	case "C":
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
	delta := 0

	switch out {
	case Draw:
		delta = 0
	case Win:
		delta = 1
	case Lose:
		delta = 2
	}

	me := (int(op) + delta) % 3

	if me == 0 {
		me = 3
	}

	return Shape(me)
}
