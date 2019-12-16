// Package nfa implements Non-Deterministic Finite Automaton(NFA).
package nfa

import (
	"github.com/8ayac/dfa-regex-engine/dfa/dfarule"
	"github.com/8ayac/dfa-regex-engine/nfa/nfarule"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// NFA represents a Non-Deterministic Finite Automaton.
type NFA struct {
	I     utils.State     // initial state
	F     mapset.Set      // accept states
	Rules nfarule.RuleMap // transition function
}

// NewNFA returns a new NFA.
func NewNFA(init utils.State, accepts mapset.Set, rules nfarule.RuleMap) *NFA {
	return &NFA{
		I:     init,
		F:     accepts,
		Rules: rules,
	}
}

// allStates returns a set of the all "From State" in Rule.
func (nfa *NFA) allStates() mapset.Set {
	states := mapset.NewSet()
	for key := range nfa.Rules {
		states.Add(key.From)
	}
	return states
}

// AllSymbol returns a set of the all "Symbol" in Rule.
func (nfa *NFA) AllSymbol() mapset.Set {
	symbols := mapset.NewSet()
	for key := range nfa.Rules {
		symbols.Add(key.C)
	}
	return symbols
}

// CalcDst returns, according to the transition function, a set of states
// to which transition is executed when c is received in the state of argument q.
func (nfa *NFA) CalcDst(q utils.State, c rune) (mapset.Set, bool) {
	s, ok := nfa.Rules[nfarule.NewRuleArgs(q, c)]
	if ok {
		return s, true
	}
	return nil, false
}

// ToWithoutEpsilon update ε-NFA to NFA whose no epsilon transitions.
func (nfa *NFA) ToWithoutEpsilon() {
	if nfa.F.IsSubset(nfa.epsilonClosure(nfa.I)) {
		nfa.F.Add(nfa.I)
	}
	nfa.Rules = nfa.removeEpsilonRule()
}

// removeEpsilonRule returns a new RuleMap removing epsilon transitions
// from original RuleMap.
func (nfa *NFA) removeEpsilonRule() (newRule nfarule.RuleMap) {
	newRule = nfarule.RuleMap{}
	states, sym := nfa.allStates(), nfa.AllSymbol()
	sym.Remove('ε')

	for q := range states.Iter() {
		for c := range sym.Iter() {
			q := q.(utils.State)
			c := c.(rune)
			for mid := range nfa.epsilonClosure(q).Iter() {
				dst := nfa.epsilonExpand(mid.(utils.State), c)
				s, ok := newRule[nfarule.NewRuleArgs(q, c)]
				if !ok {
					s = mapset.NewSet()
				}
				newRule[nfarule.NewRuleArgs(q, c)] = s.Union(dst)
			}
		}
	}

	for k := range newRule {
		if newRule[k].N() == 0 {
			delete(newRule, k)
		}
	}

	return
}

// epsilonExpand returns the state set, which is a result of simulating the transitions like 'ε*->symbol->ε*'.
func (nfa *NFA) epsilonExpand(state utils.State, symbol rune) (final mapset.Set) {
	first := nfa.epsilonClosure(state)

	second := mapset.NewSet()
	for q := range first.Iter() {
		if dst, ok := nfa.CalcDst(q.(utils.State), symbol); ok {
			second = second.Union(dst)
		}
	}

	final = mapset.NewSet()
	for q := range second.Iter() {
		dst := nfa.epsilonClosure(q.(utils.State))
		final = final.Union(dst)
	}

	return
}

// epsilonClosure returns a set of reachable states with epsilon transitions only.
func (nfa *NFA) epsilonClosure(state utils.State) (reachable mapset.Set) {
	reachable = mapset.NewSet(state)

	modified := true
	for modified {
		modified = false
		for q := range reachable.Iter() {
			dst, ok := nfa.CalcDst(q.(utils.State), 'ε')
			if !ok || reachable.IsSuperset(dst) {
				continue
			}
			reachable = reachable.Union(dst)
			modified = true
		}
	}
	return
}

// subsetConstruction implements Subset Construction.
// Returns the data for constructing the equivalent DFA from the NFA given in the argument.
// For details: https://en.wikipedia.org/wiki/Powerset_construction
func (nfa *NFA) SubsetConstruction() (dI utils.State, dF mapset.Set, dRules dfarule.RuleMap) {
	I := nfa.I
	F := nfa.F
	Rules := nfa.Rules

	dI = utils.NewState(0)
	dF = mapset.NewSet()
	dRules = dfarule.RuleMap{}

	dStates := DFAStatesMap{}
	dStates[mapset.NewSet(I)] = utils.NewState(0)

	queue := mapset.NewSet(mapset.NewSet(I))
	for queue.N() != 0 {
		dstate := queue.Pop().(mapset.Set) // the state set which can be reached from a NFA state.

		if F.Intersect(dstate).N() > 0 {
			dF.Add(dStates.getState(dstate))
		}

		Sigma := nfa.AllSymbol()
		for c := range Sigma.Iter() {
			dnext := mapset.NewSet()
			for q := range dstate.Iter() {
				d, ok := Rules[nfarule.NewRuleArgs(q.(utils.State), c.(rune))]
				if ok {
					dnext = dnext.Union(d)
				}
			}

			if dnext.N() == 0 {
				continue
			}

			if !dStates.haveKey(dnext) {
				queue.Add(dnext)
				dStates[dnext] = utils.NewState(len(dStates))
			}

			for k := range dStates {
				if k.Equal(dnext) {
					dnext = k // Swap to avoid problems with pointers
				}
			}
			dRules[dfarule.NewRuleArgs(dStates[dstate], c.(rune))] = dStates[dnext]
		}
	}

	return
}

// DFAStatesMap associates subsets of the NFA state set with the states of the DFA.
type DFAStatesMap map[mapset.Set]utils.State

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
