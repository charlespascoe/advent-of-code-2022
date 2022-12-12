package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	state, movements, err := parseInput(lines)
	if err != nil {
		log.Fatalf("Invalid input: %s", err)
	}

	fmt.Printf("State:\n%s\n", state)

	for _, mvmt := range movements {
		fmt.Println(mvmt)

		// This is the part 1 solution
		// state.ExecuteIndividualMoves(mvmt)

		// This is the part 2 solution
		state.ExecuteBulkMoves(mvmt)
	}

	fmt.Printf("Final State:\n%s\n", state)
	fmt.Printf("Top Boxes: %s\n", state.TopBoxes())
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), nil
	}
}

type State struct {
	Stacks [][]string
}

func (s State) TopBoxes() string {
	var topBoxes []string

	for _, stack := range s.Stacks {
		topBoxes = append(topBoxes, stack[len(stack)-1])
	}

	return strings.Join(topBoxes, "")
}

func (s *State) ExecuteIndividualMoves(move Movement) {
	for i := 0; i < move.Count; i++ {
		var box string
		// Note that From and To are zero-indexed
		s.Stacks[move.From], box = pop(s.Stacks[move.From])
		s.Stacks[move.To] = append(s.Stacks[move.To], box)
	}
}

func (s *State) ExecuteBulkMoves(move Movement) {
	var boxes []string
	s.Stacks[move.From], boxes = popN(s.Stacks[move.From], move.Count)
	s.Stacks[move.To] = append(s.Stacks[move.To], boxes...)
}

func (s State) String() string {
	var lines []string

	for i, stack := range s.Stacks {
		line := []byte(strconv.Itoa(i+1) + ": ")

		for _, box := range stack {
			line = append(line, []byte(fmt.Sprintf(" [%s]", box))...)
		}

		lines = append(lines, string(line))
	}

	return strings.Join(lines, "\n")
}

type Movement struct {
	// From and To are zero-indexed
	Count, From, To int
}

func (m Movement) String() string {
	return fmt.Sprintf("move %d from %d to %d", m.Count, m.From+1, m.To+1)
}

func parseInput(lines []string) (State, []Movement, error) {
	split := -1

	for i, line := range lines {
		if line == "" {
			split = i
			break
		}
	}

	var state State
	var mvmts []Movement
	var err error

	if split <= 0 {
		return state, mvmts, errors.New("couldn't find movement start")
	}

	if state, err = parseState(lines[:split]); err != nil {
		return state, mvmts, err
	}

	if mvmts, err = parseMovement(lines[split+1:]); err != nil {
		return state, mvmts, err
	}

	return state, mvmts, nil
}

func parseState(lines []string) (State, error) {
	boxRe := regexp.MustCompile(`\[([A-Z])\]`)

	// Get number of stacks from last line
	stacks := len(regexp.MustCompile(`\d+`).FindAllString(lines[len(lines)-1], -1))

	state := State{
		Stacks: make([][]string, stacks),
	}

	// Skipping len(lines) - 1 since it's just the stack numbers
	for i := len(lines) - 2; i >= 0; i-- {
		for s := 0; s < len(state.Stacks); s++ {
			// This is a weird way of doing this; improve it?
			box := lines[i][4*s : 4*s+3]

			if match := boxRe.FindStringSubmatch(box); len(match) > 0 {
				state.Stacks[s] = append(state.Stacks[s], match[1])
			}
		}
	}

	return state, nil
}

func parseMovement(lines []string) ([]Movement, error) {
	mvRe := regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)

	var cmds []Movement

	for _, line := range lines {
		match := mvRe.FindStringSubmatch(line)

		if len(match) == 0 {
			return nil, fmt.Errorf("invalid move command: %s", line)
		}

		count, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}

		from, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}

		to, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, err
		}

		// Note that From and To are zero-indexed
		cmds = append(cmds, Movement{count, from - 1, to - 1})
	}

	return cmds, nil
}

func pop(s []string) ([]string, string) {
	v := s[len(s)-1]
	return s[:len(s)-1], v
}

func popN(s []string, n int) (rest, popped []string) {
	return s[:len(s)-n], s[len(s)-n:]
}
