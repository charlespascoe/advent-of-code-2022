package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var verbosity = flag.Int("v", 0, "verbose output")

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	blueprints, err := readInput(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	// fmt.Printf("Blueprints: %v\n", blueprints)

	for _, blueprint := range blueprints {
		sim := NewSimulator(blueprint)

		schedule, finalState := sim.solve()

		schedule.print()

		fmt.Printf("Optimal Build Schedule: %s\n", schedule)
		fmt.Printf("Final State: %s\n", finalState)
		// fmt.Printf("Robots: %v\n", sim.robots)
		// fmt.Printf("Resources: %v\n", sim.resources)
		// fmt.Println()

		// fmt.Printf("\nBlueprint %d Geodes: %d\n", blueprint.Number, geodes)
		break
	}
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
