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

var part2 = flag.Bool("part2", false, "print part 2 solution")

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	lines, moves, err := readInput(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	m := NewMap(lines)
	var nav Navigator

	if *part2 {
		nav = m.CubeNavigator()
	} else {
		nav = m.MapNavigator()
	}

	actions := regexp.MustCompile(`\D|\d+`).FindAllString(moves, -1)

	for _, action := range actions {
		count := 0
		switch action {
		case "L":
			nav.Turn(Left)
			continue
		case "R":
			nav.Turn(Right)
			continue
		default:
			count = MustAtoi(action)
		}

		for count > 0 && nav.Move() {
			count--
		}
	}

	pos := nav.Pos()
	dir := nav.Dir()

	password := 1000*(pos.Row+1) + 4*(pos.Col+1) + int(dir.Turn(Left))

	fmt.Printf("Map:\n%s\n\nMoves: %s\nPosition: %d, %d\nPassword: %d\n", m, moves, pos.Row+1, pos.Col+1, password)
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
