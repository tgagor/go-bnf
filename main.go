package main

import (
	"flag"
	"fmt"
	"go-bnf/cmd"
	"os"
)

var BuildVersion string // Will be set dynamically at build time.
var appName string = "bnf"

var grammarFile string
var inputFile string
var help bool
var version bool
var lineByLine bool
var validateGrammar bool
var parseAsAST bool
var startRule string

func main() {

	// parse arguments
	flag.StringVar(&grammarFile, "g", "", "Path to the BNF grammar file")
	flag.StringVar(&inputFile, "i", "", "Path to the input file to verify against the grammar")
	flag.BoolVar(&lineByLine, "l", false, "Verify input line by line otherwise as a whole")
	flag.BoolVar(&validateGrammar, "v", false, "Only validate the grammar file")
	flag.BoolVar(&parseAsAST, "p", false, "Parse input and show AST")
	flag.StringVar(&startRule, "s", "", "Override the start rule for the grammar")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()

	if help {
		fmt.Println("Usage: bnf -g <grammar-file> -i <input-file> [options]")
		fmt.Println("   or: cat <input-file> | bnf -g <grammar-file> [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version {
		fmt.Printf("%s version: %s\n", appName, BuildVersion)
		os.Exit(0)
	}

	if grammarFile == "" {
		fmt.Println("Usage: bnf -g <grammar-file> -i <input-file> [options]")
		fmt.Println("   or: cat <input-file> | bnf -g <grammar-file> [options]")
		os.Exit(1)
	}

	cli := cmd.New(BuildVersion, appName, grammarFile, inputFile, lineByLine, validateGrammar, parseAsAST, startRule, nil)
	if err := cli.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
