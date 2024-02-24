package ir

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/prv"
)

type Meta struct {
	Structs StructsMap

	Symbols SymbolsMap
}

type SymbolsMap struct {
	// Maps symbol name into its props reference
	Props map[string]PropsRef
}

// Maps prop key into its value for a specific symbol
type PropsRef map[string]ast.PropValue

// NewPropsRef argument must contain at least one element
func NewPropsRef(props []ast.Prop) PropsRef {
	r := make(PropsRef, len(props))
	for _, p := range props {
		r[p.Key] = p.Value
	}
	return r
}

func (p PropsRef) Export() bool {
	v := p["export"]
	if v == nil {
		return false
	}
	if v.Kind() != prv.Bool {
		return false
	}
	return v.(ast.PropValueBool).Val
}

func (p PropsRef) ExtLink() bool {
	v := p["linkage"]
	if v == nil {
		return false
	}
	if v.Kind() != prv.String {
		return false
	}
	return v.(ast.PropValueString).Val == "external"
}

func (p PropsRef) LinkName() string {
	v := p["link.name"]
	if v == nil {
		return ""
	}
	if v.Kind() != prv.String {
		return ""
	}
	return v.(ast.PropValueString).Val
}
