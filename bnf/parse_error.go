package bnf

import "fmt"

type ParseError struct {
	Pos    int
	Line   int
	Column int

	RuleStack []string

	Expected []string
	Found    string
}

func (err *ParseError) Error() string {
	return fmt.Sprintf(
		"  Parse error at line %d, col %d\n  While matching rule: %s\n  Expected: %v\n  Found: %s\n",
		err.Line, err.Column, err.RuleStack, err.Expected, err.Found,
	)
}
