package bnf

import "slices"

type Node interface {
	// Match(input string, pos int) []int
	match(ctx *Context, pos int) []int
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

func (g *Grammar) SetStart(name string) {
	g.Start = name
}

func (g *Grammar) Match(input string) bool {
	if g.Start == "" {
		panic("start rule not defined")
	}
	return g.MatchFrom(g.Start, input)
}

func (g *Grammar) MatchFrom(start string, input string) bool {
	rule, ok := g.Rules[start]
	if !ok {
		panic("unknown start rule: " + start)
	}

	ctx := NewContext(input)
	matches := ctx.Match(rule.Expr, 0)
	return slices.Contains(matches, len(input))
}

func (g *Grammar) MatchPrefix(input string) bool {
	start := g.Rules[g.Start]
	ctx := NewContext(input)
	return len(ctx.Match(start.Expr, 0)) > 0
}
