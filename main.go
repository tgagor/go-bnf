package main

import (
	"bnf-test/parser"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, f := range args {
		fmt.Println("Parsing:", f)
		p := parser.New(f)
		err := p.Load()
		if err != nil {
			fmt.Println("Parsing error:", err)
			os.Exit(1)
		}
	}
}
