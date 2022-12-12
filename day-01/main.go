package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"log"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	curElf := 0
	topElves := make([]int, 3)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			addMax(topElves, curElf)

			curElf = 0
			continue
		}

		cal, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("Not a valid number: %s", err)
		}

		curElf += cal
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scan error: %s", err)
	}

	fmt.Printf("Max: %v, Sum: %d\n", topElves, sum(topElves))
}

func addMax(top []int, cal int) {
	c := cal

	for i := 0; i < len(top); i++ {
		if c > top[i] {
			top[i], c = c, top[i]
		}
	}
}

func sum(nums []int) int {
	s := 0

	for _, x := range nums {
		s += x
	}

	return s
}
