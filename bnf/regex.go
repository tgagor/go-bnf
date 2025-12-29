package bnf

import "regexp"

type Regex struct {
	Re *regexp.Regexp
}

func (r *Regex) Match(input string, pos int) []int {
	loc := r.Re.FindStringIndex(input[pos:])
	if loc != nil && loc[0] == 0 {
		return []int{pos + loc[1]}
	}
	return nil
}
