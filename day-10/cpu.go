package main

type Instruction struct{ Cycles, Delta int }

type CPU struct {
	X      int
	work   int
	prog   []Instruction
	pc     int
	Result int
}

func NewCPU(prog []Instruction) *CPU {
	return &CPU{
		X:    1,
		prog: prog,
	}
}

func (cpu *CPU) Done() bool {
	return cpu.pc >= len(cpu.prog)
}

func (cpu *CPU) Tick(cycles int) {
	if cycles%40 == 20 {
		cpu.Result += cycles * cpu.X
	}

	cpu.work++

	if cpu.work == cpu.prog[cpu.pc].Cycles {
		cpu.X += cpu.prog[cpu.pc].Delta
		cpu.pc++
		cpu.work = 0
	}
}
