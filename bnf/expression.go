package bnf

import (
	"slices"
	"strings"
)

type Expression struct {
	Name string
	// Patterns []*Pattern
	Patterns []string
	// Expressions []any
}

func (b *BNF) Match(expName string, val string) bool {
	exp := b.Symbols[expName]
	return slices.Contains(exp.Patterns, val)
}

func (b *BNF) clean(symbol string) string {
	if strings.HasPrefix(symbol, "\"") && strings.HasSuffix(symbol, "\"") {
		return strings.Trim(symbol, "\"")
	} else if strings.HasPrefix(symbol, "<") && strings.HasSuffix(symbol, ">") {
		val, ok := b.Symbols[symbol]
		if ok {
			b.clean()
		}
	}

	return ""
}
