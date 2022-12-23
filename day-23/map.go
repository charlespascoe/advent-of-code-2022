package main

import (
	"strings"
)

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

func (m *Map) Dimensions() (width, height int) {
	topLeft, bottomRight := m.Bounds()

	return (bottomRight.X + 1 - topLeft.X), (bottomRight.Y + 1 - topLeft.Y)
}

func (m *Map) Empty() int {
	width, height := m.Dimensions()

	return (width * height) - len(m.Elves)
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
