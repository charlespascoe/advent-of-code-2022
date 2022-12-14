package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const MaxInt64 = 0x7fffffffffffffff

func main() {
	input := flag.String("i", "input.txt", "program input")
	part2 := flag.Bool("part2", false, "print part 2 solution")
	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	m, start, end := parseMap(lines)

	fmt.Println("Map:")
	render(m, start, end, nil)
	fmt.Println()

	if !*part2 { // Part1
		path := ShortestPath(m, start, end)

		// Path includes both endpoints, and so is one longer than the number of
		// steps
		fmt.Printf("Shortest Path: %d steps\n", len(path)-1)
		render(m, start, end, path)
	} else { // Part 2
		newStart, path := ShortestPathFromLowestPoint(m, end)

		// Path includes both endpoints, and so is one longer than the number of
		// steps
		fmt.Printf("Best start: %v\nShortest Path: %d steps\n", newStart, len(path)-1)
		render(m, newStart, end, path)
	}
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}


func render(m Map, start, end Pos, path []Pos) {
	var buf [][]byte

	for _, row := range m {
		var bufRow []byte
		for _, cell := range row {
			chr := byte(cell.Height) + 'a'
			if len(path) > 0 {
				chr = '.'
			}
			bufRow = append(bufRow, chr)
		}
		buf = append(buf, bufRow)
	}

	buf[start.Row][start.Col] = 'S'
	buf[end.Row][end.Col] = 'E'

	prev := start
	for _, step := range path {
		buf[prev.Row][prev.Col] = prev.ArrowTo(step)
		prev = step
	}

	for _, row := range buf {
		fmt.Printf("%s\n", row)
	}
}
