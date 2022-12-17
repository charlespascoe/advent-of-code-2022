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
	// part2 := flag.Bool("part2", false, "print part 2 solution")
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

	s := &Solver{
		Map:   m,
		Start: "AA",
	}

	s.Solve()

	// instScore := InstantaneousScore(m)

	// var bestCost int

	// s := &Solver{
	// 	Map:      m,
	// 	ValvesOn: make(map[string]bool, len(m.Valves)),
	// 	BestCost: instScore,
	// }

	// if !*part2 {
	// 	// Part 1 //
	// 	// bestCost = BestSolution(m, "AA", make([]string, 0, len(m.Valves)), make(map[string]bool), false, false, 30, 0, instScore)
	// 	s.Solve("AA", 30, 0, 0)
	// 	// s.SolveMulti(Workers{
	// 	// 	{Pos: "AA", Arrival: 30},
	// 	// }, 30, 0, 0)
	// 	bestCost = s.BestCost
	// } else {
	// 	// Part 2 //
	// 	// bestCost = BestSolution(m, "AA", make([]string, 0, len(m.Valves)), make(map[string]bool), true, false, 26, CostOfUnopendValves(m, nil, 4), instScore)
	// 	s.WithElephant = true
	// 	trainingCost := s.CostOfUnopendValves(4)
	// 	// s.Solve("AA", 26, 0)
	// 	// s.SolveWithElephant()
	// 	s.SolveMulti(Workers{
	// 		{Pos: "AA", Arrival: 26},
	// 		{Pos: "AA", Arrival: 26},
	// 	}, 26, 0, 0)
	// 	bestCost = trainingCost + s.BestCost
	// }

	// fmt.Printf("Instantaneous Score: %d, Best Cost: %d, resulting in score: %d\n", instScore, bestCost, instScore-bestCost)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

// type Worker struct {
// 	Pos string
// 	// Arrival time in minutes remaining
// 	Arrival int
// }

// type Workers []Worker

// // Next TODO: Description.
// func (w Workers) Next() (Worker, int) {
// 	var next Worker
// 	idx := 0
// 	arrival := 0

// 	for i, wrkr := range w {
// 		if wrkr.Arrival >= arrival {
// 			next = wrkr
// 			idx = i
// 			arrival = wrkr.Arrival
// 		}
// 	}

// 	return next, idx
// }

// type Solver struct {
// 	Map      *Map
// 	ValvesOn map[string]bool
// 	BestCost int
// 	// LastUnopenedUpdate int

// 	WithElephant bool
// 	IsElepahnt   bool
// }

// func (s *Solver) SolveWithElephant() {
// 	me := make(map[string]bool, len(s.Map.Valves))
// 	el := make(map[string]bool, len(s.Map.Valves))

// 	best := s.BestCost
// 	fmt.Printf("STart: %d\n", best)

// 	vwr := make([]Valve, 0)
// 	for _, v := range s.Map.Valves {
// 		if v.Rate > 0 {
// 			vwr = append(vwr, v)
// 		}
// 	}

// 	for x := 0; x < 1<<(len(vwr)-1); x++ {
// 		for i, v := range vwr {
// 			me[v.Name] = (x & (1 << i)) != 0
// 			el[v.Name] = (x & (1 << i)) == 0
// 		}

// 		fmt.Printf("me: %v\nel: %v\n", me, el)

// 		s.BestCost = best

// 		s.ValvesOn = me
// 		s.Solve("AA", 26, 0, 0)
// 		fmt.Printf("best: %d\n", s.BestCost)
// 		cumCost := s.BestCost
// 		s.BestCost = best
// 		fmt.Printf("best 2: %d\n", s.BestCost)
// 		s.ValvesOn = el
// 		s.Solve("AA", 26, cumCost, 0)

// 		if s.BestCost < best {
// 			best = s.BestCost
// 		}
// 	}

// 	s.BestCost = best
// }

// func (s *Solver) Solve(pos string, minutesLeft, cumCost int, depth int) {
// 	if depth > 7 {
// 		fmt.Printf("Depth: %d\n", depth)
// 	}
// 	attempted := 0

