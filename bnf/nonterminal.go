package bnf

import "fmt"

type nonTerminal struct {
	Name string
	Rule *Rule // will be set up in 2nd pass
}

func (n *nonTerminal) match(ctx *context, pos int) ([]int, error) {
	if n.Rule == nil {
		return nil, fmt.Errorf("NonTerminal without Rule: %s", n.Name)
	}

	// debug / call stack tracking
	// only here as only NonTerminal prodive meaningful rules
	// that users can understand
	ctx.push(n.Name)
	defer ctx.pop()

	return ctx.Match(n.Rule.Expr, pos)
}

func (n *nonTerminal) Expect() []string {
	return []string{n.Name}
}
