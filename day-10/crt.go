package main

import (
	"fmt"
	"strings"
)

type CRT struct {
	x       *int
	display [][]rune
}

func NewCRT(x *int) *CRT {
	crt := &CRT{x: x}

	for i := 0; i < 6; i++ {
		line := make([]rune, 0, 40)
		for j := 0; j < cap(line); j++ {
			line = append(line, ' ')
		}

		crt.display = append(crt.display, line)
	}

	return crt
}

func (crt *CRT) String() string {
	lines := make([]string, 0, len(crt.display)+2)
	top := strings.Repeat("░", len(crt.display[0]) + 4)
	lines = append(lines, top)

	for _, l := range crt.display {
		lines = append(lines, fmt.Sprintf("░░%s░░", string(l)))
	}

	lines = append(lines, top)

	return strings.Join(lines, "\n")
}

func (crt *CRT) Tick(cycles int) {
	line, col := (cycles-1)/40, (cycles-1)%40

	diff := col - *crt.x

	if -1 <= diff && diff <= 1 {
		crt.display[line][col] = '█'
	}
}
