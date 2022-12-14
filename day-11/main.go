package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var part2 = flag.Bool("part2", false, "print part 2 solution")

func main() {
	flag.Parse()

	lines, err := readLines("test.txt")

	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	_, statements, err := ParseStatements(0, lines)
	if err != nil {
		log.Fatalf("Invalid input: %s", err)
	}

	for _, stmt := range statements {
		fmt.Printf("%s\n", stmt)
	}

	monkeys, err := ParseMonkeys(statements)
	if err != nil {
		log.Fatalf("Couldn't parse monkey data: %s", err)
	}

	fmt.Printf("%v\n", monkeys)
	rounds := 20

	if *part2 {
		rounds = 10_000
	}

	for i := 0; i < rounds; i++ {
		for _, monkey := range monkeys {
			monkey.Round()
		}

		if !*part2 {
			fmt.Printf("\nAfter round %d:\n%v\n", i+1, monkeys)
		}
	}

	fmt.Println()

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].InspectCount < monkeys[j].InspectCount
	})

	for _, monkey := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", monkey.Num, monkey.InspectCount)
	}
}

func readLines(path string) ([]Line, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return toLines(strings.TrimSpace(string(data))), nil
	}
}
