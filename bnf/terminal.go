package bnf

import "fmt"

type terminal struct {
	Value string
}

func (t *terminal) match(ctx *context, pos int) ([]MatchResult, error) {
	if pos+len(t.Value) > len(ctx.input) {
		return nil, nil
	}
	if ctx.input[pos:pos+len(t.Value)] == t.Value {
		node := &ASTNode{
			Type:  "TERMINAL",
			Value: t.Value,
		}
		return []MatchResult{{
			End:   pos + len(t.Value),
			Nodes: []*ASTNode{node},
		}}, nil
	}
	return nil, nil
}

func (t *terminal) Expect() []string {
	return []string{fmt.Sprintf("%q", t.Value)}
}
