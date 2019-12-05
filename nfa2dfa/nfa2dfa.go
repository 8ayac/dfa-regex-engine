// Package nfa2dfa implements function to convert a NFA into a DFA.
package nfa2dfa

import (
	"github.com/8ayac/dfa-regex-engine/dfa"
	"github.com/8ayac/dfa-regex-engine/nfa"
)

// ToDFA converts a NFA into a DFA which recognizes the same formal language.
func ToDFA(nfa *nfa.NFA) *dfa.DFA {
	nfa.ToWithoutEpsilon()
	I, F, Delta := nfa.SubsetConstruction()
	return dfa.NewDFA(I, F, Delta)
}
