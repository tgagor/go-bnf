package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Paths don't need to exist for New
	cli := New("1.0.0", "test-app", "../tests/simple.bnf", "../tests/input_match.txt", true, false, false, "")

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

	cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "")

	err := cli.Run()
	assert.NoError(t, err)
}

func TestRun_LineByLine(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "tests", "simple.bnf")
	inputFile := filepath.Join("..", "tests", "input_multiline.txt")

	cli := New("0.0.1", "test", grammarFile, inputFile, true, false, false, "")

	err := cli.Run()
	assert.NoError(t, err)
}

func TestRun_Mismatch(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "tests", "simple.bnf")
	inputFile := filepath.Join("..", "tests", "input_mismatch.txt")

	cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "")

	err := cli.Run()
	// Current implementation: Run() returns nil even on mismatch,
	// but prints error to stdout. We just check strictly that it doesn't crash/error.
	assert.NoError(t, err)
}

func TestRun_Postal(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "postal.bnf")

	for i := range 4 {
		inputFile := filepath.Join("..", "examples", fmt.Sprintf("postal%d.txt", i+1))

		cli := New("0.0.1", "test", grammarFile, inputFile, false, false, false, "")

		err := cli.Run()
		assert.NoError(t, err)
	}
}

func TestRun_Numbers(t *testing.T) {
	t.Parallel()
	grammarFile := filepath.Join("..", "examples", "numbers.bnf")
	inputFile := filepath.Join("..", "examples", "numbers.test")

	cli := New("0.0.1", "test", grammarFile, inputFile, true, false, false, "")

	err := cli.Run()
	assert.NoError(t, err)
}
