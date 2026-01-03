package bnf

import "strings"

// ASTNode represents a node in the parsed output tree (the Parse Tree).
type ASTNode struct {
	Type     string     // Rule name or "TERMINAL"
	Value    string     // The actual text matched
	Children []*ASTNode // Child nodes
}

func (n *ASTNode) String() string {
	if n.Type == "TERMINAL" {
		return n.Value
	}
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(n.Type)
	for _, c := range n.Children {
		sb.WriteString(" ")
		sb.WriteString(c.String())
	}
	sb.WriteString(")")
	return sb.String()
}
