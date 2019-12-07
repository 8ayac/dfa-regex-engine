// Package dfaregex provides a DFA regex engine(DFA engine).
package dfaregex

import (
	"github.com/8ayac/dfa-regex-engine/dfa"
	"github.com/8ayac/dfa-regex-engine/nfa2dfa"
	"github.com/8ayac/dfa-regex-engine/parser"
	"github.com/8ayac/dfa-regex-engine/utils"
)

// Regexp has a DFA and regexp string.
type Regexp struct {
	regexp string
	d      *dfa.DFA
}

// NewRegexp return a new Regexp.
func NewRegexp(re string) *Regexp {
	psr := parser.NewParser(re)
	ast := psr.GetAST()
	frg := ast.Assemble(utils.NewContext())
	nfa := frg.Build()
	d := nfa2dfa.ToDFA(nfa)
	d.Minimize()

	return &Regexp{
		regexp: re,
		d:      d,
	}
}

// Compile is a wrapper function of NewRegexp().
func Compile(re string) *Regexp {
	return NewRegexp(re)
}

// Match returns whether the input string matches the regular expression.
func (re *Regexp) Match(s string) bool {
	rt := re.d.GetRuntime()
	return rt.Matching(s)
}
