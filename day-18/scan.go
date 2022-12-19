package main

type ScanFunc func(Vector)

type ScanSlice [][]bool

func CreateScanSlice(width, height int) ScanSlice {
	slice := make(ScanSlice, 0, width)

	for len(slice) < cap(slice) {
		slice = append(slice, make([]bool, height))
	}

	return slice
}

func (slice ScanSlice) Count() int {
	count := 0

	for _, line := range slice {
		for _, cell := range line {
			if cell {
				count++
			}
		}
	}

	return count
}

func (slice ScanSlice) Set(x, y int, val bool) {
	slice[x][y] = slice[x][y] || val
}

type MoveFunc func(drop Droplet, pos Vector, i int) (Vector, int)

func MoveX(forward bool) MoveFunc {
	return func(droplet Droplet, pos Vector, i int) (Vector, int) {
		rem := 0
		pos.X, rem = move(i, droplet.width, forward)
		return pos, rem
	}
}

func MoveY(forward bool) MoveFunc {
	return func(droplet Droplet, pos Vector, i int) (Vector, int) {
		rem := 0
		pos.Y, rem = move(i, droplet.length, forward)
		return pos, rem
	}
}

func MoveZ(forward bool) MoveFunc {
	return func(droplet Droplet, pos Vector, i int) (Vector, int) {
		rem := 0
		pos.Z, rem = move(i, droplet.height, forward)
		return pos, rem
	}
}

func move(i, limit int, forward bool) (val, remaining int) {
	val = i % limit
	if !forward {
		val = limit - val - 1
	}
	remaining = i / limit
	return
}
