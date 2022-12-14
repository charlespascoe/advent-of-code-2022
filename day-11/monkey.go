package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Monkey struct {
	Num          int
	Items        []int
	InspectCount int

	op      Operation
	next    NextMonkey
	monkeys []*Monkey
	modulus int
}

func (m *Monkey) String() string {
	return fmt.Sprintf("Monkey %d items: %v", m.Num, m.Items)
}

func (m *Monkey) Round() {
	for len(m.Items) > 0 {
		worry := m.Items[0]
		worry = m.op(worry)

		if *part2 {
			worry %= m.modulus
		} else {
			worry /= 3
		}

		m.InspectCount++

		next := m.next.Eval(worry)

		m.Items = m.Items[1:]
		m.monkeys[next].Catch(worry)
	}
}

func (m *Monkey) Catch(item int) {
	m.Items = append(m.Items, item)
}

var monkeyRe = regexp.MustCompile(`^Monkey (\d+):`)

func ParseMonkeys(stmts []Statement) ([]*Monkey, error) {
	var monkeys []*Monkey

	modulus := 1

	for i, stmt := range stmts {
		match := monkeyRe.FindStringSubmatch(stmt.Text)

		if match == nil {
			return nil, fmt.Errorf("statement '%s': not in expected format", stmt.Text)
		}

		x, err := strconv.Atoi(match[1])
		if err != nil {
			// This shouldn't happen due to regexp
			panic(err)
		}

		if x != i {
			return nil, fmt.Errorf("statement '%s': monkey index out-of-order", stmt.Text)
		}

		monkey, err := parseMonkey(stmt)
		if err != nil {
			return nil, fmt.Errorf("%s %s", stmt.Text, err)
		}

		monkey.Num = x

		modulus *= monkey.next.Divisor

		monkeys = append(monkeys, monkey)
	}

	for _, monkey := range monkeys {
		// TODO: Better way
		monkey.monkeys = monkeys
		monkey.modulus = modulus
	}

	return monkeys, nil
}

func parseMonkey(stmt Statement) (*Monkey, error) {
	monkey := &Monkey{}
	props := stmt.NestedProperties()

	if startItems, ok := props["Starting items"]; ok {
		for _, item := range strings.Split(startItems.Text, ", ") {
			x, err := strconv.Atoi(item)
			if err != nil {
				return nil, fmt.Errorf("invalid starting item '%s': %s", item, err)
			}

			monkey.Items = append(monkey.Items, x)
		}
	} else {
		return nil, errors.New("starting items missing")
	}

	if opStmt, ok := props["Operation"]; ok {
		op, err := ParseOperation(opStmt.Text)
		if err != nil {
			return nil, fmt.Errorf("invalid operation '%s': %s", opStmt.Text, err)
		}

		monkey.op = op
	} else {
		return nil, errors.New("operation missing")
	}

	if testStmt, ok := props["Test"]; ok {
		next, err := ParseNextMonkey(testStmt)
		if err != nil {
			return nil, fmt.Errorf("invalid test '%s': %s", testStmt.Text, err)
		}

		monkey.next = next
	} else {
		return nil, errors.New("test missing")
	}

	return monkey, nil
}
