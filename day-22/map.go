package main

import (
	"fmt"
	"strings"
)

type Vector struct {
	Row, Col int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.Row + other.Row, v.Col + other.Col}
}

func (v Vector) InBounds(bounds Vector) bool {
	return 0 <= v.Row && v.Row < bounds.Row && 0 <= v.Col && v.Col < bounds.Col
}

func (v Vector) Mod(mod Vector) Vector {
	return Vector{v.Row % mod.Row, v.Col % mod.Col}
}

type Direction int

const (
	Up    Direction = 0
	Right Direction = 1
	Down  Direction = 2
	Left  Direction = 3 // Equivalent to -1 mod 4
)

// Only use Left or Right
func (dir Direction) Turn(turn Direction) Direction {
	return Direction((int(dir) + int(turn)) % 4)
}

func (dir Direction) AsVector() Vector {
	return [4]Vector{
		{-1, 0}, // Up
		{0, 1},  // Right
		{1, 0},  // Down
		{0, -1}, // Left
	}[dir]
}

func (dir Direction) String() string {
	return string("^>v<"[dir])
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

func (m *Map) MapNavigator() *MapNavigator {
	return &MapNavigator{
		tileMap: m,
		dir:     Right,
		pos:     m.start,
	}
}

func (m *Map) get(r, c int) *Tile {
	inBounds := (0 <= r && r < len(m.tiles)) && (0 <= c && c < len(m.tiles[0]))

	if !inBounds {
		return nil
	}

	return m.tiles[r][c]
}

func BuildMap(lines []string) *Map {
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
	nav.tileMap.get(nav.pos.Row, nav.pos.Col).Mark = nav.dir.String()[0]

	var tile *Tile
	next := nav.pos
	dir := nav.dir.AsVector()

	for tile == nil {
		next = next.Add(dir)

		if !next.InBounds(nav.tileMap.bounds) {
			next = next.Add(nav.tileMap.bounds).Mod(nav.tileMap.bounds)
		}

		tile = nav.tileMap.get(next.Row, next.Col)
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

type Navigator interface {
	Turn(turn Direction)
	Move() bool
	Pos() Vector
	Dir() Direction
}
