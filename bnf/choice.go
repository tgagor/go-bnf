package bnf

type choice struct {
	Options []node
}

func (c *choice) match(ctx *context, pos int) ([]MatchResult, error) {
	var results []MatchResult
	for _, opt := range c.Options {
		matches, err := ctx.Match(opt, pos)
		if err != nil {
			return nil, err
		}
		results = append(results, matches...)
	}
	return results, nil
}

func (c *choice) Expect() []string {
	var out []string
	for _, o := range c.Options {
		out = append(out, o.Expect()...)
	}
	return out
}
