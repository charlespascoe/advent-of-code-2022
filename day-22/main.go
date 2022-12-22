package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// var verbose = flag.Bool("v", false, "verbose output")
// var part2 = flag.Bool("part2", false, "print part 2 solution")

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	lines, moves, err := readInput(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	m, start := BuildMap(lines)
	tile := start

	dir := Right
	actions := regexp.MustCompile(`\D+|\d+`).FindAllString(moves, -1) 

	for _, action := range actions {
		count := 0
		switch action {
		case "L":
			dir = dir.Turn(Left)
			continue
		case "R":
			dir = dir.Turn(Right)
			continue
		default:
			count = MustAtoi(action)
		}

		for i := 0; i < count; i++ {
			tile.Mark = dir.String()[0]

			next := tile.Move(dir)

			if next.Wall {
				break
			}

			tile = next
		}
	}

	password := 1000*(tile.Row + 1) + 4*(tile.Col+1) + int(dir.Turn(Left))

	fmt.Printf("Map:\n%s\n\nMoves: %s\nPassword: %d\n", m, moves, password)
}

func readInput(path string) ([]string, string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, "", err
	} else {
		mapStr, moves, found := strings.Cut(string(data), "\n\n")
		if !found {
			return nil, "", errors.New("invalid input")
		}

		return strings.Split(mapStr, "\n"), strings.TrimSpace(moves), nil
	}
}

func MustAtoi(str string) int {
	if x, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return x
	}
}
