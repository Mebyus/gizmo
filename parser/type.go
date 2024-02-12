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
		name, err := p.scopedIdentifier()
		if err != nil {
			return nil, err
		}
		return ast.TypeName{Name: name}, nil
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
	return nil, fmt.Errorf("other type specifiers not implemented %s", p.tok.Short())
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
