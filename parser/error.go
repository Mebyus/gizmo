package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

type UnexpectedTokenError struct {
	Token token.Token
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("%s: unexpected token %s", e.Token.Pos, e.Token.Kind)
}

func (p *Parser) unexpected() error {
	return &UnexpectedTokenError{Token: p.tok}
}
