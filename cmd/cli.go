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

	GrammarFile string
	VerifyFile  string
	LineByLine  bool
}

func New(buildVersion, appName, grammarFile, verifyFile string, lineByLine bool) *CLI {
	return &CLI{
		BuildVersion: buildVersion,
		AppName:      appName,
		GrammarFile:  grammarFile,
		VerifyFile:   verifyFile,
		LineByLine:   lineByLine,
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
	fmt.Println("Grammar loaded.")

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
		match, err := g.Match(l)
		if match {
			fmt.Println(" -> matched")
		} else {
			fmt.Printf("\n%s\n\n", err.(*bnf.ParseError).Pretty(l))
		}
	}

	return nil
}
