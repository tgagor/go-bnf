package bnf

import (
	"fmt"
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
	return LoadGrammarString(`
<number>         ::= <non_zero_digit> <digit>+ | <digit> | <float>+
<non_zero_digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
<digit>          ::= <non_zero_digit> | "0"
<separator>      ::= "," | "."
<float>          ::= ( <non_zero_digit> <digit>+ | <digit> ) <separator> <digit>+
`)
}

func TestExpected_FilteredToTerminals(t *testing.T) {
	g := buildFloatGrammar()

	ok, err := g.Match("00")
	assert.False(t, ok)

	pe := err.(*ParseError)

	terms := filterTerminals(pe.Expected)
	expected := []string{`","`, `"."`}
	for i := 1; i <= 9; i++ {
		expected = append(expected, fmt.Sprintf(`"%d"`, i))
	}
	assert.ElementsMatch(t, expected, terms)
}
