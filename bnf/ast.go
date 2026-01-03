package bnf

import (
	"fmt"
	"regexp"
)

// GrammarAST represents the raw AST of a BNF grammar before it is built into a Grammar object.
type GrammarAST struct {
	Rules []*RuleAST
}

// RuleAST represents a single rule in the GrammarAST.
type RuleAST struct {
	Name string
	Expr ExprAST
}

// ExprAST is a marker interface for all AST expression nodes.
type ExprAST any

type (
	// ChoiceAST represents a choice between multiple alternatives (e.g., "a" | "b").
	ChoiceAST struct {
		Options []ExprAST
	}

	// SeqAST represents a sequence of expressions (e.g., "a" "b").
	SeqAST struct {
		Elements []ExprAST
	}

	// RepeatAST represents a repetition of an expression (e.g., "a"*, "a"+, "a"?).
	RepeatAST struct {
		Node ExprAST
		Min  int
		Max  int // -1 = infinity
	}

	// IdentAST represents a reference to another rule.
	IdentAST struct {
		Name string
	}

	// StringAST represents a literal string terminal.
	StringAST struct {
		Value string
	}

	// RegexAST represents a regular expression pattern terminal.
	RegexAST struct {
		Pattern string
	}
)

// BuildGrammar converts a raw GrammarAST into a functional Grammar object with linked rules.
func BuildGrammar(ast *GrammarAST) (*Grammar, error) {
	rules := map[string]*Rule{}

	if len(ast.Rules) == 0 {
		return nil, fmt.Errorf("empty grammar")
	}

	// 1. create rules
	for _, r := range ast.Rules {
		rules[r.Name] = &Rule{Name: r.Name}
	}

	// 2. build expr
	for _, r := range ast.Rules {
		expr, err := buildNode(r.Expr, rules)
		if err != nil {
			return nil, err
		}
		rules[r.Name].Expr = expr
	}

	return &Grammar{
		Start: ast.Rules[0].Name,
		Rules: rules,
	}, nil
}

func buildNode(e ExprAST, rules map[string]*Rule) (node, error) {
	switch t := e.(type) {
	case *StringAST:
		return &terminal{Value: t.Value}, nil

	case *RegexAST:
		re, err := regexp.Compile(t.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern %q: %w", t.Pattern, err)
		}
		return &Regex{Re: re}, nil

	case *IdentAST:
		return &nonTerminal{Name: t.Name, Rule: rules[t.Name]}, nil

	case *SeqAST:
		var elems []node
		for _, e := range t.Elements {
			n, err := buildNode(e, rules)
			if err != nil {
				return nil, err
			}
			elems = append(elems, n)
		}
		return &sequence{Elements: elems}, nil

	case *ChoiceAST:
		var opts []node
		for _, o := range t.Options {
			n, err := buildNode(o, rules)
			if err != nil {
				return nil, err
			}
			opts = append(opts, n)
		}
		return &choice{Options: opts}, nil

	case *RepeatAST:
		n, err := buildNode(t.Node, rules)
		if err != nil {
			return nil, err
		}
		return &repeat{
			Node: n,
			Min:  t.Min,
		}, nil
	}
	return nil, fmt.Errorf("unknown AST type: %T", e)
}
