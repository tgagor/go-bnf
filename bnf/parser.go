package bnf

import (
	"fmt"
	"io"
)

// Parser converts a stream of BNF tokens into a GrammarAST.
type Parser struct {
	lx   *Lexer
	look Token
	peek Token
}

// NewParser creates a new Parser for the given reader.
func NewParser(r io.Reader) (*Parser, error) {
	lx := NewLexer(r)
	look, err := lx.Next()
	if err != nil {
		return nil, err
	}
	peek, err := lx.Next()
	if err != nil {
		return nil, err
	}
	return &Parser{
		lx:   lx,
		look: look,
		peek: peek,
	}, nil
}

func (p *Parser) eat(t TokenType) (Token, error) {
	if p.look.Type != t {
		return Token{}, fmt.Errorf("unexpected token: %s, expected: %d", p.look.Text, t)
	}
	tok := p.look
	p.look = p.peek
	var err error
	p.peek, err = p.lx.Next()
	if err != nil {
		return Token{}, err
	}
	return tok, nil
}

// ParseGrammar parses the input into a complete GrammarAST.
func (p *Parser) ParseGrammar() (*GrammarAST, error) {
	var rules []*RuleAST
	for p.look.Type != EOF {
		// process rule otherwise
		r, err := p.parseRule()
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return &GrammarAST{Rules: rules}, nil
}

func (p *Parser) parseRule() (*RuleAST, error) {
	tok, err := p.eat(IDENT)
	if err != nil {
		return nil, err
	}
	name := tok.Text

	if _, err := p.eat(ASSIGN); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &RuleAST{Name: name, Expr: expr}, nil
}

func (p *Parser) parseExpr() (ExprAST, error) {
	left, err := p.parseSeq()
	if err != nil {
		return nil, err
	}

	options := []ExprAST{left}

	for {
		// multiline alternative
		if p.look.Type == PIPE {
			if _, err := p.eat(PIPE); err != nil {
				return nil, err
			}
			r, err := p.parseSeq()
			if err != nil {
				return nil, err
			}
			options = append(options, r)
			continue
		}

		// STOP: next rule starts
		if p.isRuleStart() || p.look.Type == EOF {
			break
		}

		break
	}

	if len(options) == 1 {
		return left, nil
	}
	return &ChoiceAST{Options: options}, nil
}

func (p *Parser) parseSeq() (ExprAST, error) {
	var elems []ExprAST

	for {
		// STOP conditions
		if p.look.Type == PIPE || p.isRuleStart() || p.look.Type == EOF {
			break
		}

		switch p.look.Type {
		case IDENT, STRING, REGEX, LPAREN:
			e, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			elems = append(elems, e)
		default:
			return singleOrSeq(elems), nil
		}
	}

	return singleOrSeq(elems), nil
}

func singleOrSeq(elems []ExprAST) ExprAST {
	if len(elems) == 1 {
		return elems[0]
	}
	return &SeqAST{Elements: elems}
}

func (p *Parser) parseFactor() (ExprAST, error) {
	atom, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	switch p.look.Type {
	case STAR:
		p.eat(STAR)
		return &RepeatAST{Node: atom, Min: 0, Max: -1}, nil
	case PLUS:
		p.eat(PLUS)
		return &RepeatAST{Node: atom, Min: 1, Max: -1}, nil
	case QMARK:
		p.eat(QMARK)
		return &RepeatAST{Node: atom, Min: 0, Max: 1}, nil
		// TODO: consider adding OptionalAST
		// return &OptionalAST{Node: atom}
	}

	return atom, nil
}

func (p *Parser) parseAtom() (ExprAST, error) {
	switch p.look.Type {
	case IDENT:
		tok, err := p.eat(IDENT)
		if err != nil {
			return nil, err
		}
		return &IdentAST{Name: tok.Text}, nil
	case NT_IDENT:
		tok, err := p.eat(NT_IDENT)
		if err != nil {
			return nil, err
		}
		return &IdentAST{Name: tok.Text}, nil
	case STRING:
		tok, err := p.eat(STRING)
		if err != nil {
			return nil, err
		}
		return &StringAST{Value: tok.Text}, nil
	case REGEX:
		tok, err := p.eat(REGEX)
		if err != nil {
			return nil, err
		}
		return &RegexAST{Pattern: tok.Text}, nil
	case LPAREN:
		if _, err := p.eat(LPAREN); err != nil {
			return nil, err
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.eat(RPAREN); err != nil {
			return nil, err
		}
		return e, nil
	}
	return nil, fmt.Errorf("unexpected token in atom: %v", p.look)
}

func (p *Parser) isRuleStart() bool {
	return p.look.Type == IDENT && p.peek.Type == ASSIGN
}
