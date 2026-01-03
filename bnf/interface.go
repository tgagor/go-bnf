package bnf

type MatchResult struct {
	End   int
	Nodes []*ASTNode
}

type node interface {
	match(ctx *context, pos int) ([]MatchResult, error)
	Expect() []string // for error reporting
}
