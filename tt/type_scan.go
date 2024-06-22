package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/tt/scp"
	"github.com/mebyus/gizmo/tt/sym"
	"github.com/mebyus/gizmo/tt/typ"
)

// TypeLinkKind describes how one type uses another in its definition.
// In general link type is determined by memory layout of type usage.
// When one type needs to know the size of another to form its
// memory layout it means a direct link between them. Indirect link is
// formed when only pointer of another type is included.
//
// Sources of direct links (from type B to type A):
//   - named type of A
//   - named type of array of A's
//   - struct or union field of A
//   - struct or union field of array of A's
//
// Sources of indirect links (from type B to type A):
//   - named type of chunk of A's
//   - struct or union field of pointer to A
//   - struct or union field of chunk of A's
//   - struct or union field of array pointer to A's
//   - any inclusion with more than one levels of indirections to A
type TypeLinkKind uint8

const (
	// Zero value of TypeLinkKind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where TypeLinkKind is left unspecified.
	linkEmpty TypeLinkKind = iota

	// Direct inclusion. In the example below type B directly
	// includes type A:
	//
	//	type A struct {
	//		name: str,
	//	}
	//
	//	type B struct {
	//		a:     A,
	//		title: str,
	//	}
	//
	// When there are combinations of direct and indirect inclusions inside one type,
	// then direct inclusion takes precedence over indirect one:
	//
	//	type B struct {
	//		a:     A,
	//		a2:    *A,
	//		title: str,
	//	}
	//
	// Array inclusions also lead to direct links:
	//
	//	type B struct {
	//		a:     [4]A,
	//		title: str,
	//	}
	//
	linkDirect

	// Inclusion via indirect means (pointer, chunk, array pointer).
	// In the example below type B includes type A via pointer:
	//
	//	type A struct {
	//		name: str,
	//	}
	//
	//	type B struct {
	//		a:     *A,
	//		title: str,
	//	}
	//
	// Indirect link is also formed via chunks:
	//
	//	type B struct {
	//		a:     []A,
	//		title: str,
	//	}
	//
	// When there is an inclusion with several levels of indirection it
	// also counts as indirect link:
	//
	//	type B struct {
	//		a:     []*A,
	//		title: str,
	//	}
	//
	linkIndirect
)

var linkText = [...]string{
	linkEmpty: "<nil>",

	linkDirect:   "direct",
	linkIndirect: "indirect",
}

func (k TypeLinkKind) String() string {
	return linkText[k]
}

type TypeLink struct {
	Symbol *Symbol
	Kind   TypeLinkKind
}

type TypeLinkSet map[*Symbol]TypeLinkKind

func NewTypeLinkSet() TypeLinkSet {
	return make(TypeLinkSet)
}

func (s TypeLinkSet) Add(symbol *Symbol, kind TypeLinkKind) {
	k, ok := s[symbol]
	if !ok {
		s[symbol] = kind
		return
	}
	if k == linkIndirect && kind == linkDirect {
		s[symbol] = linkDirect
		return
	}
}

func (s TypeLinkSet) Elems() []TypeLink {
	if len(s) == 0 {
		return nil
	}

	elems := make([]TypeLink, 0, len(s))
	for symbol, kind := range s {
		elems = append(elems, TypeLink{
			Symbol: symbol,
			Kind:   kind,
		})
	}
	return elems
}

// TypeContext carries information and objects which are needed to perform
// top-level unit types hoisting, cycled definitions detection, etc.
type TypeContext struct {
	links TypeLinkSet

	members MembersList

	// keeps track of current link type kind when descending/ascending
	// nested type specifiers
	kind TypeLinkKind
}

func NewTypeContext() *TypeContext {
	return &TypeContext{
		links: NewTypeLinkSet(),
		kind:  linkDirect,
	}
}

func (c *TypeContext) push(kind TypeLinkKind) TypeLinkKind {
	k := c.kind
	if k == linkDirect {
		// only direct link can be changed by descending nested type specifiers
		c.kind = kind
	}
	return k
}

func (m *Merger) shallowScanTypes() error {
	g := NewTypeGraphBuilder(len(m.types))
	// TODO: graph based scanning
	for _, s := range m.types {
		def := s.Def.(*TempTypeDef)
		ctx := NewTypeContext()
		t, err := m.shallowScanType(ctx, def)
		if err != nil {
			return err
		}

		if ctx.links[s] == linkDirect {
			return fmt.Errorf("%s: symbol \"%s\" directly references itself", s.Pos.String(), s.Name)
		}

		g.Add(s, ctx.links.Elems())

		s.Def = t
	}

	g.Scan()
	return fmt.Errorf("remove this error after finishing type graph")

	// return nil
}

