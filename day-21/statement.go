package main

import (
	"regexp"
	"strconv"
)

var exprRe = regexp.MustCompile(`^(\w+): (?:(\d+)|(\w+) ([-+*/]) (\w+))`)

type OperationStatement struct {
	Left, Op, Right string
}

type LiteralStatement struct {
	Val int
}

type StatementMap map[string]any

func BuildStatementMap(lines []string) StatementMap {
	exprs := make(StatementMap, len(lines))

	for _, line := range lines {
		match := exprRe.FindStringSubmatch(line)
		if match == nil {
			panic("invalid input")
		}

		name := match[1]

		if match[2] != "" {
			exprs[name] = LiteralStatement{
				Val: MustAtoi(match[2]),
			}
		} else {
			exprs[name] = OperationStatement{
				Left:  match[3],
				Op:    match[4],
				Right: match[5],
			}
		}
	}

	return exprs
}

func MustAtoi(str string) int {
	if x, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return x
	}
}
