package parser

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) ident() ast.Identifier {
	return ast.Identifier(p.tok)
}

func (p *Parser) basic() ast.BasicLiteral {
	return ast.BasicLiteral(p.tok)
}

func (p *Parser) expect(k token.Kind) error {
	if p.tok.Kind == k {
		return nil
	}
	return p.unexpected(p.tok)
}

func (p *Parser) parseUnitBlock() (unit ast.UnitBlock, err error) {
	if p.tok.Kind != token.Unit {
		return
	}
	p.advance() // skip "unit"
	if p.tok.Kind != token.Identifier {
		err = p.unexpected(p.tok)
		return
	}
	name := p.ident()
	p.advance() // skip unit name identifier

	block, err := p.block()
	if err != nil {
		return
	}

	return ast.UnitBlock{
		Name:       name,
		Statements: block.Statements,
	}, nil
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

// func (p *Parser) topLevelConst(public bool) error {
// 	stmt, err := p.tryParseConstStatement()
// 	if err != nil {
// 		return err
// 	}
// 	if stmt == nil {
// 		return p.unexpected(p.tok)
// 	}
// 	p.sc.Constants = append(p.sc.Constants, ast.TopLevelConst{
// 		Public: public,
// 		Const:  stmt,
// 	})
// 	return nil
// }

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
