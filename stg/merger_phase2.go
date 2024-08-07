package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/enums/smk"
)

// Phase 2 merger indexing coagulates info from AST nodes that were skipped during
// phase 1. It also performs additional scans on the tree to determine:
//
//   - top-level symbol resolution
//   - symbol hoisting information inside the unit
//   - used and defined types indexing
//   - type checking of various operations
func (m *Merger) merge() error {
	err := m.bindMethods()
	if err != nil {
		return err
	}
	err = m.shallowScanSymbols()
	if err != nil {
		return err
	}
	err = m.bindTypes()
	if err != nil {
		return err
	}
	// TODO: we need to tiebreak circular hoisting problem for types
	// and constants. For example consider these unit-level declarations:
	//
	//	type Custom struct {
	//		buf: [SIZE]u8,
	//		n: int,
	//	}
	//
	//	SIZE := 8 + 8;
	//
	// In order to properly create Custom struct we first need to scan
	// and evaluate (at compile time) expression in SIZE constant definition.
	//
	// Our solution to this problem will be to identify custom integer types and
	// unit-level integer constants first. Then process them before other types and
	// constants.
	//
	// We should probably introduce separate type kinds for signed and unsigned custom
	// integers instead of deesignating them as just regular custom types. For
	// architecture dependant integer types (int and uint) we should introduce
	// bit flag (along with flags field to store it) to indicate integer types which
	// size depends on architecture.
	//
	// Compile time constant integers (even with explicit fixed size) should have
	// separate types from regular runtime integer types with designated bit flag
	// set to indicate constant nature of values of such types. In future we can use
	// such flag to specify (at STG level) function parameters which must be known
	// at compile time.
	//
	// There is a separate problem with walrus constants syntax. Consider this example:
	//
	//	a, err := calc_a();
	//	b, err := cals_b();
	//
	// Since second statement tries to alter "err" symbol the statement is invalid with the currently
	// implied rule that constants (both compile-time and runtime) are defined via walrus syntax.
	// Rust partially solves it by allowing "let" construct to shadow symbols inside the same block.
	// This solution is clumsy in my opinion, first of all it may even alter the symbol type mid-block
	// and second it is only a somewhat mouthful way to circumvent the rule of immutability by default.
	//
	// Problem described above is an ergonomic problem of the language:
	//
	//	immutability by default vs.
	//	convenient symbol reuse where it is logically sound, without mouthful syntax
	//
	// Maybe we should eliminate syntactic difference between constants and variables,
	// and make detection whether symbol mutates or not automatic. Consider example below:
	//
	//	a := 10; // compile-time constant
	//	b := a + 2; // variable, mutates later in the code
	//	c := b - 3; // runtime constant, does not mutate later in the code (can be optimized into compile-time constant)
	//	...
	//	b = 1; // symbol mutation automatically turns it into a variable
	//
	// This is too much inference for my personal taste. It will be difficult to implement
	// and confusing to read.
	//
	// Perhaps we should ditch the concept of runtime constants for now. And use simple syntax
	// for distinguishing constants and variables like this:
	//
	//	$a := 10;
	//	b := 3;
	//
	err = m.scanConstants()
	if err != nil {
		return err
	}
	err = m.shallowScanFuns()
	if err != nil {
		return err
	}
	err = m.shallowScanMethods()
	if err != nil {
		return err
	}
	err = m.scanFuns()
	if err != nil {
		return err
	}
	err = m.scanMethods()
	if err != nil {
		return err
	}
	m.unit.Scope.WarnUnused(&Context{m: m})
	return nil
}

// Checks that each method has a custom type receiver defined in unit.
// Binds each method to its own receiver (based on receiver name).
func (m *Merger) bindMethods() error {
	if len(m.nodes.Meds) != len(m.unit.Meds) {
		panic(fmt.Sprintf("inconsistent number of method nodes (%d) and symbols (%d)",
			len(m.nodes.Meds), len(m.unit.Meds)))
	}

	for j := 0; j < len(m.unit.Meds); j += 1 {
		i := astIndexSymDef(j)
		med := m.nodes.Med(i)

		rname := med.Receiver.Name.Lit
		pos := med.Receiver.Name.Pos

		r := m.unit.Scope.sym(rname)
		if r == nil {
			return fmt.Errorf("%s: unresolved symbol \"%s\" in method receiver", pos.String(), rname)
		}
		if r.Kind != smk.Type {
			return fmt.Errorf("%s: only custom types can have methods, but \"%s\" is not a type",
				pos.String(), rname)
		}

		// method names are guaranteed to be unique for specific receiver due
		// to earlier scope symbol name check
		m.nodes.bindMethod(rname, i)
	}
	return nil
}

func (m *Merger) scanConstants() error {
	for _, s := range m.unit.Lets {
		c, err := m.scanCon(s)
		if err != nil {
			return err
		}
		s.Def = c

		if c.Type != nil {
			s.Type = c.Type
		} else {
			// TODO: this won't work if expression has other constants
			// which are not yet processed
			s.Type = c.Expr.Type()
		}
	}
	return nil
}

