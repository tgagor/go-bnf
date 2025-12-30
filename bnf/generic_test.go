package bnf

func match(node Node, input string, pos int) []int {
	ctx := &Context{
		input: input,
		memo:  make(map[memoKey]*memoEntry),
	}
	return ctx.Match(node, pos)
}
