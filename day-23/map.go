package main

import (
	"math/bits"
	"strings"
)

type Vector struct {
	// X moves right, Y moves down
	X, Y int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

var (
	North = Vector{0, -1}
	South = Vector{0, 1}
	East  = Vector{1, 0}
	West  = Vector{-1, 0}
)

var adjacent = []Vector{
	North.Add(West),
	North,
	North.Add(East),
	East,
	South.Add(East),
	South,
	South.Add(West),
	West,
}

// moveMask can be rotated to test adjacency bits:
//
//	00000111 N (maskRotation = 0)
//	00011100 E (maskRotation = 2)
//	01110000 S (maskRotation = 4)
//	11000001 W (maskRotation = 6)
const moveMask = 0b00000111

// maskOffsets is a list of bit offsets to rotate moveMask to test if the elf
// can move in that direction; it is listed in the order they are checked.
var maskOffsets = []int{
	0, // North
	4, // South
	6, // West
	2, // East
}

type Elf struct {
	m        *Map
	Pos      Vector
	NextPos  Vector
	WillMove bool
}

func (elf *Elf) PlanMove(round int) {
	elf.WillMove = false
	elf.NextPos = elf.Pos

	var adjMask byte

	for i, adjVec := range adjacent {
		pos := elf.Pos.Add(adjVec)
		other := elf.m.Get(pos)

		if other != nil && other.Pos == pos {
			adjMask |= 1 << i
		}
	}

	if adjMask == 0 {
		// No neighbours; don't move
		return
	}

	for i := 0; i < 4; i++ {
		maskRotation := maskOffsets[(i+round)%4]
		sideMask := adjMask & bits.RotateLeft8(moveMask, maskRotation)

		if sideMask == 0 {
			// No elves on this side, move in this direction
			dir := adjacent[maskRotation+1]
			elf.NextPos = elf.Pos.Add(dir)
			elf.WillMove = true
			break
		}
	}

	if elf.WillMove {
		other := elf.m.Get(elf.NextPos)
		if other != nil {
			// Another elf is trying to move to this spot (we wouldn't be moving
			// here if this elf was actually here), so we all need to stop
			other.WillMove = false
			elf.WillMove = false
			elf.NextPos = elf.Pos
		} else {
			elf.m.Set(elf.NextPos, elf)
		}
	}
}

func (elf *Elf) Move() bool {
	if !elf.WillMove {
		if elf.Pos != elf.NextPos {
			elf.m.Delete(elf.NextPos)
		}

		return false
	}

	elf.m.Delete(elf.Pos)
	elf.Pos = elf.NextPos
	elf.WillMove = false
	return true
}

type Map struct {
	Elves []*Elf
	loc   map[Vector]*Elf
}

func NewMap(lines []string) *Map {
	m := &Map{
		loc: make(map[Vector]*Elf),
	}

	for row, line := range lines {
		for col, char := range line {
			if char == '#' {
				elf := &Elf{
					m:   m,
					Pos: Vector{col, row},
				}

				m.Elves = append(m.Elves, elf)
				m.loc[elf.Pos] = elf
			}
		}
	}

	return m
}

func (m *Map) Get(pos Vector) *Elf {
	return m.loc[pos]
}

func (m *Map) Set(pos Vector, elf *Elf) {
	m.loc[pos] = elf
}

func (m *Map) Delete(pos Vector) {
	delete(m.loc, pos)
}

func (m *Map) Bounds() (topLeft, bottomRight Vector) {
	minX := m.Elves[0].Pos.X
	maxX := minX
	minY := m.Elves[0].Pos.Y
	maxY := minY

	for _, elf := range m.Elves {
		minX, maxX = MinMax(elf.Pos.X, minX, maxX)
		minY, maxY = MinMax(elf.Pos.Y, minY, maxY)
	}

	return Vector{minX, minY}, Vector{maxX, maxY}
}

func (m *Map) CountEmpty() int {
	topLeft, bottomRight := m.Bounds()

	area := (bottomRight.X + 1 - topLeft.X) * (bottomRight.Y + 1 - topLeft.Y)

	return area - len(m.Elves)
}

func (m *Map) String() string {
	var str strings.Builder

	topLeft, bottomRight := m.Bounds()

	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		if str.Len() > 0 {
			str.WriteByte('\n')
		}

		for x := topLeft.X; x <= bottomRight.X; x++ {
			if m.Get(Vector{x, y}) != nil {
				str.WriteByte('#')
			} else {
				str.WriteByte('.')
			}
		}
	}

	return str.String()
}

func MinMax(x, curMin, curMax int) (min, max int) {
	min = curMin
	max = curMax

	if x < min {
		min = x
	}

	if x > max {
		max = x
	}

	return
}
