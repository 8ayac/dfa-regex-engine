// Package nfabuilder implements some structures and functions to construct NFA.
package nfabuilder

import (
	"github.com/8ayac/dfa-regex-engine/nfa"
	"github.com/8ayac/dfa-regex-engine/nfa/nfarule"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// Fragment represents a fragment of NFA to construct a larger NFA.
type Fragment struct {
	I     utils.State // initial state
	F     mapset.Set  // accept states
	Rules nfarule.RuleMap
}

// NewFragment returns a new Fragment.
func NewFragment() *Fragment {
	return &Fragment{
		I:     utils.NewState(0),
		F:     mapset.NewSet(),
		Rules: nfarule.RuleMap{},
	}
}

// AddRule add a new transition rule to the Fragment.
// Rule concept: State(from) -->[Symbol(c)]--> State(next)
func (frg *Fragment) AddRule(from utils.State, c rune, next utils.State) {
	r := frg.Rules
	_, ok := r[nfarule.NewRuleArgs(from, c)]
	if ok {
		r[nfarule.NewRuleArgs(from, c)].Add(next)
	} else {
		r[nfarule.NewRuleArgs(from, c)] = mapset.NewSet(next)
	}
}

// CreateSkeleton returns a nfa fragment which has
// same transition rule as original fragment has.
// The initial state and accept state is set to default.
func (frg *Fragment) CreateSkeleton() (Skeleton *Fragment) {
	Skeleton = NewFragment()
	Skeleton.Rules = frg.Rules
	return
}

// MergeRule returns a new NFA fragment into which the
// transition rules of original fragment and the fragment
// given in the argument are merged.
func (frg *Fragment) MergeRule(frg2 *Fragment) (synthesizedFrg *Fragment) {
	synthesizedFrg = frg.CreateSkeleton()
	for k, v := range frg2.Rules {
		_, ok := synthesizedFrg.Rules[k]
		if !ok {
			synthesizedFrg.Rules[k] = mapset.NewSet()
		}
		synthesizedFrg.Rules[k] = synthesizedFrg.Rules[k].Union(v)
	}
	return
}

// Build converts NFA fragments into a NFA, and returns it.
func (frg *Fragment) Build() *nfa.NFA {
	return nfa.NewNFA(frg.I, frg.F, frg.Rules)
}
