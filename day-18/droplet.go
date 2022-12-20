package main

import (
	"fmt"
	"strings"
)

type Vector struct {
	X, Y, Z int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// TODO: is there a better way of doing this?
var adjacent = []Vector{
	{-1, 0, 0},
	{1, 0, 0},
	{0, -1, 0},
	{0, 1, 0},
	{0, 0, -1},
	{0, 0, 1},
}

const (
	Air   = 0
	Rock  = 1
	Water = 2
)

type Droplet struct {
	cubes                 [][][]int
	width, length, height int
	// Total internal and external surface area
	surfaceArea  int
	surroundedBy int
}

func BuildDroplet(input []Vector) Droplet {
	var droplet Droplet
	var maxX, maxY, maxZ int

	for _, pos := range input {
		maxX = Max(maxX, pos.X)
		maxY = Max(maxY, pos.Y)
		maxZ = Max(maxZ, pos.Z)
	}

	droplet.width = maxX + 1
	droplet.length = maxY + 1
	droplet.height = maxZ + 1

	for len(droplet.cubes) < droplet.width {
		row := make([][]int, 0, droplet.length)
		for len(row) < cap(row) {
			row = append(row, make([]int, droplet.height))
		}
		droplet.cubes = append(droplet.cubes, row)
	}

	droplet.surroundedBy = Air

	for _, pos := range input {
		adj := droplet.CountAdjacent(pos, Rock)

		// -adj to remove the previously exposed surface area, +(6-adj) to add
		// the new surface area of the cube
		droplet.surfaceArea += 6 - 2*adj

		droplet.cubes[pos.X][pos.Y][pos.Z] = Rock
	}

	return droplet
}

func (d Droplet) CountAdjacent(pos Vector, val int) int {
	count := 0

	for _, adj := range adjacent {
		if d.CubeAt(pos.Add(adj)) == val {
			count++
		}
	}

	return count
}

func (d Droplet) CubeAt(pos Vector) int {
	inBounds := (0 <= pos.X && pos.X < d.width) && (0 <= pos.Y && pos.Y < d.length) && (0 <= pos.Z && pos.Z < d.height)

	if !inBounds {
		return d.surroundedBy
	}

	return d.cubes[pos.X][pos.Y][pos.Z]
}

func (d *Droplet) Set(pos Vector, val int) {
	d.cubes[pos.X][pos.Y][pos.Z] = val
}

func (d *Droplet) Flood() {
	d.surroundedBy = Water

	scans := [][]MoveFunc{
		{MoveX(true), MoveY(true), MoveZ(true)},  // Bottom to Top
		{MoveX(true), MoveY(true), MoveZ(false)}, // Top to Bottom
		{MoveX(true), MoveZ(true), MoveY(true)},  // Left to Right
		{MoveX(true), MoveZ(true), MoveY(false)}, // Right to Left
		{MoveY(true), MoveZ(true), MoveX(true)},  // Front to Back
		{MoveY(true), MoveZ(true), MoveX(false)}, // Back to Front
	}

	for _, scan := range scans {
		d.Scan(func(pos Vector) {
			if d.CubeAt(pos) == Air && d.CountAdjacent(pos, Water) > 0 {
				d.Set(pos, Water)
			}
		}, scan...)
	}
}

func (d *Droplet) CalcExternalSurfaceArea() int {
	d.Flood()

	surfaceArea := 0

	d.Scan(func(pos Vector) {
		if d.CubeAt(pos) == Rock {
			surfaceArea += d.CountAdjacent(pos, Water)
		}
	}, MoveX(true), MoveY(true), MoveZ(true))

	return surfaceArea
}

func (d Droplet) Scan(scan ScanFunc, moves ...MoveFunc) {
	limit := d.width * d.length * d.height
	for i := 0; i < limit; i++ {
		pos := d.calcPos(i, moves)
		scan(pos)
	}
}

func (drop Droplet) calcPos(i int, moves []MoveFunc) Vector {
	var pos Vector
	for _, move := range moves {
		pos, i = move(drop, pos, i)
	}
	return pos
}

func (d Droplet) Print() {
	var str strings.Builder

	d.Scan(func(pos Vector) {
		if pos.X == 0 && pos.Y == 0 {
			fmt.Printf("%s\n\n", str.String())
			str = strings.Builder{}
		} else if pos.X == 0 {
			str.WriteRune('\n')
		}

		switch d.CubeAt(pos) {
		case Air:
			str.WriteRune('.')
		case Rock:
			str.WriteRune('#')
		case Water:
			str.WriteRune('+')
		default:
			panic(fmt.Sprintf("unknown type: %d", d.CubeAt(pos)))
		}
	}, MoveX(true), MoveY(true), MoveZ(true))
}
