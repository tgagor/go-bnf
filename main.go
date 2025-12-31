package main

import (
	"go-bnf/cmd"
	"fmt"
	"os"
	"flag"
)

var BuildVersion string // Will be set dynamically at build time.
var appName string = "bnf"

var grammarFile string
var inputFile string
var help bool
var version bool
var lineByLine bool

func main() {

	// parse arguments
	flag.StringVar(&grammarFile, "g", "", "Path to the BNF grammar file")
	flag.StringVar(&inputFile, "i", "", "Path to the input file to verify against the grammar")
	flag.BoolVar(&lineByLine, "l", false, "Verify input line by line otherwise as a whole")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&version, "v", false, "Show version")
	flag.Parse()

	if help {
		fmt.Println("Usage: bnf -g <grammar-file> -i <input-file>")
		fmt.Println("   or: cat <input-file> | bnf -g <grammar-file>")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version {
		fmt.Printf("%s version: %s\n", appName, BuildVersion)
		os.Exit(0)
	}

	if grammarFile == "" {
		fmt.Println("Usage: bnf -g <grammar-file> -i <input-file>")
		fmt.Println("   or: cat <input-file> | bnf -g <grammar-file>")
		os.Exit(1)
	}

	cli := cmd.New(BuildVersion, appName, grammarFile, inputFile, lineByLine)
	if err := cli.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
