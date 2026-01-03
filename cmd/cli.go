package cmd

import (
	"bufio"
	"fmt"
	"go-bnf/bnf"
	"io"
	"os"
	"strings"
)

type CLI struct {
	BuildVersion string
	AppName      string

	GrammarFile     string
	VerifyFile      string
	LineByLine      bool
	ValidateGrammar bool
	ParseAsAST      bool
	StartRule       string
}

func New(buildVersion, appName, grammarFile, verifyFile string, lineByLine, validateGrammar, parseAsAST bool, startRule string) *CLI {
	return &CLI{
		BuildVersion:    buildVersion,
		AppName:         appName,
		GrammarFile:     grammarFile,
		VerifyFile:      verifyFile,
		LineByLine:      lineByLine,
		ValidateGrammar: validateGrammar,
		ParseAsAST:      parseAsAST,
		StartRule:       startRule,
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

func (cli *CLI) Run() error {
	fmt.Println("Parsing grammar file:", cli.GrammarFile)
	g, err := bnf.LoadGrammarFile(cli.GrammarFile)
	if err != nil {
		fmt.Println("Parsing error:", err)
		os.Exit(1)
	}

	if cli.StartRule != "" {
		g.SetStart(cli.StartRule)
	}

	// Always validate grammar structure (undefined rules, etc.)
	if err := g.ValidateGrammar(); err != nil {
		fmt.Printf("Grammar validation error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Grammar loaded and validated.")

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

	fmt.Println("Loading input...")
	var tokens []string
	if cli.LineByLine {
		fmt.Println("Checking line by line...")
		tokens, err = loadFile(cli.VerifyFile, true)
	} else {
		fmt.Println("Checking whole input...")
		tokens, err = loadFile(cli.VerifyFile, false)
	}
	if err != nil {
		return err
	}

	for _, l := range tokens {
		fmt.Printf("Checking: %s", l)
		if cli.ParseAsAST {
			node, err := g.Parse(l)
			if err == nil {
				fmt.Println(" -> matched")
				fmt.Println("AST Tree:")
				fmt.Println(node.String())
			} else {
				cli.reportError(l, err)
			}
		} else {
			match, err := g.Match(l)
			if match {
				fmt.Println(" -> matched")
			} else {
				cli.reportError(l, err)
			}
		}
	}

	return nil
}

func (cli *CLI) reportError(input string, err error) {
	if err == nil {
		fmt.Println(" -> not matched (no error details)")
	} else if pe, ok := err.(*bnf.ParseError); ok {
		if pe == nil {
			fmt.Println(" -> not matched (nil ParseError)")
		} else {
			fmt.Printf("\n%s\n\n", pe.Pretty(input))
		}
	} else {
		fmt.Printf(" -> error: %v\n", err)
	}
}
