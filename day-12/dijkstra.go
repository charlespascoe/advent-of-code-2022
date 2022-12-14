package main

func ShortestPath(m Map, start, end Pos) []Pos {
	computeShortestPaths(m, start, func(fromHeight, toHeight int) bool {
		return toHeight <= (fromHeight + 1)
	})

	return reverse(buildBackPath(m, start, end))
}

func ShortestPathFromLowestPoint(m Map, end Pos) (Pos, []Pos) {
	computeShortestPaths(m, end, func(fromHeight, toHeight int) bool {
		// We're going backwards, so we need to only allow movements that we can
		// do forwards, which is why these are swapped
		return fromHeight <= (toHeight + 1)
	})

	bestDist := MaxInt64
	var bestPos Pos

	for _, row := range m {
		for _, cell := range row {
			if cell.Height == 0 && cell.BestDist < bestDist {
				bestDist = cell.BestDist
				bestPos = cell.Pos
			}
		}
	}

	return bestPos, buildBackPath(m, end, bestPos)
}

func computeShortestPaths(m Map, start Pos, canMove func(fromHeight, toHeight int) bool) {
	cell := m.Get(start)
	cell.BestDist = 0
	cell.BestStep = start

	toVisit := &Heap[*Cell, Pos]{}
	toVisit.Add(cell)

	for toVisit.Count() > 0 {
		cell := toVisit.PopMin()
		cell.Visited = true

		for _, nextPos := range m.MovesFrom(cell.Pos) {
			next := m.Get(nextPos)

			if next.Visited || !canMove(cell.Height, next.Height) {
				continue
			}

			if next.BestDist > cell.BestDist+1 {
				next.BestDist = cell.BestDist + 1
				next.BestStep = cell.Pos
			}

			toVisit.Add(next)
		}
	}
}

func reverse(path []Pos) []Pos {
	res := make([]Pos, 0, len(path))
	for i := len(path) - 1; i >= 0; i-- {
		res = append(res, path[i])
	}
	return res
}

// buildBackPath returns the shortest path starting from the end working back to
// the start. It includes both endpoints.
func buildBackPath(m Map, start, end Pos) []Pos {
	pos := end
	pathLen := m.Get(end).BestDist

	path := make([]Pos, 0, pathLen)

	for pos != start {
		path = append(path, pos)
		pos = m.Get(pos).BestStep
	}

	path = append(path, start)

	return path
}
