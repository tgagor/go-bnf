package bnf_test

import (
	"bnf-test/bnf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {

	// // Arrange
	// input := []string{
	// 	"<digit>",
	// 	"<digit>",
	// 	"<digit>",
	// 	"<digit>",
	// }
	// expected := []string{
	// 	"0",
	// 	"1",
	// 	"2",
	// 	"9",
	// }

	b, err := bnf.FromFile("../examples/numbers.bnf")

	assert.NotNil(t, b)
	assert.Nil(t, err)

	// Assert
	// for i, input := range input {
	// 	assert.Equal(t, expected[i])
	// }
}

func TestParser(t *testing.T) {
	b, _ := bnf.FromFile("../examples/numbers.bnf")
	s := b.GetSymbols()
	assert.NotNil(t, b)
	assert.Equal(t, len(s), 6)

	assert.Contains(t, s, "<digit>")
	assert.Contains(t, s, "<number>")
	assert.Contains(t, s, "<hex>")
	assert.Contains(t, s, "<hexnumber>")

	assert.NotNil(t, b.GetSymbol("<digit>"))
	assert.Nil(t, b.GetSymbol("digit"))
	assert.Contains(t, b.GetSymbol("<digit>").Patterns, "\"1\"")
	assert.Equal(t, len(b.GetSymbol("<digit>").Patterns), 10)
}
