package bnf

import "fmt"

type terminal struct {
	Value string
}

func (t *terminal) match(ctx *context, pos int) ([]int, error) {
	if pos+len(t.Value) > len(ctx.input) {
		return nil, nil
	}
	if ctx.input[pos:pos+len(t.Value)] == t.Value {
		return []int{pos + len(t.Value)}, nil
	}
	return nil, nil
}

func (t *terminal) Expect() []string {
	return []string{fmt.Sprintf("%q", t.Value)}
}
