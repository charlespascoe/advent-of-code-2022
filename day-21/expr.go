package main

import (
	"fmt"
	"strings"
)

type ExprNode interface {
	Simplify() ExprNode
	String() string
}

type Operation struct {
	Name        string
	Left, Right ExprNode
	Op          Operator
	opStr       string
}

func (op Operation) Simplify() ExprNode {
	op.Left = op.Left.Simplify()
	op.Right = op.Right.Simplify()

	left, lok := op.Left.(Literal)
	right, rok := op.Right.(Literal)
	if lok && rok {
		return Literal{
			Name: op.Name,
			Val:  op.Op.Apply(left.Val, right.Val),
		}
	}

	return op
}

func (op Operation) String() string {
	var str strings.Builder

	str.WriteString(fmt.Sprintf("%s: %s\n    ", op.Name, op.opStr))
	str.WriteString(strings.Join(strings.Split(op.Left.String(), "\n"), "\n    "))
	str.WriteString("\n    ")
	str.WriteString(strings.Join(strings.Split(op.Right.String(), "\n"), "\n    "))

	return str.String()
}

type Literal struct {
	Name string
	Val  int
}

func (con Literal) Simplify() ExprNode {
	return con
}

func (con Literal) String() string {
	return fmt.Sprintf("%s: %d", con.Name, con.Val)
}

type Unknown struct{}

func (unk Unknown) Simplify() ExprNode {
	return unk
}

func (unk Unknown) String() string {
	return "Unknown"
}

func BuiltAST(exprs StatementMap, key string) ExprNode {
	switch expr := exprs[key].(type) {
	case OperationStatement:
		return Operation{
			Name:  key,
			Left:  BuiltAST(exprs, expr.Left),
			Right: BuiltAST(exprs, expr.Right),
			Op:    GetOperation(expr.Op),
			opStr: expr.Op,
		}

	case LiteralStatement:
		return Literal{
			Name: key,
			Val:  expr.Val,
		}

	case Unknown:
		return Unknown{}

	default:
		panic(fmt.Sprintf("unexpected type %T", exprs[key]))
	}
}
