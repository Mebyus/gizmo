package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) topType(traits ast.Traits) error {
	p.advance() // skip "type"

	if p.tok.Kind != token.Identifier {
		return p.unexpected(p.tok)
	}

	name := p.word()
	p.advance() // skip name identifier

	spec, err := p.DefTypeSpecifier()
	if err != nil {
		return err
	}

	t := ast.TopType{
		Name:   name,
		Spec:   spec,
		Traits: traits,
	}
	p.atom.Types = append(p.atom.Types, t)
	return nil
}

// variant of generic type specifier parsing for usage in context
// of explicit named type definition via construct
//
//	type MyType <TypeSpecifier>
func (p *Parser) DefTypeSpecifier() (ast.TypeSpec, error) {
	t, err := p.typeSpecifier()
	if err != nil {
		return nil, err
	}
	if t.Kind() == tps.Tuple {
		return nil, fmt.Errorf("%s: tuples are not allowed to be named types", t.Pin().String())
	}
	return t, nil
}

func (p *Parser) typeSpecifier() (ast.TypeSpec, error) {
	switch p.tok.Kind {
	case token.Identifier:
		if p.next.Kind == token.Period {
			return p.importType()
		}
		return p.typeName()
	case token.Asterisk:
		if p.next.Kind == token.Any {
			return p.rawMemoryPointerType()
		}
		return p.pointerType()
	case token.ArrayPointer:
		return p.arrayPointerType()
	case token.Struct:
		return p.structType()
	case token.Chunk:
		return p.chunkType()
	case token.LeftSquare:
		return p.arrayType()
	case token.Enum:
		return p.enumType()
	case token.Fun:
		return p.fnType()
	case token.Union:
		return p.unionType()
	case token.Bag:
		return p.bagType()
	case token.LeftParentheses:
		return p.tupleType()
	default:
		return nil, fmt.Errorf("other type specifiers not implemented (start from %s at %s)",
			p.tok.Kind.String(), p.tok.Pos.String())
	}
}

func (p *Parser) tupleType() (ast.TupleType, error) {
	pos := p.pos()

	fields, err := p.functionParams()
	if err != nil {
		return ast.TupleType{}, err
	}

	return ast.TupleType{
		Pos:    pos,
		Fields: fields,
	}, nil
}

func (p *Parser) bagType() (ast.BagType, error) {
	pos := p.pos()

	p.advance() // skip "bag"

	if p.tok.Kind != token.LeftCurly {
		return ast.BagType{}, p.unexpected(p.tok)
	}
	p.advance() // skip "{"

	var methods []ast.BagMethodSpec
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"

			return ast.BagType{
				Pos:     pos,
				Methods: methods,
			}, nil
		}

		method, err := p.bagMethodSpec()
		if err != nil {
			return ast.BagType{}, err
		}
		methods = append(methods, method)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightCurly {
			// will be skipped at next iteration
		} else {
			return ast.BagType{}, p.unexpected(p.tok)
		}
	}
}

func (p *Parser) bagMethodSpec() (ast.BagMethodSpec, error) {
	if p.tok.Kind != token.Identifier {
		return ast.BagMethodSpec{}, p.unexpected(p.tok)
	}
	name := p.word()
	p.advance() // skip method name

	signature, err := p.functionSignature()
	if err != nil {
		return ast.BagMethodSpec{}, err
	}

	return ast.BagMethodSpec{
		Name: name,

		Params: signature.Params,
		Result: signature.Result,
		Never:  signature.Never,
	}, nil
}

func (p *Parser) unionType() (ast.UnionType, error) {
	pos := p.pos()

	p.advance() // skip "union"

	fields, err := p.structFields()
	if err != nil {
		return ast.UnionType{}, err
	}

	return ast.UnionType{
		Pos:    pos,
		Fields: fields,
	}, nil
}

func (p *Parser) fnType() (ast.FunctionType, error) {
	pos := p.pos()
	p.advance() // skip "fn"

	signature, err := p.functionSignature()
	if err != nil {
		return ast.FunctionType{}, nil
	}

	return ast.FunctionType{
		Pos:       pos,
		Signature: signature,
	}, nil
}

func (p *Parser) importType() (ast.ImportType, error) {
	unit := p.word()
	p.advance() // skip unit name
	p.advance() // skip "."

	if p.tok.Kind != token.Identifier {
		return ast.ImportType{}, p.unexpected(p.tok)
	}
	name := p.word()
	p.advance() // skip type name

	return ast.ImportType{
		Unit: unit,
		Name: name,
	}, nil
}

