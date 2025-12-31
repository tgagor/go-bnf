package cmd

import (
	"bufio"
	"fmt"
	"go-bnf/bnf"
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

func loadByLine(file string) ([]string, error) {
	var f *os.File
	var err error
	if file != "" {
		f, err = os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	} else {
		// read from stdin
		f = os.Stdin
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func loadWhole(file string) (string, error) {
	var f *os.File
	var err error
	if file != "" {
		f, err = os.Open(file)
		if err != nil {
			return "", err
		}
		defer f.Close()
	} else {
		// read from stdin
		f = os.Stdin
	}

	content, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (cli *CLI) Run() error {
	fmt.Println("Parsing grammar file:", cli.GrammarFile)
	g, err := bnf.LoadGrammarFile(cli.GrammarFile)
	if err != nil {
		fmt.Println("Parsing error:", err)
		os.Exit(1)
	}
	fmt.Println("Grammar loaded.")

	fmt.Println("Loading input...")
	if !cli.LineByLine {
		whole, err := loadWhole(cli.VerifyFile)
		if err != nil {
			return err
		}
		fmt.Printf("Checking whole input...")
		match, err := g.Match(whole)
		if match {
			fmt.Println(" -> matched")
		} else {
			fmt.Printf("\n%s\n", err.(*bnf.ParseError).Pretty(whole))
		}
		return nil
	}

	lines, err := loadByLine(cli.VerifyFile)

	for _, l := range lines {
		fmt.Printf("Checking: %s", l)
		match, err := g.Match(l)
		if match {
			fmt.Println(" -> matched")
		} else {
			fmt.Printf("\n%s\n", err.(*bnf.ParseError).Pretty(l))
		}
	}

	return nil
}
