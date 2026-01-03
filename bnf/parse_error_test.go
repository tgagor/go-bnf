package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPretty_WidthFromExpected(t *testing.T) {
	err := &ParseError{
		Line:     1,
		Column:   3,
		Expected: []string{`"Main St"`},
		Width:    expectedWidth([]string{`"Main St"`}),
	}

	out := err.Pretty("42 Elm St")
	assert.Contains(t, out, "  ^^^^^^^")
}

func buildFloatGrammar() *Grammar {
	g, err := LoadGrammarString(`
<number>         ::= <non_zero_digit> <digit>+ | <digit> | <float>+
<non_zero_digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
<digit>          ::= <non_zero_digit> | "0"
<separator>      ::= "," | "."
<float>          ::= ( <non_zero_digit> <digit>+ | <digit> ) <separator> <digit>+
`)
	if err != nil {
		panic(err)
	}
	return g
}

func TestExpected_FilteredToTerminals(t *testing.T) {
	g := buildFloatGrammar()

	ok, err := g.Match("00")
	assert.False(t, ok)

	pe := err.(*ParseError)

	terms := filterTerminals(pe.Expected)
	expected := []string{`","`, `"."`}
	assert.ElementsMatch(t, expected, terms)
}
