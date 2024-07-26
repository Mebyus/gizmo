package genc

import (
	"io"

	"github.com/mebyus/gizmo/tt"
)

func GenUnit(w io.Writer, u *tt.Unit) error {
	var g Builder
	g.prefix = "ku_"
	g.tprefix = "Ku"
	g.specs = make(map[*tt.Type]string)
	g.Gen(u)
	_, err := w.Write(g.Bytes())
	return err
}
