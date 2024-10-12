package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/stg/scp"
)

// LinkKind describes how one type uses another in its definition.
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
type LinkKind uint8

const (
	// Zero value of TypeLinkKind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where TypeLinkKind is left unspecified.
	linkEmpty LinkKind = iota

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

func (k LinkKind) String() string {
	return linkText[k]
}

// Link describes a dependency connection between two symbols.
type Link struct {
	Symbol *Symbol
	Kind   LinkKind
}

type LinkSet map[*Symbol]LinkKind

func NewLinkSet() LinkSet {
	return make(LinkSet)
}

func (s LinkSet) Has(symbol *Symbol) bool {
	_, ok := s[symbol]
	return ok
}

func (s LinkSet) Add(symbol *Symbol, kind LinkKind) {
	switch kind {
	case linkEmpty:
		panic("empty kind")
	case linkDirect:
		s.AddDirect(symbol)
	case linkIndirect:
		s.AddIndirect(symbol)
	default:
		panic(fmt.Sprintf("unexpected %s (%d) link kind", kind, kind))
	}
}

// Adds direct link to set.
func (s LinkSet) AddDirect(symbol *Symbol) {
	s[symbol] = linkDirect
}

// Adds indirect link to set.
func (s LinkSet) AddIndirect(symbol *Symbol) {
	k := s[symbol]
	if k == linkDirect {
		// direct link takes priority over indirect one
		return
	}
	s[symbol] = linkIndirect
}

func (s LinkSet) Elems() []Link {
	if len(s) == 0 {
		return nil
	}

	elems := make([]Link, 0, len(s))
	for symbol, kind := range s {
		elems = append(elems, Link{
			Symbol: symbol,
			Kind:   kind,
		})
	}
	return elems
}

// SymbolContext carries information and objects which are needed to perform
// unit level symbols hoisting, cycled definitions detection, etc.
type SymbolContext struct {
	// only used during type inspection
	members map[ /* member name */ string]struct{}

	Links LinkSet

	// Unit level symbol from which inspect scan started.
	Symbol *Symbol

	// keeps track of current link type kind when descending/ascending
	// nested type specifiers
	//
	// only used during type inspection
	Kind LinkKind
}

func NewSymbolContext(symbol *Symbol) *SymbolContext {
	return &SymbolContext{
		Links:  NewLinkSet(),
		Symbol: symbol,
		Kind:   linkDirect,
	}
}

func (c *SymbolContext) push(kind LinkKind) LinkKind {
	k := c.Kind
	if k == linkDirect {
		// only direct link can be changed by descending nested type specifiers
		c.Kind = kind
	}
	return k
}

func (m *Merger) inspect() error {
	g := NewGraphBuilder(len(m.unit.Types))

	for _, s := range m.unit.Types {
		i := s.Def.(astIndexSymDef)
		ctx := NewSymbolContext(s)
		err := m.shallowScanType(ctx, i)
		if err != nil {
			return err
		}

		if ctx.Links[s] == linkDirect {
			return fmt.Errorf("%s: symbol \"%s\" directly references itself", s.Pos, s.Name)
		}

		g.Add(s, ctx.Links.Elems())
	}

	for _, s := range m.unit.Constants {
		ctx := NewSymbolContext(s)
		err := m.inspectConstant(ctx)
		if err != nil {
			return err
		}

		if ctx.Links.Has(s) {
			return fmt.Errorf("%s: symbol \"%s\" references itself", s.Pos, s.Name)
		}

		g.Add(s, ctx.Links.Elems())
	}

	m.graph = g.Scan()
	return nil
}

// performs preliminary top-level type definition scan in order to obtain data
// necessary for constructing dependency graph between the named types
func (m *Merger) shallowScanType(ctx *SymbolContext, i astIndexSymDef) error {
	spec := m.nodes.Type(i).Spec
	kind := spec.Kind()

	var err error
	switch kind {
	case tps.Name:
		err = m.shallowScanNamedType(ctx, spec.(ast.TypeName))
	case tps.Struct:
		err = m.shallowScanStructType(ctx, spec.(ast.StructType))
	case tps.Enum:
		// TODO: implement enum scan
	case tps.Bag:
		panic("not implemented")
	default:
		panic(fmt.Sprintf("unexpected %s type specifier", kind.String()))
	}

	return err
}

func (m *Merger) shallowScanStructType(ctx *SymbolContext, spec ast.StructType) error {
	methods := m.nodes.MethodsByReceiver[ctx.Symbol.Name]
	ctx.members = make(map[string]struct{}, len(spec.Fields)+len(methods))

	for _, field := range spec.Fields {
		err := m.shallowScanStructField(ctx, field)
		if err != nil {
			return err
		}
	}

	for _, i := range methods {
		err := m.shallowScanTypeMethod(ctx, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Merger) shallowScanTypeMethod(ctx *SymbolContext, i astIndexSymDef) error {
	method := m.nodes.Method(i)
	name := method.Name.Lit
	pos := method.Name.Pos

	_, ok := ctx.members[name]
	if ok {
		return fmt.Errorf("%s: field with name \"%s\" is already present in struct", pos.String(), name)
	}

	ctx.members[name] = struct{}{}
	return nil
}

func (m *Merger) shallowScanStructField(ctx *SymbolContext, field ast.FieldDefinition) error {
	name := field.Name.Lit
	pos := field.Name.Pos

	_, ok := ctx.members[name]
	if ok {
		return fmt.Errorf("%s: field with name \"%s\" is already present in struct", pos.String(), name)
	}

	err := m.shallowScanTypeSpecifier(ctx, field.Type)
	if err != nil {
		return err
	}

	ctx.members[name] = struct{}{}
	return nil
}

func (m *Merger) shallowScanTypeSpecifier(ctx *SymbolContext, spec ast.TypeSpec) error {
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
	case tps.Struct, tps.Enum, tps.Union, tps.Bag:
		return fmt.Errorf("%s: anonymous nested %ss are not allowed", spec.Pin(), spec.Kind())
	default:
		panic(fmt.Sprintf("not implemented for %s", spec.Kind().String()))
	}
}

func (m *Merger) shallowScanArrayType(ctx *SymbolContext, spec ast.ArrayType) error {
	// save previous link kind for later
	k := ctx.push(linkDirect)

	err := m.inspectExp(ctx, spec.Size)
	if err != nil {
		return err
	}
	err = m.shallowScanTypeSpecifier(ctx, spec.ElemType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.Kind = k
	return nil
}

func (m *Merger) shallowScanArrayPointerType(ctx *SymbolContext, spec ast.ArrayPointerType) error {
	// save previous link kind for later
	k := ctx.push(linkIndirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.ElemType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.Kind = k
	return nil
}

func (m *Merger) shallowScanChunkType(ctx *SymbolContext, spec ast.ChunkType) error {
	// save previous link kind for later
	k := ctx.push(linkIndirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.ElemType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.Kind = k
	return nil
}

func (m *Merger) shallowScanPointerType(ctx *SymbolContext, spec ast.PointerType) error {
	// save previous link kind for later
	k := ctx.push(linkIndirect)

	err := m.shallowScanTypeSpecifier(ctx, spec.RefType)
	if err != nil {
		return err
	}

	// restore previous link kind
	ctx.Kind = k
	return nil
}

func (m *Merger) shallowScanNamedType(ctx *SymbolContext, spec ast.TypeName) error {
	name := spec.Name.Lit
	pos := spec.Name.Pos

	// in shallow scan ref count increase is not needed
	symbol := m.unit.Scope.lookup(name, pos.Num)
	if symbol == nil {
		return fmt.Errorf("%s: undefined symbol \"%s\"", pos.String(), name)
	}
	if symbol.Kind != smk.Type {
		return fmt.Errorf("%s: symbol \"%s\" is not a type", pos.String(), name)
	}

	if symbol.Scope.Kind == scp.Unit {
		ctx.Links.Add(symbol, ctx.Kind)
	}
	return nil
}
