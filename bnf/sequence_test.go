package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Parallel()

	seq := &Sequence{
		Elements: []Node{
			&Terminal{"a"},
			&Terminal{"b"},
			&Terminal{"c"},
		},
	}

	assert.Equal(t, []int{3}, match(seq, "abc", 0)) // matching 3
	assert.Nil(t, match(seq, "axc", 0))             // not matching b
	assert.Nil(t, match(seq, "ab", 0))              // no c
}

func TestSequenceMultiPath(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa") "b"
	choice := &Choice{
		Options: []Node{
			&Terminal{"a"},
			&Terminal{"aa"},
		},
	}
	seq := &Sequence{
		Elements: []Node{
			choice,
			&Terminal{"b"},
		},
	}

	assert.Equal(t, []int{2}, match(seq, "ab", 0))   // matches 2
	assert.Equal(t, []int{3}, match(seq, "aab", 0))  // matches 3
	assert.Equal(t, []int{3}, match(seq, "aabc", 0)) // matches only 3
	assert.Equal(t, []int{2}, match(seq, "abc", 0))  // matches only 2
	assert.Nil(t, match(seq, "acb", 0))              // no match
	assert.Nil(t, match(seq, "", 0))                 // no match
}
