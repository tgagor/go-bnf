package bnf

// raw AST, without links yet
type GrammarAST struct {
	Rules []*RuleAST
}

type RuleAST struct {
	Name string
	Expr ExprAST
}

type ExprAST any

type (
	ChoiceAST struct {
		Options []ExprAST
	}

	SeqAST struct {
		Elements []ExprAST
	}

	RepeatAST struct {
		Node ExprAST
		Min  int
		Max  int // -1 = infinity
	}

	IdentAST struct {
		Name string
	}

	StringAST struct {
		Value string
	}
)

func BuildGrammar(ast *GrammarAST) *Grammar {
	rules := map[string]*Rule{}

	// 1. create rules
	for _, r := range ast.Rules {
		rules[r.Name] = &Rule{Name: r.Name}
	}

	// 2. build expr
	for _, r := range ast.Rules {
		rules[r.Name].Expr = buildNode(r.Expr, rules)
	}

	return &Grammar{
		Start: ast.Rules[0].Name,
		Rules: rules,
	}
}

func buildNode(e ExprAST, rules map[string]*Rule) node {
	switch t := e.(type) {
	case *StringAST:
		return &terminal{Value: t.Value}

	case *IdentAST:
		return &nonTerminal{Name: t.Name, Rule: rules[t.Name]}

	case *SeqAST:
		var elems []node
		for _, e := range t.Elements {
			elems = append(elems, buildNode(e, rules))
		}
		return &sequence{Elements: elems}

	case *ChoiceAST:
		var opts []node
		for _, o := range t.Options {
			opts = append(opts, buildNode(o, rules))
		}
		return &choice{Options: opts}

	case *RepeatAST:
		return &repeat{
			Node: buildNode(t.Node, rules),
			Min:  t.Min,
		}
	}
	panic("unknown AST")
}
