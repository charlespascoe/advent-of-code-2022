package main

import (
	"fmt"
	"sort"
)

// Saving TODO: Description.
type Saving struct {
	V1, V2 Valve
	Savings int
}

// Route TODO: Description.
type Route []string

// Has TODO: Description.
func (r Route) Has(name string) bool {
	for _, v := range r {
		if v == name {
			return true
		}
	}

	return false
}

// Index TODO: Description.
func (r Route) Index(name string) int {
	for i, v := range r {
		if v == name {
			return i
		}
	}

	return -1
}

// Last TODO: Description.
func (r Route) Last() string {
	return r[len(r)-1]
}

// Reversed TODO: Description.
func (r Route) Reversed() Route {
	rr := make(Route, 0, len(r))

	for i := len(r)-1; i >= 0; i-- {
		rr = append(rr, r[i])
	}

	return rr
}

// // Insert TODO: Description.
// func (r Route) Insert(pos int, name string) Route {
// 	nr := make(Route, 0, len(r)+1)
// 	nr = append(nr, r[:pos]...)
// 	nr = append(nr, name)
// 	nr = append(nr, r[pos:]...)
// 	return nr
// }

// Add TODO: Description.
func (r Route) Add(name, adj string) (Route, bool) {
	if r[0] == adj {
		return append(Route{name}, r...), true
	} else if r[len(r)-1] == adj {
		return append(r, name), true
	}

	return r, false
}

// // IsInterior TODO: Description.
func (r Route) IsInterior(name string) bool {
	if len(r) < 3 {
		return false
	}

	for _, v := range r[1:len(r)-1] {
		if v == name {
			return true
		}
	}

	return false
}


type Routes []Route

func (r Routes) RouteWith(name string) int {
	for i, route := range r {
		if route.Has(name) {
			return i
		}
	}

	return -1
}

// Remove TODO: Description.
func (r Routes) Remove(indexes... int) Routes {
	routes := make(Routes, 0, len(r))

	for i, route := range r {
		exclude := false

		for _, ex := range indexes {
			if ex == i {
				exclude = true
				break
			}
		}

		if !exclude {
			routes = append(routes, route)
		}
	}

	return routes
}

// Solver TODO: Description.
type Solver struct {
	Map *Map
	Start string
}

// Solve TODO: Description.
func (s *Solver) Solve() {
	var savings []Saving

	for i, v1 := range s.Map.Valves {
		if v1.Name == s.Start || v1.Rate == 0 {
			continue
		}
		for _, v2 := range s.Map.Valves[:i] {
			if v2.Name == s.Start || v1.Rate == 0 {
				continue
			}

			savings = append(savings, Saving{
				V1: v1,
				V2: v2,
				Savings: s.Map.GetDist(s.Start, v1.Name) + s.Map.GetDist(s.Start, v2.Name) - s.Map.GetDist(v1.Name, v2.Name),
			})
		}
	}

	// Sort savings into decending order
	sort.Slice(savings, func(i, j int) bool {
		return savings[i].Savings > savings[j].Savings
	})

	var rs Routes

	for _, sav := range savings {
		v1 := sav.V1.Name
		v2 := sav.V2.Name
		r1 := rs.RouteWith(v1)
		r2 := rs.RouteWith(v2)

		if r1 == -1 && r2 == -1 {
			// No route with either; create new route
			rs = append(rs, Route{v1, v2})
			continue
		}

		if r1 >= 0 && r2 == -1 {
			var added bool
			rs[r1], added = rs[r1].Add(v2, v1)

			if added {
				// We added this link to the start/end of the route
				continue
			}

			// If added == false, it means v1 was an interior node - we can only
			// add nodes if they are at the start or end
		}

		if r1 == -1 && r2 >= 0 {
			var added bool
			rs[r2], added = rs[r2].Add(v1, v2)

			if added {
				continue
			}
		}

		if r1 == r2 {
			// Already in a route
			continue
		}

		if r1 >= 0 && r2 >= 0 && !rs[r1].IsInterior(v1) && !rs[r2].IsInterior(v2) {
			if rs[r1].Last() != v1 {
				rs[r1] = rs[r1].Reversed()
			}

			if rs[r2][0] != v2 {
				rs[r2] = rs[r2].Reversed()
			}

			if len(rs) <= 2 {
				// We're looking for two routes
				continue
			}

			// We can merge these
			fmt.Printf("R1: %v, V1: %s, Int: %t\n", rs[r1], v1, rs[r1].IsInterior(v1))
			fmt.Printf("R2: %v, V2: %s, Int: %t\n", rs[r2], v2, rs[r2].IsInterior(v2))
			merged := append(rs[r1], rs[r2]...)
			rs = append(rs.Remove(r1, r2), merged)
		}
	}

	fmt.Printf("Routes: %v\n", rs)

	for _, route := range rs {
		best := ""
		bestDist := 1000000
		for _, step := range route {
			if dist := s.Map.GetDist(s.Start, step); dist < bestDist {
				best = step
				bestDist = dist
			}
		}

		fmt.Printf("Route: %v, Best start: %s\n", route, best)
	}
}
