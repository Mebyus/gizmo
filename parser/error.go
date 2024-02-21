package parser

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

var ErrTODO = fmt.Errorf("stub error (please implement me)")

type UnexpectedTokenError struct {
	Token token.Token
}

func (u *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token %s at %s", u.Token.Kind.String(), u.Token.Pos.String())
}

func (p *Parser) unexpected(tok token.Token) error {
	return &UnexpectedTokenError{Token: tok}
}
