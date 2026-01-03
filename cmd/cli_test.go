package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Paths don't need to exist for New
	cli := New("1.0.0", "test-app", "../tests/simple.bnf", "../tests/input_match.txt", true, false, false, "", nil)

	assert.NotNil(t, cli)
	assert.Equal(t, "1.0.0", cli.BuildVersion)
	assert.Equal(t, "test-app", cli.AppName)
	assert.Equal(t, "../tests/simple.bnf", cli.GrammarFile)
	assert.Equal(t, "../tests/input_match.txt", cli.VerifyFile)
	assert.True(t, cli.LineByLine)
}

func TestRun_Success(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "tests", "simple.bnf")
	inputFile := filepath.Join("..", "tests", "input_match.txt")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "", out)

	err := cli.Run()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "-> matched")
}

func TestRun_LineByLine(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "tests", "simple.bnf")
	inputFile := filepath.Join("..", "tests", "input_multiline.txt")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, inputFile, true, false, false, "", out)

	err := cli.Run()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "matched")
}

func TestRun_Mismatch(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "tests", "simple.bnf")
	inputFile := filepath.Join("..", "tests", "input_mismatch.txt")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "", out)

	err := cli.Run()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "Parse error")
}

func TestRun_Postal(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "postal.bnf")

	for i := range 4 {
		inputFile := filepath.Join("..", "examples", fmt.Sprintf("postal%d.txt", i+1))
		out := &bytes.Buffer{}
		cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "", out)

		err := cli.Run()
		assert.NoError(t, err)
		assert.Contains(t, out.String(), "-> matched")
	}
}

func TestRun_Numbers(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "numbers.bnf")
	inputFile := filepath.Join("..", "examples", "numbers.test")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, inputFile, true, false, false, "", out)

	err := cli.Run()
	assert.NoError(t, err)
}

func TestRun_ValidateOnly(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "numbers.bnf")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, "", false, true, false, "", out)

	err := cli.Run()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "Grammar loaded and validated.")
}

func TestRun_ParseAST(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "numbers.bnf")
	inputFile := filepath.Join("..", "examples", "numbers.test")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, inputFile, true, false, true, "", out)

	err := cli.Run()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "AST Tree:")
}

func TestRun_StartRuleOverride(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "numbers.bnf")
	inputFile := filepath.Join("..", "tests", "digit.txt")
	out := &bytes.Buffer{}

	cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "digit", out)

	err := cli.Run()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "-> matched")
}

func TestRun_InvalidGrammar(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "tests", "invalid.bnf")

	cli := New("0.0.1", "test", grammarFile, "", false, true, false, "", nil)

	err := cli.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined rule")
}
