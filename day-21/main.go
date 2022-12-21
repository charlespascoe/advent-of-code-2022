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
	part2 := flag.Bool("part2", false, "print part 2 solution")

	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	exprs := buildExprMap(lines)
	var result int

	if !*part2 {
		// Part 1 //
		result = exprs.Eval("root")
	}

	fmt.Printf("Result: %d\n", result)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
