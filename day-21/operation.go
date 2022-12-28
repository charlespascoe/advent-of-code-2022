package main

import (
	"fmt"
)

type Operator interface {
	Apply(left, right int) int
	FindUnknown(left bool, known, result int) int
}

type AddOperator struct{}

func (AddOperator) Apply(left, right int) int {
	return left + right
}

func (AddOperator) FindUnknown(left bool, known, result int) int {
	// Either result = unknown + known or known + unknown;
	// either way, unknown = result - known
	return result - known
}

type SubOperator struct{}

func (SubOperator) Apply(left, right int) int {
	return left - right
}

func (SubOperator) FindUnknown(left bool, known, result int) int {
	if left {
		// result = unknown - known => unknown = result + known
		return result + known
	} else {
		// result = known - unknown => unknown = known - result
		return known - result
	}
}

type MultOperator struct{}

func (MultOperator) Apply(left, right int) int {
	return left * right
}

func (MultOperator) FindUnknown(left bool, known, result int) int {
	// Either result = unknown * known or known * unknown;
	// either way, unknown = result / known
	return result / known
}

type DivOperator struct{}

func (DivOperator) Apply(left, right int) int {
	return left / right
}

func (DivOperator) FindUnknown(left bool, known, result int) int {
	if left {
		// result = unknown / known => unknown = result * known
		return result * known
	} else {
		// result = known / unknown => unknown = known / result
		return known / result
	}
}

func GetOperation(str string) Operator {
	switch str {
	case "+":
		return AddOperator{}
	case "-":
		return SubOperator{}
	case "*":
		return MultOperator{}
	case "/":
		return DivOperator{}
	default:
		panic(fmt.Sprintf("unknown operator: '%s'", str))
	}
}
