package main

import (
	"fmt"
	"regexp"
)

type Point struct {
	// Note: +ve X is right, +ve Y is down
	X, Y int
}

// ManhattanDist returns the Manhattan distance to the given point
func (p Point) ManhattanDist(other Point) int {
	return Abs(p.X-other.X) + Abs(p.Y-other.Y)
}

type Sensor struct {
	Pos Point

	// Beacon is the closest beacon to this sensor
	Beacon Point

	// Range is the Manhattan distance to the nearest beacon; no other beacon
	// can be found within this range.
	Range int
}

// InRange returns true if the given position falls within the range of the
// sensor
func (s Sensor) InRange(pos Point) bool {
	return s.Pos.ManhattanDist(pos) <= s.Range
}

// RangeRightEdge returns the x value of the rightmost edge of the sensor range
// for the given y value
func (s Sensor) RangeRightEdge(y int) int {
	xdelta := s.Range - Abs(s.Pos.Y-y)
	if xdelta < 0 {
		panic("y value not in range")
	}

	return s.Pos.X + xdelta
}

// XBounds returns the minimum and maximum X values of the scan ranges seen so
// far.
func XBounds(s Sensor, minX, maxX int) (newMinX, newMaxX int) {
	newMinX, _ = MinMax(s.Pos.X-s.Range, minX)
	_, newMaxX = MinMax(s.Pos.X+s.Range, maxX)
	return
}

var sensorRe = regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)

func ParseSensors(lines []string) ([]Sensor, error) {
	var sensors []Sensor

	for i, line := range lines {
		var sensor Sensor
		err := ReadGroups(line, sensorRe,
			&sensor.Pos.X,
			&sensor.Pos.Y,
			&sensor.Beacon.X,
			&sensor.Beacon.Y,
		)
		if err != nil {
			return nil, fmt.Errorf("line %d does not match expected format: %s", i+1, err)
		}

		sensor.Range = sensor.Pos.ManhattanDist(sensor.Beacon)

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}
