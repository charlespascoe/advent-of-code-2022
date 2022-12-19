package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	inputFile := flag.String("i", "input.txt", "program input")
	part2 := flag.Bool("part2", false, "print part 2 solution")
	flag.Parse()

	input, err := readFile(*inputFile)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	fmt.Printf("Jets: %d\n", len(input))

	sim := NewSimulator(input)

	limit := 2022
	if !*part2 {
		for sim.RocksSettled() < limit {
			sim.Step()
		}
	} else {
		limit = 1_000_000_000_000

		// TODO: Implement
	}

	fmt.Printf("Height after %d rocks: %d\n", sim.RocksSettled(), sim.Height())
}

func Print(c Cave, pos int, s Rock, msg string) {
	c = c.Copy().EnsureRows(pos + len(s))
	s.WriteTo(c.RowsFrom(pos))
	fmt.Printf("\n%s:\n%s\n", msg, c)
}

func readFile(path string) (string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(string(data)), nil
	}
}
