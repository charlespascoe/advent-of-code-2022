package main

import "strings"

type Cave []byte

func (c Cave) Copy() Cave {
	cp := make(Cave, len(c))
	copy(cp, c)
	return cp
}

func (c Cave) EnsureRows(length int) Cave {
	for len(c) < length {
		c = append(c, 0)
	}

	return c
}

func (c Cave) RowsFrom(row int) []byte {
	if row >= len(c) {
		return nil
	}

	return c[row:]
}

func (c Cave) String() string {
	var out strings.Builder

	for i := len(c) - 1; i >= 0; i-- {
		row := c[i]
		for j := 6; j >= 0; j-- {
			if row&(1<<j) > 0 {
				out.WriteRune('#')
			} else {
				out.WriteRune('.')
			}
		}

		out.WriteRune('\n')
	}

	return out.String()
}
