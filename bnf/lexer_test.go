package bnf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	input := `
Expr ::= Term ("+" Term)*
Term ::= "a"
`
	lx := NewLexer(strings.NewReader(input))

	for {
		tok := lx.Next()
		t.Log(tok)
		if tok.Type == EOF {
			break
		}
	}
}

func TestBNF_EndToEnd_ExprGrammar(t *testing.T) {
	grammarText := `
Expr ::= Term ("+" Term)*
Term ::= "a"
`

	g := LoadGrammarString(grammarText)

	tests := []struct {
		input string
		want  bool
	}{
		{"a", true},
		{"a+a", true},
		{"a+a+a", true},
		{"", false},
		{"+a", false},
		{"a+", false},
		{"a++a", false},
		{"b", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := g.Match(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBNF_Choice(t *testing.T) {
	grammarText := `
S ::= "a" | "b"
`

	g := LoadGrammarString(grammarText)

	assert.True(t, g.Match("a"))
	assert.True(t, g.Match("b"))
	assert.False(t, g.Match("ab"))
	assert.False(t, g.Match(""))
}
