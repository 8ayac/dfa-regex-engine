// Package dfa implements Deterministic Finite Automaton(DFA).
package dfa

import (
	"github.com/8ayac/dfa-regex-engine/dfa/dfarule"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// DFA represents a Deterministic Finite Automaton.
type DFA struct {
	I     utils.State     // initial state
	F     mapset.Set      // accepts states
	Rules dfarule.RuleMap // transition function
}

// NewDFA returns a new dfa.
func NewDFA(init utils.State, accepts mapset.Set, rules dfarule.RuleMap) *DFA {
	return &DFA{
		I:     init,
		F:     accepts,
		Rules: rules,
	}
}

// Minimize minimizes the DFA.
func (dfa *DFA) Minimize() {
	states := mapset.NewSet(dfa.I)
	for _, v := range dfa.Rules {
		states.Add(v)
	}
	n := states.N()

	eqMap := map[utils.State]utils.State{}
	for i := 0; i < n; i++ {
		q1 := utils.NewState(i)
		for j := i + 1; j < n; j++ {
			q2 := utils.NewState(j)
			if !dfa.isEquivalent(q1, q2) {
				continue
			}
			if _, ok := eqMap[q2]; ok {
				continue
			}
			dfa.mergeState(q1, q2)
		}
	}
}

func (dfa *DFA) replaceState(to, from utils.State) {
	rules := dfa.Rules
	for arg, dst := range rules {
		if dst == from {
			rules[arg] = to
		}
	}
}

func (dfa *DFA) deleteState(q utils.State) {
	rules := dfa.Rules
	for arg := range rules {
		if arg.From == q {
			delete(rules, arg)
		}
	}
}

func (dfa *DFA) mergeState(to, from utils.State) {
	dfa.replaceState(to, from)
	dfa.deleteState(from)
}

func (dfa *DFA) isEquivalent(q1, q2 utils.State) bool {
	if !((dfa.F.Contains(q1) && dfa.F.Contains(q2)) ||
		(!dfa.F.Contains(q1) && !dfa.F.Contains(q2))) {
		return false
	}

	rules := dfa.Rules
	for k := range rules {
		if k.From != q1 {
			continue
		}
		if rules[dfarule.NewRuleArgs(q1, k.C)] != rules[dfarule.NewRuleArgs(q2, k.C)] {
			return false
		}
	}
	return true
}

// Runtime has a pointer to d and saves current state for
// simulating d transitions.
type Runtime struct {
	d   *DFA
	cur utils.State
}

// GetRuntime returns a new Runtime for simulating d transitions.
func (dfa *DFA) GetRuntime() *Runtime {
	return NewRuntime(dfa)
}

// NewRuntime returns a new runtime for DFA.
func NewRuntime(d *DFA) (r *Runtime) {
	r = &Runtime{
		d: d,
	}
	r.cur = d.I
	return
}

// transit execute a transition with a symbol, and returns whether
// the transition is success (or not).
func (r *Runtime) transit(c rune) bool {
	key := dfarule.NewRuleArgs(r.cur, c)
	_, ok := r.d.Rules[key]
	if ok {
		r.cur = r.d.Rules[key]
		return true
	}
	return false
}

// isAccept returns whether current status is in accept states.
func (r *Runtime) isAccept() bool {
	accepts := r.d.F
	if accepts.Contains(r.cur) {
		return true
	}
	return false
}

// Matching returns whether the string given is accepted (or not) by
// simulating the all transitions.
func (r *Runtime) Matching(str string) bool {
	r.cur = r.d.I
	for _, c := range []rune(str) {
		if !r.transit(c) {
			return false // if the transition failed, the input "str" is rejected.
		}
	}
	return r.isAccept()
}
