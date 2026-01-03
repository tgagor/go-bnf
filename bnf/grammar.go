package bnf

import (
	"fmt"
)

type Rule struct {
	Name string
	Expr node
}

type Grammar struct {
	Rules map[string]*Rule
	Start string
}

func (g *Grammar) Resolve() error {
	for _, rule := range g.Rules {
		resolveNode(rule.Expr, g.Rules)
	}

	return nil
}

func resolveNode(n node, rules map[string]*Rule) {
	switch t := n.(type) {
	case *nonTerminal:
		t.Rule = rules[t.Name]
	case *sequence:
		for _, e := range t.Elements {
			resolveNode(e, rules)
		}
	case *choice:
		for _, o := range t.Options {
			resolveNode(o, rules)
		}
	case *repeat:
		resolveNode(t.Node, rules)
		// case *Optional:
		// 	resolveNode(t.Node, rules)
	}
}

func (g *Grammar) SetStart(name string) {
	g.Start = name
}

// Validate checks if the input matches the grammar start rule.
// It is an alias for Match (or Match behaves as Validate).
func (g *Grammar) Validate(input string) (bool, error) {
	return g.Match(input)
}

func (g *Grammar) Match(input string) (bool, error) {
	if g.Start == "" {
		return false, fmt.Errorf("start rule not defined")
	}
	return g.MatchFrom(g.Start, input)
}

func (g *Grammar) MatchFrom(start string, input string) (bool, error) {
	rule, ok := g.Rules[start]
	if !ok {
		return false, fmt.Errorf("unknown start rule: %s", start)
	}

	ctx := NewContext(input)
	matches, err := ctx.Match(rule.Expr, 0)
	if err != nil {
		return false, err
	}

	for _, m := range matches {
		if m.End == len(input) {
			return true, nil
		}
	}

	return false, ctx.error
}

func (g *Grammar) MatchPrefix(input string) bool {
	start := g.Rules[g.Start]
	ctx := NewContext(input)
	matches, err := ctx.Match(start.Expr, 0)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

// Parse parses the input and returns the AST.
func (g *Grammar) Parse(input string) (*ASTNode, error) {
	if g.Start == "" {
		return nil, fmt.Errorf("start rule not defined")
	}
	rule, ok := g.Rules[g.Start]
	if !ok {
		return nil, fmt.Errorf("unknown start rule: %s", g.Start)
	}

	ctx := NewContext(input)
	matches, err := ctx.Match(rule.Expr, 0)
	if err != nil {
		return nil, err
	}

	for _, m := range matches {
		if m.End == len(input) {
			// Success!
			// If result has multiple nodes (e.g. sequence at top level), wrap them.
			// But since we matched a Rule (which is a NonTerminal implied or explicit?),
			// checking how MatchFrom works: it matches `rule.Expr`.
			// `rule.Expr` for "Start" rule.
			// Currently `MatchFrom` accesses `rule.Expr` directly.
			// `nonTerminal.match` wraps result in `ASTNode` with Type=Name.
			// But here we invoke `rule.Expr` directly which might be a `sequence` or `choice`.
			// So we miss the top-level "Start Rule" wrapper node if we don't wrap it manually.

			if len(m.Nodes) == 1 && m.Nodes[0].Type == g.Start {
				return m.Nodes[0], nil
			}
			// Wrap it
			return &ASTNode{
				Type:     g.Start,
				Children: m.Nodes,
			}, nil
		}
	}

	return nil, ctx.error
}