// 	// if s.WithElephant && !s.IsElepahnt {
// 	// 	s.IsElepahnt = true
// 	// 	s.Solve("AA", 26, cumCost - s.CostOfUnopendValves(26 - minutesLeft))
// 	// 	s.IsElepahnt = false
// 	// }

// 	// Note: Valves are ordered by rate, so this loop always tries the
// 	// highest-rate valves first
// 	for _, v := range s.Map.Valves {
// 		if v.Rate == 0 || s.ValvesOn[v.Name] {
// 			continue
// 		}

// 		attempted++

// 		// The extra +1 is to account for the one minute to open the valve
// 		timeToOpen := s.Map.GetDist(pos, v.Name) + 1

// 		if timeToOpen >= minutesLeft {
// 			// Not enough time; total our cost for this solution and see how it
// 			// compares
// 			// solutionCost := cumCost + CostOfUnopendValves(m, valvesOn)
// 			cost := s.CostOfUnopendValves(minutesLeft)
// 			solutionCost := cumCost + cost

// 			if solutionCost < s.BestCost {
// 				s.BestCost = solutionCost
// 			}

// 			continue
// 		}

// 		// The additional cost from our total score from the start to the point
// 		// at which we'd have the valve open
// 		cost := s.CostOfUnopendValves(timeToOpen)

// 		if cumCost+cost >= s.BestCost {
// 			// This wouldn't be a better solution than the best we've seen;
// 			// skip it
// 			continue
// 		}

// 		s.ValvesOn[v.Name] = true
// 		s.Solve(v.Name, minutesLeft-timeToOpen, cumCost+cost, depth+1)
// 		s.ValvesOn[v.Name] = false
// 	}

// 	if attempted == 0 && cumCost < s.BestCost {
// 		// There were no more valves to open, which means cumCost won't get any
// 		// higher
// 		s.BestCost = cumCost
// 	}
// }

// func (s *Solver) SolveMulti(workers Workers, lastUpdate int, cumCost int, depth int) {
// 	// fmt.Printf("Depth: %d\n", depth)
// 	attempted := 0

// 	worker, wi := workers.Next()
// 	pos := worker.Pos
// 	minutesLeft := worker.Arrival

// 	// fmt.Printf("Workers: %v\n", workers)
// 	// Arrived and finished turning on valve; increment cost so far
// 	// fmt.Printf("W: %d, L: %d, A: %d, D: %d\n", wi, lastUpdate, minutesLeft, lastUpdate-minutesLeft)
// 	if minutesLeft < 0 {
// 		minutesLeft = 0
// 	}

// 	cumCost += s.CostOfUnopendValves(lastUpdate - minutesLeft)

// 	if minutesLeft <= 0 && cumCost < s.BestCost {
// 		s.BestCost = cumCost
// 		return
// 	}

// 	s.ValvesOn[pos] = true
// 	defer func() {
// 		s.ValvesOn[pos] = false
// 	}()

// 	// Note: Valves are ordered by rate, so this loop always tries the
// 	// highest-rate valves first
// 	for _, v := range s.Map.Valves {
// 		if v.Rate == 0 || s.ValvesOn[v.Name] {
// 			continue
// 		}

// 		attempted++

// 		// The extra +1 is to account for the one minute to open the valve
// 		timeToOpen := s.Map.GetDist(pos, v.Name) + 1
// 		if timeToOpen <= 0 {
// 			fmt.Printf("F: %s T: %s TTO: %d\n", pos, v.Name, timeToOpen)
// 			panic("INVALID TTO")
// 		}

// 		if timeToOpen < minutesLeft {
// 			// The additional cost from our total score from the start to the point
// 			// at which we'd have the valve open
// 			cost := v.Rate*timeToOpen //s.CostOfUnopendValves(timeToOpen)

// 			if cumCost+cost >= s.BestCost {
// 				// This wouldn't be a better solution than the best we've seen;
// 				// skip it
// 				continue
// 			}
// 			// Not enough time
// 			// continue
// 			// Not enough time; total our cost for this solution and see how it
// 			// compares
// 			// solutionCost := cumCost + CostOfUnopendValves(m, valvesOn)
// 			// cost := s.CostOfUnopendValves(minutesLeft)
// 			// solutionCost := cumCost + cost

