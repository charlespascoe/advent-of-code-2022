package main

import (
	"errors"
	"regexp"
	"strconv"
)

var opRe = regexp.MustCompile(`new = (old|\d+)\s+([*+])\s+(old|\d+)`)

type Operation func(worry int) int

func ParseOperation(str string) (Operation, error) {
	match := opRe.FindStringSubmatch(str)

	if match == nil {
		return nil, errors.New("not in expected format")
	}

	return parseArithmeticOperation(
		match[2],
		parseOperand(match[1]),
		parseOperand(match[3]),
	), nil
}

func parseArithmeticOperation(str string, a, b Operation) Operation {
	switch str {
	case "*":
		return multiplyOperation(a, b)
	case "+":
		return additionOperation(a, b)
	default:
		panic("unknown operation")
	}
}

func parseOperand(str string) Operation {
	if str == "old" {
		return oldValueOperand
	} else {
		x, err := strconv.Atoi(str)
		if err != nil {
			// This shouldn't happen due to regexp
			panic(err)
		}

		return constantOperand(x)
	}
}

func multiplyOperation(a, b Operation) Operation {
	return func(worry int) int {
		return a(worry) * b(worry)
	}
}

func additionOperation(a, b Operation) Operation {
	return func(worry int) int {
		return a(worry) + b(worry)
	}
}

func constantOperand(x int) Operation {
	return func(_ int) int {
		return x
	}
}

func oldValueOperand(worry int) int {
	return worry
}
