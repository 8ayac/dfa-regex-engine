// Package node implements some AST nodes.
package node

import (
	"fmt"

	"github.com/8ayac/dfa-regex-engine/nfa/nfabuilder"
	"github.com/8ayac/dfa-regex-engine/utils"
	mapset "github.com/8ayac/golang-set"
)

// String to identify the type of Node.
const (
	TypeCharacter = "Character"
	TypeUnion     = "Union"
	TypeConcat    = "Concat"
	TypeStar      = "Star"
)

// Node is the interface Node implements.
type Node interface {
	// SubtreeString returns a string to which converts
	// a subtree with Node at the top.
	SubtreeString() string

	// Assemble returns a NFA fragment assembled with a Node.
	Assemble(*utils.Context) *nfabuilder.Fragment
}

// Character represents the Character node.
type Character struct {
	Ty string
	V  rune
}

func (c *Character) String() string {
	return c.SubtreeString()
}

// NewCharacter returns a new Character node.
func NewCharacter(r rune) *Character {
	return &Character{
		Ty: TypeCharacter,
		V:  r,
	}
}

/*
Assemble returns a NFA fragment assembled with Character node.
The fragment assembled from a Character node is like below:
	q1(Initial State) -- [Character.V] --> q2(Accept state)
*/
func (c *Character) Assemble(ctx *utils.Context) *nfabuilder.Fragment {
	newFrg := nfabuilder.NewFragment()

	q1 := utils.NewState(ctx.Increment())
	q2 := utils.NewState(ctx.Increment())
	newFrg.AddRule(q1, c.V, q2)

	newFrg.I = q1
	newFrg.F = newFrg.F.Union(mapset.NewSet(q2))

	return newFrg
}

// SubtreeString returns a string to which converts
// a subtree with the Character node at the top.
func (c *Character) SubtreeString() string {
	return fmt.Sprintf("\x1b[32m%s('%s')\x1b[32m", c.Ty, string(c.V))
}

// Union represents the Union node.
type Union struct {
	Ty   string
	Ope1 Node
	Ope2 Node
}

func (u *Union) String() string {
	return u.SubtreeString()
}

// NewUnion returns a new Union node.
func NewUnion(ope1, ope2 Node) *Union {
	return &Union{
		Ty:   TypeUnion,
		Ope1: ope1,
		Ope2: ope2,
	}
}

/*
Assemble returns a NFA fragment assembled with Union node.
The fragment assembled from a Union node is like below:

	I'(new initial state) -- ['ε'] --> frg1
    	              	 `- ['ε'] --> frg2

	+ frg1(fragment assembled with Union.Ope1): I1 -- [???] --> F1
	+ frg2(fragment assembled with Union.Ope2): I2 -- [???] --> F2
*/
func (u *Union) Assemble(ctx *utils.Context) *nfabuilder.Fragment {
	frg1 := u.Ope1.Assemble(ctx)
	frg2 := u.Ope2.Assemble(ctx)

	newFrg := frg1.MergeRule(frg2)
	newFrg.I = utils.NewState(ctx.Increment())
	newFrg.F = newFrg.F.Union(frg1.F)
	newFrg.F = newFrg.F.Union(frg2.F)
	newFrg.AddRule(newFrg.I, 'ε', frg1.I)
	newFrg.AddRule(newFrg.I, 'ε', frg2.I)

	return newFrg
}

// SubtreeString returns a string to which converts
// a subtree with the Union node at the top.
func (u *Union) SubtreeString() string {
	return fmt.Sprintf("\x1b[36m%s(%s, %s\x1b[36m)\x1b[0m", u.Ty, u.Ope1.SubtreeString(), u.Ope2.SubtreeString())
}

// Concat represents the Concat node.
type Concat struct {
	Ty   string
	Ope1 Node
	Ope2 Node
}

func (c *Concat) String() string {
	return c.SubtreeString()
}

// NewConcat returns a new Concat node.
func NewConcat(ope1, ope2 Node) *Concat {
	return &Concat{
		Ty:   TypeConcat,
		Ope1: ope1,
		Ope2: ope2,
	}
}

/*
Assemble returns a NFA fragment assembled with Concat node.
The fragment assembled from a Concat node is like below:

	frg1 -- ['ε']　--> frg2

	+ frg1(fragment assembled with Concat.Ope1): I1 -- [???] --> F1
	+ frg2(fragment assembled with Concat.Ope2): I2 -- [???] --> F2
*/
func (c *Concat) Assemble(ctx *utils.Context) *nfabuilder.Fragment {
	frg1 := c.Ope1.Assemble(ctx)
	frg2 := c.Ope2.Assemble(ctx)

	newFrg := frg1.MergeRule(frg2)
	newFrg.I = frg1.I
	newFrg.F = newFrg.F.Union(frg2.F)
	for q := range frg1.F.Iter() {
		newFrg.AddRule(q.(utils.State), 'ε', frg2.I)
	}

	return newFrg
}

// SubtreeString returns a string to which converts
// a subtree with the Concat node at the top.
func (c *Concat) SubtreeString() string {
	return fmt.Sprintf("\x1b[31m%s(%s, %s\x1b[31m)\x1b[0m", c.Ty, c.Ope1.SubtreeString(), c.Ope2.SubtreeString())
}

// Star represents the Star node.
type Star struct {
	Ty  string
	Ope Node
}

func (s *Star) String() string {
	return s.SubtreeString()
}

// NewStar returns a new Star node.
func NewStar(ope Node) *Star {
	return &Star{
		Ty:  TypeStar,
		Ope: ope,
	}
}

/*
Assemble returns a NFA fragment assembled with Star node.
The fragment assembled from a Star node is like below:

	(new state1) -- ['ε'] --> I1 -----> F1 -- ['ε'] --> (new state2)
	   \					  ↑--['ε']-´						 ↑
		\														 /
		 `-------------------------['ε']------------------------´

	+ frg1(fragment assembled with Ope): I1 -- [???] --> F1

Note: Accept states of new fragment is "(new state1)", "(new state2)" and "I1".
*/
func (s *Star) Assemble(ctx *utils.Context) *nfabuilder.Fragment {
	orgFrg := s.Ope.Assemble(ctx)
	newFrg := orgFrg.CreateSkeleton()
	newState1 := utils.NewState(ctx.Increment())
	newState2 := utils.NewState(ctx.Increment())

	newFrg.I = newState1
	newFrg.F.Add(orgFrg.I)
	newFrg.F.Add(newState2)

	newFrg.AddRule(newState1, 'ε', newState2)
	newFrg.AddRule(newState1, 'ε', orgFrg.I)
	for q := range orgFrg.F.Iter() {
		newFrg.AddRule(q.(utils.State), 'ε', newState2)
		newFrg.AddRule(q.(utils.State), 'ε', orgFrg.I)
	}

	return newFrg
}

// SubtreeString returns a string to which converts
// a subtree with the Star node at the top.
func (s *Star) SubtreeString() string {
	return fmt.Sprintf("\x1b[33m%s(%s\x1b[33m)\x1b[0m", s.Ty, s.Ope.SubtreeString())
}
