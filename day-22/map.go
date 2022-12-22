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

type Map struct {
	tiles [][]*Tile
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

func (m *Map) LinkFlat() {
	prev := make([]*Tile, len(m.tiles[0]))
	top := make([]*Tile, len(m.tiles[0]))

	for _, row := range m.tiles {
		var first, last *Tile

		for i, tile := range row {
			if tile != nil {
				if first == nil {
					first = tile
				}

				if prev[i] != nil {
					prev[i].Neighbours[Down] = tile
					tile.Neighbours[Up] = prev[i]
				} else {
					top[i] = tile
				}

				if last != nil {
					last.Neighbours[Right] = tile
					tile.Neighbours[Left] = last
				}
			} else {
				if last != nil {
					last.Neighbours[Right] = first
					first.Neighbours[Left] = last
				}

				if prev[i] != nil {
					prev[i].Neighbours[Down] = top[i]
					top[i].Neighbours[Up] = prev[i]
				}
			}

			last = tile
		}

		if last != nil {
			last.Neighbours[Right] = first
			first.Neighbours[Left] = last
		}

		prev = row
	}

	for i, tile := range prev {
		if tile != nil {
			tile.Neighbours[Down] = top[i]
			top[i].Neighbours[Up] = tile
		}
	}

}

func (m *Map) checkAllPointers() {
	invalid := 0

	for _, row := range m.tiles {
		for _, tile := range row {
			if tile == nil {
				continue
			}

			x := 0

			for _, n := range tile.Neighbours {
				if n == nil {
					invalid++
					x++
				}
			}

			if x > 0 {
				tile.Mark = '0' + byte(x)
			}
		}
	}

	if invalid > 0 {
		fmt.Printf("Invalid: %d\nMap:\n%s\n", invalid, m)
		panic("invalid map")
	}
}

func BuildMap(lines []string) (*Map, *Tile) {
	maxLen := 0

	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	tiles := make([][]*Tile, 0, len(lines))
	var start *Tile

	// TODO: Refactor this; it's too nested and messy
	for r, line := range lines {
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
		}

		tiles = append(tiles, row)
	}

	return &Map{tiles}, start
}
