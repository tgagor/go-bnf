package bnf

type Choice struct {
	Options []Node
}

func (c *Choice) Match(input string, pos int) []int {
	var results []int

	for _, option := range c.Options {
		matches := option.Match(input, pos)
		if len(matches) > 0 {
			results = append(results, matches...)
		}
	}

	if len(results) == 0 {
		return nil
	}
	return results
}