func (p *Parser) typeName() (ast.TypeName, error) {
	name := p.word()
	p.advance() // skip type name identifier

	return ast.TypeName{Name: name}, nil
}

func (p *Parser) enumType() (ast.EnumType, error) {
	pos := p.pos()

	p.advance() // skip "enum"

	base, err := p.typeName()
	if err != nil {
		return ast.EnumType{}, err
	}

	if p.tok.Kind != token.LeftCurly {
		return ast.EnumType{}, p.unexpected(p.tok)
	}
	p.advance() // skip "{"

	var entries []ast.EnumEntry
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"

			return ast.EnumType{
				Pos:     pos,
				Base:    base,
				Entries: entries,
			}, nil
		}

		entry, err := p.enumEntry()
		if err != nil {
			return ast.EnumType{}, err
		}
		entries = append(entries, entry)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightCurly {
			// will be skipped at next iteration
		} else {
			return ast.EnumType{}, p.unexpected(p.tok)
		}
	}
}

func (p *Parser) enumEntry() (ast.EnumEntry, error) {
	if p.tok.Kind != token.Identifier {
		return ast.EnumEntry{}, p.unexpected(p.tok)
	}
	name := p.word()
	p.advance() // skip entry name identifier

	if p.tok.Kind != token.Assign {
		// entry without explicit assigned value
		return ast.EnumEntry{Name: name}, nil
	}

	p.advance() // skip "="

	expr, err := p.exp()
	if err != nil {
		return ast.EnumEntry{}, err
	}

	return ast.EnumEntry{
		Name:       name,
		Expression: expr,
	}, nil
}

func (p *Parser) rawMemoryPointerType() (ast.AnyPointerType, error) {
	pos := p.pos()

	p.advance() // skip "*"
	p.advance() // skip "any"

	return ast.AnyPointerType{
		Pos: pos,
	}, nil
}

func (p *Parser) pointerType() (ast.PointerType, error) {
	pos := p.pos()

	p.advance() // skip "*"

	ref, err := p.typeSpecifier()
	if err != nil {
		return ast.PointerType{}, err
	}

	return ast.PointerType{
		Pos:     pos,
		RefType: ref,
	}, nil
}

func (p *Parser) arrayType() (ast.ArrayType, error) {
	pos := p.pos()

	p.advance() // skip "["

	size, err := p.exp()
	if err != nil {
		return ast.ArrayType{}, err
	}

	if p.tok.Kind != token.RightSquare {
		return ast.ArrayType{}, p.unexpected(p.tok)
	}
	p.advance() // skip "]"

	elem, err := p.typeSpecifier()
	if err != nil {
		return ast.ArrayType{}, err
	}

	return ast.ArrayType{
		Pos:      pos,
		ElemType: elem,
		Size:     size,
	}, nil
}

func (p *Parser) arrayPointerType() (ast.ArrayPointerType, error) {
	pos := p.pos()

	p.advance() // skip "[*]"

	elem, err := p.typeSpecifier()
	if err != nil {
		return ast.ArrayPointerType{}, err
	}

	return ast.ArrayPointerType{
		Pos:      pos,
		ElemType: elem,
	}, nil
}

func (p *Parser) structType() (ast.StructType, error) {
	pos := p.pos()

	p.advance() // skip "struct"

	fields, err := p.structFields()
	if err != nil {
		return ast.StructType{}, err
	}

	return ast.StructType{
		Pos:    pos,
		Fields: fields,
	}, nil
}

func (p *Parser) structFields() ([]ast.FieldDefinition, error) {
	if p.tok.Kind != token.LeftCurly {
		return nil, p.unexpected(p.tok)
	}
	p.advance() // skip "{"

	var fields []ast.FieldDefinition
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"
			return fields, nil
		}

		field, err := p.field()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightCurly {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected(p.tok)
		}
	}
}

func (p *Parser) chunkType() (ast.ChunkType, error) {
	pos := p.pos()

	p.advance() // skip "[]"

	elem, err := p.typeSpecifier()
	if err != nil {
		return ast.ChunkType{}, err
	}

	return ast.ChunkType{
		Pos:      pos,
		ElemType: elem,
	}, nil
}