// 			// if solutionCost < s.BestCost {
// 			// 	s.BestCost = solutionCost
// 			// }

// 			// continue
// 		}

// 		// s.ValvesOn[v.Name] = true
// 		workers[wi].Pos = v.Name
// 		workers[wi].Arrival = minutesLeft - timeToOpen
// 		// s.Solve(v.Name, minutesLeft-timeToOpen, cumCost+cost)
// 		s.SolveMulti(workers, minutesLeft, cumCost, depth+1)
// 		workers[wi] = worker
// 		// s.ValvesOn[v.Name] = false
// 	}

// 	if attempted == 0 && cumCost < s.BestCost {
// 		// There were no more valves to open, which means cumCost won't get any
// 		// higher
// 		s.BestCost = cumCost
// 	}
// }

// // CostOfUnopenedValves TODO: Description.
// func (s *Solver) CostOfUnopendValves(dur int) int {
// 	if dur < 0 {
// 		panic("dur cannot be negative")
// 	}
// 	sum := 0

// 	for _, v := range s.Map.Valves {
// 		if !s.ValvesOn[v.Name] {
// 			sum += v.Rate * dur
// 		}
// 	}

// 	return sum
// }

// // func BestSolution(m *Map, pos string, path []string, valvesOn
// // map[string]bool, withElephant, isElephant bool, minutesLeft, cumCost,
// // bestCost int) int {
// // 	// elapsed := 30 - minutesLeft

// // 	if withElephant && !isElephant {
// // 		newBestCost := BestSolution(m, "AA", path, valvesOn, withElephant, true, 26, cumCost-CostOfUnopendValves(m, valvesOn, 26-minutesLeft), bestCost)

// // 		if cumCost+newBestCost < bestCost {
// // 			bestCost = newBestCost
// // 		}
// // 	}

// // 	attempted := 0

// // 	// Note: Valves are ordered by rate, so this loop always tries the
// // 	// highest-rate valves first
// // 	for _, v := range m.Valves {
// // 		if v.Rate == 0 || valvesOn[v.Name] {
// // 			continue
// // 		}

// // 		attempted++

// // 		// The extra +1 is to account for the one minute to open the valve
// // 		timeToOpen := m.GetDist(pos, v.Name) + 1

// // 		if timeToOpen >= minutesLeft {
// // 			// Not enough time; total our cost for this solution and see how it
// // 			// compares
// // 			// solutionCost := cumCost + CostOfUnopendValves(m, valvesOn)
// // 			cost := CostOfUnopendValves(m, valvesOn, minutesLeft)
// // 			solutionCost := cumCost + cost

// // 			if solutionCost < bestCost {
// // 				bestCost = solutionCost
// // 			}

// // 			continue
// // 		}

// // 		// The additional cost from our total score from the start to the point
// // 		// at which we'd have the valve open
// // 		cost := CostOfUnopendValves(m, valvesOn, timeToOpen)

// // 		if cumCost+cost >= bestCost {
// // 			// This wouldn't be a better solution than the best we've seen;
// // 			// skip it
// // 			continue
// // 		}

// // 		valvesOn[v.Name] = true
// // 		newBestCost := BestSolution(m, v.Name, append(path, v.Name), valvesOn, withElephant, isElephant, minutesLeft-timeToOpen, cumCost+cost, bestCost)
// // 		valvesOn[v.Name] = false

// // 		if newBestCost < bestCost {
// // 			bestCost = newBestCost
// // 		}
// // 	}

// // 	if attempted == 0 && cumCost < bestCost {
// // 		// There were no more valves to open, which means cumCost won't get any
// // 		// higher
// // 		bestCost = cumCost
// // 	}

// // 	return bestCost
// // }

// // func CostOfUnopendValves(m *Map, valvesOn map[string]bool, dur int) int {
// // 	sum := 0

// // 	for _, v := range m.Valves {
// // 		if !valvesOn[v.Name] {
// // 			sum += v.Rate * dur
// // 		}
// // 	}

// // 	return sum
// // }

// func InstantaneousScore(m *Map) int {
// 	score := 0

// 	for _, v := range m.Valves {
// 		score += 30 * v.Rate
// 	}

// 	return score
// }
