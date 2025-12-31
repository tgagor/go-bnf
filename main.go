package main

import (
	"go-bnf/bnf"
	"bufio"
	"fmt"
	"os"
	"strings"
)

var BuildVersion string // Will be set dynamically at build time.
var appName string = "bnf"

func loadExamples(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func main() {
	args := os.Args[1:]
	grammarFile := args[0]
	examplesFile := args[1]

	fmt.Println("Parsing:", grammarFile)
	g, err := bnf.LoadGrammarFile(grammarFile)
	if err != nil {
		fmt.Println("Parsing error:", err)
		os.Exit(1)
	}
	fmt.Println("Grammar loaded.")

	fmt.Println("Loading examples...")
	examples, err := loadExamples(examplesFile)

	for _, e := range examples {
		fmt.Printf("Checking: %s", e)
		if g.Match(e) {
			fmt.Println(" -> matched")
		} else {
			fmt.Println(" -> not matched")
		}
	}

}
