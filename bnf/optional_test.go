package bnf

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionalTerminal(t *testing.T) {
	g := &optional{
		Node: &terminal{Value: "a"},
	}

	assert.Equal(t, []int{0, 1}, testMatch(g, "a", 0))
	assert.Equal(t, []int{0}, testMatch(g, "", 0))
}

func TestOptionalSequence(t *testing.T) {
	// "a"? "b"
	seq := &sequence{
		Elements: []node{
			&optional{Node: &terminal{Value: "a"}},
			&terminal{Value: "b"},
		},
	}

	assert.Equal(t, []int{2}, testMatch(seq, "ab", 0)) // matched optional "a" and "b"
	assert.Equal(t, []int{1}, testMatch(seq, "b", 0))  // skipped optional "a", matched "b"
	assert.Nil(t, testMatch(seq, "a", 0))              // no "b" to match
}

func TestOptionalPlus_Terminal(t *testing.T) {
	// a?+ should work as (a?)+
	// not as a(+?)
	node := &repeat{
		Node: &optional{
			Node: &terminal{Value: "a"},
		},
		Min: 1,
	}

	// Helper to extract ends
	ends := func(input string) []int {
		ctx := NewContext(input)
		res, _ := ctx.Match(node, 0)
		var out []int
		for _, r := range res {
			out = append(out, r.End)
		}
		return out
	}

	assert.True(t, slices.Contains(ends("a"), 1))
	assert.True(t, slices.Contains(ends("aa"), 2))
	assert.False(t, slices.Contains(ends(""), 0))
}
