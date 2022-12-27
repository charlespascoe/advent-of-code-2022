package main

import (
	"fmt"
	"math"
	"strings"
)

type Navigator interface {
	Turn(turn Direction)
	Move() bool
	Pos() Vector
	Dir() Direction
}

type Tile struct {
	Pos  Vector
	Wall bool
	Mark byte
}

type Map struct {
	tiles     [][]*Tile
	tileCount int
	bounds    Vector
	start     Vector
}

func NewMap(lines []string) *Map {
	maxLen := 0

	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	tileCount := 0
	tiles := make([][]*Tile, 0, len(lines))
	var start *Tile

	// TODO: Refactor this; it's too nested and messy
	for r, line := range lines {
		row := make([]*Tile, maxLen)

		for i, char := range line {
			if char == '.' {
				row[i] = &Tile{Pos: Vector{Row: r, Col: i}}

				if start == nil {
					start = row[i]
				}

				tileCount++
			} else if char == '#' {
				row[i] = &Tile{Pos: Vector{Row: r, Col: i}, Wall: true}
				tileCount++
			}
		}

		tiles = append(tiles, row)
	}

	return &Map{
		tiles:     tiles,
		tileCount: tileCount,
		bounds: Vector{
			Row: len(tiles),
			Col: len(tiles[0]),
		},
		start: start.Pos,
	}
}

func (m *Map) MapNavigator() *MapNavigator {
	return &MapNavigator{
		tileMap: m,
		dir:     Right,
		pos:     m.start,
	}
}

func (m *Map) CubeNavigator() *CubeNavigator {
	edge := int(math.Sqrt(float64(m.tileCount / 6)))

	faceWidth := len(m.tiles[0]) / edge
	faceHeight := len(m.tiles) / edge
	cubeNet := make(CubeNet, 0, faceHeight)

	for r := 0; r < faceHeight; r++ {
		row := make([]*Face, 0, faceWidth)

		for c := 0; c < faceWidth; c++ {
			row = append(row, m.buildFace(r, c, edge))
		}

		cubeNet = append(cubeNet, row)
	}

	cube := &Cube{Length: edge}

	cube.Fold(cubeNet, cubeNet.Start(), Identity)

	return &CubeNavigator{
		cube:     cube,
		dir:      Right,
		rotation: Identity,
		pos:      Vector{0, 0}, // Top-left corner of the start face
	}
}

func (m *Map) String() string {
	var str strings.Builder

	for i, row := range m.tiles {
		if str.Len() > 0 {
			str.WriteByte('\n')
		}

		str.WriteString(fmt.Sprintf("%2d: ", i))

		for _, tile := range row {
			if tile == nil {
				str.WriteByte(' ')
			} else if tile.Wall {
				str.WriteByte('#')
			} else if tile.Mark != 0 {
				str.WriteByte(tile.Mark)
			} else {
				str.WriteByte('.')
			}
		}
	}

	return str.String()
}

func (m *Map) buildFace(faceRow, faceCol, length int) *Face {
	if m.get(Vector{faceRow * length, faceCol * length}) == nil {
		return nil
	}

	face := &Face{}

	for r := 0; r < length; r++ {
		face.Tiles = append(face.Tiles, m.tiles[faceRow*length+r][faceCol*length:(faceCol+1)*length])
	}

	return face
}

func (m *Map) get(pos Vector) *Tile {
	if !pos.InBounds(m.bounds) {
		return nil
	}

	return m.tiles[pos.Row][pos.Col]
}

type MapNavigator struct {
	tileMap *Map
	dir     Direction
	pos     Vector
}

func (nav *MapNavigator) Turn(turn Direction) {
	nav.dir = nav.dir.Turn(turn)
}

func (nav *MapNavigator) Move() bool {
	// Mark the tile that we're about to leave (see Map.String())
	nav.tileMap.get(nav.pos).Mark = nav.dir.String()[0]

	var tile *Tile
	next := nav.pos
	dir := nav.dir.AsVector()

	for tile == nil {
		next = next.Add(dir)

		if !next.InBounds(nav.tileMap.bounds) {
			next = next.Add(nav.tileMap.bounds).Mod(nav.tileMap.bounds)
		}

		tile = nav.tileMap.get(next)
	}

	if tile.Wall {
		return false
	}

	nav.pos = next

	return true
}

func (nav *MapNavigator) Pos() Vector {
	return nav.pos
}

func (nav *MapNavigator) Dir() Direction {
	return nav.dir
}
