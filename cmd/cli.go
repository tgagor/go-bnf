package cmd

import (
	"bufio"
	"fmt"
	"github.com/tgagor/go-bnf/bnf"
	"io"
	"os"
	"strings"
)

// CLI handles the command-line interface execution flow, including grammar loading and input verification.
type CLI struct {
	BuildVersion string
	AppName      string

	GrammarFile     string
	VerifyFile      string
	LineByLine      bool
	ValidateGrammar bool
	ParseAsAST      bool
	StartRule       string
	Output          io.Writer
}

// New creates a new CLI instance with the specified configuration.
func New(buildVersion, appName, grammarFile, verifyFile string, lineByLine, validateGrammar, parseAsAST bool, startRule string, output io.Writer) *CLI {
	if output == nil {
		output = os.Stdout
	}
	return &CLI{
		BuildVersion:    buildVersion,
		AppName:         appName,
		GrammarFile:     grammarFile,
		VerifyFile:      verifyFile,
		LineByLine:      lineByLine,
		ValidateGrammar: validateGrammar,
		ParseAsAST:      parseAsAST,
		StartRule:       startRule,
		Output:          output,
	}
}

func loadFile(file string, lineByLine bool) ([]string, error) {
	var r io.ReadCloser
	var err error
	if file == "" { // no file -> read from stdin
		r, err = io.NopCloser(os.Stdin), nil
	} else {
		r, err = os.Open(file)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if lineByLine {
		var lines []string
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lines = append(lines, strings.TrimSpace(scanner.Text()))
		}
		return lines, scanner.Err()
	} else {
		content, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return []string{string(content)}, nil
	}
}

// Run executes the CLI logic: loads the grammar, validates it, and matches/parses the provided input.
func (cli *CLI) Run() error {
	fmt.Fprintln(cli.Output, "Parsing grammar file:", cli.GrammarFile)
	g, err := bnf.LoadGrammarFile(cli.GrammarFile)
	if err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	if cli.StartRule != "" {
		g.SetStart(cli.StartRule)
	}

	// Always validate grammar structure (undefined rules, etc.)
	if err := g.ValidateGrammar(); err != nil {
		return fmt.Errorf("grammar validation error: %w", err)
	}
	fmt.Fprintln(cli.Output, "Grammar loaded and validated.")

	if cli.ValidateGrammar {
		// If only validation was requested, we are done
		return nil
	}

	if cli.VerifyFile == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat stdin: %w", err)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no input provided: use -i <file> or pipe data to stdin")
		}
	}

	fmt.Fprintln(cli.Output, "Loading input...")
	var tokens []string
	if cli.LineByLine {
		fmt.Fprintln(cli.Output, "Checking line by line...")
		tokens, err = loadFile(cli.VerifyFile, true)
	} else {
		fmt.Fprintln(cli.Output, "Checking whole input...")
		tokens, err = loadFile(cli.VerifyFile, false)
	}
	if err != nil {
		return err
	}

	for _, l := range tokens {
		fmt.Fprintf(cli.Output, "Checking: %s", l)
		if cli.ParseAsAST {
			node, err := g.Parse(l)
			if err == nil {
				fmt.Fprintln(cli.Output, " -> matched")
				fmt.Fprintln(cli.Output, "AST Tree:")
				fmt.Fprintln(cli.Output, node.String())
			} else {
				cli.reportError(l, err)
			}
		} else {
			match, err := g.Match(l)
			if match {
				fmt.Fprintln(cli.Output, " -> matched")
			} else {
				cli.reportError(l, err)
			}
		}
	}

	return nil
}

func (cli *CLI) reportError(input string, err error) {
	if err == nil {
		fmt.Fprintln(cli.Output, " -> not matched (no error details)")
	} else if pe, ok := err.(*bnf.ParseError); ok {
		if pe == nil {
			fmt.Fprintln(cli.Output, " -> not matched (nil ParseError)")
		} else {
			fmt.Fprintf(cli.Output, "\n%s\n\n", pe.Pretty(input))
		}
	} else {
		fmt.Fprintf(cli.Output, " -> error: %v\n", err)
	}
}
