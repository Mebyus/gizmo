package format

import (
	"io"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
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

	b := Format(atom, holder.Tokens())
	_, err = w.Write(b)
	return err
}

func Format(atom ast.Atom, tokens []token.Token) []byte {
	f := New()
	return f.Format(atom, tokens)
}

type Builder struct {
	g source.Builder
}

func New() *Builder {
	return &Builder{}
}

func (g *Builder) Format(atom ast.Atom, tokens []token.Token) []byte {
	for _, top := range atom.Nodes {
		g.TopLevel(top)
	}
	return g.g.Bytes()
}
