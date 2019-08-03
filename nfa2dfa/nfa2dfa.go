// Package nfa2dfa implements function to convert a NFA into a DFA.
package nfa2dfa

import (
	"github.com/8ayac/dfa-regex-engine/dfa"
	"github.com/8ayac/dfa-regex-engine/dfa/dfarule"
	"github.com/8ayac/dfa-regex-engine/nfa"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// DFAStatesMap associates subsets of the NFA state set with the states of the DFA.
type DFAStatesMap map[mapset.Set]utils.State

// ToDFA converts a NFA into a DFA which recognizes the same formal language.
func ToDFA(nfa *nfa.NFA) *dfa.DFA {
	nfa.ToWithoutEpsilon()
	I, F, Delta := subsetConstruction(nfa)
	return dfa.NewDFA(I, F, Delta)
}

// subsetConstruction implements Subset Construction.
// Returns the data for constructing the equivalent DFA from the NFA given in the argument.
// For details: https://en.wikipedia.org/wiki/Powerset_construction
func subsetConstruction(nfa *nfa.NFA) (dInit utils.State, dAccepts mapset.Set, dRules dfarule.RuleMap) {
	dInit = utils.NewState(0)
	dAccepts = mapset.NewSet()
	dRules = dfarule.RuleMap{}

	Sigma := nfa.AllSymbol()
	dStates := DFAStatesMap{}
	dStates[mapset.NewSet(nfa.Init)] = utils.NewState(0)

	queue := mapset.NewSet(mapset.NewSet(nfa.Init))
	for queue.N() != 0 {
		nDst := queue.Pop().(mapset.Set) // the state set which can be reached from a NFA state.

		if nfa.Accepts.Intersect(nDst).N() > 0 {
			dAccepts.Add(dStates.getState(nDst))
		}

		for c := range Sigma.Iter() {
			dDst := mapset.NewSet()
			for q := range nDst.Iter() {
				d, ok := nfa.CalcDst(q.(utils.State), c.(rune))
				if ok {
					dDst = dDst.Union(d)
				}
			}

			if dDst.N() == 0 {
				continue
			}

			if !dStates.haveKey(dDst) {
				queue.Add(dDst)
				dStates[dDst] = utils.NewState(len(dStates))
			}

			dNext := dStates.getState(dDst)
			dFrom := dStates.getState(nDst)
			dRules[dfarule.NewRuleArgs(dFrom, c.(rune))] = dNext
		}
	}

	return
}

// getState returns the state associated with key.
// If there is no corresponding state, it returns empty state struct.
func (dm DFAStatesMap) getState(key mapset.Set) utils.State {
	if dm.haveKey(key) {
		for k := range dm {
			if k.Equal(key) {
				return dm[k]
			}
		}
	}
	return utils.State{}
}

// haveKey returns whether DFAStatesMap has the set given as the argument "key".
// If it has, returns true.
func (dm DFAStatesMap) haveKey(key mapset.Set) bool {
	for k := range dm {
		if k.Equal(key) {
			return true
		}
	}
	return false
}
