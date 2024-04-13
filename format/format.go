package format

import (
	"io"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/token"
)

func FormatFile(w io.Writer, path string) error {
	lx, err := lexer.FromFile(path)
	if err != nil {
		return err
	}
	holder := NewHolder(lx)
	parser := parser.New(holder)

	_, err = parser.Header()
	if err != nil {
		return err
	}

	atom, err := parser.Parse()
	if err != nil {
		return err
	}

	Format(atom, holder.Tokens())
	return nil
}

func Format(atom ast.Atom, tokens []token.Token) {

}
