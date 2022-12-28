package main

const MaxInt64 = 0x7fffffffffffffff

type DistMatrix map[string]int

func (dm DistMatrix) DistTo(name string) int {
	if dist, ok := dm[name]; ok {
		return dist
	} else {
		return MaxInt64
	}
}

type Map struct {
	Valves      []Valve
	ValveLookup map[string]Valve
	Dists       map[string]DistMatrix
}

func NewMap(valves []Valve) *Map {
	vl := make(map[string]Valve, len(valves))

	for _, v := range valves {
		vl[v.Name] = v
	}

	return &Map{
		Valves:      valves,
		ValveLookup: vl,
		Dists:       make(map[string]DistMatrix),
	}
}

func (m *Map) GetDist(from, to string) int {
	if dist, ok := m.Dists[from][to]; ok {
		return dist
	} else {
		return MaxInt64
	}
}

func (m *Map) ComputeShortestPaths() {
	for _, valve := range m.Valves {
		distMatrix := make(DistMatrix)
		distMatrix[valve.Name] = 0

		toVisit := []Valve{valve}
		visited := make(map[string]bool, len(m.Valves))

		for len(toVisit) > 0 {
			var v Valve
			v, toVisit = popMin(toVisit, distMatrix)
			visited[v.Name] = true

			for _, next := range v.Tunnels {
				if visited[next] {
					continue
				}

				if distMatrix.DistTo(next) > distMatrix.DistTo(v.Name)+1 {
					distMatrix[next] = distMatrix.DistTo(v.Name) + 1
				}

				toVisit = append(toVisit, m.ValveLookup[next])
			}
		}

		m.Dists[valve.Name] = distMatrix
	}
}

func popMin(valves []Valve, distMatrix DistMatrix) (Valve, []Valve) {
	// It's not worth using a heap for this problem

	var min Valve
	pos := 0
	best := MaxInt64

	for i, v := range valves {
		if dist := distMatrix.DistTo(v.Name); dist < best {
			min = v
			pos = i
			best = dist
		}
	}

	// I could just mutate 'valves', but that could cause other problems
	res := make([]Valve, 0, len(valves)-1)
	res = append(res, valves[:pos]...)
	res = append(res, valves[pos+1:]...)

	return min, res
}
