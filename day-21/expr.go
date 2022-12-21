package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var exprRe = regexp.MustCompile(`^(\w+): (?:(\d+)|(\w+) ([-+*/]) (\w+))`)

type Operator struct {
	Left, Op, Right string
}

type Constant struct {
	Val int
}

func EvalExpression(lines []string) int {
	exprs := make(map[string]any, len(lines))
	for _, line := range lines {
		match := exprRe.FindStringSubmatch(line)
		if match == nil {
			panic("invalid input")
		}

		name := match[1]

		if match[2] != "" {
			exprs[name] = Constant{
				Val: MustAtoi(match[2]),
			}
		} else {
			exprs[name] = Operator{
				Left:  match[3],
				Op:    match[4],
				Right: match[5],
			}
		}
	}

	return evalExpr(exprs, "root")
}

func evalExpr(exprs map[string]any, key string) int {
	switch val := exprs[key].(type) {
	case Constant:
		return val.Val
	case Operator:
		return evalOp(val, exprs)
	default:
		panic(fmt.Sprintf("unexpected type %T", exprs[key]))
	}
}

func evalOp(op Operator, exprs map[string]any) int {
	left := evalExpr(exprs, op.Left)
	right := evalExpr(exprs, op.Right)

	switch op.Op {
	case "-":
		return left - right
	case "+":
		return left + right
	case "*":
		return left * right
	case "/":
		// Just in case the result isn't an int
		if left % right != 0 {
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
