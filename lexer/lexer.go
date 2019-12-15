// Package lexer implements lexer for a simple regular expression.
package lexer

import (
	"github.com/8ayac/dfa-regex-engine/token"
)

// Lexer has a slice of symbols to analyze.
type Lexer struct {
	s []rune // string to be analyzed
}

// NewLexer returns a new Lexer.
// This constructor create a sequence of symbols from
// the string given in the argument and hold it.
func NewLexer(s string) *Lexer {
	return &Lexer{
		s: []rune(s),
	}
}

// Scan returns the token list to which converted from
// the symbol slice held in Lexer struct.
func (l *Lexer) Scan() (tokenList []token.Token) {
	for i := 0; i < len(l.s); i++ {
		switch l.s[i] {
		case '\x00':
			tokenList = append(tokenList, token.NewToken(l.s[i], token.EOF))
		case '|':
			tokenList = append(tokenList, token.NewToken(l.s[i], token.UNION))
		case '(':
			tokenList = append(tokenList, token.NewToken(l.s[i], token.LPAREN))
		case ')':
			tokenList = append(tokenList, token.NewToken(l.s[i], token.RPAREN))
		case '*':
			tokenList = append(tokenList, token.NewToken(l.s[i], token.STAR))
		case '+':
			tokenList = append(tokenList, token.NewToken(l.s[i], token.PLUS))
		case '\\':
			tokenList = append(tokenList, token.NewToken(l.s[i+1], token.CHARACTER))
			i++
		default:
			tokenList = append(tokenList, token.NewToken(l.s[i], token.CHARACTER))
		}
	}
	return
}
