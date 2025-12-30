package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoice(t *testing.T) {
	t.Parallel()

	// A ::= "a" | "ab"
	c := &Choice{
		Options: []Node{
			&Terminal{"a"},
			&Terminal{"ab"},
		},
	}

	assert.Equal(t, []int{1, 2}, match(c, "ab", 0)) // matching both
	assert.Equal(t, []int{1}, match(c, "a", 0))     // matching 1
	assert.Nil(t, match(c, "ab", 1))                // no match
	assert.Nil(t, match(c, "b", 0))                 // not matching b
}

func TestChoiceDuplicate(t *testing.T) {
	t.Parallel()

	// A ::= "a" | "a"
	c := &Choice{
		Options: []Node{
			&Terminal{"a"},
			&Terminal{"a"},
		},
	}

	assert.Equal(t, []int{1, 1}, match(c, "ab", 0)) // matching same
	assert.Equal(t, []int{1, 1}, match(c, "a", 0))  // matching same
	assert.Nil(t, match(c, "b", 0))                 // not matching b
}

func TestChoiceMultiPath(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa") | "b"
	inner := &Choice{
		Options: []Node{
			&Terminal{"a"},
			&Terminal{"aa"},
		},
	}
	c := &Choice{
		Options: []Node{
			inner,
			&Terminal{"b"},
		},
	}

	assert.Equal(t, []int{1}, match(c, "ab", 0))      // matches 2
	assert.Equal(t, []int{1, 2}, match(c, "aab", 0))  // matches 3
	assert.Equal(t, []int{1, 2}, match(c, "aabc", 0)) // matches only 3
	assert.Equal(t, []int{1}, match(c, "abc", 0))     // only a matches
	assert.Equal(t, []int{1}, match(c, "acb", 0))     // only a matches
	assert.Nil(t, match(c, "cba", 0))                 // only a matches
	assert.Nil(t, match(c, "", 0))                    // no match
}
