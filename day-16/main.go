package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	valves, err := ParseValves(lines)
	if err != nil {
		log.Fatalf("Couldn't read input: %s", err)
	}

	// Sort valves into descending order by rate
	sort.Slice(valves, func(i, j int) bool {
		return valves[i].Rate > valves[j].Rate
	})

	fmt.Printf("Valves: %v\n", valves)

	m := NewMap(valves)

	m.ComputeShortestPaths()

	fmt.Printf("Map: %v\n", m)

	// Part 1 //

	start := m.ValveLookup["AA"]

	var nonzeroValves []Valve
	var unopened int

	for i, valve := range m.Valves {
		if valve.Rate == 0 {
			nonzeroValves = m.Valves[:i]
			break
		}

		unopened |= 1 << i
	}

	startTime := time.Now()
	score := Solve(m, nonzeroValves, start, unopened, 30, 0, 0)

	fmt.Printf("Part 1 Score: %d, Dur: %s\n", score, time.Now().Sub(startTime))

	// Part 2 //

	bound := (unopened >> 1) + 1

	best := 0

	startTime = time.Now()

	for i := 1; i < bound; i++ {
		score := Solve(m, nonzeroValves, start, i, 26, 0, 0)
		score += Solve(m, nonzeroValves, start, unopened^i, 26, 0, 0)

		if score > best {
			best = score
		}
	}

	fmt.Printf("Part 2 Score: %d, Dur: %s\n", best, time.Now().Sub(startTime))
}

func Solve(m *Map, valves []Valve, curValve Valve, unopened int, timeRem, score, best int) int {
	if timeRem <= 0 {
		return score
	}

	score += timeRem * curValve.Rate

	if unopened == 0 && score > best {
		return score
	}

	for i, valve := range valves {
		if unopened&(1<<i) == 0 {
			continue
		}

		nextScore := Solve(
			m,
			valves,
			valve,
			unopened&^(1<<i),
			timeRem-m.GetDist(curValve.Name, valve.Name)-1,
			score,
			best,
		)

		if nextScore > best {
			best = nextScore
		}
	}

	return best
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
