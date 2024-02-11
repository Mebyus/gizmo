package gencpp

import (
	"io"

	"github.com/mebyus/gizmo/ast"
)

func Gen(w io.Writer, atom ast.UnitAtom) error {
	builder := NewBuilder(0)
	builder.UnitAtom(atom)

	_, err := io.Copy(w, builder)
	return err
}
