package ir

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/prv"
)

type Meta struct {
	Structs StructsMap

	Symbols SymbolsMap
}

type SymbolsMap struct {
	// Maps symbol name into its props reference
	Props map[string]*PropsRef
}

// Maps prop key into its value for a specific symbol
type PropsRef struct {
	val map[string]ast.PropValue

	// what symbol name should be used when linking compiled objects together
	linkName string

	// symbol should be exported from compiled object (available for external linkage)
	export bool

	// should be linked from external object when building
	external bool

	// should be found in asm code inside the same unit
	asm bool

	// symbol marked as test function
	test bool
}

// NewPropsRef argument must contain at least one element
func NewPropsRef(props []ast.Prop) *PropsRef {
	r := &PropsRef{val: make(map[string]ast.PropValue)}

	for _, p := range props {
		if p.Key != "" {
			switch p.Key {
			case "link.name":
				r.linkName = propValueString(p.Value)
			default:
				r.val[p.Key] = p.Value
			}
			continue
		}

		if p.Value.Kind() != prv.Tag {
			panic(fmt.Sprintf("unexpected prop value: %s", p.Value.Kind().String()))
		}

		tags := p.Value.(ast.PropValueTags).Tags
		for _, tag := range tags {
			switch tag.Lit {
			case "export":
				r.export = true
			case "external":
				r.external = true
			case "test":
				r.test = true
			case "asm":
				r.asm = true
			}
		}
	}

	if len(r.val) == 0 {
		r.val = nil
	}
	return r
}

func (p *PropsRef) Asm() bool {
	if p == nil {
		return false
	}
	return p.asm
}

func (p *PropsRef) Test() bool {
	if p == nil {
		return false
	}
	return p.test
}

func (p *PropsRef) Export() bool {
	if p == nil {
		return false
	}
	return p.export
}

func (p *PropsRef) External() bool {
	if p == nil {
		return false
	}
	return p.external
}

func (p *PropsRef) LinkName() string {
	if p == nil {
		return ""
	}
	return p.linkName
}

func propValueString(v ast.PropValue) string {
	if v == nil {
		return ""
	}
	if v.Kind() != prv.String {
		return ""
	}
	return v.(ast.PropValueString).Val
}
