package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "verbose output")
	part2   = flag.Bool("part2", false, "print part 2 solution")
)

func main() {
	inputFile := flag.String("i", "input.txt", "program input")
	flag.Parse()

	input, err := readInput(*inputFile)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	var zero *ListNode
	var prev *ListNode
	var list []*ListNode

	for _, val := range input {
		if *part2 {
			val *= 811589153
		}

		node := &ListNode{Val: val}

		if prev != nil {
			node.Prev = prev
			prev.Next = node
		}

		if val == 0 {
			zero = node
		}

		prev = node
		list = append(list, node)
	}

	// Join ends
	prev.Next = list[0]
	list[0].Prev = prev

	if *verbose {
		fmt.Printf("List: %s\n", zero)
	}

	mixCount := 1
	if *part2 {
		mixCount = 10
	}

	for i := 0; i < mixCount; i++ {
		mix(list)
	}

	if *verbose {
		fmt.Printf("Mixed List: %s\n", zero)
	}

	result := 0

	n := zero
	for i := 0; i < 3; i++ {
		n = n.Walk(1000 % len(input))
		result += n.Val
	}

	fmt.Printf("Result: %d\n", result)
}

func mix(list []*ListNode) {
	for _, node := range list {
		node.Move(node.Val, len(list))
	}
}

type ListNode struct {
	Val  int
	Prev *ListNode
	Next *ListNode
}

func (node *ListNode) Extract() (prev, next *ListNode) {
	prev = node.Prev
	next = node.Next
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	node.Next = nil
	node.Prev = nil
	return
}

func (node *ListNode) InsertAfter(other *ListNode) {
	node.Prev = other
	node.Next = other.Next

	other.Next = node
	node.Next.Prev = node
}

func (node *ListNode) InsertBefore(other *ListNode) {
	node.Next = other
	node.Prev = other.Prev

	other.Prev = node
	node.Prev.Next = node
}

func (node *ListNode) Move(steps, listLen int) {
	// Reduce steps to only walk as far as we need to; when the node is removed,
	// there will only be (listLen - 1) elements left
	steps %= listLen - 1

	n := node.Walk(steps)

	if n == node {
		return
	}

	node.Extract()

	if steps > 0 {
		node.InsertAfter(n)
	} else {
		node.InsertBefore(n)
	}
}

func (node *ListNode) Walk(steps int) *ListNode {
	n := node

	for i := 0; i > steps; i-- {
		n = n.Prev
	}

	for i := 0; i < steps; i++ {
		n = n.Next
	}

	return n
}

func (node *ListNode) String() string {
	var str strings.Builder

	str.WriteString(strconv.Itoa(node.Val))

	n := node.Next

	for n != node {
		str.WriteString(", ")
		str.WriteString(strconv.Itoa(n.Val))
		n = n.Next
	}

	return str.String()
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
