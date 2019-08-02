// Package dfa implements Deterministic Finite Automaton(DFA).
package dfa

import (
	"github.com/8ayac/dfa-regex-engine/dfa/dfarule"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// DFA represents a Deterministic Finite Automaton.
type DFA struct {
	Init    utils.State     // initial state
	Accepts mapset.Set      // accepts states
	Rules   dfarule.RuleMap // transition function
}

// NewDFA returns a new dfa.
func NewDFA(init utils.State, accepts mapset.Set, rules dfarule.RuleMap) *DFA {
	return &DFA{
		Init:    init,
		Accepts: accepts,
		Rules:   rules,
	}
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
	r.cur = d.Init
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
	accepts := r.d.Accepts
	if accepts.Contains(r.cur) {
		return true
	}
	return false
}

// Matching returns whether the string given is accepted (or not) by
// simulating the all transitions.
func (r *Runtime) Matching(str string) bool {
	for _, c := range []rune(str) {
		if !r.transit(c) {
			return false // if the transition failed, the input "str" is rejected.
		}
	}
	return r.isAccept()
}
