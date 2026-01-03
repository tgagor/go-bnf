package bnf

import (
	"fmt"
	"strings"
)

// ASTNode represents a node in the parsed output tree (the Parse Tree).
type ASTNode struct {
	Type     string     // Rule name or "TERMINAL"
	Value    string     // The actual text matched
	Children []*ASTNode // Child nodes
}

func (n *ASTNode) String() string {
	var sb strings.Builder
	n.format(&sb, 0)
	return sb.String()
}

func (n *ASTNode) format(sb *strings.Builder, level int) {
	indent := strings.Repeat("  ", level)

	// Terminals are leaf nodes
	if n.Type == "TERMINAL" || n.Type == "REGEX" {
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("%q", n.Value))
		return
	}

	sb.WriteString(indent)
	sb.WriteString("(")
	sb.WriteString(n.Type)

	if len(n.Children) == 0 {
		if n.Value != "" {
			sb.WriteString(" ")
			sb.WriteString(fmt.Sprintf("%q", n.Value))
		}
		sb.WriteString(")")
		return
	}

	for _, child := range n.Children {
		sb.WriteString("\n")
		child.format(sb, level+1)
	}
	sb.WriteString("\n")
	sb.WriteString(indent)
	sb.WriteString(")")
}
