package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type Sensor struct {
	Pos Vector

	// Beacon is the closest beacon to this sensor
	Beacon Vector

	// Range is the distance to the nearest beacon; no other beacon can be found
	// within this range.
	Range int
}

// InRange returns true if the given position falls within the range of the
// sensor
func (s Sensor) InRange(pos Vector) bool {
	return pos.Sub(s.Pos).Magnitude() <= s.Range
}

// RangeRightEdge returns the x value of the rightmost edge of the sensor range
// for the given y value
func (s Sensor) RangeRightEdge(y int) int {
	xdelta := s.Range - abs(s.Pos.Y-y)

	if xdelta < 0 {
		panic("y value not in range")
	}

	return s.Pos.X + xdelta
}

// XBounds returns the minimum and maximum X values of the scan ranges seen so
// far.
func XBounds(s Sensor, minX, maxX int) (newMinX, newMaxX int) {
	newMinX = minX
	newMaxX = maxX

	mag := s.Beacon.Sub(s.Pos).Magnitude()

	if s.Pos.X-mag < newMinX {
		newMinX = s.Pos.X - mag
	}

	if s.Pos.X+mag > newMaxX {
		newMaxX = s.Pos.X + mag
	}

	return
}

var sensorRe = regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)

func ParseSensors(lines []string) ([]Sensor, error) {
	var sensors []Sensor

	for i, line := range lines {
		match := sensorRe.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("line %d does not match expected format", i+1)
		}

		// Don't worry about errors - just rely on regular expression
		pos := Vector{
			X: MustAtoi(match[1]),
			Y: MustAtoi(match[2]),
		}

		beacon := Vector{
			X: MustAtoi(match[3]),
			Y: MustAtoi(match[4]),
		}

		sensors = append(sensors, Sensor{
			Pos:    pos,
			Beacon: beacon,
			Range:  beacon.Sub(pos).Magnitude(),
		})
	}

	return sensors, nil
}

func MustAtoi(str string) int {
	if x, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return x
	}
}
