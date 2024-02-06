package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/token"
)

func (p *Parser) typeSpecifier() (ast.TypeSpecifier, error) {
	if p.tok.Kind == token.Identifier {
		name, err := p.scopedIdentifier()
		if err != nil {
			return nil, err
		}
		return ast.TypeName{Name: name}, nil
	}
	return nil, fmt.Errorf("other type specifiers not implemented")
}
