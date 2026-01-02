package bnf

type nonTerminal struct {
	Name string
	Rule *Rule // will be set up in 2nd pass
}

func (n *nonTerminal) match(ctx *Context, pos int) []int {
	if n.Rule == nil {
		panic("NonTerminal without Rule: " + n.Name)
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
