package main

import (
	"fmt"
)

const MaxInt64 = 0x7fffffffffffffff

type Simulator struct {
	blueprint Blueprint
}

func NewSimulator(blueprint Blueprint) *Simulator {
	return &Simulator{
		blueprint: blueprint,
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
	for schedule.State.time < 24 {
		lastGeodeMiner = schedule
		schedule = sim.optimiseStep(schedule, GeodeResource)
	}

	schedule = lastGeodeMiner

	for schedule.State.time > 24 {
		// The last state ih the schedule has time >= 24, which is too late to
		// have any effect on the result, so pop it
		schedule = schedule.Prev
	}

	// Pop step 0 (time == 0, the initial ore miner we already had)
	// schedule = schedule[1:]

	finalState := schedule.State

	// Simulate final state for resource counts
	finalState = sim.simulateTime(finalState, 24-finalState.time)

	return schedule, finalState
}

type BuildStep struct {
	Robot Resource
	State SimState
	Prev  *BuildStep
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

func (sim *Simulator) buildNextRobot(state SimState, robot Resource) (final SimState, gatherTime Resources) {
	cost := sim.robotCost(robot)
	maxWait := 0

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
