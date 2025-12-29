package bnf_test

import (
	"bnf-test/bnf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildSimpleGrammar() *bnf.Grammar {
	// S ::= "a"* "b"
	return &bnf.Grammar{
		Start: "S",
		Rules: map[string]*bnf.Rule{
			"S": {
				Name: "S",
				Expr: &bnf.Sequence{
					Elements: []bnf.Node{
						&bnf.Repeat{
							Node: &bnf.Terminal{"a"},
							Min:  0,
						},
						&bnf.Terminal{"b"},
					},
				},
			},
		},
	}
}

func TestGrammar(t *testing.T) {
	t.Parallel()
	g := buildSimpleGrammar()

	assert.True(t, g.Match("b"))    // true
	assert.True(t, g.Match("ab"))   // true
	assert.True(t, g.Match("aaab")) // true
	assert.False(t, g.Match("aaa")) // false
	assert.False(t, g.Match("ba"))  // false
}

func TestGrammarPrefix(t *testing.T) {
	t.Parallel()
	// S ::= "a"* "b"
	g := buildSimpleGrammar()

	assert.True(t, g.MatchPrefix("b"))    // true
	assert.True(t, g.MatchPrefix("ab"))   // true
	assert.True(t, g.MatchPrefix("aaab")) // true
	assert.False(t, g.MatchPrefix("aaa")) // false
	assert.True(t, g.MatchPrefix("ba"))   // matches b prefix
}

func buildLeftRecursiveGrammar() *bnf.Grammar {
	// Expr ::= Expr "+" Term | Term
	// Term ::= "a"
	expr := &bnf.Rule{Name: "Expr"}
	term := &bnf.Rule{Name: "Term"}

	// Term ::= "a"
	term.Expr = &bnf.Terminal{Value: "a"}

	// Expr ::= Expr "+" Term | Term
	expr.Expr = &bnf.Choice{
		Options: []bnf.Node{
			&bnf.Sequence{
				Elements: []bnf.Node{
					&bnf.NonTerminal{Name: "Expr", Rule: expr},
					&bnf.Terminal{Value: "+"},
					&bnf.NonTerminal{Name: "Term", Rule: term},
				},
			},
			&bnf.NonTerminal{Name: "Term", Rule: term},
		},
	}

	return &bnf.Grammar{
		Start: "Expr",
		Rules: map[string]*bnf.Rule{
			"Expr": expr,
			"Term": term,
		},
	}
}

func TestRecursiveGrammar(t *testing.T) {
	t.Parallel()
	g := buildLeftRecursiveGrammar()

	assert.True(t, g.Match("a"))     // true
	assert.True(t, g.Match("a+a"))   // true
	assert.True(t, g.Match("a+a+a")) // true
	assert.False(t, g.Match(""))     // false
	assert.False(t, g.Match("+a"))   // false
	assert.False(t, g.Match("a+"))   // false
}
