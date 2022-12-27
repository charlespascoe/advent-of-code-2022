package main

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
