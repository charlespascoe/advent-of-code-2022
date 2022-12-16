package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
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

	instScore := InstantaneousScore(m)

	bestCost := BestSolution(m, "AA", make([]string, 0, len(m.Valves)), make(map[string]bool), 30, 0, instScore)

	fmt.Printf("Instantaneous Score: %d, Best Cost: %d, resulting in score: %d\n", instScore, bestCost, instScore-bestCost)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func BestSolution(m *Map, pos string, path []string, valvesOn map[string]bool, minutesLeft, cumCost, bestCost int) int {
	elapsed := 30 - minutesLeft

	attempted := 0

	// Note: Valves are ordered by rate, so this loop always tries the
	// highest-rate valves first
	for _, v := range m.Valves {
		if v.Rate == 0 || valvesOn[v.Name] {
			continue
		}

		attempted++

		// fmt.Printf("Path: %v, Trying %s...\n", path, v.Name)

		// The extra +1 is to account for the one minute to open the valve
		timeToOpen := m.GetDist(pos, v.Name) + 1

		// fmt.Printf("TTO: %d, Rem: %d\n", timeToOpen, minutesLeft)

		if timeToOpen >= minutesLeft {
			// Not enough time; total our cost for this solution and see how it
			// compares
			// solutionCost := cumCost + CostOfUnopendValves(m, valvesOn)
			cost :=  CostOfUnopendValves(m, valvesOn)
			solutionCost := cumCost + cost

			// fmt.Printf(">>>>>>>>>> CumCost: %d, Cost: %d, Sum: %d, BestCost: %d\n", cumCost, cost, cumCost + cost, bestCost)

			if solutionCost < bestCost {
				bestCost = solutionCost
			}

			continue
		}

		// The cost from our total score from the start to the point at which
		// we'd have the valve open
		cost := v.Rate * (elapsed + timeToOpen)


		if cumCost+cost >= bestCost {
			// This wouldn't be a better solution than the best we've seen;
			// skip it
			continue
		}

		valvesOn[v.Name] = true
		newBestCost := BestSolution(m, v.Name, append(path, v.Name), valvesOn, minutesLeft-timeToOpen, cumCost+cost, bestCost)
		valvesOn[v.Name] = false

		if newBestCost < bestCost {
			bestCost = newBestCost
		}
	}

	if attempted == 0 && cumCost < bestCost {
		// There were no more valves to open, which means cumCost won't get any
		// higher
		bestCost = cumCost
	}

	return bestCost
}

func CostOfUnopendValves(m *Map, valvesOn map[string]bool) int {
	sum := 0

	for _, v := range m.Valves {
		if !valvesOn[v.Name] {
			sum += v.Rate * 30
		}
	}

	return sum
}

func InstantaneousScore(m *Map) int {
	score := 0

	for _, v := range m.Valves {
		score += 30 * v.Rate
	}

	return score
}
