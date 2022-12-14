package main

var (
	Down      = Vector{0, 1}
	DownLeft  = Vector{-1, 1}
	DownRight = Vector{1, 1}
)

type Vector struct {
	// Note: +ve X is right, +ve Y is down
	X, Y int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

func (v Vector) Sub(other Vector) Vector {
	return Vector{v.X - other.X, v.Y - other.Y}
}

func (v Vector) Unit() Vector {
	if v.X != 0 && v.Y != 0 {
		panic("cannot produce diagonal unit vector")
	}

	return Vector{unit(v.X), unit(v.Y)}
}

func unit(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}
