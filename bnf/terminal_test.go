package bnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminal(t *testing.T) {
	n := &Terminal{
		Value: "abc",
	}

	assert.Equal(t, []int{3}, testMatch(n, "abcdef", 0)) // matching 3
	assert.Nil(t, testMatch(n, "abcdef", 1))             // not matching
	assert.Nil(t, testMatch(n, "ab", 0))                 // not matching
	assert.Nil(t, testMatch(n, "", 0))                   // not matching
}
