package main

type Device interface {
	Tick(cycles int)
}

type Clock struct {
	cycle   int
	devices []Device
}

func NewClock(devices ...Device) *Clock {
	return &Clock{
		devices: devices,
	}
}

func (c *Clock) Tick() {
	c.cycle++

	for _, d := range c.devices {
		d.Tick(c.cycle)
	}
}
