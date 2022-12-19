package main

type Simulator struct {
	cave     Cave
	input    string
	rock     Rock
	rockIdx  int
	inputIdx int
	pos      int
}

func NewSimulator(input string) *Simulator {
	return &Simulator{
		input: input,
		pos:   3,
		rock:  rocks[0],
	}
}

func (sim *Simulator) Step() {
	jet := sim.input[sim.inputIdx%len(sim.input)]

	if movedSh := sim.rock.Shift(jet == '<'); !movedSh.OverlapsWith(sim.cave.RowsFrom(sim.pos)) {
		// There's nothing blocking us to move this way
		sim.rock = movedSh
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
	sim.rock = rocks[sim.rockIdx%len(rocks)]
	sim.pos = len(sim.cave) + 3
}

func (sim *Simulator) RocksSettled() int {
	return sim.rockIdx
}

func (sim *Simulator) Height() int {
	return len(sim.cave)
}
