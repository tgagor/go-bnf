package bnf

import "slices"

type Node interface {
	Match(input string, pos int) []int
}

type Rule struct {
	Name string
	Expr Node
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

func resolveNode(n Node, rules map[string]*Rule) {
	switch t := n.(type) {
	case *NonTerminal:
		t.Rule = rules[t.Name]
	case *Sequence:
		for _, e := range t.Elements {
			resolveNode(e, rules)
		}
	case *Choice:
		for _, o := range t.Options {
			resolveNode(o, rules)
		}
	case *Repeat:
		resolveNode(t.Node, rules)
		// case *Optional:
		// 	resolveNode(t.Node, rules)
	}
}

func (g *Grammar) Match(input string) bool {
	startRule, ok := g.Rules[g.Start]
	if !ok {
		panic("start rule not found: " + g.Start)
	}

	matches := startRule.Expr.Match(input, 0)
	return slices.Contains(matches, len(input))
}

func (g *Grammar) MatchPrefix(input string) bool {
	start := g.Rules[g.Start]
	return len(start.Expr.Match(input, 0)) > 0
}
