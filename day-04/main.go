package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	ovCount := 0
	fcCount := 0

	for _, line := range lines {
		elf1, elf2, found := strings.Cut(line, ",")
		if !found {
			log.Fatalf("Couldn't split assignments in string '%s'", line)
		}

		rng1, err := parseRange(elf1)
		if err != nil {
			log.Fatalf("Invalid range '%s': %s", elf1, err)
		}

		rng2, err := parseRange(elf2)
		if err != nil {
			log.Fatalf("Invalid range '%s': %s", elf2, err)
		}

		if overlap(rng1, rng2) {
			ovCount++

			if fullyContained(rng1, rng2) {
				fcCount++
			}
		}
	}

	fmt.Printf("Overlap count: %d, Fully Contained Count: %d\n", ovCount, fcCount)
}

type Range struct {
	Start, End int
}

func (r Range) Has(x int) bool {
	return r.Start <= x && x <= r.End
}

func parseRange(str string) (rng Range, err error) {
	s, e, found := strings.Cut(str, "-")
	if !found {
		return Range{}, errors.New("range sep not found")
	}

	if rng.Start, err = strconv.Atoi(s); err != nil {
		return
	}

	if rng.End, err = strconv.Atoi(e); err != nil {
		return
	}

	return
}

func overlap(r1, r2 Range) bool {
	return r1.Has(r2.Start) || r1.Has(r2.End) || r2.Has(r1.Start) || r2.Has(r1.End)
}

func fullyContained(r1, r2 Range) bool {
	return (r1.Has(r2.Start) && r1.Has(r2.End)) || (r2.Has(r1.Start) && r2.Has(r1.End))
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
