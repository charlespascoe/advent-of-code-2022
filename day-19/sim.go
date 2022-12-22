package main

import (
	"fmt"
)

const MaxInt64 = 0x7fffffffffffffff

type Simulator struct {
	blueprint Blueprint
	maxTime int
}

func NewSimulator(blueprint Blueprint, maxTime int) *Simulator {
	return &Simulator{
		blueprint: blueprint,
		maxTime: maxTime,
	}
}

func (sim *Simulator) solve() (*BuildStep, SimState) {
	schedule := InitialBuildStep()

	// Find the optimal path to the first of each robot (excluding an ore miner,
	// as we already start with one). optimiseStep() will add any intermediate
	// robots required to build the given robot as quickly as possible.
	for _, resource := range resources[1:] {
		schedule = sim.optimiseStep(schedule, resource)
	}

	lastGeodeMiner := schedule

	// Use remaining time to optimially build as many geode miner robots as
	// possible
	for schedule.State.time < sim.maxTime {
		lastGeodeMiner = schedule
		schedule = sim.optimiseStep(schedule, GeodeResource)
	}

	schedule = lastGeodeMiner

	for schedule.State.time > sim.maxTime {
		// The last state ih the schedule has time >= 24, which is too late to
		// have any effect on the result, so pop it
		schedule = schedule.Prev
	}

	// Pop step 0 (time == 0, the initial ore miner we already had)
	// schedule = schedule[1:]

	finalState := schedule.State

	// Simulate final state for resource counts
	finalState = sim.simulateTime(finalState, sim.maxTime-finalState.time)

	return schedule, finalState
}

func (sim *Simulator) solve2() (*BuildStep, SimState) {
	start := InitialBuildStep()
	schedule := start

	for _, resource := range resources[1:] {
		schedule = sim.buildNextRobot2(schedule, resource)


// 		fmt.Printf("\033[36mStarting with: %s\033[0m\n", next)
// 		improved := sim.optimiseTo(next, start)

// 		for improved != nil {
// 			fmt.Printf("\033[36mResult: %d, compared to %d\033[0m\n", improved.State.time, next.State.time)
// 			if improved.State.time <= next.State.time {
// 				next = improved
// 			}

// 			fmt.Printf("\033[36mStarting with: %s\033[0m\n", next)
// 			improved = sim.optimiseTo(next, start)
// 		}

// 		schedule = next
	}

	schedule.print()

	next := schedule

	for {
		fmt.Printf("\033[36mStarting with: %s\033 (time: %d)[0m\n", next, next.State.time)
		next.Reset()
		improved := sim.optimiseTo(next, start)

		for improved != nil {
			fmt.Printf("\033[36mResult: %d, compared to %d\033[0m\n", improved.State.time, next.State.time)
			if improved.State.time < next.State.time {
				fmt.Print("IMPROVED\n")
				next = improved
			}

			fmt.Printf("\033[36mStarting with: %s\033[0m\n", next)
			improved = sim.optimiseTo(next, start)
		}

		if next.State.time >= 24 {
			break
		}

		fmt.Print("\n\033[42;30m=== NEW SOLUTION ===\033[0m\n")
		fmt.Printf("%s (time: %d)\n\n", next, next.State.time)
		schedule = next
		next = sim.buildNextRobot2(schedule, GeodeResource)
	}

	finalState := sim.simulateTime(schedule.State, 24-schedule.State.time)

	return schedule, finalState
}

func (sim *Simulator) solve3() SimState {
	start := NewSimState()
	state := start
	for _, resource := range resources[1:] {
		state, _ = sim.buildNextRobot(state, resource)
	}

	if state.time < sim.maxTime {
		state = sim.simulateTime(state, sim.maxTime-state.time)
	}

	// return sim.FindBest(start, schedule, final)
	return sim.FindBest(start, start)
}

// func (sim *Simulator) FindBest(prev *BuildStep, bestStep *BuildStep, bestState SimState) (*BuildStep, SimState) {
// 	var available []Resource

// 	for i, resource := range resources {
// 		if prev.State.robots[resource] == 0 {
// 			break
// 		}

// 		if resource != GeodeResource {
// 			available = resources[:i+2]
// 		} else {
// 			available = resources
// 		}
// 	}

// 	final := prev.State

// 	if final.time < 24 {
// 		final = sim.simulateTime(final, 24-final.time)
// 	}

// 	if final.time < bestState.time || final.Geodes() > bestState.Geodes() {
// 		bestStep = prev
// 		bestState = final
// 	}

// 	for _, resource := range available {
// 		next := sim.buildNextRobot2(prev, resource)

// 		if next.State.time >= 24 {
// 			continue
// 		}

// 		step, state := sim.FindBest(next, bestStep, bestState)

