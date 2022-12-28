package main

import (
	"fmt"
)

type Simulator struct {
	cave     Cave
	input    string
	rock     Rock
	rockIdx  int
	inputIdx int
	pos      int
	newRock  bool
}

func NewSimulator(input string) *Simulator {
	return &Simulator{
		input:   input,
		pos:     3,
		rock:    rocks[0],
		newRock: true,
	}
}

func (sim *Simulator) Step() {
	sim.newRock = false
	jet := sim.input[sim.inputIdx%len(sim.input)]

	if movedRock := sim.rock.Shift(jet == '<'); !movedRock.OverlapsWith(sim.cave.RowsFrom(sim.pos)) {
		// There's nothing blocking us to move this way
		sim.rock = movedRock
	}

	sim.inputIdx++

	if sim.pos > 0 && !sim.rock.OverlapsWith(sim.cave.RowsFrom(sim.pos-1)) {
		sim.pos--
		return
	}

	// Rock has stopped
	sim.cave = sim.cave.EnsureRows(sim.pos + len(sim.rock))
	sim.rock.WriteTo(sim.cave.RowsFrom(sim.pos))

	// Next rock
	sim.rockIdx++
	sim.newRock = true
	sim.rock = rocks[sim.rockIdx%len(rocks)]
	sim.pos = len(sim.cave) + 3
}

func (sim *Simulator) HasSameState(other *Simulator) bool {
	r1, i1, n1 := sim.State()
	r2, i2, n2 := other.State()

	// We have to see if the top of the caves are the same; this way, we know
	// all future rocks will fall and stop in the same way
	topEqual := sim.cave.TopUntilBlocked().Equals(other.cave.TopUntilBlocked())

	return r1 == r2 && i1 == i2 && n1 == n2 && topEqual
}

func (sim *Simulator) State() (rockIdx, inputIdx int, newRock bool) {
	return sim.rockIdx % len(rocks), sim.inputIdx % len(sim.input), sim.newRock
}

func (sim *Simulator) StateString() string {
	r, i, n := sim.State()
	return fmt.Sprintf("R: %d, I: %d, N: %t", r, i, n)
}

func (sim *Simulator) RocksSettled() int {
	return sim.rockIdx
}

func (sim *Simulator) Height() int {
	return len(sim.cave)
}
