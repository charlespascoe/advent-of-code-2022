package main

import (
	"fmt"
	"strings"
)

type Cave struct {
	space      [][]rune
	minX, maxX int
	empty      rune
	solid      rune
}

func NewCave() *Cave {
	return &Cave{
		minX:  500,
		maxX:  500,
		empty: ' ',
		solid: '█',
	}
}

func (c Cave) Width() int {
	return c.maxX + 1 - c.minX
}

func (c *Cave) Depth() int {
	return len(c.space)
}

func (c *Cave) DrawStructure(s Structure) {
	prev := s[0]

	for _, point := range s[1:] {
		c.DrawBetween(prev, point)
		prev = point
	}
}

func (c *Cave) DrawBetween(p1, p2 Vector) {
	d := p2.Sub(p1).Unit()

	for p := p1; p != p2; p = p.Add(d) {
		c.Draw(p, c.solid)
	}

	c.Draw(p2, c.solid)
}

func (c *Cave) AddFloor() {
	c.space = append(c.space, []rune(strings.Repeat(string(c.empty), 1000)))
	c.space = append(c.space, []rune(strings.Repeat(string(c.solid), 1000)))
}

func (c *Cave) Draw(p Vector, v rune) {
	for len(c.space) <= p.Y {
		c.space = append(c.space, []rune(strings.Repeat(string(c.empty), 1000)))
	}

	c.space[p.Y][p.X] = v

	if p.X < c.minX {
		c.minX = p.X
	}

	if p.X > c.maxX {
		c.maxX = p.X
	}
}

func (c *Cave) IsEmpty(p Vector) bool {
	return c.Get(p) == c.empty
}

func (c *Cave) Get(p Vector) rune {
	if p.Y >= len(c.space) {
		return c.empty
	}

	// TODO: Check X?
	return c.space[p.Y][p.X]
}

func (c *Cave) String() string {
	if c.Width() > 200 {
		return fmt.Sprintf("... (omitted due to being %d wide) ...", c.Width())
	}

	var lines []string

	for _, row := range c.space {
		lines = append(lines, string(row[c.minX-1:c.maxX+2]))
	}

	return strings.Join(lines, "\n")
}
