package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/tpk"
	"github.com/mebyus/gizmo/source"
)

// Member represents something that can be selected from expression by "." chain.
type Member interface {
	Member()
}

// This is dummy implementation of Member interface.
//
// Used for embedding into other (non-dummy) member nodes.
type NodeM struct{}

func (NodeM) Member() {}

type Field struct {
	NodeM

	// position where this field is defined.
	Pos source.Pos

	// Field name. Cannot be empty.
	Name string

	// Type of value which this member holds.
	Type *Type
}

type Method struct {
	NodeM

	// Symbol which defines the method.
	Symbol *Symbol
}

// Member tries to find member by name which belongs
// to this type. Returns error if member not found.
// Panics on types which cannot have members.
func (t *Type) Member(name ast.Identifier) (Member, error) {
	if t.Kind != tpk.Custom {
		return nil, fmt.Errorf("%s: %s type cannot have members", name.Pos, t.Kind)
	}

	// try to find method by name
	def := t.Def.(CustomTypeDef)
	for i := range len(def.Methods) {
		s := def.Methods[i]
		methodSymbolName := def.Symbol.Name + "." + name.Lit
		if s.Name == methodSymbolName {
			return Method{Symbol: s}, nil
		}
	}

	b := def.Base
	if b.Kind != tpk.Struct {
		return nil, fmt.Errorf("%s: %s type cannot have members", name.Pos, t.Kind)
	}

	bdef := b.Def.(*StructTypeDef)
	for i := range len(bdef.Fields) {
		field := bdef.Fields[i]
		if field.Name == name.Lit {
			return field, nil
		}
	}

	return nil, fmt.Errorf("%s: type %s does not have field or method \"%s\"", name.Pos, t.Kind, name.Lit)
}