// performs preliminary top-level type definition scan in order to obtain data
// necessary for constructing dependency graph between the named types
func (m *Merger) shallowScanType(ctx *TypeContext, def *TempTypeDef) (*Type, error) {
	name := def.top.Name.Lit
	kind := def.top.Spec.Kind()

	var err error
	// TODO: mechanism for other types dependency shallow scanning and hoisting
	switch kind {
	case tps.Name:
		err = m.shallowScanNamedType(ctx, def.top.Spec.(ast.TypeName))
	case tps.Struct:
		err = m.shallowScanStructType(ctx, def.top.Spec.(ast.StructType), def.methods)
	case tps.Enum:
		// TODO: implement enum scan
	case tps.Bag:
		// TODO: implement bag scan
	default:
		panic(fmt.Sprintf("unexpected %s type specifier", kind.String()))
	}

	if err != nil {
		return nil, err
	}

	return &Type{
		Kind: typ.Named,
		Name: name,
	}, nil
}

func (m *Merger) shallowScanStructType(ctx *TypeContext, spec ast.StructType, methods []ast.Method) error {
	ctx.members.Init(len(spec.Fields) + len(methods))

	for _, field := range spec.Fields {
		err := m.shallowScanStructField(ctx, field)
		if err != nil {
			return err
		}
	}

	for _, method := range methods {
		err := m.shallowScanMethod(ctx, method)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Merger) shallowScanMethod(ctx *TypeContext, method ast.Method) error {
	name := method.Name.Lit
	pos := method.Name.Pos

	r := ctx.members.Find(name)
	if r != nil {
		return fmt.Errorf("%s: field with name \"%s\" is already present in struct", pos.String(), name)
	}

	ctx.members.Add(Member{
		Pos:  pos,
		Name: name,
		Kind: MemberMethod,
	})

	return nil
}

func (m *Merger) shallowScanStructField(ctx *TypeContext, field ast.FieldDefinition) error {
	name := field.Name.Lit
	pos := field.Name.Pos

	r := ctx.members.Find(name)
	if r != nil {
		return fmt.Errorf("%s: field with name \"%s\" is already present in struct", pos.String(), name)
	}

	err := m.shallowScanTypeSpecifier(ctx, field.Type)
	if err != nil {
		return err
	}

	ctx.members.Add(Member{
		Pos:  pos,
		Name: name,
		Kind: MemberField,
	})

	return nil
}

func (m *Merger) shallowScanTypeSpecifier(ctx *TypeContext, spec ast.TypeSpecifier) error {
	switch spec.Kind() {
	case tps.Name:
		return m.shallowScanNamedType(ctx, spec.(ast.TypeName))
	case tps.Pointer:
		return m.shallowScanPointerType(ctx, spec.(ast.PointerType))
	case tps.Chunk:
		return m.shallowScanChunkType(ctx, spec.(ast.ChunkType))
	case tps.Array:
		return m.shallowScanArrayType(ctx, spec.(ast.ArrayType))
	case tps.ArrayPointer:
		return m.shallowScanArrayPointerType(ctx, spec.(ast.ArrayPointerType))
	case tps.Struct:
		return fmt.Errorf("%s: anonymous nested structs are not allowed", spec.Pin().String())
	case tps.Enum:
		return fmt.Errorf("%s: anonymous nested enums are not allowed", spec.Pin().String())
	case tps.Union:
		return fmt.Errorf("%s: anonymous nested unions are not allowed", spec.Pin().String())
	default:
		panic(fmt.Sprintf("not implemented for %s", spec.Kind().String()))
	}
}

func (m *Merger) shallowScanArrayType(ctx *TypeContext, spec ast.ArrayType) error {
	// save previous link kind for later
	k := ctx.push(linkDirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.ElemType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.kind = k
	return nil
}

func (m *Merger) shallowScanArrayPointerType(ctx *TypeContext, spec ast.ArrayPointerType) error {
	// save previous link kind for later
	k := ctx.push(linkIndirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.ElemType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.kind = k
	return nil
}

func (m *Merger) shallowScanChunkType(ctx *TypeContext, spec ast.ChunkType) error {
	// save previous link kind for later
	k := ctx.push(linkIndirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.ElemType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.kind = k
	return nil
}

func (m *Merger) shallowScanPointerType(ctx *TypeContext, spec ast.PointerType) error {
	// save previous link kind for later
	k := ctx.push(linkIndirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.RefType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.kind = k
	return nil
}

func (m *Merger) shallowScanNamedType(ctx *TypeContext, spec ast.TypeName) error {
	name := spec.Name.Lit
	pos := spec.Name.Pos

	// in shallow scan ref count increase is not needed
	symbol := m.unit.Scope.lookup(name, pos.Num)
	if symbol == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if symbol.Kind != sym.Type {
		return fmt.Errorf("%s: symbol \"%s\" is not a type", pos.String(), name)
	}

	if symbol.Scope.Kind == scp.Unit {
		ctx.links.Add(symbol, ctx.kind)
	}
	return nil
}
