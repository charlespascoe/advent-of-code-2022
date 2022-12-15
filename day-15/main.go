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
	yScan := flag.Int("y", 2_000_000, "the y level to scan (part 1)")
	scanBounds := flag.Int("b", 4_000_000, "x/y bounds for scanning (part 2)")
	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	sensors, err := ParseSensors(lines)
	if err != nil {
		log.Fatalf("Invalid input: %s", err)
	}

	// --- Part 1 --- //

	minX := sensors[0].Pos.X
	maxX := minX

	// Find the minimum and maximum X values that we need to iterate over
	for _, sensor := range sensors {
		minX, maxX = XBounds(sensor, minX, maxX)
	}

	fmt.Printf("Min X: %d, Max X: %d\n", minX, maxX)

	inRange := 0

	for x := minX; x <= maxX; x++ {
		for _, sensor := range sensors {
			pos := Vector{x, *yScan}
			if sensor.Beacon != pos && sensor.InRange(pos) {
				inRange++
				break
			}
		}
	}

	fmt.Printf("Points along y=%d in range: %d\n", *yScan, inRange)

	// --- Part 2 --- //

	for y := 0; y <= *scanBounds; y++ {
	xloop:
		for x := 0; x <= *scanBounds; x++ {
			for _, sensor := range sensors {
				if sensor.InRange(Vector{x, y}) {
					x = sensor.RangeRightEdge(y)
					continue xloop
				}
			}

			fmt.Printf("Position: %d %d, tuning frequency: %d\n", x, y, 4_000_000*x+y)
		}
	}
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
