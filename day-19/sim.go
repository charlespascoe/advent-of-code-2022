package main

import (
	"fmt"
)

const MaxInt64 = 0x7fffffffffffffff

type BuildSchedule struct {
	Blueprint *Blueprint
	Robot     Resource
	Next      *BuildSchedule
	Prev      *BuildSchedule
	Best      int
}

func (sch *BuildSchedule) Cost() Cost {
	return sch.Blueprint.Costs[sch.Robot]
}

func (sch *BuildSchedule) InsertBefore(robot Resource) *BuildSchedule {
	prev := &BuildSchedule{
		Blueprint: sch.Blueprint,
		Robot:     robot,
		Next:      sch,
		Prev:      sch.Prev,
	}

	if sch.Prev != nil {
		sch.Prev.Next = prev
	}

	sch.Prev = prev

	return prev
}

func (sch *BuildSchedule) PopPrevious() {
	if sch.Prev == nil {
		panic("nothing before")
	}

	sch.Prev = sch.Prev.Prev

	if sch.Prev != nil {
		sch.Prev.Next = sch
	}
}

func (sch *BuildSchedule) Head() *BuildSchedule {
	if sch.Prev == nil {
		return sch
	} else {
		return sch.Prev.Head()
	}
}

func (sch *BuildSchedule) String() string {
	if sch == nil {
		return ""
	}

	return fmt.Sprintf("%s, %s", sch.Robot, sch.Next)
}

func InitialBuildSchedule(blueprint *Blueprint) *BuildSchedule {
	var head *BuildSchedule
	var prev *BuildSchedule

	// The [1:] is to skip "ore", as we always start with an ore miner
	for _, resource := range resources[1:] {
		sch := &BuildSchedule{
			Blueprint: blueprint,
			Robot:     resource,
			Best:      MaxInt64,
		}

		if head == nil {
			head = sch
		} else {
			sch.Prev = prev
			prev.Next = sch
		}

		prev = sch
	}

	return head
}

type Simulator struct {
	blueprint Blueprint
	schedule  *BuildSchedule
}

func NewSimulator(blueprint Blueprint) *Simulator {
	return &Simulator{
		blueprint: blueprint,
		schedule:  InitialBuildSchedule(&blueprint),
	}
}

func (sim *Simulator) solve() {
	init := NewSimState(&sim.blueprint)
	state := init

	fmt.Printf("Initial Schedule: %s\n\n", init.nextRobot)

	// for state.nextRobot != nil {
	for state.time < 24 {
		if state.nextRobot.Next == nil {
			state.nextRobot.Next = &BuildSchedule{
				Blueprint: &sim.blueprint,
				Robot: "geode",
				Prev: state.nextRobot,
			}
		}

		vprintf(2, "\n============\nStart Time: %d, Next Robot: %s\n", state.time, state.nextRobot.Robot)

		state = sim.optimiseStep(state, 24)

		// state = sim.simUntil(state, state.nextRobot.Next)
	}

	fmt.Printf("State: %v\n", state)
	fmt.Printf("Final Schedule: %s\n", init.nextRobot)
}

func (sim *Simulator) optimiseStep(start SimState, mustBeat int) SimState {
	// prevBest := step.Best

	nextRobot := start.nextRobot
	state := sim.simUntil(start.copy(), nextRobot)

	bestState := start
	best := state.time
	vprintf(3, "BEST state: %v, BEST Seq: %s\n", state, start.nextRobot.Head())

	if best >= mustBeat {
		return state
	}

	count := 0

	for count < 3 {
		diff := nextRobot.Cost().Sub(state.prevResources)
		critRes := ZeroValues(diff)

		vprintf(3, "Resources: %v, Cost: %v, Critical: %v (diff: %v)\n", state.resources, nextRobot.Cost(), critRes, diff)

		if len(critRes) == 0 {
			break
		}

		improved := bestState.copy()
		inserted := nextRobot.InsertBefore(critRes[0])

		if improved.nextRobot == nextRobot {
			improved.nextRobot = inserted
		}

		vprintf(3, "Improved: %v\n", improved)
		improved = sim.optimiseStep(improved, best)
		state = sim.simUntil(improved.copy(), nextRobot)
		vprintf(3, "Inserted %s (now %s), Best Time: %d, Run Time: %d\n", critRes[0], nextRobot.Head(), best, state.time)
		vprintf(3, "Out State: %v\n", state)

		if state.time <= best {
			best = state.time
			bestState = improved
		} else {
			nextRobot.PopPrevious()
		}

		// count++
	}

	// if count >= 3 {
	// 	vprintf(3, "HAD TO STOP\n")
	// }

	vprintf(3, "\n")

	return sim.simUntil(bestState, bestState.nextRobot)
}

func (sim *Simulator) simUntil(state SimState, nextRobot *BuildSchedule) SimState {
	for state.nextRobot != nextRobot {
		vprintf(3, "Looking for next robot!\n")
		state = sim.step(state)
	}

	if nextRobot != nil {
		for state.nextRobot == nextRobot {
			vprintf(3, "Looking for AFTER robot!\n")
			state = sim.step(state)
		}
	}

	// Stop just after building the robot

	return state
}

