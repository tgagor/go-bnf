package bnf

import "fmt"

type nonTerminal struct {
	Name string
	Rule *Rule // will be set up in 2nd pass
}

func (n *nonTerminal) match(ctx *context, pos int) ([]MatchResult, error) {
	if n.Rule == nil {
		return nil, fmt.Errorf("NonTerminal without Rule: %s", n.Name)
	}

	ctx.push(n.Name)
	defer ctx.pop()

	matches, err := ctx.Match(n.Rule.Expr, pos)
	if err != nil {
		return nil, err
	}

	var results []MatchResult
	for _, m := range matches {
		// Wrap children in a new ASTNode representing this rule
		node := &ASTNode{
			Type:     n.Name,
			Children: m.Nodes,
			// Value is usually empty for non-terminals unless we want capturing
		}
		results = append(results, MatchResult{
			End:   m.End,
			Nodes: []*ASTNode{node},
		})
	}

	return results, nil
}

func (n *nonTerminal) Expect() []string {
	return []string{n.Name}
}
