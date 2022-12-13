package main

// TreeMap is a grid of values in row-major order. TreeMap.Values[0] is the
// top-left corner.
type TreeMap struct {
	Values []int
	Rows   int
	Cols   int
}

func (tm TreeMap) Count() int {
	return len(tm.Values)
}

func (tm TreeMap) String() string {
	out := make([]byte, 0, len(tm.Values)+tm.Rows)

	for pos, val := range tm.Values {
		if pos > 0 && pos%tm.Cols == 0 {
			out = append(out, '\n')
		}

		if val > 9 {
			panic("TreeMap.String() only supports values 0-9")
		}

		// Characters '0' to '9' are all single bytes in UTF-8, so this is safe
		out = append(out, byte(val)+'0')
	}

	return string(out)
}

func loadTreeMap(input []string) TreeMap {
	tm := TreeMap{
		Rows: len(input),
		Cols: len(input[0]),
	}

	tm.Values = make([]int, 0, tm.Rows*tm.Cols)

	for _, row := range input {
		for _, char := range row {
			tm.Values = append(tm.Values, int(char-'0'))
		}
	}

	return tm
}
