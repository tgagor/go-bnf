package main

import (
	"bnf-test/bnf"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, f := range args {
		fmt.Println("Parsing:", f)
		b := bnf.New(f)
		err := b.Load()
		if err != nil {
			fmt.Println("Parsing error:", err)
			os.Exit(1)
		}

		for _, s := range b.GetSymbols() {
			fmt.Printf("Parsing line: %s = ", s)
			for _, p := range b.GetSymbol(s).Patterns {
				fmt.Printf("%+v ", p)
			}
			fmt.Println()
		}
	}
}
