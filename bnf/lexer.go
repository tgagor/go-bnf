package bnf

import (
	"bufio"
	"io"
	"strings"
)

type TokenType int

const (
	EOF      TokenType = iota
	IDENT              // digit
	NT_IDENT           // <digit>
	STRING             // "a string"
	ASSIGN             // ::=
	PIPE               // |
	STAR               // *
	PLUS               // +
	QMARK              // ?
	LPAREN             // (
	RPAREN             // )
)

type Token struct {
	Type TokenType
	Text string
}

type Lexer struct {
	r *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

func (l *Lexer) Next() Token {
	// 1. skip whitespace
	for {
		ch, _, err := l.r.ReadRune()
		if err == io.EOF {
			return Token{Type: EOF}
		}

		if isWhitespace(ch) {
			continue
		}

		// line comment
		if ch == ';' || ch == '#' {
			l.skipUntilEOL()
			continue
		}

		if ch == '/' {
			next, _, err := l.r.ReadRune()
			if err == nil && next == '/' {
				l.skipUntilEOL()
				continue
			}
			l.r.UnreadRune()
		}

		// ch is meaningful
		l.r.UnreadRune()
		break
	}

	ch, _, err := l.r.ReadRune()
	if err == io.EOF {
		return Token{Type: EOF}
	}

	// 2. identificator
	if isIdentStart(ch) {
		var sb strings.Builder
		sb.WriteRune(ch)

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil || !isIdentPart(ch) {
				if err == nil {
					l.r.UnreadRune()
				}
				break
			}
			sb.WriteRune(ch)
		}

		return Token{
			Type: IDENT,
			Text: sb.String(),
		}
	}

	// 2.5. <identificator>
	if ch == '<' {
		var sb strings.Builder

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil {
				panic("unterminated <identifier>")
			}
			if ch == '>' {
				break
			}
			if !isIdentPart(ch) {
				panic("invalid character in <identifier>: " + string(ch))
			}
			sb.WriteRune(ch)
		}

		if sb.Len() == 0 {
			panic("empty <identifier>")
		}

		return Token{
			Type: IDENT, // normalize NT_IDENT as <> are just a sugar here
			Text: sb.String(),
		}
	}

	// 3. string literal
	if ch == '"' || ch == '\'' {
		quote := ch
		var sb strings.Builder

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil {
				panic("unterminated string literal")
			}

			if ch == quote {
				break
			}

			if ch == '\\' {
				esc, _, err := l.r.ReadRune()
				if err != nil {
					panic("unterminated escape sequence")
				}
				switch esc {
				case '"':
					sb.WriteRune('"')
				case '\'':
					sb.WriteRune('\'')
				case '\\':
					sb.WriteRune('\\')
				case 'n':
					sb.WriteRune('\n')
				case 't':
					sb.WriteRune('\t')
				default:
					panic("unknown escape sequence: \\" + string(esc))
				}
				continue
			}

			sb.WriteRune(ch)
		}

		return Token{
			Type: STRING,
			Text: sb.String(),
		}
	}

	// 4. ASSIGN ::=
	if ch == ':' {
		ch2, _, err := l.r.ReadRune()
		if err == nil && ch2 == ':' {
			ch3, _, err := l.r.ReadRune()
			if err == nil && ch3 == '=' {
				return Token{Type: ASSIGN, Text: "::="}
			}
		}
		panic("expected ::=")
	}

	// 5. pojedyncze symbole
	switch ch {
	case '|':
		return Token{Type: PIPE, Text: "|"}
	case '*':
		return Token{Type: STAR, Text: "*"}
	case '+':
		return Token{Type: PLUS, Text: "+"}
	case '?':
		return Token{Type: QMARK, Text: "?"}
	case '(':
		return Token{Type: LPAREN, Text: "("}
	case ')':
		return Token{Type: RPAREN, Text: ")"}
	}

	panic("unexpected character: " + string(ch))
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isIdentStart(ch rune) bool {
	return ch == '_' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

func isIdentPart(ch rune) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '-'
}

func (l *Lexer) skipUntilEOL() error {
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			return err
		}

		// Linux new line
		if ch == '\n' {
			return nil
		}
		// Windows, or old Mac style new line
		if ch == '\r' {
			// check if it's not \r\n
			next, _, err := l.r.ReadRune()
			if err == nil && next != '\n' {
				l.r.UnreadRune()
			}
			return nil
		}
	}
}
