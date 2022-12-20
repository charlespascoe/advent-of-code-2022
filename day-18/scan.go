package main

type ScanFunc func(Vector)

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
