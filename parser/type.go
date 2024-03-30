package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) topLevelType() (ast.TopLevel, error) {
	p.advance() // skip "type"

	if p.tok.Kind != token.Identifier {
		return nil, p.unexpected(p.tok)
	}
	if p.next.Kind == token.LeftParentheses {
		return p.topPrototype()
	}

	name := p.idn()
	p.advance() // skip name identifier

	spec, err := p.typeSpecifier()
	if err != nil {
		return nil, err
	}

	return ast.TopType{
		Name: name,
		Spec: spec,
	}, nil
}

func (p *Parser) topPrototype() (ast.TopPrototype, error) {
	name := p.idn()
	p.advance() // skip name identifier

	if p.tok.Kind != token.LeftParentheses {
		panic("unexpected token")
	}

	params, err := p.protoParams()
	if err != nil {
		return ast.TopPrototype{}, err
	}

	if p.tok.Kind != token.Struct {
		return ast.TopPrototype{}, fmt.Errorf("prototypes of %s are not allowed %s", p.tok.Kind.String(), p.tok.Pos.String())
	}
	spec, err := p.structType()
	if err != nil {
		return ast.TopPrototype{}, err
	}

	return ast.TopPrototype{
		Name:   name,
		Spec:   spec,
		Params: params,
	}, nil
}

func (p *Parser) protoParams() ([]ast.TypeParam, error) {
	p.advance() // skip "("

	var params []ast.TypeParam
	for {
		if p.tok.Kind == token.RightParentheses {
			if len(params) == 0 {
				return nil, fmt.Errorf("no params in prototype %s", p.pos().String())
			}

			p.advance() // skip ")"
			return params, nil
		}

		param, err := p.protoParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightParentheses {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected(p.tok)
		}
	}
}

func (p *Parser) protoParam() (ast.TypeParam, error) {
	err := p.expect(token.Identifier)
	if err != nil {
		return ast.TypeParam{}, err
	}
	name := p.idn()
	p.advance() // skip param name identifier

	err = p.expect(token.Colon)
	if err != nil {
		return ast.TypeParam{}, err
	}
	p.advance() // skip ":"

	// TODO: design data type to represent constraint
	if p.tok.Kind != token.Type {
		return ast.TypeParam{}, p.unexpected(p.tok)
	}
	p.advance() // skip "type"

	return ast.TypeParam{
		Name:       name,
		Constraint: nil,
	}, nil
}

func (p *Parser) typeSpecifier() (ast.TypeSpecifier, error) {
	switch p.tok.Kind {
	case token.Identifier:
		return p.typeNameOrInstance()
	case token.Asterisk:
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
	case token.Fn:
		return p.fnType()
	case token.Union:
		return p.unionType()
	default:
		return nil, fmt.Errorf("other type specifiers not implemented (start from %s at %s)",
			p.tok.Kind.String(), p.tok.Pos.String())
	}
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

func (p *Parser) typeName() (ast.TypeName, error) {
	name, err := p.scopedIdentifier()
	if err != nil {
		return ast.TypeName{}, err
	}
	return ast.TypeName{Name: name}, nil
}

func (p *Parser) typeNameOrInstance() (ast.TypeSpecifier, error) {
	name, err := p.scopedIdentifier()
	if err != nil {
		return nil, err
	}
	if p.tok.Kind != token.LeftDoubleSquare {
		return ast.TypeName{Name: name}, nil
	}

	args, err := p.templateArgs()
	if err != nil {
		return nil, err
	}
	return ast.TemplateInstanceType{
		Params: args,
		Name:   name,
	}, nil
}

func (p *Parser) templateArgs() ([]ast.TypeSpecifier, error) {
	p.advance() // skip "[["

	var args []ast.TypeSpecifier
	for {
		if p.tok.Kind == token.RightDoubleSquare {
			if len(args) == 0 {
				return nil, fmt.Errorf("no args in template instance %s", p.pos().String())
			}

			p.advance() // skip "]]"
			return args, nil
		}

		arg, err := p.typeSpecifier()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.tok.Kind == token.Comma {
			p.advance() // skip ","
		} else if p.tok.Kind == token.RightDoubleSquare {
			// will be skipped at next iteration
		} else {
			return nil, p.unexpected(p.tok)
		}
	}
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
