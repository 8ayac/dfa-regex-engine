// Package dfarule implements the transition function of DFA.
package dfarule

import (
	"fmt"
	"reflect"

	"github.com/8ayac/dfa-regex-engine/utils"
)

// RuleMap represents a transition function of d.
// The key is a pair like "(from state, input symbol)".
// The value is a destination state when "input symbol"
// is received in "from state".
type RuleMap map[RuleArgs]utils.State

func (r RuleMap) String() string {
	s := ""

	keys := reflect.ValueOf(r).MapKeys()
	for i, k := range keys {
		from := k.FieldByName("From").Interface().(utils.State)
		c := k.FieldByName("C").Interface().(rune)
		dst := r[NewRuleArgs(from, c)]
		s += fmt.Sprintf("s%d\t--['%c']-->\t%s", from, c, dst)
		if i+1 < len(keys) {
			s += "\n"
		}
	}
	return s
}

// RuleArgs is a key for the map as transition function of d.
type RuleArgs struct {
	From utils.State // from state
	C    rune        // input symbol
}

// NewRuleArgs returns a new RuleArgs.
func NewRuleArgs(from utils.State, in rune) RuleArgs {
	return RuleArgs{
		From: from,
		C:    in,
	}
}
