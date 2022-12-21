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
}

func (sch *BuildSchedule) Cost() Resources {
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

func (sch *BuildSchedule) Pop() {
	if sch.Prev != nil {
		sch.Prev.Next = sch.Next
	}

	if sch.Next != nil {
		sch.Next.Prev = sch.Prev
	}
}

func (sch *BuildSchedule) Copy() *BuildSchedule {
	// cp := new(BuildSchedule)
	// *cp = *sch
	out := sch.copyNode()

	next := out
	node := out.Prev.copyNode()

	for node != nil {
		node.Next = next
		next.Prev = node
		node = node.Prev
	}

	prev := out
	node = out.Next.copyNode()

	for node != nil {
		node.Prev = prev
		prev.Next = node
		node = node.Next
	}

	return out
}

func (sch *BuildSchedule) copyNode() *BuildSchedule {
	if sch == nil {
		return nil
	}

	cp := new(BuildSchedule)
	*cp = *sch
	return cp
}

type Robot struct {
	blueprint *Blueprint
	Resource  Resource
	Prev      *Robot
}

func (robot *Robot) Cost() Resources {
	return robot.blueprint.Costs[robot.Resource]
}

func (robot *Robot) BuildSchedule() []Resource {
	if robot == nil {
		return nil
	}

	return append(robot.Prev.BuildSchedule(), robot.Resource)
}

// func (sch *BuildSchedule) PopPrevious() {
// 	if sch.Prev == nil {
// 		panic("nothing before")
// 	}

// 	sch.Prev = sch.Prev.Prev

// 	if sch.Prev != nil {
// 		sch.Prev.Next = sch
// 	}
// }

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

// func InitialBuildSchedule(blueprint *Blueprint) *BuildSchedule {
// 	var head *BuildSchedule
// 	var prev *BuildSchedule

// 	// The [1:] is to skip "ore", as we always start with an ore miner
// 	for _, resource := range resources[1:] {
// 		sch := &BuildSchedule{
// 			Blueprint: blueprint,
// 			Robot:     resource,
// 			Best:      MaxInt64,
// 		}

// 		if head == nil {
// 			head = sch
// 		} else {
// 			sch.Prev = prev
// 			prev.Next = sch
// 		}

// 		prev = sch
// 	}

// 	return head
// }

type Simulator struct {
	blueprint Blueprint
	schedule  *BuildSchedule
}

func NewSimulator(blueprint Blueprint) *Simulator {
	return &Simulator{
		blueprint: blueprint,
		// schedule:  InitialBuildSchedule(&blueprint),
	}
}

func (sim *Simulator) solve() {
	init := NewSimState()
	state := init

	// fmt.Printf("Initial Schedule: %s\n\n", init.nextRobot)

	// for state.nextRobot != nil {
	// for state.time < 24 {
	// 	if state.nextRobot.Next == nil {
	// 		state.nextRobot.Next = &BuildSchedule{
	// 			Blueprint: &sim.blueprint,
	// 			Robot:     GeodeResource,
	// 			Prev:      state.nextRobot,
	// 		}
	// 	}

	// 	vprintf(2, "\n============\nStart Time: %d, Next Robot: %s\n", state.time, state.nextRobot.Robot)

	// 	state = sim.optimiseStep(state)

	// 	// state = sim.simUntil(state, state.nextRobot.Next)
	// }

	var robot *Robot

	for _, resource := range resources[1:] {
		robot = sim.newRobot(resource, robot)
		fmt.Printf("Chain (before): %s\n", robot.BuildSchedule())
		state = sim.optimiseStep(state.NextRobot(robot), 24)
		fmt.Printf("Chain (after): %s\n", robot.BuildSchedule())
	}

	for state.time < 24 {
		robot = sim.newRobot(GeodeResource, robot)
		fmt.Printf("Chain (before): %s\n", robot.BuildSchedule())
		state = sim.optimiseStep(state.NextRobot(robot), 24)
		fmt.Printf("Chain (after): %s\n", robot.BuildSchedule())
	}

	fmt.Printf("State: %v\n", state)
	fmt.Printf("Final Schedule: %s\n", state.nextRobot.BuildSchedule())
}

// func (sim *Simulator) optimiseStep(start SimState, mustBeat int) SimState {
// 	// prevBest := step.Best

// 	nextRobot := start.nextRobot
// 	state := sim.simUntil(start, nextRobot)

// 	bestState := start
// 	best := state.time
// 	vprintf(3, "BEST state: %v, BEST Seq: %s\n", state, start.nextRobot.Head())

// 	if best >= mustBeat {
// 		return state
// 	}

// 	count := 0

// 	for count < 3 {
// 		diff := nextRobot.Cost().Sub(state.prevResources)
// 		critRes := diff.ZeroValues()

// 		vprintf(3, "Resources: %v, Cost: %v, Critical: %v (diff: %v)\n", state.resources, nextRobot.Cost(), critRes, diff)

// 		if len(critRes) == 0 {
// 			break
// 		}

// 		improved := bestState
// 		inserted := nextRobot.InsertBefore(critRes[0])

// 		if improved.nextRobot == nextRobot {
// 			improved.nextRobot = inserted
// 		}

// 		vprintf(3, "Improved: %v\n", improved)
// 		improved = sim.optimiseStep(improved, best)
// 		state = sim.simUntil(improved, nextRobot)
// 		vprintf(3, "Inserted %s (now %s), Best Time: %d, Run Time: %d\n", critRes[0], nextRobot.Head(), best, state.time)
// 		vprintf(3, "Out State: %v\n", state)

// 		if state.time <= best {
// 			best = state.time
// 			bestState = improved
// 		} else {
// 			nextRobot.PopPrevious()
// 		}

// 		// count++
// 	}

// 	// if count >= 3 {
// 	// 	vprintf(3, "HAD TO STOP\n")
// 	// }

// 	vprintf(3, "\n")

// 	return sim.simUntil(bestState, bestState.nextRobot)
// }

var counter = 0

func (sim *Simulator) optimiseStep(state SimState, mustBeat int) SimState {
	counter++
	x := counter
	defer func() {
		counter--
	}()
	fmt.Printf("\n\n=== OPTIMISE %d: %v\n", x, state.nextRobot.BuildSchedule())
	defer fmt.Printf("\n\n=== END OPTIMISE %d: %v\n", x, state.nextRobot.BuildSchedule())
	start := state
	nextRobot := start.nextRobot
	end := sim.simUntil2(start)
	var newRobot *Robot

	if end.time > mustBeat {
		return end
	}

	critRes := nextRobot.Cost().EqualValues(end.prevResources)
	vprintf(3, "Critical Resources for %s robot: %s, Cost: %s, Resources: %s\n", nextRobot.Resource, critRes, nextRobot.Cost(), end.prevResources)

	for _, resource := range critRes {
		robot := sim.newRobot(resource, nextRobot.Prev)
		fmt.Printf("SEQ1: %v\n", robot.BuildSchedule())

		// Note we're using the original state - we're only applying one change
		// at a time
		intermediate := sim.optimiseStep(state.NextRobot(robot), end.time).NextRobot(nextRobot)
		if intermediate.time > end.time {
			// Getting to this intermediate robot ends up taking even more time,
			// so skip it
			continue
		}
		fmt.Printf("SEQ1A: %v\n", intermediate.nextRobot.BuildSchedule())

		result := sim.simUntil2(intermediate)

		if result.time <= end.time {
			newRobot = robot
			start = intermediate
			end = result
		}
	}

	if newRobot != nil && newRobot != nextRobot {
		fmt.Printf("SEQ2: %v\n", newRobot.BuildSchedule())
		fmt.Printf("SEQ3: %v\n", nextRobot.BuildSchedule())
		nextRobot.Prev = newRobot
		fmt.Printf("SEQ4: %v\n", nextRobot.BuildSchedule())
		// fmt.Printf("SEQ4: %v\n", start.nextRobot.BuildSchedule())

		abc := sim.optimiseStep(start, end.time)
		if abc.time <= end.time {
			fmt.Println("BETTERR!")
			end = abc
		} else {
			fmt.Println("WORSE!")
		}
	}

	return end
}

func (sim *Simulator) newRobot(resource Resource, prev *Robot) *Robot {
	return &Robot{
		Resource:  resource,
		Prev:      prev,
		blueprint: &sim.blueprint,
	}
}

func (sim *Simulator) simUntil2(state SimState) SimState {
	for state.nextRobot != nil {
		state = sim.step(state)
	}

	return state
}

// func (sim *Simulator) simUntil(state SimState, nextRobot *BuildSchedule) SimState {
func (sim *Simulator) simUntil(state SimState, nextRobot *Robot) SimState {
	panic("Not implemented")
	for state.nextRobot != nextRobot && state.time < 24 {
		vprintf(3, "Looking for next robot!\n")
		state = sim.step(state)
	}

	if nextRobot != nil {
		for state.nextRobot == nextRobot && state.time < 24 {
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

	state.prevResources = state.resources

	newRobot := false

	if state.nextRobot != nil && state.haveResources(state.nextRobot.Cost()) {
		newRobot = true
	}

	state = state.mineResources()

	vprintf(4, "Time: %d\n", state.time)

	if newRobot {
		state = state.consumeResources(state.nextRobot.Cost())
		state.robots[state.nextRobot.Resource]++
		vprintf(4, "\t++ Built %s robot ++\n", state.nextRobot.Resource)
		// state.nextRobot = state.nextRobot.Next
		state.nextRobot = nil
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
	time int
	// nextRobot *BuildSchedule
	nextRobot *Robot
	resources Resources
	robots    Resources

	// TODO: Get rid of this
	prevResources Resources
}

func NewSimState() SimState {
	return SimState{
		// nextRobot: InitialBuildSchedule(blueprint),
		robots: Resources{
			OreResource: 1,
		},
	}
}

func (state SimState) String() string {
	return fmt.Sprintf("Time: %d, Resources: %v, Robots: %v", state.time, state.resources, state.robots)
}

func (state SimState) NextRobot(robot *Robot) SimState {
	state.nextRobot = robot
	return state
}

func (state SimState) haveResources(cost Resources) bool {
	for resource, count := range cost {
		if state.resources[resource] < count {
			return false
		}
	}

	return true
}

func (state SimState) consumeResources(cost Resources) SimState {
	for resource, count := range cost {
		state.resources[resource] -= count
	}

	return state
}

func (state SimState) mineResources() SimState {
	for resource, robots := range state.robots {
		// map[] is a reference type, so this works despite not using a pointer
		// receiver
		state.resources[resource] += robots
	}

	return state
}
