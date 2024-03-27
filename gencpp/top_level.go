package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/ir"
)

func (g *Builder) TopLevel(node ast.TopLevel) {
	if g.cfg.SourceLocationComments {
		g.comment("gizmo.source = " + node.Pin().String())
	}

	switch node.Kind() {
	case toplvl.Fn:
		g.TopFn(node.(ast.TopFunctionDefinition))
	case toplvl.Type:
		g.TopType(node.(ast.TopType))
	case toplvl.Const:
		g.TopConst(node.(ast.TopConst))
	case toplvl.Declare:
		g.TopDeclare(node.(ast.TopFunctionDeclaration))
	case toplvl.Var:
		g.TopVar(node.(ast.TopVar))
	case toplvl.Method:
		g.Method(node.(ast.Method))
	case toplvl.Blue:
		g.TopFunctionTemplate(node.(ast.TopBlueprint))
	case toplvl.Proto:
		g.TopTypeTemplate(node.(ast.TopPrototype))
	case toplvl.Pmb:
		g.MethodTemplate(node.(ast.ProtoMethodBlueprint))
	default:
		g.write(fmt.Sprintf("<%s node not implemented>", node.Kind().String()))
		g.nl()
	}
}

func (g *Builder) TopTypeTemplate(top ast.TopPrototype) {
	g.write("template")
	g.templateParams(top.Params)
	g.nl()

	g.TopStructType(top.Name, top.Spec.(ast.StructType))
}

func (g *Builder) TopType(top ast.TopType) {
	switch top.Spec.Kind() {
	case tps.Struct:
		g.TopStructType(top.Name, top.Spec.(ast.StructType))
	case tps.Enum:
		g.TopEnumType(top.Name, top.Spec.(ast.EnumType))
	case tps.Function:
		g.TopFunctionType(top.Name, top.Spec.(ast.FunctionType))
	case tps.Union:
		g.TopUnionType(top.Name, top.Spec.(ast.UnionType))
	default:
		g.write(fmt.Sprintf("<top %s type not implemented>", top.Spec.Kind()))
		g.nl()
	}
}

func (g *Builder) TopUnionType(name ast.Identifier, spec ast.UnionType) {
	g.write("union ")
	g.Identifier(name)
	g.space()

	// methodsKey := g.symName(name)
	g.structFields(spec.Fields, nil)

	g.semi()
	g.nl()
}

func (g *Builder) TopFunctionType(name ast.Identifier, spec ast.FunctionType) {
	signature := spec.Signature
	g.write("typedef ")
	g.functionReturnType(signature.Result)
	g.space()
	g.Identifier(name)
	g.functionParams(signature.Params)
	g.semi()
	g.nl()
}

func (g *Builder) TopStructType(name ast.Identifier, spec ast.StructType) {
	g.write("struct ")
	g.Identifier(name)
	g.space()

	methodsKey := g.symName(name)
	g.structFields(spec.Fields, g.meta.Structs.Methods[methodsKey])

	g.semi()
	g.nl()
}

// Output a function definition of the form (example):
//
//	extern "C" i32 custom_link_name(s: str, k: int) noexcept {
//		return original_source_name(s, k);
//	}
//
// This function with fixed linkname just calls the original function
// with the same arguments
func (g *Builder) topFnLinkNameBind(top ast.TopFunctionDefinition, linkName string, props *ir.PropsRef) {
	g.write(`extern "C" `)
	if top.Definition.Head.Signature.Never {
		g.write("[[noreturn]] ")
	}
	if !props.Export() {
		g.write("static ")
	}

	g.functionReturnType(top.Definition.Head.Signature.Result)
	g.nl()

	g.write(linkName)
	g.functionParams(top.Definition.Head.Signature.Params)
	g.write(" noexcept {")
	g.nl()

	g.inc()
	g.indent()
	if top.Definition.Head.Signature.Result != nil {
		g.write("return ")
	}
	g.Identifier(top.Definition.Head.Name)

	params := top.Definition.Head.Signature.Params
	g.write("(")
	if len(params) != 0 {
		g.Identifier(params[0].Name)
		for _, param := range params[1:] {
			g.write(", ")
			g.Identifier(param.Name)
		}
	}
	g.write(");")
	g.nl()
	g.dec()

	g.write("}")
	g.nl()
}

func (g *Builder) topFnExport(top ast.TopFunctionDefinition, linkName string, props *ir.PropsRef) {
	if top.Definition.Head.Signature.Never {
		g.write("[[noreturn]] ")
	}
	g.write("static ")
	g.FunctionDefinition(top.Definition)
	g.nl()

	g.comment("gizmo.export.bind = " + top.Definition.Head.Name.Lit)
	g.topFnLinkNameBind(top, linkName, props)
}

func (g *Builder) TopFn(top ast.TopFunctionDefinition) {
	symName := g.symName(top.Definition.Head.Name)
	props := g.meta.Symbols.Props[symName]

	linkName := props.LinkName()
	export := props.Export()
	if export && linkName != "" {
		g.topFnExport(top, linkName, props)
		return
	}

	if top.Definition.Head.Signature.Never {
		g.write("[[noreturn]] ")
	}
	g.write("static ")

	g.FunctionDefinition(top.Definition)
}

func (g *Builder) TopConst(top ast.TopConst) {
	g.ConstInit(top.ConstInit)
	g.semi()
	g.nl()
}

func (g *Builder) topDeclareAsmLink(top ast.TopFunctionDeclaration, linkName string) {
	g.comment("gizmo.link.obj = asm")
	g.topDeclareBind(top, linkName)
}

