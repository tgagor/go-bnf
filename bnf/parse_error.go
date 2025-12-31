package bnf

import (
	"fmt"
	"strings"
)

type ParseError struct {
	Pos    int
	Line   int
	Column int

	RuleStack []string

	Expected []string
	Found    string

	Width int // number of characters to highlight
}

func (err *ParseError) Error() string {
	return fmt.Sprintf(
		"  Parse error at line %d, col %d\n  While matching rule: %s\n  Expected: %v\n  Found: %s\n",
		err.Line, err.Column, err.RuleStack, err.Expected, err.Found,
	)
}

func (e *ParseError) Pretty(input string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"Parse error at line %d, column %d, rule %s\n\n",
		e.Line, e.Column, e.RuleStack,
	))

	line := extractLine(input, e.Pos)
	sb.WriteString(line)
	sb.WriteByte('\n')

	// caret line
	for i := 1; i < e.Column; i++ {
		sb.WriteByte(' ')
	}
	for i := 0; i < max(1, e.Width); i++ {
		sb.WriteByte('^')
	}

	if e.Found != "" {
		sb.WriteString("\nFound: ")
		sb.WriteString(e.Found)
		sb.WriteString(" expected one of: ")
		terms := filterTerminals(e.Expected)
		if len(terms) == 0 {
			terms = e.Expected // fallback
		}
		sb.WriteString(strings.Join(terms, ", "))
	}

	return sb.String()
}

func extractLine(input string, pos int) string {
	runes := []rune(input)

	start := pos
	for start > 0 && runes[start-1] != '\n' {
		start--
	}

	end := pos
	for end < len(runes) && runes[end] != '\n' {
		end++
	}

	return string(runes[start:end])
}

func filterTerminals(expected []string) []string {
	var out []string
	seen := map[string]bool{}

	for _, e := range expected {
		// terminal = string literal
		if len(e) >= 2 && e[0] == '"' && e[len(e)-1] == '"' {
			if !seen[e] {
				seen[e] = true
				out = append(out, e)
			}
		}
	}
	return out
}