// func (sim *Simulator) BestSolution() (Solution, int) {
//	best := sim.solution.Copy()
//	geodes := sim.resources["geode"]

//	vprintf(2, "\nStart\n\n")

//	for {
//		sim.RunIteration()
//		vprintf(2, "Resources: %v\n", sim.resources)

//		if g := sim.resources["geode"]; g > geodes {
//			best = sim.solution.Copy()
//			geodes = g
//		}

//		if !sim.Improve() {
//			break
//		}
//	}

//	return best, geodes
//}

// func (sim *Simulator) Improve() bool {
// 	if next, ok := sim.schedule.Next(); ok {
// 		cost := sim.blueprint.Costs[next]

// 		missing := cost.Sub(sim.resources).AsSlice()
// 		vprintf(2, "Unable to build %s robot, missing: %v\n", next, missing)

// 		if missing[0].Count < 0 {
// 			vprintf(2, "Count is negative; we have the resources, but not enough time\n")
// 			return false
// 		}

// 		sim.solution[missing[0].Resource]++
// 	} else {
// 		vprintf(2, "End of schedule, adding geode miner\n")
// 		sim.solution["geode"]++
// 	}

// 	vprintf(2, "Next Iteration Robots: %v\n\n", sim.solution)

// 	return true
// }

// func (sim *Simulator) RunIteration() {
//	sim.reset()

//	for sim.time < 24 {
//		sim.Step()
//	}
//}

func (sim *Simulator) step(state SimState) SimState {
	state.time++

	state.prevResources = copyResourceMap(state.resources)

	newRobot := false

	if state.nextRobot != nil && state.haveResources(state.nextRobot.Cost()) {
		newRobot = true
	}

	state.mineResources()

	vprintf(4, "Time: %d\n", state.time)

	if newRobot {
		state.consumeResources(state.nextRobot.Cost())
		state.robots[state.nextRobot.Robot]++
		vprintf(4, "\t++ Built %s robot ++\n", state.nextRobot.Robot)
		state.nextRobot = state.nextRobot.Next
	}

	vprintf(4, "\tResources:\n\t\t%v\n\tRobots:\n\t\t%v\n\n", state.resources, state.robots)

	return state
}

// func (sim *Simulator) buildRobot() (Resource, Cost, bool) {
//	// Refactor this something along the lines of:
//	//
//	//   - Start with a really bad solution, e.g. make one robot of each
//	//     required type, and calculate how long it would take to build the
//	//     desired robot
//	//   - Find the resource that is needed the most and increase its robot
//	//     count by one
//	//   - Recalculate solution to see if it reduces the time (always prioritise
//	//     building the simplest types first)
//	//   - Repeat until you have an optimal solution

//	if sim.nextRobot != nil {
//		next := sim.nextRobot.Robot
//		cost := sim.blueprint.Costs[next]

//		if sim.haveResources(cost) {
//			sim.consumeResources(cost)
//			sim.nextRobot = sim.nextRobot.Next
//			return next, true
//		}
//	}

//	return "", false
//}

// func (sim *Simulator) reset() {
// 	sim.time = 0
// 	sim.nextRobot = sim.schedule
// 	sim.resources = make(map[Resource]int)
// 	sim.robots = map[Resource]int{
// 		"ore": 1,
// 	}
// }

type SimState struct {
	time      int
	nextRobot *BuildSchedule
	resources map[Resource]int
	robots    map[Resource]int

	// TODO: Get rid of this
	prevResources map[Resource]int
}

func NewSimState(blueprint *Blueprint) SimState {
	return SimState{
		nextRobot: InitialBuildSchedule(blueprint),
		resources: make(map[Resource]int),
		robots: map[Resource]int{
			"ore": 1,
		},
	}
}

func (state SimState) String() string {
	return fmt.Sprintf("Time: %d, Resources: %v, Robots: %v", state.time, state.resources, state.robots)
}

func (state SimState) copy() SimState {
	state.resources = copyResourceMap(state.resources)
	state.robots = copyResourceMap(state.robots)

	return state
}

func (state SimState) copyWithNext(robot Resource) SimState {
	state = state.copy()
	state.nextRobot = state.nextRobot.InsertBefore(robot)
	return state
}

func (state SimState) haveResources(cost Cost) bool {
	for resource, count := range cost {
		if state.resources[resource] < count {
			return false
		}
	}

	return true
}

func (state SimState) consumeResources(cost Cost) {
	for resource, count := range cost {
		state.resources[resource] -= count
	}
}

func (state SimState) mineResources() {
	for resource, robots := range state.robots {
		// map[] is a reference type, so this works despite not using a pointer
		// receiver
		state.resources[resource] += robots
	}
}

func copyResourceMap(m map[Resource]int) map[Resource]int {
	res := make(map[Resource]int, len(m))

	for resource, count := range m {
		res[resource] = count
	}

	return res
}