// 		if state.Geodes() > bestState.Geodes() {
// 			bestStep = step
// 			bestState = state
// 		}
// 	}

// 	return bestStep, bestState
// }

func (sim *Simulator) FindBest(prev SimState, bestState SimState) SimState {
	// var available []Resource

	final := prev

	if final.time < sim.maxTime {
		final = sim.simulateTime(final, sim.maxTime-final.time)
	}

	if final.time < bestState.time || final.Geodes() > bestState.Geodes() {
		bestState = final
	}

	for _, resource := range resources {
		if resource != OreResource && prev.robots[resource-1] == 0 {
			// We always need at least one of the previous robot
			continue
		}

		next, _ := sim.buildNextRobot(prev, resource)

		if next.time >= sim.maxTime {
			continue
		}

		bestState = sim.FindBest(next, bestState)

// 		if state.Geodes() > bestState.Geodes() {
// 			bestState = state
// 		}
	}

	return bestState
}

// func (sim *Simulator) FindBest(state SimState, bestState SimState) SimState {
// 	if state.time == sim.maxTime {
// 		// fmt.Printf("State: %s\n", state)
// 		if state.Geodes() > bestState.Geodes() {
// 			bestState = state
// 		}

// 		return bestState
// 	}

// 	state.time++

// 	if state.HasResources(sim.blueprint.Costs[GeodeResource]) {
// 		state.resources = state.resources.Add(state.robots).Sub(sim.blueprint.Costs[GeodeResource])
// 		state.robots[GeodeResource]++
// 		return sim.FindBest(state, bestState)
// 	}

// 	// fmt.Printf("state: %s\n", state)
// 	for resource, cost := range sim.blueprint.Costs {
// 		if Resource(resource) == GeodeResource {
// 			break
// 		}
// 		// fmt.Printf("cost for %s: %s\n", Resource(resource), cost)

// 		if state.HasResources(cost) {
// 			next := state
// 			// fmt.Printf("Resources BEFORE: %s\n", next.resources)
// 			next.resources = next.resources.Add(next.robots).Sub(cost)
// 			// fmt.Printf("Resources AFTER: %s\n", next.resources)
// 			next.robots[resource]++

// 			bestState = sim.FindBest(next, bestState)
// 		}
// 	}

// 	state.resources = state.resources.Add(state.robots)
// 	return sim.FindBest(state, bestState)

// // 	// var available []Resource

// // 	final := prev

// // 	if final.time < sim.maxTime {
// // 		final = sim.simulateTime(final, sim.maxTime-final.time)
// // 	}

// // 	if final.time < bestState.time || final.Geodes() > bestState.Geodes() {
// // 		bestState = final
// // 	}

// // 	for _, resource := range resources {
// // 		if resource != OreResource && prev.robots[resource-1] == 0 {
// // 			// We always need at least one of the previous robot
// // 			continue
// // 		}

// // 		next, _ := sim.buildNextRobot(prev, resource)

// // 		if next.time >= sim.maxTime {
// // 			continue
// // 		}

// // 		state := sim.FindBest(next, bestState)

// // 		if state.Geodes() > bestState.Geodes() {
// // 			bestState = state
// // 		}
// // 	}

// // 	return bestState
// }

type BuildStep struct {
	Robot      Resource
	State      SimState
	Prev       *BuildStep
	GatherTime Resources
	Cost Resources
}

func InitialBuildStep() *BuildStep {
	return &BuildStep{
		Robot: OreResource,
		State: SimState{
			robots: Resources{
				OreResource: 1,
			},
		},
	}
}

func (step *BuildStep) Reset() {
	if step == nil {
		return
	}

	step.GatherTime = step.Cost
	step.Prev.Reset()
}

func (step *BuildStep) String() string {
	if step.Prev == nil {
		return step.Robot.String()
	} else {
		return fmt.Sprintf("%s, %s", step.Prev, step.Robot)
	}
}

func (step *BuildStep) print() {
	if step.Prev != nil {
		step.Prev.print()
	}

	fmt.Printf("=== Time %d ===\n", step.State.time)
	fmt.Printf("Resources: %s\n", step.State.resources)
	fmt.Printf("Robots: %s\n", step.State.robots)
	fmt.Printf("Just got %s robot\n", step.Robot)
	fmt.Println()
}

var counter = 0

