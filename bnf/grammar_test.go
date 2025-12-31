package bnf

import (
	"fmt"
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

	ok, err := g.Match("b")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("ab")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("aaab")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("aaa")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.Match("ba")
	assert.False(t, ok)
	assert.Error(t, err)
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

func buildComplexGrammar() *Grammar {
	// Term ::= "a"
	// Expr ::= Term ("+" Term)*
	expr := &Rule{Name: "Expr"}
	term := &Rule{Name: "Term"}

	// Term ::= "a"
	term.Expr = &Terminal{Value: "a"}
	// exprNT := &NonTerminal{Name: "Expr", Rule: expr}
	termNT := &NonTerminal{Name: "Term", Rule: term}

	// Expr ::= Term ("+" Term)*
	expr.Expr = &Sequence{
		Elements: []Node{
			termNT,
			&Repeat{
				Node: &Sequence{
					Elements: []Node{
						&Terminal{Value: "+"},
						termNT,
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

func TestComplexGrammar(t *testing.T) {
	t.Parallel()
	g := buildComplexGrammar()
	fmt.Println("Grammar:", g.Rules["Expr"].Expr)

	ok, err := g.Match("a")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("a+a")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("a+a+a")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.Match("+a")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.Match("a+")
	assert.False(t, ok)
	assert.Error(t, err)
}

func buildLeftRecursiveGrammar() *Grammar {
	// Expr ::= Expr "+" Term | Term
	// Term ::= "a"
	expr := &Rule{Name: "Expr"}
	term := &Rule{Name: "Term"}

	exprNT := &NonTerminal{Name: "Expr", Rule: expr}
	termNT := &NonTerminal{Name: "Term", Rule: term}

	// Term ::= "a"
	term.Expr = &Terminal{Value: "a"}

	// Expr ::= Expr "+" Term | Term
	expr.Expr = &Choice{
		Options: []Node{
			&Sequence{
				Elements: []Node{
					exprNT,
					&Terminal{Value: "+"},
					termNT,
				},
			},
			termNT,
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

	ok, err := g.Match("a")
	assert.True(t, ok)
	assert.NoError(t, err)

	// assert.True(t, g.Match("a+a"))   // left recursion detected
	// assert.True(t, g.Match("a+a+a")) // left recursion detected

	ok, err = g.Match("")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.Match("+a")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.Match("a+")
	assert.False(t, ok)
	assert.Error(t, err)
}

func TestRepeatAlone(t *testing.T) {
	term := &Terminal{Value: "a"}

	expr := &Sequence{
		Elements: []Node{
			term,
			&Repeat{
				Node: &Sequence{
					Elements: []Node{
						&Terminal{Value: "+"},
						term,
					},
				},
				Min: 0,
			},
		},
	}

	assert.Equal(t, []int{1, 3, 5}, match(expr, "a+a+a", 0))
}
