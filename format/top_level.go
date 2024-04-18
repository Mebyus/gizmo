package format

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
)

func (f *Formatter) TopLevel(node ast.TopLevel) {
	switch node.Kind() {
	case toplvl.Fn:
		f.TopFn(node.(ast.TopFunctionDefinition))
	case toplvl.Type:
		f.TopType(node.(ast.TopType))
	case toplvl.Const:
		f.TopConst(node.(ast.TopConst))
	case toplvl.Declare:
		f.TopDeclare(node.(ast.TopFunctionDeclaration))
	case toplvl.Var:
		f.TopVar(node.(ast.TopVar))
	case toplvl.Method:
		f.Method(node.(ast.Method))
	default:
		f.g.Str(fmt.Sprintf("<%s node not implemented>\n", node.Kind().String()))
	}
}

func (f *Formatter) TopFn(top ast.TopFunctionDefinition) {
	// f.FunctionDefinition(top.Definition)
}

func (f *Formatter) TopConst(top ast.TopConst) {
	// f.ConstInit(top.ConstInit)
}

func (f *Formatter) TopDeclare(top ast.TopFunctionDeclaration) {
	// f.FunctionDeclaration(top.Declaration)
}

func (f *Formatter) TopVar(top ast.TopVar) {
	// f.VarInit(top.VarInit)
}

func (f *Formatter) Method(top ast.Method) {
	f.g.Str("<method>\n")
}

func (f *Formatter) TopType(top ast.TopType) {

}
