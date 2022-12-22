package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var verbosity = flag.Int("v", 0, "verbose output")

func main() {
	input := flag.String("i", "input.txt", "program input")
	part2 := flag.Bool("part2", false, "print part 2 solution")
	flag.Parse()

	blueprints, err := readInput(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	// fmt.Printf("Blueprints: %v\n", blueprints)
	if *part2 {
		solvePart2(blueprints)
	} else {
		solvePart1(blueprints)
	}
}

func solvePart1(blueprints []Blueprint) {
	result := 0

	var wg sync.WaitGroup
	c := make(chan int)

	for _, blueprint := range blueprints {
		wg.Add(1)

		go func(sim *Simulator) {
			defer wg.Done()

			finalState := sim.solve3()

			fmt.Printf("Final State: %s\n", finalState)

			geodes := finalState.resources[GeodeResource]
			score := sim.blueprint.Number * geodes
			fmt.Printf("Blueprint %d geodes: %d, quality score: %d\n", sim.blueprint.Number, geodes, score)

			c <- score
		}(NewSimulator(blueprint, 24))
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for val := range c {
		result += val
	}

	fmt.Printf("\nResult: %d\n", result)
}

func solvePart2(blueprints []Blueprint) {
	result := 1

	var wg sync.WaitGroup
	c := make(chan int)

	for _, blueprint := range blueprints[:3] {
		wg.Add(1)

		go func(sim *Simulator) {
			defer wg.Done()

			finalState := sim.solve3()

			fmt.Printf("Final State: %s\n", finalState)

			geodes := finalState.resources[GeodeResource]
			// score := sim.blueprint.Number * geodes
			fmt.Printf("Blueprint %d geodes: %d\n", sim.blueprint.Number, geodes)

			c <- geodes
		}(NewSimulator(blueprint, 32))
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for val := range c {
		result *= val
	}

	fmt.Printf("\nResult: %d\n", result)
}

func readInput(path string) ([]Blueprint, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		blueprints := make([]Blueprint, 0, len(lines))

		for _, line := range lines {
			blueprints = append(blueprints, ParseBlueprint(line))
		}

		return blueprints, nil
	}
}

func vprintf(level int, format string, args ...any) {
	if level <= *verbosity {
		fmt.Printf(format, args...)
	}
}
