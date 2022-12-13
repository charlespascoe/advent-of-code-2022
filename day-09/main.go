package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"flag"
)

func main() {
	part2 := flag.Bool("part2", false, "print part 2 solution")
	flag.Parse()

	lines, err := readLines("input.txt")
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	ropeLen := 2
	if *part2 {
		ropeLen = 10
	}

	head := NewRope(ropeLen)

	visited := make(map[Vector]bool)

	for _, line := range lines {
		dir, count, err := readInstruction(line)
		if err != nil {
			panic(err)
		}

		for i := 0; i < count; i++ {
			head.Move(dir)
			visited[head.Tail().Pos] = true
		}
	}

	fmt.Printf("Visited: %d\n", len(visited))
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

type Vector struct {
	// Note: +ve X is right, +ve Y is up
	X, Y int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

func (v Vector) Sub(other Vector) Vector {
	return Vector{v.X - other.X, v.Y - other.Y}
}

// AdjacentTo returns true if the given vector resolves to an adjacent (or the
// same) point from the origin. Diagonal is also considered adjacent.
func (v Vector) AdjacentTo(other Vector) bool {
	delta := v.Sub(other)
	return (-1 <= delta.X && delta.X <= 1) && (-1 <= delta.Y && delta.Y <= 1)
}

var (
	Up    = Vector{0, 1}
	Down  = Vector{0, -1}
	Right = Vector{1, 0}
	Left  = Vector{-1, 0}
)

func readInstruction(line string) (dir Vector, count int, err error) {
	l, r, found := strings.Cut(line, " ")
	if !found {
		err = errors.New("invalid instruction")
		return
	}

	count, err = strconv.Atoi(r)
	if err != nil {
		return
	}

	switch l {
	case "U":
		dir = Up
	case "D":
		dir = Down
	case "R":
		dir = Right
	case "L":
		dir = Left
	default:
		err = fmt.Errorf("unknown direction '%s'", l)
	}

	return
}

type RopeNode struct {
	Num  int
	Pos  Vector
	Next *RopeNode
}

func (n *RopeNode) Move(dir Vector) {
	if dir == (Vector{}) {
		// No movement, just stop
		return
	}

	n.Pos = n.Pos.Add(dir)
	n.Next.Follow(n.Pos)
}

func (n *RopeNode) Follow(parent Vector) {
	if n == nil {
		return
	}

	if n.Pos.AdjacentTo(parent) {
		// We're already adjacent (or at the same point), so don't need to move
		return
	}

	move := Vector{}

	if parent.X < n.Pos.X {
		move = move.Add(Left)
	} else if parent.X > n.Pos.X {
		move = move.Add(Right)
	}

	if parent.Y < n.Pos.Y {
		move = move.Add(Down)
	} else if parent.Y > n.Pos.Y {
		move = move.Add(Up)
	}

	n.Pos = n.Pos.Add(move)

	n.Next.Follow(n.Pos)
}

func (n *RopeNode) Tail() *RopeNode {
	if n.Next == nil {
		return n
	} else {
		return n.Next.Tail()
	}
}

func NewRope(length int) *RopeNode {
	head := &RopeNode{Num: 0}
	n := head

	for i := 1; i < length; i++ {
		n.Next = &RopeNode{Num: n.Num + 1}
		n = n.Next
	}

	return head
}
