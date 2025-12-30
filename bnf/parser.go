package bnf

import (
	"io"
	"os"
	"strings"
)

type Parser struct {
	lx   *Lexer
	look Token
}

func NewParser(r io.Reader) *Parser {
	lx := NewLexer(r)
	return &Parser{
		lx:   lx,
		look: lx.Next(),
	}
}

func (p *Parser) eat(t TokenType) Token {
	if p.look.Type != t {
		panic("unexpected token: " + p.look.Text)
	}
	tok := p.look
	p.look = p.lx.Next()
	return tok
}

func (p *Parser) ParseGrammar() *GrammarAST {
	var rules []*RuleAST
	for p.look.Type != EOF {
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

	for p.look.Type == PIPE {
		p.eat(PIPE)
		options = append(options, p.parseSeq())
	}

	if len(options) == 1 {
		return left
	}
	return &ChoiceAST{Options: options}
}

func (p *Parser) parseSeq() ExprAST {
	var elems []ExprAST
	for p.look.Type == IDENT || p.look.Type == STRING || p.look.Type == LPAREN {
		elems = append(elems, p.parseFactor())
	}

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
	}

	return atom
}

func (p *Parser) parseAtom() ExprAST {
	switch p.look.Type {
	case IDENT:
		return &IdentAST{Name: p.eat(IDENT).Text}
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

func LoadGrammarFile(path string) (*Grammar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := NewParser(f)
	ast := p.ParseGrammar()
	return BuildGrammar(ast), nil
}

func LoadGrammarString(s string) *Grammar {
	p := NewParser(strings.NewReader(s))
	ast := p.ParseGrammar()
	return BuildGrammar(ast)
}
