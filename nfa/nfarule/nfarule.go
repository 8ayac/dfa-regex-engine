// Package nfarule implements the transition function of NFA.
package nfarule

import (
	"fmt"
	"reflect"

	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// RuleMap represents a transition function of NFA.
// The key is a pair like "(from state, input symbol)".
// The value is a set of transition destination states
// when "input symbol" is received in "from state".
type RuleMap map[RuleArgs]mapset.Set

func (r RuleMap) String() string {
	s := ""

	keys := reflect.ValueOf(r).MapKeys()
	for i, k := range keys {
		from := k.FieldByName("From").Interface().(utils.State)
		c := k.FieldByName("C").Interface().(rune)
		dst := r[NewRuleArgs(from, c)]
		s += fmt.Sprintf("%s\t--['%c']-->\t%s", from, c, dst)
		if i+1 < len(keys) {
			s += "\n"
		}
	}
	return s
}

// RuleArgs is a key for the map as transition function of NFA.
type RuleArgs struct {
	From utils.State // from state
	C    rune        // input symbol
}

// NewRuleArgs returns a new RuleArgs.
func NewRuleArgs(from utils.State, c rune) RuleArgs {
	return RuleArgs{
		From: from,
		C:    c,
	}
}
