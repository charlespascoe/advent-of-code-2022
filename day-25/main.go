package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

// var verbose = flag.Bool("v", false, "verbose output")
// var part2 = flag.Bool("part2", false, "print part 2 solution")

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	sum := 0
	for _, line := range lines {
		sum += snafuToDec(line)
	}

	fmt.Printf("Sum: %d, Sum in SNAFU format: %s\n", sum, decToSnafu(sum))
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

const snafuDigits = "=-012"

func snafuToDec(snafu string) int {
	pos := 1
	res := 0

	for _, digChar := range reverseString(snafu) {
		dig := strings.IndexRune(snafuDigits, digChar) - 2
		res += pos * dig
		pos *= 5
	}

	return res
}

func decToSnafu(num int) string {
	var str []byte

	for num > 0 {
		idx := (num + 2) % 5

		if idx < 2 {
			num += 5
		}

		str = append(str, snafuDigits[idx])
		num /= 5
	}

	return reverseString(string(str))
}

func reverseString(str string) string {
	res := make([]rune, utf8.RuneCountInString(str))
	i := len(res)

	for _, char := range str {
		i--
		res[i] = char
	}

	return string(res)
}
