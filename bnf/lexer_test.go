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
	t.Parallel()

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
			got, err := g.Match(tt.input)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func TestBNF_Choice(t *testing.T) {
	t.Parallel()

	grammarText := `
S ::= "a" | "b"
`

	g := LoadGrammarString(grammarText)

	ok, err := g.Match("a")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("b")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.Match("ab")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.Match("")
	assert.False(t, ok)
	assert.Error(t, err)

}

func TestBNF_Numbers(t *testing.T) {
	t.Parallel()

	grammarText := `
<non_null_digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
<digit> ::= "0" | <non_null_digit>
<number> ::= <digit> | <non_null_digit> <number>
`

	g := LoadGrammarString(grammarText)

	ok, err := g.MatchFrom("number", "0") // single zero is fine
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("number", "01") // can't start with zero
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.MatchFrom("number", "11")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err =g.MatchFrom("number", "111")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err =g.MatchFrom("number", "1234567890")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err =g.MatchFrom("number", "")
	assert.False(t, ok) // not a number
	assert.Error(t, err)

}

func TestLexerNewlines(t *testing.T) {
	// All styles of new lines: Unix (\n), Windows (\r\n), old Mac (\r)
	input := "A ::= \"a\"\r\nB ::= \"b\"\rC ::= \"c\"\nD ::= \"d\""
	lx := NewLexer(strings.NewReader(input))

	var count int
	for {
		tok := lx.Next()
		if tok.Type == NEWLINE {
			count++
		}
		if tok.Type == EOF {
			break
		}
	}

	assert.Equal(t, 3, count)
}

func TestLexer_StringQuotes(t *testing.T) {
	t.Parallel()

	l := NewLexer(strings.NewReader(`"a" 'b' "c'd" 'e"f'`))

	a := l.Next()
	assert.Equal(t, STRING, a.Type)
	assert.Equal(t, "a", a.Text)

	b := l.Next()
	assert.Equal(t, STRING, b.Type)
	assert.Equal(t, "b", b.Text)

	c := l.Next()
	assert.Equal(t, STRING, c.Type)
	assert.Equal(t, "c'd", c.Text)

	d := l.Next()
	assert.Equal(t, STRING, d.Type)
	assert.Equal(t, `e"f`, d.Text)
}

func TestLexer_EmptyString(t *testing.T) {
	t.Parallel()

	l := NewLexer(strings.NewReader(`"" '' "a" 'b'`))

	tok := l.Next()
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "", tok.Text)

	tok = l.Next()
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "", tok.Text)

	tok = l.Next()
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "a", tok.Text)

	tok = l.Next()
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "b", tok.Text)
}

func TestPostalAddress(t *testing.T) {
	g, _ := LoadGrammarFile("../examples/postal.bnf")

	ok := []string{
		"John Smith\n123 Main St\nSpringfield, MA 02139\n",
		"J. Doe Jr.\n42 Elm St\nBoston, NY 10001\n",
		"John A. Jane Smith Sr.\n42 Elm St Apt12\nBoston, MA 10001\n",
	}

	bad := []string{
		"John Smith 123 Main St\nSpringfield, MA 02139\n",
		"John Smith\n123 Main St\nSpringfield, MA\n",
		"John Smith\n123 Main St Apt\nSpringfield, MA 02139\n",
	}

	for _, s := range ok {
		m, err := g.MatchFrom("postal-address", s)
		assert.True(t, m, s)
		assert.NoError(t, err)
	}

	for _, s := range bad {
		m, err := g.MatchFrom("postal-address", s)
		assert.False(t, m, s)
		assert.Error(t, err)
	}
}
