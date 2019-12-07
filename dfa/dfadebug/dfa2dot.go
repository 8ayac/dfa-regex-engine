// Package dfadebug provides some features to debug a dfa.
package dfadebug

import (
	"fmt"
	"github.com/8ayac/dfa-regex-engine/dfa"
	"github.com/awalterschulze/gographviz"
	"os"
)

type CommonNodeAttrs map[string]string

func NewCommonNodeAttrs() CommonNodeAttrs {
	return CommonNodeAttrs{
		"fontname": "meiryo",
		"fontsize": "18",
	}
}

type CommonEdgeAttrs map[string]string

func NewCommonEdgeAttrs() CommonEdgeAttrs {
	return CommonEdgeAttrs{
		"fontname":   "meiryo",
		"fontsize":   "18",
		"len":        "1.5",
		"labelfloat": "false",
	}
}

// DFA2dot outputs
func DFA2dot(d dfa.DFA) {
	const GRAPH_NAME = "DFA"
	g := gographviz.NewGraph()

	// General
	_ = g.SetName(GRAPH_NAME)
	_ = g.SetDir(true)
	_ = g.AddAttr(GRAPH_NAME, "rankdir", "LR")

	// For initial state
	dummyAttrs := NewCommonNodeAttrs()
	dummyAttrs["shape"] = "point"
	_ = g.AddNode(GRAPH_NAME, "\"\"", dummyAttrs)
	_ = g.AddNode(GRAPH_NAME, d.Init.String(), NewCommonNodeAttrs())

	initEdgeAttrs := NewCommonEdgeAttrs()
	initEdgeAttrs["len"] = "2"
	_ = g.AddEdge("\"\"", d.Init.String(), true, initEdgeAttrs)

	// For transition rules
	rules := d.Rules
	for _, v := range rules {
		attrs := NewCommonNodeAttrs()
		if d.Accepts.Contains(v) {
			attrs["shape"] = "doublecircle"
		}
		_ = g.AddNode(GRAPH_NAME, v.String(), attrs)
	}

	for arg, dst := range rules {
		attrs := NewCommonEdgeAttrs()
		attrs["label"] = fmt.Sprintf("\"'%c'\"", arg.C)
		_ = g.AddEdge(arg.From.String(), dst.String(), true, attrs)
	}

	// Output DOT
	file, err := os.Create(`dfa.dot`)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(g.String()))
	if err != nil {
		panic(err)
	}
}
