package main

import (
	"fmt"
	"strings"
)

type Direction int

// Only use Left or Right
func (dir Direction) Turn(turn Direction) Direction {
	return Direction((int(dir) + int(turn)) % 4)
}

func (dir Direction) String() string {
	switch dir {
	case Up:
		return "^"
	case Right:
		return ">"
	case Down:
		return "v"
	case Left:
		return "<"
	default:
		return "?"
	}
}

const (
	Up    Direction = 0
	Right Direction = 1
	Down  Direction = 2
	Left  Direction = 3 // Equivalent to -1 mod 4
)

type Tile struct {
	Row int
	Col int
	Neighbours [4]*Tile
	Wall       bool
	Mark byte
}

func (tile *Tile) Move(dir Direction) *Tile {
	return tile.Neighbours[dir]
}

type Map [][]*Tile

func (m Map) String() string {
	var str strings.Builder

	for i, row := range m {
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

func BuildMap(lines []string) (Map, *Tile) {
	maxLen := 0

	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	tiles := make([][]*Tile, 0, len(lines))
	bottom := make([]*Tile, maxLen)
	var start *Tile

	// TODO: Refactor this; it's too nested and messy
	for r, line := range lines {
		var first, last *Tile
		row := make([]*Tile, maxLen)

		for i, char := range line {
			if char == '.' {
				row[i] = &Tile{Row: r, Col: i}

				if start == nil {
					start = row[i]
				}
			} else if char == '#' {
				row[i] = &Tile{Row: r, Col: i, Wall: true}
			}

			if row[i] != nil {
				if first == nil {
					first = row[i]
				}

				bottom[i] = row[i]

				if last != nil {
					last.Neighbours[Right] = row[i]
					row[i].Neighbours[Left] = last
				}
			} else {
				if last != nil {
					last.Neighbours[Right] = first
					first.Neighbours[Left] = last
				}
			}

			last = row[i]
		}

		if last != nil {
			last.Neighbours[Right] = first
			first.Neighbours[Left] = last
		}

		tiles = append(tiles, row)
	}

	prev := make([]*Tile, maxLen)

	for _, row := range tiles {
		for i, tile := range row {
			if tile != nil {
				if prev[i] == nil {
					bottom[i].Neighbours[Down] = tile
					tile.Neighbours[Up] = bottom[i]
				} else {
					prev[i].Neighbours[Down] = tile
					tile.Neighbours[Up] = prev[i]
				}
			}
		}

		prev = row
	}

	invalid := 0

	for x := 0; x < len(tiles); x++ {
		for y := 0; y < len(tiles[x]); y++ {
			if tiles[x][y] == nil {
				continue
			}
			for _, n := range tiles[x][y].Neighbours {
				if n == nil {
					invalid++
				}
			}
		}
	}

	fmt.Printf("Invalid: %d\n", invalid)

	return tiles, start
}
