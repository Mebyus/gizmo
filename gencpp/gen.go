package gencpp

import (
	"io"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ir"
)

func Gen(w io.Writer, atom ast.UnitAtom) error {
	sm := ir.IndexMethods(atom)

	builder := NewBuilder(0)
	builder.sm = sm
	builder.UnitAtom(atom)

	_, err := io.Copy(w, builder)
	return err
}
