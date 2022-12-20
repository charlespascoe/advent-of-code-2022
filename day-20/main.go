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
	debug := flag.Bool("d", false, "debug output")
	flag.Parse()

	input, err := readInput(*inputFile)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	fmt.Printf("Input: %v\n", input)

	indices := make([]int, 0, len(input))
	for i := range input {
		indices = append(indices, i)
	}

	printnums(input, indices)

	done := 0
	i := 0
	for done < len(indices) {
		for indices[i] != done {
			i = mod(i+1, len(input))
		}

		x := input[indices[i]]
		from := i
		to := x + i
		move(from, to, indices)
		done++

		if *debug {
			fmt.Println()
			fmt.Printf("Move val %d at %d to %d (%d), Done %d\n", x, from, to, mod(to, len(input)), done)
			printnums(input, indices, from, mod(to, len(input)))
		}
	}

	zeroPos := 0

	for _, idx := range indices {
		if input[idx] == 0 {
			zeroPos = idx
			break
		}
	}

	fmt.Printf("Zero pos: %d\n", zeroPos)
	result := getval(input, indices, zeroPos+1000) + getval(input, indices, zeroPos+2000) + getval(input, indices, zeroPos+3000)
	fmt.Printf("Result: %d\n", result)
}

func getval(input, indices []int, pos int) int {
	return input[indices[mod(pos, len(indices))]]
}

func newpos(i, val, n int) (pos, ni int) {
	pos = mod(i + val, n)
	ni = i

	if val < 0 {
		pos = mod(pos - 1, n)

		if pos > i {
			ni = i - 1
		}
	}

	return
}

func move(from, to int, nums []int)  {
	dir := 1
	if to <= from {
		dir = -1
	}

	n := len(nums)

	for i := from; i != to; i += dir {
		nums[mod(i, n)], nums[mod(i+dir, n)] = nums[mod(i+dir, n)], nums[mod(i, n)]
	}
}

func printnums(nums, indices []int, special ...int) {
	s := toset(special)
	for i := 0; i < len(nums); i++ {
		v := i
		if c := s[i]; c > 0 {
			fmt.Printf("\033[%dm%2d\033[0m | ", c+30, v)
		} else {
			fmt.Printf("%2d | ", v)
		}
	}
	fmt.Println()
	for i := 0; i < len(nums); i++ {
		// fmt.Printf("%2d : ", indices[i])
		v := indices[i]
		if c := s[i]; c > 0 {
			fmt.Printf("\033[%dm%2d\033[0m | ", c+30, v)
		} else {
			fmt.Printf("%2d : ", v)
		}
	}
	fmt.Println()
	for i := 0; i < len(nums); i++ {
		// fmt.Printf("%2d : ", nums[indices[i]])
		v := nums[indices[i]]
		if c := s[i]; c > 0 {
			fmt.Printf("\033[%dm%2d\033[0m | ", c+30, v)
		} else {
			fmt.Printf("%2d : ", v)
		}
	}

	fmt.Println()
}

func toset(nums []int) map[int]int {
	m := make(map[int]int)
	for i, x := range nums {
		m[x] = i + 1
	}
	return m
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
