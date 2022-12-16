package main

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
)

func ReadLines(path string) ([]string, error) {
	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
	}
}

func MinMax(a, b int) (min, max int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func ReadGroups(str string, re *regexp.Regexp, groups ...any) error {
	match := re.FindSubmatch([]byte(str))
	if match == nil {
		return errors.New("does not match expected pattern")
	}

	if len(groups) >= len(match) {
		panic("too many groups")
	}

	for i, group := range groups {
		if err := json.Unmarshal(match[i+1], group); err != nil {
			return err
		}
	}

	return nil
}
