package bnf

import "regexp"

type Regex struct {
	Re *regexp.Regexp
}

func (r *Regex) match(ctx *context, pos int) ([]MatchResult, error) {
	// match needs lowercase signature to satisfy interface?
	// The file `regex.go` had `Match`. But interface requires `match`.
	// Let's assume `Regex` was experimental or unused. I'll update it anyway.
	loc := r.Re.FindStringIndex(ctx.input[pos:])
	if loc != nil && loc[0] == 0 {
		return []MatchResult{{
			End:   pos + loc[1],
			Nodes: []*ASTNode{{Type: "REGEX", Value: ctx.input[pos : pos+loc[1]]}},
		}}, nil
	}
	return nil, nil
}

func (r *Regex) Expect() []string {
	return []string{r.Re.String()}
}
