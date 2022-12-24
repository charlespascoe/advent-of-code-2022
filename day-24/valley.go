package main

import (
	"fmt"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// cardinalVectors is in the same order as the list of Directions
var cardinalVectors = [...]Vector{
	{0, -1}, // Up
	{0, 1},  // Down
	{-1, 0}, // Left
	{1, 0},  // Right
}

var dirStrings = [...]string{
	"^",
	"v",
	"<",
	">",
}

func GetDirection(char rune) Direction {
	switch char {
	case '^':
		return Up
	case 'v':
		return Down
	case '<':
		return Left
	case '>':
		return Right
	default:
		panic(fmt.Sprintf("unknown direction: %c", char))
	}
}

type Vector struct {
	// X moves right, Y moves down
	X, Y int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

func (v Vector) Mult(scalar int) Vector {
	return Vector{v.X * scalar, v.Y * scalar}
}

func (v Vector) Mod(mod Vector) Vector {
	return Vector{v.X % mod.X, v.Y % mod.Y}
}

type BlizzardSet struct {
	start   [][]bool
	invMove Vector
	size    Vector
}

func NewBlizzardSet(dir Direction, rows, cols int) *BlizzardSet {
	start := make([][]bool, 0, rows)

	for len(start) < cap(start) {
		start = append(start, make([]bool, cols))
	}

	size := Vector{cols, rows}

	// Multiplying by -1 inverts the direction
	invMove := cardinalVectors[dir].Mult(-1).Add(size).Mod(size)

	return &BlizzardSet{
		start:   start,
		invMove: invMove,
		size:    size,
	}
}

func (bls *BlizzardSet) Add(row, col int) {
	bls.start[row][col] = true
}

func (bls BlizzardSet) IsEmpty(pos Vector, time int) bool {
	startPos := pos.Add(bls.invMove.Mult(time)).Mod(bls.size)

	return !bls.start[startPos.Y][startPos.X]
}

type Valley struct {
	blizzards  [4]*BlizzardSet
	size Vector
}

func NewValley(lines []string) *Valley {
	rowCount := len(lines) - 2
	colCount := len(lines[0]) - 2

	fmt.Printf("Rows: %d, Cols: %d\n", rowCount, colCount)

	val := &Valley{
		size: Vector{X: colCount, Y: rowCount},
	}

	for i := 0; i < 4; i++ {
		val.blizzards[i] = NewBlizzardSet(Direction(i), rowCount, colCount)
	}

	for r, row := range lines[1 : 1+rowCount] {
		for c, char := range row[1 : 1+colCount] {
			if char != '.' {
				dir := GetDirection(char)
				val.blizzards[dir].Add(r, c)
			}
		}
	}

	return val
}

func (val *Valley) IsEmpty(pos Vector, time int) bool {
	for _, bl := range val.blizzards {
		if !bl.IsEmpty(pos, time) {
			return false
		}
	}

	return true
}

func (val *Valley) StringAtTime(time int) string {
	var str strings.Builder

	for y := 0; y < val.size.Y; y++ {
		if str.Len() > 0 {
			str.WriteRune('\n')
		}

		for x := 0; x < val.size.X; x++ {
			count := 0
			dir := Up
			for i, bl := range val.blizzards {
				if !bl.IsEmpty(Vector{x, y}, time) {
					count++
					dir = Direction(i)
				}
			}

			if count == 0 {
				str.WriteByte('.')
			} else if count == 1 {
				str.WriteString(dirStrings[dir])
			} else {
				str.WriteByte('0'+byte(count))
			}
		}
	}

	return str.String()
}

func (val *Valley) String() string {
	return val.StringAtTime(0)
}
