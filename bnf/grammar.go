package bnf

import (
	"slices"
)

type Node interface {
	match(ctx *context, pos int) []int
	Expect() []string // for error reporting, node types expected at this point
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

func (g *Grammar) Match(input string) (bool, error) {
	if g.Start == "" {
		panic("start rule not defined")
	}
	return g.MatchFrom(g.Start, input)
}

func (g *Grammar) MatchFrom(start string, input string) (bool, error) {
	rule, ok := g.Rules[start]
	if !ok {
		panic("unknown start rule: " + start)
	}

	ctx := NewContext(input)
	matches := ctx.Match(rule.Expr, 0)
	if slices.Contains(matches, len(input)) {
		return true, nil
	}

	return false, ctx.error
}

func (g *Grammar) MatchPrefix(input string) bool {
	start := g.Rules[g.Start]
	ctx := NewContext(input)
	return len(ctx.Match(start.Expr, 0)) > 0
}
