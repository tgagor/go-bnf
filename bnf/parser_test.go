package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBNF_MultilineOr(t *testing.T) {
	t.Parallel()
	grammar := `
<DIGIT> ::= "0"
          | "1"
          | "2"
`

	g := LoadGrammarString(grammar)
	assert.NotNil(t, g)

	ok, err := g.MatchFrom("DIGIT", "0")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("DIGIT", "1")
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = g.MatchFrom("DIGIT", "3")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.MatchFrom("DIGIT", "10")
	assert.False(t, ok)
	assert.Error(t, err)

	ok, err = g.MatchFrom("DIGIT", "")
	assert.False(t, ok)
	assert.Error(t, err)
}
