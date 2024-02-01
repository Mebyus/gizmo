package parser

// func (p *Parser) topLevelImport(public bool) error {
// 	p.advance() // skip "import"

// 	var origin ast.Identifier
// 	if p.tok.IsIdent() {
// 		origin = ast.Identifier(p.tok)
// 		p.advance() // skip import origin
// 	}

// 	err := p.expect(token.LeftCurly)
// 	if err != nil {
// 		return err
// 	}
// 	p.advance() // skip "{"

// 	block := ast.ImportBlock{
// 		Origin: origin,
// 		Public: public,
// 	}

// 	for {
// 		if p.tok.Kind == token.RightCurly {
// 			p.advance() // skip "}"
// 			p.sc.Imports = append(p.sc.Imports, block)
// 			return nil
// 		}

// 		spec, err := p.importSpec()
// 		if err != nil {
// 			return err
// 		}
// 		block.Specs = append(block.Specs, spec)
// 	}
// }

// func (p *Parser) importSpec() (ast.Import, error) {
// 	err := p.expect(token.Identifier)
// 	if err != nil {
// 		return ast.Import{}, err
// 	}
// 	name := p.ident()
// 	p.advance() // skip identifier

// 	err = p.expect(token.RightArrow)
// 	if err != nil {
// 		return ast.Import{}, err
// 	}
// 	p.advance() // skip "=>"

// 	err = p.expect(token.String)
// 	if err != nil {
// 		return ast.Import{}, err
// 	}
// 	str := ast.ImportString(p.tok)
// 	p.advance() // skip string

// 	return ast.Import{
// 		Name:   name,
// 		String: str,
// 	}, nil
// }
