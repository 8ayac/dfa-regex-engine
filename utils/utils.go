// Package utils contains utility types and functions for dfa-regex-engine.
package utils

import "fmt"

// State represents a state including in DFA or NFA.
// It has its number. The number can NOT be duplicate in same DFA or NFA.
// Basically, the number is set incrementally.
type State struct {
	N int
}

// NewState returns a new state with its number set.
func NewState(n int) State {
	return State{
		N: n,
	}
}

func (s State) String() string {
	return fmt.Sprintf("q%d", s.N)
}

// Context has a number which is basically used to create incremental stuff.
// Example incremental stuff: state number(q0, q1, q2)
type Context struct {
	N int
}

// NewContext returns a new Context.
// The default value of N is -1.
func NewContext() *Context {
	return &Context{
		N: -1,
	}
}

// Increment add 1 to N which held in Context struct,
// and returns the number.
func (ctx *Context) Increment() int {
	ctx.N++
	return ctx.N
}
