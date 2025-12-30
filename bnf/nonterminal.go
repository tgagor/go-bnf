package bnf

type NonTerminal struct {
	Name string
	Rule *Rule // will be set up in 2nd pass
}

// func (n *NonTerminal) Match(input string, pos int) []int {
// 	return []int{0}
// }

func (n *NonTerminal) match(ctx *Context, pos int) []int {
	if n.Rule == nil {
		panic("NonTerminal without Rule: " + n.Name)
	}
	return ctx.Match(n.Rule.Expr, pos)
}
