package bnf

import "fmt"

type Terminal struct {
	Value string
}

func (t *Terminal) match(ctx *Context, pos int) []int {
	if pos+len(t.Value) > len(ctx.input) {
		return nil
	}
	if ctx.input[pos:pos+len(t.Value)] == t.Value {
		return []int{pos + len(t.Value)}
	}
	return nil
}

func (t *Terminal) Expect() []string {
    return []string{fmt.Sprintf("%q", t.Value)}
}
