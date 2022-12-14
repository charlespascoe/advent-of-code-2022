package main

type Map [][]Cell

func (m Map) MovesFrom(pos Pos) []Pos {
	var res []Pos

	// Iterate over the 4 cardinal directions (U/D/L/R)
	for _, delta := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		p := pos.Move(delta[0], delta[1])

		if m.Has(p) {
			res = append(res, p)
		}
	}

	return res
}

func (m Map) Get(pos Pos) *Cell {
	if m.Has(pos) {
		return &m[pos.Row][pos.Col]
	}

	return nil
}

func (m Map) Has(pos Pos) bool {
	return (0 <= pos.Row && pos.Row < len(m)) && (0 <= pos.Col && pos.Col < len(m[0]))
}

func parseMap(lines []string) (m Map, start, end Pos) {
	m = make(Map, 0, len(lines))

	for ln, line := range lines {
		row := make([]Cell, 0, len(line))

		for cn, char := range line {
			cell := Cell{
				Pos:      Pos{ln, cn},
				BestDist: MaxInt64,
			}

			switch char {
			case 'S':
				start = cell.Pos
				cell.Height = 0
			case 'E':
				end = cell.Pos
				cell.Height = 25
			default:
				cell.Height = int(char - 'a')
			}

			row = append(row, cell)
		}

		m = append(m, row)
	}

	return
}

// --- Pos --- //

type Pos struct {
	Row, Col int
}

func (p Pos) Move(r, c int) Pos {
	return Pos{p.Row + r, p.Col + c}
}

func (p Pos) ArrowTo(to Pos) byte {
	// Assumes p and to are adjacent on one side.
	switch {
	case to.Row < p.Row:
		return '^'
	case to.Row > p.Row:
		return 'v'
	case to.Col < p.Col:
		return '<'
	case to.Col > p.Col:
		return '>'
	default:
		return '?'
	}
}

// --- Cell --- //

type Cell struct {
	Pos      Pos
	Height   int
	Visited  bool
	BestStep Pos
	BestDist int
}

func (c *Cell) Key() Pos {
	return c.Pos
}

func (c *Cell) LessThan(other *Cell) bool {
	return c.BestDist < other.BestDist
}
