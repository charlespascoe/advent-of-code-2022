package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Valve TODO: Description.
type Valve struct {
	Name    string
	Rate    int
	Tunnels []string
}

var valveRe = regexp.MustCompile(`Valve ([A-Z]+) has flow rate=(\d+); tunnels? leads? to valves? ([A-Z, ]+)`)

func ParseValves(lines []string) ([]Valve, error) {
	var valves []Valve

	for i, line := range lines {
		match := valveRe.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("line %d not in expected format: '%s'", i+1, line)
		}

		valves = append(valves, Valve{
			Name:    match[1],
			Rate:    MustAtoi(match[2]),
			Tunnels: strings.Split(match[3], ", "),
		})
	}

	return valves, nil
}
