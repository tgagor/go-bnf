package bnf_test

import (
	"bytes"
	"go-bnf/bnf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrammar_Match(t *testing.T) {
	t.Parallel()

	grammarText := `
Expr ::= Term ("+" Term)*
Term ::= "a"
`
	g, err := bnf.LoadGrammarString(grammarText)
	assert.NoError(t, err)

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
			got, err := g.Match(tt.input)
			assert.Equal(t, tt.want, got)
			if tt.want {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestGrammar_Choice(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`S ::= "a" | "b"`)
	assert.NoError(t, err)

	ok, _ := g.Match("a")
	assert.True(t, ok)

	ok, _ = g.Match("b")
	assert.True(t, ok)

	ok, _ = g.Match("ab")
	assert.False(t, ok)
}

func TestGrammar_MatchFrom(t *testing.T) {
	t.Parallel()

	grammarText := `
<non_null_digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
<digit> ::= "0" | <non_null_digit>
<number> ::= <digit> | <non_null_digit> <number>
`
	g, err := bnf.LoadGrammarString(grammarText)
	assert.NoError(t, err)

	ok, _ := g.MatchFrom("number", "123")
	assert.True(t, ok)

	ok, _ = g.MatchFrom("digit", "5")
	assert.True(t, ok)
}

func TestGrammar_Parse(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`S ::= "a" "b"`)
	assert.NoError(t, err)

	node, err := g.Parse("ab")
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "S", node.Type)

	// Test ASTNode.String()
	treeStr := node.String()
	assert.Contains(t, treeStr, "(S")
	assert.Contains(t, treeStr, "\"a\"")
	assert.Contains(t, treeStr, "\"b\"")
}

func TestGrammar_Validate(t *testing.T) {
	t.Parallel()
	g, _ := bnf.LoadGrammarString(`S ::= "a"`)
	ok, _ := g.Validate("a")
	assert.True(t, ok)
	ok, _ = g.Validate("b")
	assert.False(t, ok)
}

func TestGrammar_ValidateGrammar(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`S ::= "a"`)
	assert.NoError(t, err)
	assert.NoError(t, g.ValidateGrammar())

	g2, _ := bnf.LoadGrammarString(`S ::= <undefined>`)
	assert.Error(t, g2.ValidateGrammar())
}

func TestGrammar_SetStart(t *testing.T) {
	t.Parallel()
	g, _ := bnf.LoadGrammarString(`
A ::= "a"
B ::= "b"
`)
	g.SetStart("B")
	ok, _ := g.Match("b")
	assert.True(t, ok)
	ok, _ = g.Match("a")
	assert.False(t, ok)
}

func TestGrammar_MatchPrefix(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`S ::= "a" "b"`)
	assert.NoError(t, err)

	assert.True(t, g.MatchPrefix("abc"))
	assert.False(t, g.MatchPrefix("ac"))
}

func TestParseError_Reporting(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`S ::= "a"`)
	assert.NoError(t, err)

	_, err = g.Match("b")
	assert.Error(t, err)

	pe, ok := err.(*bnf.ParseError)
	assert.True(t, ok)

	// Test Error()
	msg := pe.Error()
	assert.Contains(t, msg, "Parse error at line")

	// Test Pretty()
	pretty := pe.Pretty("b")
	assert.Contains(t, pretty, "expected one of: \"a\"")
}

func TestGrammar_Loaders(t *testing.T) {
	// LoadGrammarString
	g, err := bnf.LoadGrammarString(`S ::= "a"`)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	// LoadGrammar
	g2, err := bnf.LoadGrammar(bytes.NewReader([]byte(`S ::= "a"`)))
	assert.NoError(t, err)
	assert.NotNil(t, g2)

	// LoadGrammarFile
	g3, err := bnf.LoadGrammarFile("../examples/postal.bnf")
	assert.NoError(t, err)
	assert.NotNil(t, g3)
}

func TestGrammar_MultilineOr(t *testing.T) {
	t.Parallel()
	grammar := `
<DIGIT> ::= "0"
          | "1"
          | "2"
`
	g, err := bnf.LoadGrammarString(grammar)
	assert.NoError(t, err)

	ok, _ := g.MatchFrom("DIGIT", "1")
	assert.True(t, ok)
}

func TestGrammar_PostalAddress(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarFile("../examples/postal.bnf")
	assert.NoError(t, err)

	ok := []string{
		"John Smith\n123 Main St\nSpringfield, MA 02139\n",
		"J. Doe Jr.\n42 Elm St\nBoston, NY 10001\n",
	}

	for _, s := range ok {
		m, err := g.Match(s)
		assert.True(t, m, s)
		assert.NoError(t, err)
	}
}

func TestGrammar_Comments(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`
# config
a ::= "a" // comment
`)
	assert.NoError(t, err)

	m, _ := g.Match("a")
	assert.True(t, m)
}

func TestGrammar_RecurrentParenthesis(t *testing.T) {
	t.Parallel()

	g, err := bnf.LoadGrammarString(`a ::= (((("1" | "2") | "3") | "4") | "c")`)
	assert.NoError(t, err)

	m, _ := g.Match("1")
	assert.True(t, m)
	m, _ = g.Match("c")
	assert.True(t, m)
}
