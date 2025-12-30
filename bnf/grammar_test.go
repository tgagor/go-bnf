package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildSimpleGrammar() *Grammar {
	// S ::= "a"* "b"
	return &Grammar{
		Start: "S",
		Rules: map[string]*Rule{
			"S": {
				Name: "S",
				Expr: &Sequence{
					Elements: []Node{
						&Repeat{
							Node: &Terminal{"a"},
							Min:  0,
						},
						&Terminal{"b"},
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

func buildLeftRecursiveGrammar() *Grammar {
	// Expr ::= Term ("+" Term)*
	// Term ::= "a"
	expr := &Rule{Name: "Expr"}
	term := &Rule{Name: "Term"}

	// Term ::= "a"
	term.Expr = &Terminal{Value: "a"}

	// Expr ::= Term ("+" Term)*
	expr.Expr = &Sequence{
		Elements: []Node{
			&NonTerminal{Name: "Term", Rule: term},
			&Repeat{
				Node: &Sequence{
					Elements: []Node{
						&Terminal{Value: "+"},
						&NonTerminal{Name: "Term", Rule: term},
					},
				},
				Min: 0,
			},
		},
	}

	return &Grammar{
		Start: "Expr",
		Rules: map[string]*Rule{
			"Expr": expr,
			"Term": term,
		},
	}
}

func TestRecursiveGrammar(t *testing.T) {
	t.Parallel()
	g := buildLeftRecursiveGrammar()

	// FIXME: those commented don't pass, but they should
	assert.True(t, g.Match("a"))     // true
	assert.True(t, g.Match("a+a"))   // true
	assert.True(t, g.Match("a+a+a")) // true
	assert.False(t, g.Match(""))     // false
	assert.False(t, g.Match("+a"))   // false
	assert.False(t, g.Match("a+"))   // false
}
