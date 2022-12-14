package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	testCondRe   = regexp.MustCompile(`divisible by (\d+)`)
	nextMonkeyRe = regexp.MustCompile(`throw to monkey (\d+)`)
)

type NextMonkey struct {
	Divisor int
	IfTrue  int
	IfFalse int
}

func (nm NextMonkey) Eval(worry int) int {
	if worry%nm.Divisor == 0 {
		return nm.IfTrue
	} else {
		return nm.IfFalse
	}
}

func ParseNextMonkey(stmt Statement) (NextMonkey, error) {
	match := testCondRe.FindStringSubmatch(stmt.Text)

	if match == nil {
		return NextMonkey{}, errors.New("not in expected format")
	}

	var err error
	var next NextMonkey

	if next.Divisor, err = strconv.Atoi(match[1]); err != nil {
		// This shouldn't happen due to regexp
		panic(err)
	}

	props := stmt.NestedProperties()

	if next.IfTrue, err = parseOutcome("If true", props); err != nil {
		return NextMonkey{}, err
	}

	if next.IfFalse, err = parseOutcome("If false", props); err != nil {
		return NextMonkey{}, err
	}

	return next, nil
}

func parseOutcome(key string, props map[string]Statement) (int, error) {
	prop, ok := props[key]
	if !ok {
		return 0, fmt.Errorf("'%s' property missing", key)
	}

	match := nextMonkeyRe.FindStringSubmatch(prop.Text)
	if match == nil {
		return 0, fmt.Errorf("'%s' value not in expected format", key)
	}

	return strconv.Atoi(match[1])
}
