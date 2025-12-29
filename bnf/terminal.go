package bnf

type Terminal struct {
	Value string
}

func (t *Terminal) Match(input string, pos int) []int {
	// if end of input, no match
	if pos+len(t.Value) > len(input) {
		return nil
	}

	// check if part matches
	if input[pos:pos+len(t.Value)] == t.Value {
		return []int{pos + len(t.Value)}
	}

	return nil
}
