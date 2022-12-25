package main

type Path struct {
	Time int
	Pos  Vector
	Prev *Path
}

func (path *Path) LessThan(other *Path) bool {
	return path.Time < other.Time
}

func (path *Path) Moves() []Vector {
	if path == nil {
		return nil
	}

	return append(path.Prev.Moves(), path.Pos)
}

type PhaseMarkers struct {
	// markers[row][col][phase] == best path node at this row+col at this phase
	// during the blizzard loop (time mod period)
	markers [][][]*Path

	// The blizzard repeat time period
	period int
}

func NewPhaseMarkers(val *Valley) PhaseMarkers {
	// The +2 is to handle the start positions, which are outside the main
	// valley coordinates
	best := make([][][]*Path, 0, val.Rows()+2)

	for len(best) < cap(best) {
		row := make([][]*Path, 0, val.Cols())

		for len(row) < cap(row) {
			row = append(row, make([]*Path, val.Period))
		}

		best = append(best, row)
	}

	return PhaseMarkers{
		markers: best,
		period:  val.Period,
	}
}

func (markers *PhaseMarkers) Get(path *Path) *Path {
	// +1 offsets the fact that one of the start positions is before the first
	// (zero) row.
	return markers.markers[path.Pos.Y+1][path.Pos.X][path.Time%markers.period]
}

func (markers *PhaseMarkers) Set(path *Path) {
	// +1 offsets the fact that one of the start positions is before the first
	// (zero) row.
	markers.markers[path.Pos.Y+1][path.Pos.X][path.Time%markers.period] = path
}

type Solver struct {
	// valley is the valley with blizzard descriptions
	valley *Valley

	// nextMoves is a min heap of next moves
	nextMoves *Heap[*Path]

	// phaseMarkers is a collection of markers on each tile for each blizzard
	// phase (from 0 through period - 1). Used to detect cycles so we can
	// abandon useless walks.
	phaseMarkers PhaseMarkers

	// Start and end positions
	start, end Vector
}

func NewSolver(valley *Valley, start, end Vector, startTime int) *Solver {
	solver := &Solver{
		valley:       valley,
		nextMoves:    NewHeap[*Path, int](),
		phaseMarkers: NewPhaseMarkers(valley),
		start:        start,
		end:          end,
	}

	solver.nextMoves.Add(&Path{
		Time: startTime,
		Pos:  solver.start,
	})

	return solver
}

func (solver *Solver) Solve() *Path {
	var path *Path

	for path == nil {
		path = solver.nextStep()
	}

	return path
}

func (solver *Solver) nextStep() *Path {
	path := solver.nextMoves.PopMin()

	if path == nil {
		// We should always have a move, otherwise it's unsolvable
		panic("out of moves")
	}

	if solver.betterFound(path) {
		// We walked in a circle
		return nil
	}

	for _, dir := range cardinalVectors {
		next := path.Pos.Add(dir)

		if next == solver.end {
			// We found the end; the use of a min-heap ensures this is
			// definitely the lowest-time path we could have found
			return &Path{
				Time: path.Time + 1,
				Pos:  next,
				Prev: path,
			}
		}

		if !next.InBounds(solver.valley.Size) {
			// We're at the edge
			continue
		}

		if !solver.valley.IsEmpty(next, path.Time+1) {
			// There's going to be a blizzard here
			continue
		}

		solver.nextMoves.Add(&Path{
			Time: path.Time + 1,
			Pos:  next,
			Prev: path,
		})
	}

	if path.Pos == solver.start || solver.valley.IsEmpty(path.Pos, path.Time+1) {
		solver.nextMoves.Add(&Path{
			Time: path.Time + 1,
			Pos:  path.Pos,
			Prev: path,
		})
	}

	return nil
}

func (solver *Solver) betterFound(path *Path) bool {
	// See if we've visited this node before at this particular phase during the
	// blizzard cycle
	previous := solver.phaseMarkers.Get(path)

	if previous != nil {
		if previous.Time > path.Time {
			// Because we use a min-heap for all moves, the first time we visit
			// this tile should always be the soonest time
			panic("unexpectedly encounting an existing solution that was worse")
		}

		return true
	}

	solver.phaseMarkers.Set(path)

	return false
}
