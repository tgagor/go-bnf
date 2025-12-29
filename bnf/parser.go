package bnf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type BNF struct {
	File    string
	Symbols map[string]*Expression
	// Expressions map[string]*Expression
}

type Pattern struct {
	Val any
}

func New(path string) *BNF {
	return &BNF{
		File:    path,
		Symbols: map[string]*Expression{},
	}
}

func (p *BNF) GetSymbols() []string {
	keys := make([]string, len(p.Symbols))
	i := 0
	for k := range p.Symbols {
		keys[i] = k
		i++
	}
	return keys
}

func (p *BNF) GetSymbol(name string) *Expression {
	return p.Symbols[name]
}

func (p *BNF) Load() error {
	file, err := os.Open(p.File)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewScanner(file)
	for reader.Scan() {
		line := reader.Text()
		p.ParseLine(line)
		// p.Expressions[rule.Name] = &rule
	}
	fmt.Printf("Expressions: %+v\n", p.Symbols)

	return nil
}

func FromFile(path string) (*BNF, error) {
	b := New(path)
	e := b.Load()
	return b, e
}

func (p *BNF) ParseLine(line string) {
	split := strings.Split(line, "::=")

	name := strings.TrimSpace(split[0])
	patterns := p.ParseExpression(name, split[1])

	p.Symbols[name] = &Expression{
		Name:     name,
		Patterns: patterns,
	}
}

func (p *BNF) ParseExpression(name string, line string) (patterns []string) {
	patternSets := strings.Split(line, "|")
	patterns = []string{}
	for _, patternStr := range patternSets {
		// clean cases
		patternStr = strings.TrimSpace(patternStr)
		for _, symbol := range strings.Split(patternStr, " ") {
			// terminal, literal expression
			if strings.HasPrefix(symbol, "\"") && strings.HasSuffix(symbol, "\"") {
				// patterns = append(patterns, &Pattern{
				// 	Val: symbol,
				// })
				patterns = append(patterns, symbol)
			} else if strings.HasPrefix(symbol, "<") && strings.HasSuffix(symbol, ">") {
				// non-terminal, pointing to another expression
				// patterns = append(patterns, &Pattern{
				// 	Val: symbol,
				// })
				patterns = append(patterns, symbol)
			} else {
				fmt.Println("Unknown expression:", symbol)
				os.Exit(2)
			}
		}
	}

	return patterns
}
