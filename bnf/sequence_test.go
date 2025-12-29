package bnf_test

import (
	"bnf-test/bnf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Parallel()

	seq := &bnf.Sequence{
		Elements: []bnf.Node{
			&bnf.Terminal{"a"},
			&bnf.Terminal{"b"},
			&bnf.Terminal{"c"},
		},
	}

	assert.Equal(t, []int{3}, seq.Match("abc", 0)) // matching 3
	assert.Nil(t, seq.Match("axc", 0))             // not matching b
	assert.Nil(t, seq.Match("ab", 0))              // no c
}

func TestSequenceMultiPath(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa") "b"
	choice := &bnf.Choice{
		Options: []bnf.Node{
			&bnf.Terminal{"a"},
			&bnf.Terminal{"aa"},
		},
	}
	seq := &bnf.Sequence{
		Elements: []bnf.Node{
			choice,
			&bnf.Terminal{"b"},
		},
	}

	assert.Equal(t, []int{2}, seq.Match("ab", 0))   // matches 2
	assert.Equal(t, []int{3}, seq.Match("aab", 0))  // matches 3
	assert.Equal(t, []int{3}, seq.Match("aabc", 0)) // matches only 3
	assert.Equal(t, []int{2}, seq.Match("abc", 0))  // matches only 2
	assert.Nil(t, seq.Match("acb", 0))              // no match
	assert.Nil(t, seq.Match("", 0))                 // no match
}
