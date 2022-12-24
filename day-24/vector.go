package main

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
	return Vector{AbsMod(v.X, mod.X), AbsMod(v.Y, mod.Y)}
}

func (v Vector) InBounds(size Vector) bool {
	return 0 <= v.X && v.X < size.X && 0 <= v.Y && v.Y < size.Y
}

func AbsMod(x, n int) int {
	return ((x % n) + n) % n
}
