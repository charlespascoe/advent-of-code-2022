package main

type Vector struct {
	// Note: +ve X is right, +ve Y is down
	X, Y int
}

func (v Vector) Sub(other Vector) Vector {
	return Vector{v.X - other.X, v.Y - other.Y}
}

// Magnitude returns the Manhattan magnitude of the vector
func (v Vector) Magnitude() int {
	return abs(v.X) + abs(v.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
