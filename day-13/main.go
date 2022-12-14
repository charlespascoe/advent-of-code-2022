package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	input := flag.String("i", "input.txt", "program input")
	flag.Parse()

	packetPairs, err := readPacketPairs(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	sumOfCorrectIndices := 0
	var packets []any

	for i, pair := range packetPairs {
		fmt.Printf("Left:  %v\n", pair[0])
		fmt.Printf("Right: %v\n", pair[1])

		packets = append(packets, pair[:]...)

		ord, err := order(pair[0], pair[1])
		if err != nil {
			log.Fatalf("Invalid packet pair %d: %s", i+1, err)
		}

		if ord < 0 {
			sumOfCorrectIndices += i + 1
			fmt.Println("\033[32m>>> In order\033[0m")
		} else {
			fmt.Println("\033[31m... Out of order\033[0m")
		}

		fmt.Println()
	}

	fmt.Printf("Sum of correct indices: %d\n", sumOfCorrectIndices)

	// ---  Part 2 --- //

	// Add divider packets
	packets = append(packets, []any{[]any{float64(2)}}, []any{[]any{float64(6)}})

	sort.Slice(packets, func(i, j int) bool {
		// We can safely ignore this error - it would have
		ord, err := order(packets[i], packets[j])
		if err != nil {
			panic(err)
		}

		return ord < 0
	})

	decoderKey := 1

	fmt.Println("\nSorted packets:")
	for i, packet := range packets {
		ps := fmt.Sprintf("%v", packet)
		fmt.Println(ps)

		if ps == "[[2]]" || ps == "[[6]]" {
			decoderKey *= i+1
		}
	}

	fmt.Printf("\nDecoder key: %d\n", decoderKey)
}

func readPacketPairs(path string) ([][2]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var packets [][2]any

	for i, packetPair := range bytes.Split(bytes.TrimSpace(data), []byte("\n\n")) {
		l, r, found := bytes.Cut(packetPair, []byte("\n"))
		if !found {
			return nil, fmt.Errorf("packet pair %d not valid", i+1)
		}

		var left any
		if err := json.Unmarshal(l, &left); err != nil {
			return nil, fmt.Errorf("invalid left packet '%s': %s", l, err)
		}

		var right any
		if err := json.Unmarshal(r, &right); err != nil {
			return nil, fmt.Errorf("invalid right packet '%s': %s", r, err)
		}

		packets = append(packets, [2]any{left, right})
	}

	return packets, nil
}

func order(left, right any) (int, error) {
	switch l := left.(type) {
	case float64:
		switch r := right.(type) {
		case float64:
			return numOrder(l, r), nil

		case []any:
			return order([]any{l}, r)

		default:
			return 0, fmt.Errorf("unknown type: %T", right)
		}

	case []any:
		switch r := right.(type) {
		case float64:
			return order(l, []any{r})

		case []any:
			return listOrder(l, r)

		default:
			return 0, fmt.Errorf("unknown type: %T", right)
		}
	default:
		return 0, fmt.Errorf("unknown type: %T", left)
	}
}

func numOrder(left, right float64) int {
	if left < right {
		return -1
	} else if left > right {
		return 1
	} else {
		return 0
	}
}

func listOrder(left, right []any) (int, error) {
	l := len(left)
	if len(right) < l {
		l = len(right)
	}

	for i := 0; i < l; i++ {
		res, err := order(left[i], right[i])

		if res != 0 || err != nil {
			return res, err
		}
	}

	return numOrder(float64(len(left)), float64(len(right))), nil
}
