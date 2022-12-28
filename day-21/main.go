package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input := flag.String("i", "input.txt", "program input")
	part2 := flag.Bool("part2", false, "print part 2 solution")
	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	lines, err := readLines(*input)
	if err != nil {
		log.Fatalf("Couldn't read input file: %s", err)
	}

	exprs := BuildStatementMap(lines)
	var result int

	if *part2 {
		exprs["humn"] = Unknown{}
	}

	ast := BuiltAST(exprs, "root")

	if *verbose {
		fmt.Printf("AST:\n%s\n", ast)
	}

	ast = ast.Simplify()

	if *verbose {
		fmt.Printf("Simplified AST:\n%s\n", ast)
	}

	if !*part2 {
		// Part 1 //
		result = (ast.(Literal)).Val
	} else {
		// Part 2 //
		root := ast.(Operation)

		var target int
		var other any

		if lit, ok := root.Left.(Literal); ok {
			target = lit.Val
			other = root.Right
		} else if lit, ok := root.Right.(Literal); ok {
			target = lit.Val
			other = root.Left
		} else {
			panic("ast not simplified")
		}

		result = findUnknown(other, target)
	}

	fmt.Printf("Result: %d\n", result)
}

func readLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func findUnknown(node any, result int) int {
	switch n := node.(type) {
	case Operation:
		var nextResult int
		var next any

		// Because we've simplified the AST, one branch will be always literal, the
		// other will always be either Unknown or Operation
		if lit, ok := n.Left.(Literal); ok {
			nextResult = n.Op.FindUnknown(false, lit.Val, result)
			next = n.Right
		} else if lit, ok := n.Right.(Literal); ok {
			nextResult = n.Op.FindUnknown(true, lit.Val, result)
			next = n.Left
		} else {
			panic("ast not simplified")
		}

		return findUnknown(next, nextResult)

	case Unknown:
		return result

	default:
		panic(fmt.Sprintf("unexpected type: %T", node))
	}
}
