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
		g, err := bnf.LoadGrammarFile(f)
		if err != nil {
			fmt.Println("Parsing error:", err)
			os.Exit(1)
		}

		for _, s := range g.Rules {
			fmt.Printf("Parsing line: %s = ", s.Name)
			// for _, p := range g.Rules[s.Name] {
			// 	fmt.Printf("%+v ", p)
			// }
			fmt.Println()
		}
	}
}
