package bnf

import "fmt"

// raw AST, without links yet
type GrammarAST struct {
	Rules []*RuleAST
}

type RuleAST struct {
	Name string
	Expr ExprAST
}

type ExprAST any

type (
	ChoiceAST struct {
		Options []ExprAST
	}

	SeqAST struct {
		Elements []ExprAST
	}

	RepeatAST struct {
		Node ExprAST
		Min  int
		Max  int // -1 = infinity
	}

	IdentAST struct {
		Name string
	}

	StringAST struct {
		Value string
	}
)

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
