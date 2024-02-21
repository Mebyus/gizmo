package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) topLevelType() (ast.TopType, error) {
	p.advance() // skip "type"

	if p.tok.Kind != token.Identifier {
		return ast.TopType{}, p.unexpected(p.tok)
	}

	name := p.idn()
	p.advance() // skip name identifier

	spec, err := p.typeSpecifier()
	if err != nil {
		return ast.TopType{}, err
	}

	return ast.TopType{
		Name: name,
		Spec: spec,
	}, nil
}

func (p *Parser) typeSpecifier() (ast.TypeSpecifier, error) {
	if p.tok.Kind == token.Identifier {
		return p.typeName()
	}
	if p.tok.Kind == token.Asterisk {
		return p.pointerType()
	}
	if p.tok.Kind == token.ArrayPointer {
		return p.arrayPointerType()
	}
	if p.tok.Kind == token.Struct {
		return p.structType()
	}
	if p.tok.Kind == token.Chunk {
		return p.chunkType()
	}
	if p.tok.Kind == token.LeftSquare {
		return p.arrayType()
	}
	if p.tok.Kind == token.Enum {
		return p.enumType()
	}
	return nil, fmt.Errorf("other type specifiers not implemented (start from %s at %s)",
		p.tok.Kind.String(), p.tok.Pos.String())
}

func (p *Parser) typeName() (ast.TypeName, error) {
	name, err := p.scopedIdentifier()
	if err != nil {
		return ast.TypeName{}, err
	}
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
	name := p.idn()
	p.advance() // skip entry name identifier

	if p.tok.Kind != token.Assign {
		// entry without explicit assigned value
		return ast.EnumEntry{Name: name}, nil
	}

	p.advance() // skip "="

	expr, err := p.expr()
	if err != nil {
		return ast.EnumEntry{}, err
	}

	return ast.EnumEntry{
		Name:       name,
		Expression: expr,
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

	size, err := p.expr()
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
