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

	var cave Cave
	pos := len(cave) + 3
	inputIdx := 0
	rockCount := 0
	rock := rocks[0]

	limit := 2022
	if *part2 {
		limit = 1_000_000_000_000
	}

	for rockCount < limit {
		jet := input[inputIdx%len(input)]

		// Print(cave, pos, rock, fmt.Sprintf("Before move '%c'", jet))

		if movedSh := rock.Shift(jet == '<'); !movedSh.OverlapsWith(cave.RowsFrom(pos)) {
			// There's nothing blocking us to move this way
			rock = movedSh

			// Print(cave, pos, movedSh, "After Move")
		}

		if pos > 0 && !rock.OverlapsWith(cave.RowsFrom(pos-1)) {
			pos--
		} else {
			// Rock has stopped
			cave = cave.EnsureRows(pos + len(rock))
			rock.WriteTo(cave.RowsFrom(pos))
			// Print(cave, pos, rock, fmt.Sprintf("Rock stopped - cave after rock %d", rockCount+1))

			// Next rock
			rockCount++
			rock = rocks[rockCount%len(rocks)]
			pos = len(cave) + 3

			// fmt.Printf("\nCave after shape %d:\n%s\n", rockCount, cave)
		}

		inputIdx++
	}

	// fmt.Printf("\nCave after shape %d:\n%s\n", rockCount, cave)
	fmt.Printf("Height after %d rocks: %d\n", rockCount, len(cave))
}

func Print(c Cave, pos int, s Rock, msg string) {
	c = c.Copy().EnsureRows(pos+len(s))
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
