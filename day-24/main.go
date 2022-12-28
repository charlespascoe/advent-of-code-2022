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

	valley := NewValley(lines)

	for i := 0; i < 10; i++ {
		fmt.Printf("Valley at %d:\n%s\n\n", i, valley.StringAtTime(i))
	}

	fmt.Printf("Valley Dimensions: %d rows by %d columns\n", valley.Rows(), valley.Cols())
	fmt.Printf("Blizzard sequence period: %d\n", valley.Period)

	pathEnd := NewSolver(valley, valley.Start, valley.End, 0).Solve()

	fmt.Printf("Shortest path time: %d\n", pathEnd.Time)

	returnEnd := NewSolver(valley, valley.End, valley.Start, pathEnd.Time).Solve()
	fmt.Printf("Shorest return time: %d\n", returnEnd.Time)

	endAgain := NewSolver(valley, valley.Start, valley.End, returnEnd.Time).Solve()
	fmt.Printf("Shorest time back to end: %d\n", endAgain.Time)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
