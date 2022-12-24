package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// var verbose = flag.Bool("v", false, "verbose output")
// var part2 = flag.Bool("part2", false, "print part 2 solution")

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	val := NewValley(lines)

	for i := 0; i < 10; i++ {
		fmt.Printf("Valley at %d:\n%s\n\n", i, val.StringAtTime(i))
	}
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
