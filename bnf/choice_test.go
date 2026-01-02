package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoice(t *testing.T) {
	t.Parallel()

	// A ::= "a" | "ab"
	c := &choice{
		Options: []node{
			&terminal{"a"},
			&terminal{"ab"},
		},
	}

	assert.Equal(t, []int{1, 2}, testMatch(c, "ab", 0)) // matching both
	assert.Equal(t, []int{1}, testMatch(c, "a", 0))     // matching 1
	assert.Nil(t, testMatch(c, "ab", 1))                // no match
	assert.Nil(t, testMatch(c, "b", 0))                 // not matching b
}

func TestChoiceDuplicate(t *testing.T) {
	t.Parallel()

	// A ::= "a" | "a"
	c := &choice{
		Options: []node{
			&terminal{"a"},
			&terminal{"a"},
		},
	}

	assert.Equal(t, []int{1, 1}, testMatch(c, "ab", 0)) // matching same
	assert.Equal(t, []int{1, 1}, testMatch(c, "a", 0))  // matching same
	assert.Nil(t, testMatch(c, "b", 0))                 // not matching b
}

func TestChoiceMultiPath(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa") | "b"
	inner := &choice{
		Options: []node{
			&terminal{"a"},
			&terminal{"aa"},
		},
	}
	c := &choice{
		Options: []node{
			inner,
			&terminal{"b"},
		},
	}

	assert.Equal(t, []int{1}, testMatch(c, "ab", 0))      // matches 2
	assert.Equal(t, []int{1, 2}, testMatch(c, "aab", 0))  // matches 3
	assert.Equal(t, []int{1, 2}, testMatch(c, "aabc", 0)) // matches only 3
	assert.Equal(t, []int{1}, testMatch(c, "abc", 0))     // only a matches
	assert.Equal(t, []int{1}, testMatch(c, "acb", 0))     // only a matches
	assert.Nil(t, testMatch(c, "cba", 0))                 // only a matches
	assert.Nil(t, testMatch(c, "", 0))                    // no match
}
