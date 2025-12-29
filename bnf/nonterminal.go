package bnf

type NonTerminal struct {
	Name string
	Rule *Rule // will be set up in 2nd pass
}

func (n *NonTerminal) Match(input string, pos int) []int {
	return []int{0}
}
