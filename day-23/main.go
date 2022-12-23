package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	m := NewMap(lines)

	fmt.Printf("Map:\n%s\n", m)

	round := 0

	// All elves just moved into their starting positions
	moved := len(m.Elves)

	for moved > 0 {
		for _, elf := range m.Elves {
			elf.PlanMove(round)
		}

		moved = 0
		for _, elf := range m.Elves {
			if elf.Move() {
				moved++
			}
		}

		round++

		if round == 10 {
			fmt.Printf("Empty tiles after %d rounds: %d\n", round, m.CountEmpty())
		}
	}

	fmt.Printf("Number of Rounds: %d\n", round)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
