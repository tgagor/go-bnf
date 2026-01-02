package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatStar(t *testing.T) {
	t.Parallel()

	// A ::= "a"*
	r := &Repeat{
		Node: &Terminal{"a"},
		Min:  0,
	}

	// Flow for "aaa":
	// 	i=0 -> result: [0]
	// 	i=1 -> [1]
	// 	i=2 -> [2]
	// 	i=3 -> [3]
	// 	no further change -> stop
	assert.Equal(t, []int{0, 1, 2, 3}, testMatch(r, "aaa", 0)) // matches
	assert.Equal(t, []int{0, 1}, testMatch(r, "abc", 0))       // 0 for optional, but then 1 for actual match
	assert.Equal(t, []int{0}, testMatch(r, "", 0))             // 0 for optional
}

func TestRepeatStarComplex(t *testing.T) {
	t.Parallel()

	// A ::= ("a" | "aa")*
	r := &Repeat{
		Node: &Choice{
			Options: []Node{
				&Terminal{"a"},
				&Terminal{"aa"},
			},
		},
		Min: 0,
	}

	// possible paths:
	// a a a 	-> 3
	// aa a 	-> 3
	// a aa 	-> 3
	// a 		-> 1
	// aa 		-> 2
	// duplicates are fine
	assert.Equal(t, []int{0, 1, 2, 2, 3, 3, 3}, testMatch(r, "aaa", 0)) // matches
	assert.Equal(t, []int{0, 1}, testMatch(r, "a", 0))                  // 0 for optional, but then 1 for actual match
	assert.Equal(t, []int{0}, testMatch(r, "", 0))                      // 0 for optional
}

func TestRepeatPlus(t *testing.T) {
	t.Parallel()

	// A ::= "a"+
	r := &Repeat{
		Node: &Terminal{"a"},
		Min:  1,
	}

	// Flow for "aaa":
	// 	i=0 -> result: [0]
	// 	i=1 -> [1]
	// 	i=2 -> [2]
	// 	i=3 -> [3]
	// 	no further change -> stop
	assert.Equal(t, []int{1, 2, 3}, testMatch(r, "aaa", 0)) // it stops at 2nd char (no match)
	assert.Equal(t, []int{1}, testMatch(r, "abc", 0))       // it stops at 2nd char (no match)
	assert.Nil(t, testMatch(r, "", 0))                      // it stops at 2nd char (no match)
}
