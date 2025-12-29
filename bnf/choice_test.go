package bnf_test

import (
	"bnf-test/bnf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoice(t *testing.T) {
	t.Parallel()

	// A ::= "a" | "ab"
	seq := &bnf.Choice{
		Options: []bnf.Node{
			&bnf.Terminal{"a"},
			&bnf.Terminal{"ab"},
		},
	}

	assert.Equal(t, []int{1, 2}, seq.Match("ab", 0)) // matching both
	assert.Equal(t, []int{1}, seq.Match("a", 0))     // matching 1
	assert.Nil(t, seq.Match("ab", 1))                // no match
	assert.Nil(t, seq.Match("b", 0))                 // not matching b
}

func TestChoiceDuplicate(t *testing.T) {
	t.Parallel()

	// A ::= "a" | "a"
	seq := &bnf.Choice{
		Options: []bnf.Node{
			&bnf.Terminal{"a"},
			&bnf.Terminal{"a"},
		},
	}

	assert.Equal(t, []int{1, 1}, seq.Match("ab", 0)) // matching same
	assert.Equal(t, []int{1, 1}, seq.Match("a", 0))  // matching same
	assert.Nil(t, seq.Match("b", 0))                 // not matching b
}

func TestChoiceMultiPath(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa") | "b"
	choice := &bnf.Choice{
		Options: []bnf.Node{
			&bnf.Terminal{"a"},
			&bnf.Terminal{"aa"},
		},
	}
	seq := &bnf.Choice{
		Options: []bnf.Node{
			choice,
			&bnf.Terminal{"b"},
		},
	}

	assert.Equal(t, []int{1}, seq.Match("ab", 0))      // matches 2
	assert.Equal(t, []int{1, 2}, seq.Match("aab", 0))  // matches 3
	assert.Equal(t, []int{1, 2}, seq.Match("aabc", 0)) // matches only 3
	assert.Equal(t, []int{1}, seq.Match("abc", 0))     // only a matches
	assert.Equal(t, []int{1}, seq.Match("acb", 0))     // only a matches
	assert.Nil(t, seq.Match("cba", 0))                 // only a matches
	assert.Nil(t, seq.Match("", 0))                    // no match
}
