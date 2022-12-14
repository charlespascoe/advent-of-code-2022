package main

import (
	"strconv"
	"strings"
)

type Structure []Vector

func ParseStructures(lines []string) []Structure {
	var sts []Structure

	for _, line := range lines {
		var st Structure
		for _, point := range strings.Split(line, " -> ") {
			xs, ys, found := strings.Cut(point, ",")
			if !found {
				panic("invalid point")
			}

			x, err := strconv.Atoi(xs)
			if err != nil {
				panic(err)
			}

			y, err := strconv.Atoi(ys)
			if err != nil {
				panic(err)
			}

			st = append(st, Vector{x, y})
		}

		sts = append(sts, st)
	}

	return sts
}
