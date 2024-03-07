package gencpp

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/toplvl"
	"github.com/mebyus/gizmo/ast/tps"
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
	case toplvl.FnTemplate:
		g.TopFunctionTemplate(node.(ast.TopFunctionTemplate))
	case toplvl.TypeTemplate:
		g.TopTypeTemplate(node.(ast.TopTypeTemplate))
	case toplvl.MethodTemplate:
		g.MethodTemplate(node.(ast.MethodTemplate))
	default:
		g.write(fmt.Sprintf("<%s node not implemented>", node.Kind().String()))
		g.nl()
	}
}

func (g *Builder) TopTypeTemplate(top ast.TopTypeTemplate) {
	g.write("template")
	g.templateParams(top.TypeParams)
	g.nl()

	g.TopStructType(top.Name, top.Spec.(ast.StructType))
}

func (g *Builder) TopType(top ast.TopType) {
	switch top.Spec.Kind() {
	case tps.Struct:
		g.TopStructType(top.Name, top.Spec.(ast.StructType))
	case tps.Enum:
		g.TopEnumType(top.Name, top.Spec.(ast.EnumType))
	default:
		g.write(fmt.Sprintf("<top %s type not implemented>", top.Spec.Kind()))
	}
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

func (g *Builder) TopFn(top ast.TopFunctionDefinition) {
	symName := g.symName(top.Definition.Head.Name)
	props := g.meta.Symbols.Props[symName]

	linkName := props.LinkName()
	if linkName != "" {
		g.write(`extern "C" `)
	}
	if top.Definition.Head.Signature.Never {
		g.write("[[noreturn]] ")
	}
	if !props.Export() {
		g.write("static ")
	}

	if linkName != "" {
		top.Definition.Head.Name.Lit = linkName
	}
	g.FunctionDefinition(top.Definition)
}

func (g *Builder) TopConst(top ast.TopConst) {
	g.ConstInit(top.ConstInit)
	g.semi()
	g.nl()
}

func (g *Builder) TopDeclare(top ast.TopFunctionDeclaration) {
	symName := g.symName(top.Declaration.Name)
	props := g.meta.Symbols.Props[symName]

	linkName := props.LinkName()
	extLink := props.ExtLink()
	if extLink || linkName != "" {
		g.write(`extern "C" `)
	}
	if !extLink && !props.Export() {
		g.write("static ")
	}

	if linkName != "" {
		top.Declaration.Name.Lit = linkName
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

func (g *Builder) MethodTemplate(top ast.MethodTemplate) {
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

func (g *Builder) templateReceiverArgs(args []ast.Identifier) {
	g.write("<")
	g.Identifier(args[0])
	for _, param := range args[1:] {
		g.write(", ")
		g.Identifier(param)
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

func (g *Builder) templateParams(params []ast.Identifier) {
	g.write("<")
	g.typeParam(params[0])
	for _, param := range params[1:] {
		g.write(", ")
		g.typeParam(param)
	}
	g.write(">")
}

func (g *Builder) TopFunctionTemplate(top ast.TopFunctionTemplate) {
	g.write("template")
	g.templateParams(top.TypeParams)
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
