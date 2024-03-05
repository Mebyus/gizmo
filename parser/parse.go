package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// create identifier from current token
func (p *Parser) idn() ast.Identifier {
	return ast.Identifier{
		Pos: p.tok.Pos,
		Lit: p.tok.Lit,
	}
}

func (p *Parser) identifier() (ast.Identifier, error) {
	if p.tok.Kind != token.Identifier {
		return ast.Identifier{}, p.unexpected(p.tok)
	}
	identifier := p.idn()
	p.advance()
	return identifier, nil
}

// gives a copy of current parser position
func (p *Parser) pos() source.Pos {
	return p.tok.Pos
}

func (p *Parser) basic() ast.BasicLiteral {
	return ast.BasicLiteral{Token: p.tok}
}

func (p *Parser) expect(k token.Kind) error {
	if p.tok.Kind == k {
		return nil
	}
	return p.unexpected(p.tok)
}

func (p *Parser) unitBlock() (*ast.UnitBlock, error) {
	p.advance() // skip "unit"
	if p.tok.Kind != token.Identifier {
		return nil, p.unexpected(p.tok)
	}
	name := p.idn()
	p.advance() // skip unit name identifier

	if p.tok.Kind != token.LeftCurly {
		return nil, p.unexpected(p.tok)
	}

	block, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.UnitBlock{
		Name:  name,
		Block: block,
	}, nil
}

func (p *Parser) namespaceBlock() (ast.NamespaceBlock, error) {
	p.advance() // skip "namespace"
	name, err := p.scopedIdentifier()
	if err != nil {
		return ast.NamespaceBlock{}, err
	}

	if p.tok.Kind != token.LeftCurly {
		return ast.NamespaceBlock{}, p.unexpected(p.tok)
	}
	p.advance() // skip "{"

	var nodes []ast.TopLevel
	for {
		if p.tok.Kind == token.RightCurly {
			p.advance() // skip "}"
			return ast.NamespaceBlock{
				Name:  name,
				Nodes: nodes,
			}, nil
		}

		node, err := p.topLevel()
		if err != nil {
			return ast.NamespaceBlock{}, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) scopedIdentifier() (ast.ScopedIdentifier, error) {
	if p.tok.Kind != token.Identifier {
		return ast.ScopedIdentifier{}, p.unexpected(p.tok)
	}
	name := p.idn()
	p.advance() // skip identifier

	var scopes []ast.Identifier
	for {
		if p.tok.Kind != token.DoubleColon {
			return ast.ScopedIdentifier{
				Scopes: scopes,
				Name:   name,
			}, nil
		}
		p.advance() // skip "::"

		if p.tok.Kind != token.Identifier {
			return ast.ScopedIdentifier{}, p.unexpected(p.tok)
		}

		scopes = append(scopes, name)
		name = p.idn()
		p.advance() // skip identifier
	}
}

func (p *Parser) topLevel() (ast.TopLevel, error) {
	err := p.gatherProps()
	if err != nil {
		return nil, err
	}

	switch p.tok.Kind {
	case token.Type:
		return p.topLevelType()
	case token.Var:
		return p.topLevelVar()
	// case token.Import:
	// 	return p.topLevelImport(false)
	case token.Declare:
		return p.topLevelDeclare()
	case token.Fn:
		return p.topLevelFn()
	case token.Const:
		return p.topLevelConst()
	case token.Method:
		return p.topLevelMethod()
	// case token.Pub:
	// 	return p.parseTopLevelPublic()
	default:
		return nil, p.unexpected(p.tok)
	}
}

// func (p *Parser) parseInUnitMode() (err error) {
// 	for {
// 		if p.isEOF() {
// 			break
// 		}
// 		err = p.parseTopLevel()
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// func (p *Parser) parseInNoUnitMode() (err error) {
// 	for {
// 		if p.isEOF() {
// 			break
// 		}
// 		err = p.parseScriptTopLevel()
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// func (p *Parser) parseScriptTopLevel() (err error) {
// 	switch p.tok.Kind {
// 	case token.Type:
// 		return p.parseTopLevelType(false)
// 	case token.Var:
// 	case token.Import:
// 		return p.topLevelImport(false)
// 	case token.Fn:
// 		return p.parseTopLevelFunction(false)
// 	default:
// 		var stmt ast.Statement
// 		stmt, err = p.parseStatement()
// 		if err != nil {
// 			return
// 		}
// 		_ = stmt
// 		// p.stmts = append(p.stmts, stmt)
// 		return
// 	}
// 	p.advance()
// 	return
// }

func (p *Parser) topLevelVar() (ast.TopVar, error) {
	statement, err := p.varStatement()
	if err != nil {
		return ast.TopVar{}, err
	}

	return ast.TopVar{
		VarInit: statement.VarInit,
	}, nil
}

func (p *Parser) topLevelConst() (ast.TopConst, error) {
	statement, err := p.constStatement()
	if err != nil {
		return ast.TopConst{}, err
	}
	return ast.TopConst{
		ConstInit: statement.ConstInit,
	}, nil
}

func (p *Parser) topLevelMethod() (ast.TopLevel, error) {
	return p.method()
}
