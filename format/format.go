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

	b := Format(atom, holder.Tokens())
	_, err = w.Write(b)
	return err
}

func Format(atom ast.Atom, tokens []token.Token) []byte {
	g := New(tokens)
	nodes := g.Nodes(atom)
	return NewStapler(tokens, nodes).staple()
}
