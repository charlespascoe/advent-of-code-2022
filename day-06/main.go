package main

import (
	"fmt"
	"log"
	"os"
	"flag"
)

func main() {
	part2 := flag.Bool("part2", false, "print part 2 solution")
	flag.Parse()

	input, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	seqLen := 4

	if *part2 {
		seqLen = 14
	}

	for i := 0; i < len(input)-seqLen; i++ {
		if allUnique(input[i:i+seqLen]) {
			fmt.Printf("Start: %d\n", i+seqLen)
			break
		}
	}
}

// allUnique determines if the slice of lowercase letters ('a' - 'z')
// contains only unique characters
func allUnique(chars []byte) bool {
	var mask, prev uint32

	for _, c := range chars {
		// Use xor to flip the bit:
		//   - Not seen before: 0 -> 1, `mask` increases  (mask > prev)
		//   - Seen before:     1 -> 0, `mask` descreases (mask < prev)
		mask ^= 1 << (c - 'a')

		if mask < prev {
			// mask decreased - we flipped a bit 0 -> 1 -> 0 because we've
			// encountered a character twice
			return false
		}

		prev = mask
	}

	return true
}
