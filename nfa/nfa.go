// Package nfa implements Non-Deterministic Finite Automaton(NFA).
package nfa

import (
	"github.com/8ayac/dfa-regex-engine/nfa/nfarule"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// NFA represents a Non-Deterministic Finite Automaton.
type NFA struct {
	Init    utils.State     // initial state
	Accepts mapset.Set      // accept states
	Rules   nfarule.RuleMap // transition function
}

// NewNFA returns a new NFA.
func NewNFA(init utils.State, accepts mapset.Set, rules nfarule.RuleMap) *NFA {
	return &NFA{
		Init:    init,
		Accepts: accepts,
		Rules:   rules,
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
	if nfa.Accepts.IsSubset(nfa.epsilonClosure(nfa.Init)) {
		nfa.Accepts.Add(nfa.Init)
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
func (nfa *NFA) epsilonExpand(state utils.State, symbol rune) (dst mapset.Set) {
	dst = mapset.NewSet()

	orgDst, ok := nfa.CalcDst(state, symbol)
	if !ok {
		return
	}

	for q := range orgDst.Iter() {
		e := nfa.epsilonClosure(q.(utils.State))
		dst = dst.Union(e)
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
