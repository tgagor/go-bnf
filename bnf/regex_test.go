package bnf_test

import (
	"testing"

	"github.com/tgagor/go-bnf/bnf"

	"github.com/stretchr/testify/assert"
)

func TestRegex_Basic(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<word> ::= /[a-zA-Z]+/
`)
	assert.NoError(t, err)

	// Should match
	ok, err := g.Match("hello")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = g.Match("WORLD")
	assert.NoError(t, err)
	assert.True(t, ok)

	// Should not match
	ok, err = g.Match("123")
	assert.False(t, ok)

	ok, err = g.Match("")
	assert.False(t, ok)
}

func TestRegex_Digits(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<number> ::= /\d+/
`)
	assert.NoError(t, err)

	ok, err := g.Match("123")
	assert.NoError(t, err)

	assert.True(t, ok)

	ok, err = g.Match("0")
	assert.True(t, ok)

	ok, err = g.Match("abc")
	assert.False(t, ok)
}

func TestRegex_WithSequence(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<assignment> ::= /[a-z]+/ "=" /\d+/
`)
	assert.NoError(t, err)

	ok, err := g.Match("x=123")
	assert.NoError(t, err)

	assert.True(t, ok)

	ok, err = g.Match("variable=42")
	assert.True(t, ok)

	ok, err = g.Match("X=123")
	assert.False(t, ok) // uppercase not allowed

	ok, err = g.Match("x=abc")
	assert.False(t, ok) // letters not allowed in number
}

func TestRegex_WithChoice(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<token> ::= /[a-z]+/ | /\d+/
`)
	assert.NoError(t, err)

	ok, err := g.Match("hello")
	assert.NoError(t, err)

	assert.True(t, ok)

	ok, err = g.Match("123")
	assert.True(t, ok)

	ok, err = g.Match("HELLO")
	assert.False(t, ok)
}

func TestRegex_WithRepetition(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<words> ::= /[a-z]+/ (" " /[a-z]+/)*
`)
	assert.NoError(t, err)

	ok, err := g.Match("hello")
	assert.NoError(t, err)

	assert.True(t, ok)

	ok, err = g.Match("hello world")
	assert.True(t, ok)

	ok, err = g.Match("one two three")
	assert.True(t, ok)
}

func TestRegex_ComplexPattern(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<email> ::= /[a-zA-Z0-9._%+-]+/ "@" /[a-zA-Z0-9-]+/ "." /[a-zA-Z]{2,}/
`)
	assert.NoError(t, err)

	ok, err := g.Match("user@example.com")
	assert.NoError(t, err)

	assert.True(t, ok)

	ok, err = g.Match("test.user+tag@sub-domain.org")
	assert.True(t, ok)

	ok, err = g.Match("invalid@")
	assert.False(t, ok)
}

func TestRegex_InvalidPattern(t *testing.T) {
	t.Parallel()

	_, err := bnf.LoadGrammarString(`
<bad> ::= /[unclosed/
`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid regex pattern")
}

func TestRegex_EscapedSlash(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<path> ::= /[a-z]+\/[a-z]+/
`)
	assert.NoError(t, err)

	ok, err := g.Match("foo/bar")
	assert.NoError(t, err)

	assert.True(t, ok)

	ok, err = g.Match("foobar")
	assert.False(t, ok)
}

func TestRegex_Anchors(t *testing.T) {
	t.Parallel()

	// Note: The regex matcher only matches from the current position,
	// so ^ is implicit. $ should work for end-of-string.
	g, err := bnf.LoadGrammarString(`
<exact> ::= /hello$/
`)
	assert.NoError(t, err)

	ok, err := g.Match("hello")
	assert.NoError(t, err)

	assert.True(t, ok)

	// This should fail because there's more after "hello"
	ok, err = g.Match("helloworld")
	assert.False(t, ok)
}

func TestRegex_AST(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
<word> ::= /[a-z]+/
`)
	assert.NoError(t, err)

	node, err := g.Parse("hello")
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "word", node.Type)
	assert.Equal(t, 1, len(node.Children))
	assert.Equal(t, "REGEX", node.Children[0].Type)
	assert.Equal(t, "hello", node.Children[0].Value)
}
