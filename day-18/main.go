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
	flag.Parse()

	input, err := readInput(*inputFile)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	droplet := BuildDroplet(input)

	extSurfaceArea := droplet.CalcExternalSurfaceArea()
	droplet.Print()
	fmt.Printf("Total Surface Area: %d\n", droplet.surfaceArea)
	fmt.Printf("External Surface Area: %d\n", extSurfaceArea)
}

func readInput(path string) ([]Vector, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		var input []Vector
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")

		for _, line := range lines {
			input = append(input, ParsePos(line))
		}

		return input, nil
	}
}

func ParsePos(str string) Vector {
	parts := strings.SplitN(str, ",", 3)

	return Vector{
		X: MustAtoi(parts[0]),
		Y: MustAtoi(parts[1]),
		Z: MustAtoi(parts[2]),
	}
}
