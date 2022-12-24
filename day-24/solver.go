package main

import (
	"fmt"
)

type Path struct {
	Time int
	Pos  Vector
	Prev *Path
	Wait bool
}

func (path *Path) LessThan(other *Path) bool {
	// Prioritise movement over waiting
	if path.Wait != other.Wait {
		return other.Wait
	}

	return path.Time < other.Time
}

type Solver struct {
	val        *Valley
	nextMoves  *Heap[*Path]
	start, end Vector
	step int
	moves []Vector
}

func NewSolver(val *Valley) *Solver {
	solver := &Solver{
		val:       val,
		nextMoves: NewHeap[*Path, int](),
		start:     Vector{0, -1}, // Just above top-left corner
		// val.size will already have the correct Y value but be one to the
		// right
		end: val.size.Add(Vector{-1, 0}),
		moves: []Vector{
			cardinalVectors[Down],
			cardinalVectors[Right],
			cardinalVectors[Left],
			cardinalVectors[Up],
		},
	}

	solver.nextMoves.Add(&Path{
		Time: 0,
		Pos:  solver.start,
	})

	return solver
}

func (solver *Solver) Solve() {
	var path *Path

	for path == nil {
		path = solver.nextStep()
	}

	// TODO: 

	fmt.Printf("Done! %d\n", path.Time)
}

func (solver *Solver) nextStep() *Path {
	solver.step++

	// fmt.Printf("Moves: %d\n", solver.nextMoves.Count())
	path := solver.nextMoves.PopMin()

	if solver.step % 10000 == 0 {
		fmt.Printf("%d (Time: %d, Moves: %d)\n", solver.step, path.Time, solver.nextMoves.Count())
	}

	if path == nil {
		panic("out of moves")
	}

	for _, dir := range solver.moves {
		next := path.Pos.Add(dir)

		if next == solver.end {
			return &Path{
				Time: path.Time + 1,
				Pos:  next,
				Prev: path,
			}
		}

		if !next.InBounds(solver.val.size) {
			// We're at the edge
			continue
		}

		if !solver.val.IsEmpty(next, path.Time+1) {
			// There's going to be a blizzard here
			continue
		}

		solver.nextMoves.Add(&Path{
			Time: path.Time + 1,
			Pos:  next,
			Prev: path,
		})
	}

	if path.Pos == solver.start || solver.val.IsEmpty(path.Pos, path.Time+1) {
		// fmt.Printf("Wait in place\n")
		solver.nextMoves.Add(&Path{
			Time: path.Time + 1,
			Pos:  path.Pos,
			Prev: path,
			Wait: true,
		})
	}

	return nil
}
