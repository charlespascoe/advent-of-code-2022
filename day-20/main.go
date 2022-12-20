package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)


func main() {
	inputFile := flag.String("i", "input.txt", "program input")
	part2 := flag.Bool("part2", false, "print part 2 solution")

	flag.Parse()

	input, err := readInput(*inputFile)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	fmt.Printf("Input: %v\n", input)

	indices := make([]int, 0, len(input))

	for i := range input {
		indices = append(indices, i)

		if *part2 {
			input[i] *= 811589153
		}
	}

	mixCount := 1
	if *part2 {
		mixCount = 10
	}

	for i := 0; i < mixCount; i++ {
		mix(input, indices)
	}

	zeroPos := 0

	for i, idx := range indices {
		if input[idx] == 0 {
			zeroPos = i
			break
		}
	}

	fmt.Printf("Zero pos: %d (%d)\n", zeroPos, input[indices[zeroPos]])
	result := getval(input, indices, zeroPos+1000) + getval(input, indices, zeroPos+2000) + getval(input, indices, zeroPos+3000)
	fmt.Printf("Result: %d\n", result)
}

func mix(input, indices []int) {
	for i, val := range input {
		from := getpos(i, indices)
		to := (from + val) % (len(indices) - 1)
		move(from, to, indices)
	}
}

func getpos(i int, indices []int) int {
	for pos, idx := range indices {
		if idx == i {
			return pos
		}
	}

	panic("not found")
}

func getval(input, indices []int, pos int) int {
	return input[indices[mod(pos, len(indices))]]
}

func move(from, to int, nums []int) {
	dir := 1
	if to <= from {
		dir = -1
	}

	n := len(nums)

	for i := from; i != to; i += dir {
		nums[mod(i, n)], nums[mod(i+dir, n)] = nums[mod(i+dir, n)], nums[mod(i, n)]
	}
}

// mod returns the positive value of x mod n
func mod(x, n int) int {
	return ((x % n) + n) % n
}

func readInput(path string) ([]int, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		input := make([]int, 0, len(lines))

		for _, line := range lines {
			input = append(input, MustAtoi(line))
		}

		return input, nil
	}
}

func MustAtoi(str string) int {
	if x, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return x
	}
}
