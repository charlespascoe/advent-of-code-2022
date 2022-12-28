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

	height := 0

	if !*part2 {
		height = solvePart1(input)
	} else {
		height = solvePart2(input)
	}

	fmt.Printf("Height: %d\n", height)
}

func solvePart1(input string) int {
	sim := NewSimulator(input)

	limit := 2022

	for sim.RocksSettled() < limit {
		sim.Step()
	}

	return sim.Height()
}

func solvePart2(input string) int {
	// Although this is solution is functionally identical to part 1, it is
	// provided separately for clarity

	limit := 1_000_000_000_000

	tortoise := NewSimulator(input)
	hare := NewSimulator(input)

	hare.Step()
	hare.Step()
	tortoise.Step()

	for !tortoise.HasSameState(hare) {
		hare.Step()
		hare.Step()
		tortoise.Step()
	}

	tortoise = NewSimulator(input)

	for !tortoise.HasSameState(hare) {
		hare.Step()
		tortoise.Step()
	}

	hare.Step()

	for !tortoise.HasSameState(hare) {
		hare.Step()
	}

	rocksPerLoop := hare.RocksSettled() - tortoise.RocksSettled()
	heightPerLoop := hare.Height() - tortoise.Height()

	loops := (limit - tortoise.RocksSettled()) / rocksPerLoop
	remaining := (limit - tortoise.RocksSettled()) % rocksPerLoop

	newLimit := tortoise.RocksSettled() + remaining

	fmt.Printf("Tortoise: %d, %d; Hare: %d, %d\n", tortoise.RocksSettled(), tortoise.Height(), hare.RocksSettled(), hare.Height())
	fmt.Printf("Rocks/Loop: %d, Height/Loop: %d\n", rocksPerLoop, heightPerLoop)
	fmt.Printf("Tortoise Cave:\n%s\nHare Cave:\n%s\n", tortoise.cave.TopUntilBlocked(), hare.cave.TopUntilBlocked())
	fmt.Printf("Loops: %d, Remaining: %d, New Limit: %d\n", loops, remaining, newLimit)

	for tortoise.RocksSettled() < newLimit {
		tortoise.Step()
	}

	return tortoise.Height() + heightPerLoop*loops
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