func (m *Merger) scanCon(s *Symbol) (*ConstDef, error) {
	scope := m.unit.Scope
	node := m.nodes.Con(s.Def.(astIndexSymDef))

	t, err := scope.Types.Lookup(node.Type)
	if err != nil {
		return nil, err
	}

	expr := node.Expr
	if expr == nil {
		panic("nil expression in constant definition")
	}
	ctx := m.newConstCtx()
	e, err := scope.Scan(ctx, expr)
	if err != nil {
		return nil, err
	}

	if ctx.ref.Has(s) {
		return nil, fmt.Errorf("%s: init cycle (constant \"%s\" definition references itself)", s.Pos.String(), s.Name)
	}

	// TODO: check that expression type and constant type match
	refs := ctx.ref.Elems()
	if len(refs) != 0 {
		panic("referencing other symbols in constant definition is not implemented")
	}

	return &ConstDef{
		Expr: e,
		Refs: refs,
		Type: t,
	}, nil
}

func (m *Merger) shallowScanFuns() error {
	for _, s := range m.unit.Funs {
		err := m.shallowScanFun(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) shallowScanMethods() error {
	for _, s := range m.unit.Types {
		t := s.Def.(*Type)
		methods := t.Def.(CustomTypeDef).Methods
		for _, method := range methods {
			err := m.shallowScanMethod(t, method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Merger) shallowScanMethod(t *Type, s *Symbol) error {
	scope := m.unit.Scope
	node := m.nodes.Med(s.Def.(astIndexSymDef))

	params, err := newParamSymbols(scope, node.Signature.Params)
	if err != nil {
		return err
	}

	result, err := scope.Types.Lookup(node.Signature.Result)
	if err != nil {
		return err
	}

	var receiver *Type
	if node.Receiver.Ptr {
		receiver = scope.Types.storePointer(t)
	} else {
		receiver = t
	}

	pos := node.Body.Pos
	def := &MethodDef{
		Receiver: receiver,
		Signature: Signature{
			Params: params,
			Result: result,
			Never:  node.Signature.Never,
		},
		Body: Block{Pos: pos},
	}

	def.Body.Scope = NewTopScope(m.unit.Scope, &def.Body.Pos)
	for _, param := range params {
		name := param.Name
		p := def.Body.Scope.sym(name)
		if p != nil {
			return fmt.Errorf("%s: parameter \"%s\" redeclared in this function", pos.String(), name)
		}
		def.Body.Scope.Bind(param)
	}

	s.Def = def
	return nil
}

func newParamSymbols(scope *Scope, defs []ast.FieldDefinition) ([]*Symbol, error) {
	if len(defs) == 0 {
		return nil, nil
	}

	params := make([]*Symbol, 0, len(defs))
	for _, p := range defs {
		t, err := scope.Types.Lookup(p.Type)
		if err != nil {
			return nil, err
		}

		params = append(params, &Symbol{
			Pos:  p.Name.Pos,
			Name: p.Name.Lit,
			Type: t,
			Kind: smk.Param,
		})
	}
	return params, nil
}

// scan function signature, body will be scanned in a separate pass
func (m *Merger) shallowScanFun(s *Symbol) error {
	scope := m.unit.Scope
	node := m.nodes.Fun(s.Def.(astIndexSymDef))

	params, err := newParamSymbols(scope, node.Signature.Params)
	if err != nil {
		return err
	}

	pos := node.Body.Pos
	result, err := scope.Types.Lookup(node.Signature.Result)
	if err != nil {
		return err
	}

	f := &FunDef{
		Signature: Signature{
			Params: params,
			Result: result,
			Never:  node.Signature.Never,
		},
		Body: Block{Pos: pos},
	}

	f.Body.Scope = NewTopScope(m.unit.Scope, &f.Body.Pos)
	for _, param := range params {
		name := param.Name
		s := f.Body.Scope.sym(name)
		if s != nil {
			return fmt.Errorf("%s: parameter \"%s\" redeclared in this function", pos.String(), name)
		}
		f.Body.Scope.Bind(param)
	}

	s.Def = f
	return nil
}

func (m *Merger) scanFuns() error {
	for i, s := range m.unit.Funs {
		err := m.scanFun(s, astIndexSymDef(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) scanFun(s *Symbol, i astIndexSymDef) error {
	body := m.nodes.Fun(i).Body
	f := s.Def.(*FunDef)
	return m.scanFunBody(f, body.Statements)
}

func (m *Merger) scanFunBody(def *FunDef, statements []ast.Statement) error {
	if len(statements) == 0 {
		if def.Result != nil {
			return fmt.Errorf("%s: function with return type cannot have empty body", def.Body.Pos.String())
		}
		if def.Never {
			return fmt.Errorf("%s: function marked as never returning cannot have empty body", def.Body.Pos.String())
		}
		return nil
	}

	ctx := m.newFunCtx(def)
	err := def.Body.Fill(ctx, statements)
	if err != nil {
		return err
	}
	def.Refs = ctx.ref.Elems()
	return nil
}

func (m *Merger) scanMethods() error {
	for i, s := range m.unit.Meds {
		err := m.scanMethod(s, astIndexSymDef(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Merger) scanMethod(s *Symbol, i astIndexSymDef) error {
	body := m.nodes.Med(i).Body
	def := s.Def.(*MethodDef)
	return m.scanMethodBody(def, body.Statements)
}

func (m *Merger) scanMethodBody(def *MethodDef, statements []ast.Statement) error {
	if len(statements) == 0 {
		if def.Result != nil {
			return fmt.Errorf("%s: method with return type cannot have empty body", def.Body.Pos.String())
		}
		if def.Never {
			return fmt.Errorf("%s: method marked as never returning cannot have empty body", def.Body.Pos.String())
		}
		return nil
	}

	ctx := m.newMethodCtx(def)
	err := def.Body.Fill(ctx, statements)
	if err != nil {
		return err
	}
	def.Refs = ctx.ref.Elems()
	return nil
}
