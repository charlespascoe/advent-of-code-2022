package main

// Permutation describes how a cube has been rotated. Each index/value has the
// following meanings:
//
//   - 0: the face that you're looking at
//   - 1: top face
//   - 2: right face
//   - 3: bottom face
//   - 4: left face
//   - 5: rear face (opposite the face you're looking at)
//
// Think of the index of a permutation as a number written on a face of the
// cube; each face is numbered according to its start position as listed above.
// The value at that index in a permutation is the position that the numbered
// face has been moved to, e.g. p[3] == 1 means face number 3, which starts at
// the bottom, is now on top. Permutation{0,1,2,3,4,5} is the identity
// permutation (i.e. the starting position of the cube before any rotations).
type Permutation [6]int

func (perm Permutation) Inverse() Permutation {
	var inv Permutation

	for i, face := range perm {
		inv[face] = i
	}

	return inv
}

func (perm Permutation) Apply(right Permutation) Permutation {
	var res Permutation

	for i, face := range perm {
		res[i] = right[face]
	}

	return res
}

var (
	Identity        = Permutation{0, 1, 2, 3, 4, 5}
	RotateUpward    = Permutation{1, 5, 2, 0, 4, 3}
	RotateDownward  = RotateUpward.Inverse()
	RotateRight     = Permutation{2, 1, 5, 3, 0, 4}
	RotateLeft      = RotateRight.Inverse()
	RotateClockwise = Permutation{0, 2, 3, 4, 1, 5}
)

var rotate = [4]Permutation{
	RotateDownward, // To go up one face, rotate the cube down
	RotateLeft,     // To go right one face, rotate the cube to the left
	RotateUpward,   // To go down one face, rotate the cube up
	RotateRight,    // To go left one face, rotate the cube to the right
}

type Face struct {
	Tiles    [][]*Tile
	Visited  bool
	Rotation Permutation
}

func (face Face) Tile(pos Vector) *Tile {
	return face.Tiles[pos.Row][pos.Col]
}

type CubeNet [][]*Face

func (net CubeNet) Bounds() Vector {
	return Vector{Row: len(net), Col: len(net[0])}
}

func (net CubeNet) Face(pos Vector) *Face {
	return net[pos.Row][pos.Col]
}

func (net CubeNet) Start() Vector {
	for c, face := range net[0] {
		if face != nil {
			return Vector{Row: 0, Col: c}
		}
	}

	panic("top row should not be empty")
}

type Cube struct {
	Faces  [6]*Face
	Length int
}

func (cube *Cube) Bounds() Vector {
	return Vector{cube.Length, cube.Length}
}

func (cube *Cube) Fold(net CubeNet, pos Vector, rotation Permutation) {
	face := net.Face(pos)

	if face == nil || face.Visited {
		return
	}

	face.Visited = true
	face.Rotation = rotation
	// The rotation permutation will tell you what position a given numbered
	// face has moved to, therefore the inverse permutation will tell you what
	// numbered face is in that position. Position 0 is the face that we're
	// currently looking at.
	cube.Faces[rotation.Inverse()[0]] = face

	for i := 0; i < 4; i++ {
		dir := Direction(i)
		next := pos.Add(dir.AsVector())

		if !next.InBounds(net.Bounds()) {
			continue
		}

		cube.Fold(net, next, rotation.Apply(rotate[dir]))
	}
}

type CubeNavigator struct {
	cube *Cube

	// dir is the direction of movement.
	dir Direction

	// pos is the position within the cube face. pos = (0, 0) is the top-left
	// corner of the face as it originally appeared in the input file.
	pos Vector

	// rotation is the permutation that describes the combination of all
	// transformations applied to the cube so far.
	rotation Permutation
}

func (nav *CubeNavigator) Turn(turn Direction) {
	nav.dir = nav.dir.Turn(turn)
}

func (nav *CubeNavigator) Move() bool {
	rotation := nav.rotation
	dir := nav.dir
	face := nav.cube.Faces[rotation.Inverse()[0]]

	// Mark the tile that we're about to leave (see Map.String())
	face.Tile(nav.pos).Mark = dir.String()[0]

	bounds := nav.cube.Bounds()
	next := nav.pos.Add(dir.AsVector())

	if !next.InBounds(bounds) {
		// Apply mod to be within bounds
		next = next.Add(bounds).Mod(bounds)

		// Move to the neighbouring face
		rotation = rotation.Apply(rotate[nav.dir])
		face = nav.cube.Faces[rotation.Inverse()[0]]

		// Adjust the rotation, direction, and position so that this face is the
		// correct way up (i.e. Vector{0,0} is the top-left tile as per the
		// input file)
		for rotation != face.Rotation {
			rotation = rotation.Apply(RotateClockwise)
			dir = dir.Turn(Right)
			next = rotateClockwise(next, nav.cube.Length)
		}
	}

	if face.Tile(next).Wall {
		return false
	}

	nav.rotation = rotation
	nav.dir = dir
	nav.pos = next

	return true
}

func (nav *CubeNavigator) Pos() Vector {
	face := nav.cube.Faces[nav.rotation.Inverse()[0]]

	return face.Tile(nav.pos).Pos
}

func (nav *CubeNavigator) Dir() Direction {
	return nav.dir
}

func rotateClockwise(pos Vector, edge int) Vector {
	return Vector{
		Row: pos.Col,
		Col: -(pos.Row - (edge - 1)),
	}
}
