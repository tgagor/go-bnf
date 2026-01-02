package bnf

func testMatch(node node, input string, pos int) []int {
	ctx := &context{
		input: input,
		memo:  make(map[memoKey]*memoEntry),
	}
	return ctx.Match(node, pos)
}