func (g *Builder) topDeclareExtLink(top ast.TopFunctionDeclaration, linkName string) {
	g.comment("gizmo.link.obj = external")
	g.topDeclareBind(top, linkName)
}

func (g *Builder) topDeclareBind(top ast.TopFunctionDeclaration, linkName string) {
	g.comment("gizmo.link.bind = " + top.Declaration.Name.Lit)

	g.write(`extern "C" `)
	if top.Declaration.Signature.Never {
		g.write("[[noreturn]] ")
	}
	g.functionReturnType(top.Declaration.Signature.Result)
	g.nl()

	g.write(linkName)
	g.functionParams(top.Declaration.Signature.Params)
	g.write(" noexcept;")
	g.nl()
	g.nl()

	g.topDeclareExtLinkBind(top, linkName)
}

// Output a function definition of the form (example):
//
//	static i32
//	original_source_name(s: str, k: int) noexcept {
//		return custom_link_name(s, k);
//	}
//
// This function with original name just calls another function
// with fixed linkname giving it the same arguments
func (g *Builder) topDeclareExtLinkBind(top ast.TopFunctionDeclaration, linkName string) {
	if top.Declaration.Signature.Never {
		g.write("[[noreturn]] ")
	}
	g.write("static ")

	g.functionReturnType(top.Declaration.Signature.Result)
	g.nl()

	g.Identifier(top.Declaration.Name)
	g.functionParams(top.Declaration.Signature.Params)
	g.write(" noexcept {")
	g.nl()

	g.inc()
	g.indent()
	if top.Declaration.Signature.Result != nil {
		g.write("return ")
	}
	g.write(linkName)

	params := top.Declaration.Signature.Params
	g.write("(")
	if len(params) != 0 {
		g.Identifier(params[0].Name)
		for _, param := range params[1:] {
			g.write(", ")
			g.Identifier(param.Name)
		}
	}
	g.write(");")
	g.nl()
	g.dec()

	g.write("}")
	g.nl()
}

func (g *Builder) TopDeclare(top ast.TopFunctionDeclaration) {
	symName := g.symName(top.Declaration.Name)
	props := g.meta.Symbols.Props[symName]

	asm := props.Asm()
	external := props.External()
	linkName := props.LinkName()
	if external && linkName != "" {
		g.topDeclareExtLink(top, linkName)
		return
	}
	if asm && linkName != "" {
		g.topDeclareAsmLink(top, linkName)
		return
	}
	if top.Declaration.Signature.Never {
		g.write("[[noreturn]] ")
	}
	if !props.Export() {
		g.write("static ")
	}

	g.FunctionDeclaration(top.Declaration)
	g.semi()
	g.nl()
}

func (g *Builder) TopVar(top ast.TopVar) {
	g.VarInit(top.VarInit)
	g.semi()
	g.nl()
}

func (g *Builder) Method(top ast.Method) {
	g.functionReturnType(top.Signature.Result)
	g.space()
	g.Identifier(top.Receiver)
	g.write("::")
	g.Identifier(top.Name)
	g.functionParams(top.Signature.Params)
	g.write(" noexcept ")
	g.Block(top.Body)
	g.nl()
}

func (g *Builder) MethodTemplate(top ast.ProtoMethodBlueprint) {
	g.write("template")
	g.templateParams(top.TypeParams)
	g.nl()
	g.functionReturnType(top.Signature.Result)
	g.nl()
	g.Identifier(top.Receiver)
	g.templateReceiverArgs(top.TypeParams)
	g.write("::")
	g.Identifier(top.Name)
	g.functionParams(top.Signature.Params)
	g.write(" noexcept ")
	g.Block(top.Body)
	g.nl()
}

func (g *Builder) templateReceiverArgs(args []ast.TypeParam) {
	g.write("<")
	g.Identifier(args[0].Name)
	for _, param := range args[1:] {
		g.write(", ")
		g.Identifier(param.Name)
	}
	g.write(">")
}

func (g *Builder) TopEnumType(name ast.Identifier, spec ast.EnumType) {
	g.write("enum struct ")
	g.Identifier(name)
	g.space()

	if len(spec.Entries) == 0 {
		g.write("{};")
		g.nl()
		return
	}

	g.write("{")
	g.nl()

	g.inc()
	for _, entry := range spec.Entries {
		g.indent()

		g.Identifier(entry.Name)

		if entry.Expression != nil {
			g.write(" = ")
			g.Expression(entry.Expression)
		}

		g.write(",")
		g.nl()
	}
	g.dec()

	g.write("};")
	g.nl()
}

func (g *Builder) typeParam(param ast.Identifier) {
	g.write("typename ")
	g.Identifier(param)
}

func (g *Builder) templateParams(params []ast.TypeParam) {
	g.write("<")
	g.typeParam(params[0].Name)
	for _, param := range params[1:] {
		g.write(", ")
		g.typeParam(param.Name)
	}
	g.write(">")
}

func (g *Builder) TopFunctionTemplate(top ast.TopBlueprint) {
	g.write("template")
	g.templateParams(top.Params)
	g.nl()

	symName := g.symName(top.Name)
	props := g.meta.Symbols.Props[symName]

	if !props.Export() {
		g.write("static ")
	}

	g.FunctionDefinition(ast.FunctionDefinition{
		Head: ast.FunctionDeclaration{
			Signature: top.Signature,
			Name:      top.Name,
		},
		Body: top.Body,
	})
}
