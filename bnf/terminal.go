package bnf

import "fmt"

type terminal struct {
	Value string
}

func (t *terminal) match(ctx *Context, pos int) []int {
	if pos+len(t.Value) > len(ctx.input) {
		return nil
	}
	if ctx.input[pos:pos+len(t.Value)] == t.Value {
		return []int{pos + len(t.Value)}
	}
	return nil
}

func (t *terminal) Expect() []string {
	return []string{fmt.Sprintf("%q", t.Value)}
}
