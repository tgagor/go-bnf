package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	File        string
	Expressions map[string]*Expression
	// Expressions map[string]*Expression
}

type Expression struct {
	Name     string
	Patterns []*Pattern
	// Expressions []any
}

type Pattern struct {
	Val any
}

func New(path string) *Parser {
	return &Parser{
		File:        path,
		Expressions: map[string]*Expression{},
	}
}

func (p *Parser) Load() error {
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
	fmt.Printf("Expressions: %+v\n", p.Expressions)

	return nil
}

func (p *Parser) ParseLine(line string) {
	split := strings.Split(line, "::=")

	// clean name to raw string
	name := strings.TrimSpace(split[0])

	// rule.Name = strings.TrimLeft(rule.Name, "<")
	// rule.Name = strings.TrimRight(rule.Name, ">")
	// rule.Name = strings.TrimSpace(rule.Name)

	patterns := p.ParseExpression(name, split[1])

	fmt.Printf("Parsing line: %s = ", name)
	for _, p := range patterns {
		fmt.Printf("%+v ", p)
	}
	fmt.Println()

	p.Expressions[name] = &Expression{
		Name:     name,
		Patterns: patterns,
	}
}

func (p *Parser) ParseExpression(name string, line string) (patterns []*Pattern) {
	patternSets := strings.Split(line, "|")
	patterns = []*Pattern{}
	for _, patternStr := range patternSets {
		// clean cases
		patternStr = strings.TrimSpace(patternStr)
		symbols := strings.Split(patternStr, " ")
		for _, symbol := range symbols {
			// terminal, literal expression
			if strings.HasPrefix(symbol, "\"") && strings.HasSuffix(symbol, "\"") {
				patterns = append(patterns, &Pattern{
					Val: symbol,
				})
			} else if strings.HasPrefix(symbol, "<") && strings.HasSuffix(symbol, ">") {
				// non-terminal, pointing to another expression
				patterns = append(patterns, &Pattern{
					Val: symbol,
				})
			} else {
				fmt.Println("Unknown expression:", symbol)
			}
		}
	}

	return patterns
}
