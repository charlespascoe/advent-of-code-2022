package main

import (
	"strconv"
)

func MustAtoi(str string) int {
	if x, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return x
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
