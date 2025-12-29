package bnf

type Char struct {
	C byte
}

func (c *Char) Match(input string, pos int) []int {
	if pos < len(input) && input[pos] == c.C {
		return []int{pos + 1}
	}
	return nil
}
