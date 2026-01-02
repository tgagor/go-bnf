package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Parallel()

	seq := &sequence{
		Elements: []node{
			&terminal{"a"},
			&terminal{"b"},
			&terminal{"c"},
		},
	}

	assert.Equal(t, []int{3}, testMatch(seq, "abc", 0)) // matching 3
	assert.Nil(t, testMatch(seq, "axc", 0))             // not matching b
	assert.Nil(t, testMatch(seq, "ab", 0))              // no c
}

func TestSequenceMultiPath(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa") "b"
	choice := &choice{
		Options: []node{
			&terminal{"a"},
			&terminal{"aa"},
		},
	}
	seq := &sequence{
		Elements: []node{
			choice,
			&terminal{"b"},
		},
	}

	assert.Equal(t, []int{2}, testMatch(seq, "ab", 0))   // matches 2
	assert.Equal(t, []int{3}, testMatch(seq, "aab", 0))  // matches 3
	assert.Equal(t, []int{3}, testMatch(seq, "aabc", 0)) // matches only 3
	assert.Equal(t, []int{2}, testMatch(seq, "abc", 0))  // matches only 2
	assert.Nil(t, testMatch(seq, "acb", 0))              // no match
	assert.Nil(t, testMatch(seq, "", 0))                 // no match
}
