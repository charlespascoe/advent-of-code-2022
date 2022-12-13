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

	prog, err := readInstructions(lines)
	if err != nil {
		log.Fatalf("Parse program failed: %s", err)
	}

	cpu := NewCPU(prog)
	crt := NewCRT(&cpu.X)
	clock := NewClock(crt, cpu)

	for !cpu.Done() {
		clock.Tick()
	}

	fmt.Printf("Result: %d (%d cycles)\n", cpu.Result, clock.cycle)

	fmt.Printf("Display:\n%s\n", crt)
}


func readInstructions(lines []string) ([]Instruction, error) {
	var prog []Instruction

	for _, line := range lines {
		inst, err := readInstruction(line)
		if err != nil {
			return nil, fmt.Errorf("invalid instruction '%s': %s", line, err)
		}

		prog = append(prog, inst)
	}

	return prog, nil
}

func readInstruction(line string) (Instruction, error) {
	if line == "noop" {
		return Instruction{Cycles: 1}, nil
	}

	l, r, found := strings.Cut(line, " ")
	if !found || l != "addx" {
		return Instruction{}, errors.New("must be noop or addx")
	}

	delta, err := strconv.Atoi(r)
	if err != nil {
		return Instruction{}, err
	}

	return Instruction{Cycles: 2, Delta: delta}, nil
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}
