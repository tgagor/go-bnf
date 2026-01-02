package bnf

import (
	"io"
	"os"
	"strings"
)

type Parser struct {
	lx   *Lexer
	look Token
	peek Token
}

func NewParser(r io.Reader) *Parser {
	lx := NewLexer(r)
	return &Parser{
		lx:   lx,
		look: lx.Next(),
		peek: lx.Next(),
	}
}

func (p *Parser) eat(t TokenType) Token {
	if p.look.Type != t {
		panic("unexpected token: " + p.look.Text)
	}
	tok := p.look
	p.look = p.peek
	p.peek = p.lx.Next()
	return tok
}

func (p *Parser) ParseGrammar() *GrammarAST {
	var rules []*RuleAST
	for p.look.Type != EOF {
		// process rule otherwise
		rules = append(rules, p.parseRule())
	}
	return &GrammarAST{Rules: rules}
}

func (p *Parser) parseRule() *RuleAST {
	name := p.eat(IDENT).Text
	p.eat(ASSIGN)
	expr := p.parseExpr()

	return &RuleAST{Name: name, Expr: expr}
}

func (p *Parser) parseExpr() ExprAST {
	left := p.parseSeq()
	options := []ExprAST{left}

	for {
		// multiline alternative
		if p.look.Type == PIPE {
			p.eat(PIPE)
			options = append(options, p.parseSeq())
			continue
		}

		// STOP: next rule starts
		if p.isRuleStart() || p.look.Type == EOF {
			break
		}

		break
	}

	if len(options) == 1 {
		return left
	}
	return &ChoiceAST{Options: options}
}

func (p *Parser) parseSeq() ExprAST {
	var elems []ExprAST

	for {
		// STOP conditions
		if p.look.Type == PIPE || p.isRuleStart() || p.look.Type == EOF {
			break
		}

		switch p.look.Type {
		case IDENT, STRING, LPAREN:
			elems = append(elems, p.parseFactor())
		default:
			return singleOrSeq(elems)
		}
	}

	return singleOrSeq(elems)
}

func singleOrSeq(elems []ExprAST) ExprAST {
	if len(elems) == 1 {
		return elems[0]
	}
	return &SeqAST{Elements: elems}
}


func (p *Parser) parseFactor() ExprAST {
	atom := p.parseAtom()

	switch p.look.Type {
	case STAR:
		p.eat(STAR)
		return &RepeatAST{Node: atom, Min: 0, Max: -1}
	case PLUS:
		p.eat(PLUS)
		return &RepeatAST{Node: atom, Min: 1, Max: -1}
	case QMARK:
		p.eat(QMARK)
		return &RepeatAST{Node: atom, Min: 0, Max: 1}
		// TODO: consider adding OptionalAST
		// return &OptionalAST{Node: atom}
	}

	return atom
}

func (p *Parser) parseAtom() ExprAST {
	switch p.look.Type {
	case IDENT:
		return &IdentAST{Name: p.eat(IDENT).Text}
	case NT_IDENT:
		return &IdentAST{Name: p.eat(NT_IDENT).Text}
	case STRING:
		return &StringAST{Value: p.eat(STRING).Text}
	case LPAREN:
		p.eat(LPAREN)
		e := p.parseExpr()
		p.eat(RPAREN)
		return e
	}
	panic("unexpected token")
}

func (p *Parser) isRuleStart() bool {
	return p.look.Type == IDENT && p.peek.Type == ASSIGN
}

func LoadGrammarIOReader(r io.Reader) *Grammar {
	p := NewParser(r)
	ast := p.ParseGrammar()
	return BuildGrammar(ast)
}

func LoadGrammarFile(path string) (*Grammar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadGrammarIOReader(f), nil
}

func LoadGrammarString(s string) *Grammar {
	return LoadGrammarIOReader(strings.NewReader(s))
}
