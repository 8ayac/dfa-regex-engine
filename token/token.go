// Package token provides tokens for parsing the regular expressions.
package token

import "fmt"

// Type is integer to identify the type of token.
type Type int

// Each token is identified by a unique integer.
const (
	CHARACTER Type = iota
	UNION
	STAR
	PLUS
	LPAREN
	RPAREN
	EOF
)

func (k Type) String() string {
	switch k {
	case CHARACTER:
		return "CHARACTER"
	case UNION:
		return "UNION"
	case STAR:
		return "STAR"
	case PLUS:
		return "PLUS"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case EOF:
		return "EOF"
	default:
		return ""
	}
}

// Token represents a token.
type Token struct {
	V  rune // token value
	Ty Type // token type
}

func (t Token) String() string {
	return fmt.Sprintf("V -> \x1b[32m%v\x1b[0m\tKind -> \x1b[32m%v\x1b[0m", string(t.V), t.Ty)
}

// NewToken returns a new Token.
func NewToken(value rune, k Type) Token {
	return Token{
		V:  value,
		Ty: k,
	}
}
