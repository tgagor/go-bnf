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
		tok, err := lx.Next()
		if err != nil {
			t.Logf("Lexer error: %v", err)
			break
		}
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

	g, err := LoadGrammarString(grammarText)
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
			} else { // failed match generate error
				assert.Error(t, err)
			}
		})
	}
}

func TestBNF_Choice(t *testing.T) {
	t.Parallel()

	grammarText := `
S ::= "a" | "b"
`

	g, err := LoadGrammarString(grammarText)
	assert.NoError(t, err)

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

	g, err := LoadGrammarString(grammarText)
	assert.NoError(t, err)

	ok, err := g.MatchFrom("number", "0") // single zero is fine
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("number", "01") // can't start with zero
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.MatchFrom("number", "11")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("number", "111")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("number", "1234567890")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("number", "")
	assert.False(t, ok) // not a number
	assert.Error(t, err)

}

func TestLexer_StringQuotes(t *testing.T) {
	t.Parallel()

	l := NewLexer(strings.NewReader(`"a" 'b' "c'd" 'e"f'`))

	a, err := l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, a.Type)
	assert.Equal(t, "a", a.Text)

	b, err := l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, b.Type)
	assert.Equal(t, "b", b.Text)

	c, err := l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, c.Type)
	assert.Equal(t, "c'd", c.Text)

	d, err := l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, d.Type)
	assert.Equal(t, `e"f`, d.Text)
}

func TestLexer_EmptyString(t *testing.T) {
	t.Parallel()

	l := NewLexer(strings.NewReader(`"" '' "a" 'b'`))

	tok, err := l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "", tok.Text)

	tok, err = l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "", tok.Text)

	tok, err = l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "a", tok.Text)

	tok, err = l.Next()
	assert.NoError(t, err)
	assert.Equal(t, STRING, tok.Type)
	assert.Equal(t, "b", tok.Text)
}

func TestPostalAddress(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarFile("../examples/postal.bnf")
	assert.NoError(t, err)

	ok := []string{
		"John Smith\n123 Main St\nSpringfield, MA 02139\n",
		"J. Doe Jr.\n42 Elm St\nBoston, NY 10001\n",
		"John Doe III\n42 Elm St Apt12\nBoston, MA 10001\n",
	}

	bad := []string{
		"John Smith 123 Main St\nSpringfield, MA 02139\n",
		"John Smith\n123 Main St\nSpringfield, MA\n",
		"John Smith\n123 Main St Apt\nSpringfield, MA 02139\n",
	}

	for _, s := range ok {
		m, err := g.Match(s)
		assert.True(t, m, s)
		assert.NoError(t, err)
	}

	for _, s := range bad {
		m, err := g.Match(s)
		assert.False(t, m, s)
		assert.Error(t, err)
	}
}

func TestCommentAtEOF(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarString(`
a ::= "a" // eof comment`)
	assert.NoError(t, err)

	m, err := g.Match("a")
	assert.True(t, m)
	assert.NoError(t, err)
}

func TestCommentOnlyLine(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarString(`
# this is a comment
; another one
a ::= "a"
`)
	assert.NoError(t, err)

	m, err := g.Match("a")
	assert.True(t, m)
	assert.NoError(t, err)
}

func TestCommentAfterAlternative(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarString(`
a ::= "a" | "b" // alternative c
`)
	assert.NoError(t, err)

	m, err := g.Match("a")
	assert.True(t, m)
	assert.NoError(t, err)

	m, err = g.Match("b")
	assert.True(t, m)
	assert.NoError(t, err)

	m, err = g.Match("c")
	assert.False(t, m)
	assert.Error(t, err)
}

func TestCommentInsideString(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarString(`
a ::= "//" | "#" | ";"
`)
	assert.NoError(t, err)

	for _, c := range []string{"//", "#", ";"} {
		m, err := g.Match(c)
		assert.True(t, m)
		assert.NoError(t, err)

	}
}

func TestCommentWithWhitespace(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarString(`
   a ::=   "a"     # comment
           | "b"   // another
           | "c"
`)
	assert.NoError(t, err)

	for _, c := range []string{"a", "b", "c"} {
		m, err := g.Match(c)
		assert.True(t, m)
		assert.NoError(t, err)

	}
}

func TestRecurrentParenthesis(t *testing.T) {
	t.Parallel()

	g, err := LoadGrammarString(`
		a ::= (((("1" | "2") | "3") | "4") | "c")
	`)
	assert.NoError(t, err)

	tests := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"2", true},
		{"3", true},
		{"4", true},
		{"c", true},
		{"5", false},
		{"", false},
		{"12", false},
		{"cc", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := g.Match(tt.input)
			if tt.want {
				assert.True(t, got)
				assert.NoError(t, err)
			} else {
				assert.False(t, got)
				assert.Error(t, err)
			}
		})
	}
}
