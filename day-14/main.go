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

	structs := ParseStructures(lines)

	fmt.Printf("%v\n", structs)

	cave := NewCave()

	for _, st := range structs {
		cave.DrawStructure(st)
	}

	fmt.Printf("Cave:\n%s\n", cave)

	count := 0

	spawn := Vector{500, 0}

	for !simulateGrain(cave, spawn) && cave.IsEmpty(spawn) {
		count++
	}

	fmt.Printf("\nCave after simulation:\n%s\n", cave)

	fmt.Printf("Grains at rest: %d\n", count)

	cave.AddFloor()

	fmt.Printf("Cave with floor:\n%s\n", cave)

	for !simulateGrain(cave, spawn) {
		count++
	}

	fmt.Printf("Cave after simulation with floor:\n%s\n", cave)

	fmt.Printf("Grains at rest: %d\n", count)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func simulateGrain(cave *Cave, spawn Vector) bool {
	pos := spawn

mainloop:
	for pos.Y < cave.Depth() && cave.IsEmpty(spawn) {
		for _, move := range []Vector{Down, DownLeft, DownRight} {
			if cave.IsEmpty(pos.Add(move)) {
				pos = pos.Add(move)
				continue mainloop
			}
		}

		cave.Draw(pos, 'o')
		return false
	}

	return true
}
