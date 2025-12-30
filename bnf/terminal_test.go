package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminal(t *testing.T) {
	n := &Terminal{
		Value: "abc",
	}

	assert.Equal(t, []int{3}, match(n, "abcdef", 0)) // matching 3
	assert.Nil(t, match(n, "abcdef", 1))             // not matching
	assert.Nil(t, match(n, "ab", 0))                 // not matching
	assert.Nil(t, match(n, "", 0))                   // not matching
}
