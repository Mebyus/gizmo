package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
)

// TempFnDef is temporary function symbol definition in the form of AST node.
// Used as an intermidiate transfering data structure to carry
// node data between phase 1 and phase 2 merging.
type TempFnDef struct {
	nodeSymDef

	top ast.TopFun
}

func NewTempFnDef(top ast.TopFun) TempFnDef {
	return TempFnDef{top: top}
}

// Explicit interface implementation check
var _ SymDef = TempFnDef{}

type TempTypeDef struct {
	nodeSymDef

	methods []ast.Method

	// Methods map. Maps method name to index inside methods slice.
	mm map[string]int

	top ast.TopType
}

func NewTempTypeDef(top ast.TopType) *TempTypeDef {
	return &TempTypeDef{top: top}
}

func (d *TempTypeDef) addMethod(method ast.Method) error {
	name := method.Name.Lit

	if len(d.methods) == 0 {
		d.mm = make(map[string]int)
	} else {
		_, ok := d.mm[name]
		if ok {
			return errMultMethodDef(name, method.Pin())
		}
	}

	d.mm[name] = len(d.methods)
	d.methods = append(d.methods, method)
	return nil
}

func errMultMethodDef(name string, pos source.Pos) error {
	return fmt.Errorf("%s: multiple definitions of method \"%s\"", pos.String(), name)
}

// Explicit interface implementation check
var _ SymDef = &TempTypeDef{}

type TempConstDef struct {
	nodeSymDef

	top ast.TopCon
}

func NewTempConstDef(top ast.TopCon) TempConstDef {
	return TempConstDef{top: top}
}
