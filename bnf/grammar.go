package bnf

import (
	"fmt"
)

// Rule represents a single BNF rule with its name and expression.
type Rule struct {
	Name string
	Expr node
}

// Grammar represents a complete set of BNF rules and an optional start rule.
type Grammar struct {
	Rules map[string]*Rule
	Start string
}

// Resolve recursively links all non-terminal nodes in the grammar to their rule definitions.
func (g *Grammar) Resolve() error {
	for _, rule := range g.Rules {
		if err := resolveNode(rule.Expr, g.Rules); err != nil {
			return err
		}
	}

	return nil
}

func resolveNode(n node, rules map[string]*Rule) error {
	switch t := n.(type) {
	case *nonTerminal:
		rule, ok := rules[t.Name]
		if !ok {
			return fmt.Errorf("undefined rule: %s", t.Name)
		}
		t.Rule = rule
	case *sequence:
		for _, e := range t.Elements {
			if err := resolveNode(e, rules); err != nil {
				return err
			}
		}
	case *choice:
		for _, o := range t.Options {
			if err := resolveNode(o, rules); err != nil {
				return err
			}
		}
	case *repeat:
		return resolveNode(t.Node, rules)
	case *optional:
		return resolveNode(t.Node, rules)
	}
	return nil
}

// ValidateGrammar checks if the grammar is self-consistent (all referenced rules exist and a start rule is defined).
func (g *Grammar) ValidateGrammar() error {
	if g.Start == "" {
		return fmt.Errorf("no start rule defined")
	}
	if _, ok := g.Rules[g.Start]; !ok {
		return fmt.Errorf("start rule %q is not defined", g.Start)
	}
	return g.Resolve()
}

// SetStart sets the entry rule for matching and parsing.
func (g *Grammar) SetStart(name string) {
	g.Start = name
}

// Validate checks if the input matches the grammar start rule.
// It is an alias for Match (or Match behaves as Validate).
func (g *Grammar) Validate(input string) (bool, error) {
	return g.Match(input)
}

// Match checks if the entire input matches the grammar starting from the defined start rule.
func (g *Grammar) Match(input string) (bool, error) {
	if g.Start == "" {
		return false, fmt.Errorf("start rule not defined")
	}
	return g.MatchFrom(g.Start, input)
}

// MatchFrom checks if the entire input matches the grammar starting from the specified rule.
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

// MatchPrefix checks if any prefix of the input matches the grammar start rule.
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
