package bnf

type Choice struct {
	Options []Node
}

func (c *Choice) match(ctx *Context, pos int) []int {
	var results []int
	for _, opt := range c.Options {
		results = append(results, ctx.Match(opt, pos)...)
	}
	return results
}

func (c *Choice) Expect() []string {
    var out []string
    for _, o := range c.Options {
        out = append(out, o.Expect()...)
    }
    return out
}
