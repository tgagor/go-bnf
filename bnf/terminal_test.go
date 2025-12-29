package bnf_test

import (
	"bnf-test/bnf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminal(t *testing.T) {
	n := &bnf.Terminal{
		Value: "abc",
	}

	assert.Equal(t, []int{3}, n.Match("abcdef", 0)) // matching 3
	assert.Nil(t, n.Match("abcdef", 1))             // not matching
	assert.Nil(t, n.Match("ab", 0))                 // not matching
	assert.Nil(t, n.Match("", 0))                   // not matching
}