func (sim *Simulator) optimiseStep(schedule *BuildStep, robot Resource) *BuildStep {
	bestSch := schedule
	endState, gatherTime := sim.buildNextRobot(schedule.State, robot)

	maxTime := gatherTime.Max()

	for r, count := range gatherTime {
		resource := Resource(r)

		if count == 0 || count < maxTime || resource == robot {
			// Either:
			//    1) We already have enough resources (count == 0)
			//    2) it's not on the critical path (count < maxTime)
			//    3) it would try to build the same robot (e.g. ore as a
			//       critical resource for an ore miner)
			continue
		}

		// Note we're using the original build schedule - we're only
		// applying one change at a time
		sch := sim.optimiseStep(schedule, resource)

		if sch.State.time > endState.time {
			// Getting to this intermediate robot ends up taking even more time,
			// so skip it
			continue
		}

		result, _ := sim.buildNextRobot(sch.State, robot)

		// If result.time == end.time, it implies that we've improved this
		// resource (because we built a new robot without extending the build
		// time for this robot), but another resource has the same critical path
		// length.
		if result.time <= endState.time {
			bestSch = sch
			endState = result
		}
	}

	if bestSch != schedule {
		// New optimal build step inserted; we now need to optimise between this
		// new intermediate robot and the target robot
		return sim.optimiseStep(bestSch, robot)
	}

	// No other insertion was able to provide a better time than the current
	// schedule; the given robot is therefore the optimal next step
	return &BuildStep{
		Robot: robot,
		State: endState,
		Prev:  schedule,
	}
}

func (sim *Simulator) optimiseTo(step, after *BuildStep) *BuildStep {
	if step == after {
		return nil
	}

	if prev := sim.optimiseTo(step.Prev, after); prev != nil {
		return sim.buildNextRobot2(prev, step.Robot)
	}

	maxTime := step.GatherTime.Max()

	if maxTime == 0 {
		fmt.Printf("\033[31mCan't optimise %s robot\033[0m (prev: \033[34m%s\033[0m, after: \033[35m%s\033[0m)\n", step.Robot, step.Prev, after)
		// We already have enough resources, there are no improvements to be
		// made
		return nil
	}

	for r, count := range step.GatherTime {
		resource := Resource(r)
		if count < maxTime || resource == step.Robot {
			continue
		}

		fmt.Printf("\033[32mOptimising %s robot\033[0m (prev: \033[34m%s\033[0m, after: \033[35m%s\033[0m)\n", step.Robot, step.Prev, after)

		fmt.Printf("\033[33mTrying %s\033[0m\n", resource)

		// Clear it to prevent checking this resource again
		step.GatherTime[resource] = 0

		prev := sim.buildNextRobot2(step.Prev, resource)
		return sim.buildNextRobot2(prev, step.Robot)
	}

	// fmt.Printf("\033[31mNothing to optimise for %s\033[0m\n", step.Robot)
	return nil
}

func (sim *Simulator) buildNextRobot2(prev *BuildStep, robot Resource) *BuildStep {
	// state, gatherTime := sim.buildNextRobot(prev.State, robot)
	state, _ := sim.buildNextRobot(prev.State, robot)

	return &BuildStep{
		Robot:      robot,
		State:      state,
		Prev:       prev,
		// GatherTime: gatherTime,
		GatherTime: sim.robotCost(robot),
		Cost: sim.robotCost(robot),
	}
}

func (sim *Simulator) buildNextRobot(state SimState, robot Resource) (final SimState, gatherTime Resources) {
	cost := sim.robotCost(robot)
	maxWait := 0

	// fmt.Printf("Building %s robot\nCost: %s\nState:%s\n\n", robot, cost, state)

	missing := cost.Sub(state.resources)

	for resource, count := range missing {
		if count <= 0 {
			continue
		}

		gatherTime[resource] = CeilDiv(count, state.robots[resource])

		maxWait = Max(maxWait, gatherTime[resource])
	}

	// Add one minute to build the robot
	final = sim.simulateTime(state, maxWait+1)
	final.resources = final.resources.Sub(cost)
	final.robots[robot]++

	return
}

func (sim *Simulator) simulateTime(state SimState, dur int) SimState {
	state.time += dur
	state.resources = state.resources.Add(state.robots.Multiply(dur))
	return state
}

func (sim *Simulator) robotCost(robot Resource) Resources {
	return sim.blueprint.Costs[robot]
}

type SimState struct {
	time      int
	resources Resources
	robots    Resources
}

func NewSimState() SimState {
	return SimState{
		// nextRobot: InitialBuildSchedule(blueprint),
		robots: Resources{
			OreResource: 1,
		},
	}
}

func (state SimState) Geodes() int {
	return state.resources[GeodeResource]
}

func (state SimState) HasResources(res Resources) bool {
	for i, count := range state.resources {
		if count < res[i] {
			return false
		}
	}

	return true
}

func (state SimState) String() string {
	return fmt.Sprintf("Time: %d, Resources: %v, Robots: %v", state.time, state.resources, state.robots)
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func CeilDiv(a, b int) int {
	c := a / b

	if a%b != 0 {
		// Non-zero remainer, so round up
		c++
	}

	return c
}
