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
	switch p.tok.Kind {
	// case token.Type:
	// 	return p.parseTopLevelType(false)
	// case token.Var:
	// 	return p.topLevelVar(false)
	// case token.Identifier:
	// 	return p.topLevelConst(false)
	// case token.Import:
	// 	return p.topLevelImport(false)
	case token.Declare:
		return p.topLevelDeclare()
	case token.Fn:
		return p.topLevelFunction()
	case token.Const:
		return p.topLevelConst()
	// case token.Pub:
	// 	return p.parseTopLevelPublic()
	// case token.Atr:
	// 	return p.topLevelAtr()
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

// func (p *Parser) topLevelVar(public bool) error {
// 	stmt, err := p.parseVariableStatement()
// 	if err != nil {
// 		return err
// 	}
// 	p.sc.Variables = append(p.sc.Variables, ast.TopLevelVar{
// 		Public: public,
// 		Var:    stmt,
// 	})
// 	return nil
// }

func (p *Parser) topLevelConst() (ast.TopConstInit, error) {
	statement, err := p.constStatement()
	if err != nil {
		return ast.TopConstInit{}, err
	}
	return ast.TopConstInit{
		ConstInit: statement.ConstInit,
	}, nil
}

// func (p *Parser) topLevelAtr() error {
// 	if p.saved.ok {
// 		return &cer.Error{
// 			Kind: cer.OrphanedAtrBlock,
// 			Span: p.saved.Pos,
// 		}
// 	}

// 	blk, err := p.atrBlock()
// 	if err != nil {
// 		return err
// 	}
// 	p.saved.AtrBlock = blk
// 	p.saved.ok = true
// 	return nil
// }

// func (p *Parser) parseTopLevel() error {
// 	switch p.tok.Kind {
// 	case token.Type:
// 		return p.parseTopLevelType(false)
// 	case token.Var:
// 		return p.topLevelVar(false)
// 	case token.Identifier:
// 		return p.topLevelConst(false)
// 	case token.Import:
// 		return p.topLevelImport(false)
// 	case token.Fn:
// 		return p.parseTopLevelFunction(false)
// 	case token.Pub:
// 		return p.parseTopLevelPublic()
// 	case token.Atr:
// 		return p.topLevelAtr()
// 	default:
// 		return p.unexpected(p.tok)
// 	}
// }

// func (p *Parser) parseTopLevelPublic() error {
// 	p.advance() // skip "pub"
// 	switch p.tok.Kind {
// 	case token.Import:
// 		return p.topLevelImport(true)
// 	case token.Type:
// 		return p.parseTopLevelType(true)
// 	case token.Var:
// 		return p.topLevelVar(true)
// 	case token.Identifier:
// 		return p.topLevelConst(true)
// 	case token.Fn:
// 		return p.parseTopLevelFunction(true)
// 	default:
// 		return p.unexpected(p.tok)
// 	}
// }
