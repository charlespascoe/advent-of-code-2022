package main

// TODO: is there a better way of doing this?
var adjacent = []Vector{
	{-1, 0, 0},
	{1, 0, 0},
	{0, -1, 0},
	{0, 1, 0},
	{0, 0, -1},
	{0, 0, 1},
}

type Vector struct {
	X, Y, Z int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

type Droplet struct {
	cubes                 [][][]bool
	width, length, height int
	// Total internal and external surface area
	surfaceArea int
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
		row := make([][]bool, 0, droplet.length)
		for len(row) < cap(row) {
			row = append(row, make([]bool, droplet.height))
		}
		droplet.cubes = append(droplet.cubes, row)
	}

	for _, pos := range input {
		adj := droplet.CountAdjacent(pos)

		// -adj to remove the previously exposed surface area, +(6-adj) to add
		// the new surface area of the cube
		droplet.surfaceArea += 6 - 2*adj

		droplet.cubes[pos.X][pos.Y][pos.Z] = true
	}

	return droplet
}

func (d Droplet) CountAdjacent(pos Vector) int {
	count := 0

	for _, adj := range adjacent {
		if d.HasCube(pos.Add(adj)) {
			count++
		}
	}

	return count
}

func (d Droplet) HasCube(pos Vector) bool {
	inBounds := (0 <= pos.X && pos.X < d.width) && (0 <= pos.Y && pos.Y < d.length) && (0 <= pos.Z && pos.Z < d.height)

	if !inBounds {
		return false
	}

	return d.cubes[pos.X][pos.Y][pos.Z]
}
