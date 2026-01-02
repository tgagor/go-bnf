package bnf

type choice struct {
	Options []Node
}

func (c *choice) match(ctx *Context, pos int) []int {
	var results []int
	for _, opt := range c.Options {
		results = append(results, ctx.Match(opt, pos)...)
	}
	return results
}

func (c *choice) Expect() []string {
	var out []string
	for _, o := range c.Options {
		out = append(out, o.Expect()...)
	}
	return out
}
