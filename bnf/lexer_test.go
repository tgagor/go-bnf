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
			break
		}
		if tok.Type == EOF {
			break
		}
	}
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
