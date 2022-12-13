package main

// ScanFunc is a function passed to scanMap that is evaluated once for every
// value in a TreeMap. 'reset' is true when starting a new row/column.
type ScanFunc func(val int, reset bool) int

// scanMap calls the ScanFunc on each value of the TreeMap, either along each
// row (scanCol == false) or each column (scanCol == true), forwards (rev ==
// false) or backwards (rev == true).
func scanMap(tm TreeMap, scanCol, rev bool, scanFunc ScanFunc) TreeMap {
	result := TreeMap{
		Rows:   tm.Rows,
		Cols:   tm.Cols,
		Values: make([]int, tm.Count()),
	}

	// Having this variable allows the tree map to be non-square
	resetInterval := tm.Rows

	if scanCol {
		resetInterval = tm.Cols
	}

	for i := 0; i < tm.Count(); i++ {
		// reset == true if we've found the start of a new row/column
		reset := i%resetInterval == 0

		pos := scanPos(
			reverse(i, tm.Count(), rev),
			tm,
			scanCol,
		)

		result.Values[pos] = scanFunc(tm.Values[pos], reset)
	}

	return result
}

// reverse inverts each value in an integer sequence so that a full iteration is
// reversed.
//
//	i = [0..n-1]
//	reverse(i, n, false) = [0..n-1]
//	reverse(i, n, true)  = [n-1..0]
func reverse(i, n int, rev bool) int {
	if rev {
		// i = [0..n-1]
		// (-1*i - 1) mod n = [n-1..0]
		return mod(-1*i-1, n)
	} else {
		return i
	}
}

// scanPos returns the TreeMap position based on the scan index i, depending on
// the scan direction (by rows or by columns). The scan index is assumed to
// increase/decrease by 1 on each scan step. To illustrate, a 3x3 TreeMap would
// have these position indexes:
//
//	┌───┬───┬───┐
//	│ 0 │ 1 │ 2 │
//	├───┼───┼───┤
//	│ 3 │ 4 │ 5 │
//	├───┼───┼───┤
//	│ 6 │ 7 │ 8 │
//	└───┴───┴───┘
//
// Row scanning is just incrementing i, but when scanCol == true, this function
// maps the sequence [0, 1, 2, 3, 4, 5, 6, 7, 8] to [0, 3, 6, 1, 4, 7, 2, 5, 8],
// which iterates over each column.
func scanPos(i int, tm TreeMap, scanCol bool) int {
	if scanCol {
		return tm.Cols*mod(i, tm.Rows) + i/tm.Rows
	} else {
		// TreeMap is already in row-major order; incrementing i by 1 advances
		// along a row
		return i
	}
}

// maxSeen returns a ScanFunc that returns the maximum value seen so far in each
// row/column during a scan (e.g. when scanning left to right along each row, it
// returns the maximum value of all values to the left). The first value is
// always -1.
func maxSeen() ScanFunc {
	max := -1
	return func(val int, reset bool) (result int) {
		if reset {
			max = -1
		}

		result = max

		if val > max {
			max = val
		}

		return
	}
}

// visibility returns a ScanFunc that returns the visibility looking back in the
// direction of the scan (e.g. when scanning left to right along each row, it
// returns the visibility looking towards the left)
func visibility() ScanFunc {
	// Index of the scan through the row/column; zero at the start of the scan,
	// regardless of scan direction (e.g. left -> right vs right -> left)
	i := 0
	// Slice of indexes of the last seen tree of each height, e.g. seen[5] == 4
	// means that the last tree of height 5 was seen when i == 4
	var seen []int

	return func(val int, reset bool) (result int) {
		if reset {
			// We only encounter trees heights 0-9, so length of 10 is fine
			seen = make([]int, 10)
			i = 0
			// We're at an edge, there's nothing else to see
			result = 0
		} else {
			i++
			// Find the nearest index of all trees that are at least as tall as
			// this tree
			result = i - max(seen[val:])
		}

		seen[val] = i

		return
	}
}

// mod computes the positive modulus of x mod n, e.g. mod(-1, 10) == 9
func mod(x, n int) int {
	return (n + (x % n)) % n
}

func max(nums []int) int {
	max := -1

	for _, x := range nums {
		if x > max {
			max = x
		}
	}

	return max
}
