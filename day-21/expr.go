package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var exprRe = regexp.MustCompile(`^(\w+): (?:(\d+)|(\w+) ([-+*/]) (\w+))`)

type Operator struct {
	Name            string
	Left, Op, Right string
}

type Constant struct {
	Name string
	Val  int
}

type ExprMap map[string]any



func buildExprMap(lines []string) ExprMap {
	exprs := make(ExprMap, len(lines))

	for _, line := range lines {
		match := exprRe.FindStringSubmatch(line)
		if match == nil {
			panic("invalid input")
		}

		name := match[1]

		if match[2] != "" {
			exprs[name] = Constant{
				Name: name,
				Val:  MustAtoi(match[2]),
			}
		} else {
			exprs[name] = Operator{
				Name:  name,
				Left:  match[3],
				Op:    match[4],
				Right: match[5],
			}
		}
	}

	return exprs
}

func (exprs ExprMap) Eval(key string) int {
	switch val := exprs[key].(type) {
	case Constant:
		return val.Val
	case Operator:
		return exprs.evalOp(val)
	default:
		panic(fmt.Sprintf("unexpected type %T", exprs[key]))
	}
}

func (exprs ExprMap) evalOp(op Operator) int {
	left := exprs.Eval(op.Left)
	right := exprs.Eval(op.Right)

	switch op.Op {
	case "-":
		return left - right
	case "+":
		return left + right
	case "*":
		return left * right
	case "/":
		// Just in case the result isn't an int
		if left%right != 0 {
			panic("not int division")
		}

		return left / right
	default:
		panic(fmt.Sprintf("unknown operator '%s'", op.Op))
	}
}

func MustAtoi(str string) int {
	if x, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return x
	}
}
