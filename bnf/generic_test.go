package bnf

func testMatch(node Node, input string, pos int) []int {
	ctx := &context{
		input: input,
		memo:  make(map[memoKey]*memoEntry),
	}
	return ctx.Match(node, pos)
}
