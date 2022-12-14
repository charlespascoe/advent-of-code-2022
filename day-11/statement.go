package main

import (
	"fmt"
	"regexp"
	"strings"
)

var indentRe = regexp.MustCompile(`^\s*`)

type Line struct {
	Num    int
	Text   string
	Indent int
}

func toLines(input string) []Line {
	var lines []Line

	for i, line := range strings.Split(input, "\n") {
		indent := len(indentRe.FindString(line)) / 2

		lines = append(lines, Line{
			Num:    i + 1,
			Text:   strings.TrimSpace(line),
			Indent: indent,
		})
	}

	return lines
}

type Statement struct {
	Text   string
	Nested []Statement
}

func (s Statement) String() string {
	var lines []string
	for _, stmt := range s.Nested {
		for _, line := range strings.Split(stmt.String(), "\n") {
			lines = append(lines, fmt.Sprintf("+ %s", line))
		}
	}

	if len(lines) == 0 {
		return s.Text
	}

	return fmt.Sprintf("%s\n%s", s.Text, strings.Join(lines, "\n"))
}

func (s Statement) NestedProperties() map[string]Statement {
	result := make(map[string]Statement)

	for _, stmt := range s.Nested {
		key, val, _ := strings.Cut(stmt.Text, ":")
		val = strings.TrimSpace(val)

		// TODO: duplicate key?
		result[key] = Statement{
			Text:   val,
			Nested: stmt.Nested,
		}
	}

	return result
}

func ParseStatements(curIndent int, lines []Line) ([]Line, []Statement, error) {
	var result []Statement

	for len(lines) > 0 {
		line := lines[0]

		if len(line.Text) == 0 {
			lines = lines[1:]
			continue
		}

		if line.Indent > curIndent {
			if len(result) == 0 {
				return nil, nil, fmt.Errorf("line %d: indent too large", line.Num)
			}

			remainingLines, children, err := ParseStatements(curIndent+1, lines)
			if err != nil {
				return nil, nil, err
			}

			result[len(result)-1].Nested = children
			lines = remainingLines
			continue
		}

		if line.Indent < curIndent {
			break
		}

		result = append(result, Statement{
			Text: line.Text,
		})

		lines = lines[1:]
	}

	return lines, result, nil
}
