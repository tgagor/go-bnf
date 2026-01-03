package bnf_test

import (
	"fmt"
	"github.com/tgagor/go-bnf/bnf"
	"log"
)

func ExampleGrammar_Validate() {
	grammar, err := bnf.LoadGrammarString(`<number> ::= "0" | "1"`)
	if err != nil {
		log.Fatal(err)
	}
	v1, _ := grammar.Validate("1")
	fmt.Println(v1)
	v2, _ := grammar.Validate("3")
	fmt.Println(v2)
	// Output:
	// true
	// false
}

func ExampleGrammar_Parse() {
	grammar, err := bnf.LoadGrammarString(`
<S> ::= "a" <S> | "b"
`)
	if err != nil {
		log.Fatal(err)
	}

	tree, err := grammar.Parse("aab")
	if err != nil {
		log.Fatal(err)
	}

	// Tree string representation: (Type Child...)
	// Since <S> is recursive, we expect nesting.
	// Structure roughly depend on how repeated S is nested.
	// "a" <S> -> "a" ("a" "b")
	// (S "a" (S "a" (S "b")))

	fmt.Println(tree.Type)
	// Output:
	// S
}
