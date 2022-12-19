package main

type Rock []byte

func (rock Rock) Shift(left bool) Rock {
	res := make(Rock, 0, len(rock))
	var edge byte = 1
	if left {
		edge <<= 6
	}

	for _, row := range rock {
		if row&edge != 0 {
			// End bit is set; we can't move any further in this direction,
			// return the original rock
			return rock
		}

		if left {
			res = append(res, row<<1)
		} else {
			res = append(res, row>>1)
		}
	}

	return res
}

func (rock Rock) OverlapsWith(rows []byte) bool {
	for i, row := range rows {
		if i >= len(rock) {
			break
		} else if row&rock[i] > 0 {
			return true
		}
	}

	return false
}

func (rock Rock) WriteTo(rows []byte) {
	// Assumes that len(rows) >= len(rock)
	for i := range rock {
		rows[i] |= rock[i]
	}
}

// Note: The first row is closest to the floor, so these rocks are essentially
// upside down. All rocks are in their starting position.
var rocks = []Rock{
	{0b0011110},
	{
		0b0001000,
		0b0011100,
		0b0001000,
	},
	{
		0b0011100,
		0b0000100,
		0b0000100,
	},
	{
		0b0010000,
		0b0010000,
		0b0010000,
		0b0010000,
	},
	{
		0b0011000,
		0b0011000,
	},
}
