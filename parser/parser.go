// Package parser implements function to parse the regular expressions.
package parser

import (
	"fmt"
	"log"

	"github.com/8ayac/dfa-regex-engine/lexer"
	"github.com/8ayac/dfa-regex-engine/node"
	"github.com/8ayac/dfa-regex-engine/token"
)

// Parser has a slice of tokens to parse, and now looking token.
type Parser struct {
	tokens []*token.Token
	look   *token.Token
}

// NewParser returns a new Parser with the tokens to
// parse that were obtained by scanning.
func NewParser(s string) *Parser {
	p := &Parser{
		tokens: lexer.NewLexer(s).Scan(),
	}
	p.move()
	return p
}

// GetAST returns the root node of AST obtained by parsing.
func (psr *Parser) GetAST() node.Node {
	ast := psr.expression()
	return ast
}

// move updates the now looking token to the next token in token slice.
// If token slice is empty, will set token.EOF as now looking token.
func (psr *Parser) move() {
	if len(psr.tokens) == 0 {
		psr.look = token.NewToken('\x00', token.EOF)
	} else {
		psr.look = psr.tokens[0]
		psr.tokens = psr.tokens[1:]
	}
}

// moveWithValidation execute move() with validating whether
// now looking Token type is an expected (or not).
func (psr *Parser) moveWithValidation(expect token.Type) {
	if psr.look.Ty != expect {
		err := fmt.Sprintf("[syntax error] expect:\x1b[31m%s\x1b[0m actual:\x1b[31m%s\x1b[0m", expect, psr.look.Ty)
		log.Fatal(err)
	}
	psr.move()
}

// expression -> subexpr
func (psr *Parser) expression() node.Node {
	nd := psr.subexpr()
	psr.moveWithValidation(token.EOF)

	return nd
}

// subexpr -> subexpr '|' seq | seq
// (
//	subexpr  -> seq _subexpr
//	_subexpr -> '|' seq _subexpr | ε
// )
func (psr *Parser) subexpr() node.Node {
	nd := psr.seq()

	for {
		if psr.look.Ty == token.OpeUnion {
			psr.moveWithValidation(token.OpeUnion)
			nd2 := psr.seq()
			nd = node.NewUnion(nd, nd2)
		} else {
			break
		}
	}
	return nd
}

// seq -> subseq | ε
func (psr *Parser) seq() node.Node {
	if psr.look.Ty == token.LPAREN || psr.look.Ty == token.CHARACTER {
		return psr.subseq()
	}
	return node.NewCharacter('ε')

}

// subseq -> subseq star | star
// (
//	subseq  -> star _subseq
//	_subseq -> star _subseq | ε
// )
func (psr *Parser) subseq() node.Node {
	nd := psr.star()

	if psr.look.Ty == token.LPAREN || psr.look.Ty == token.CHARACTER {
		nd2 := psr.subseq()
		return node.NewConcat(nd, nd2)
	}
	return nd
}

// star -> factor '*' | factor
func (psr *Parser) star() node.Node {
	nd := psr.factor()

	if psr.look.Ty == token.OpeStar {
		psr.moveWithValidation(token.OpeStar)
		return node.NewStar(nd)
	}
	return nd
}

// factor -> '(' subexpr ')' | CHARACTER
func (psr *Parser) factor() node.Node {
	if psr.look.Ty == token.LPAREN {
		psr.moveWithValidation(token.LPAREN)
		nd := psr.subexpr()
		psr.moveWithValidation(token.RPAREN)
		return nd
	}
	nd := node.NewCharacter(psr.look.V)
	psr.moveWithValidation(token.CHARACTER)
	return nd
}
