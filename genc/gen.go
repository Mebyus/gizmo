package genc

import (
	"io"

	"github.com/mebyus/gizmo/stg"
)

func GenUnit(w io.Writer, u *stg.Unit) error {
	var g Builder
	g.prefix = "ku_"
	g.tprefix = "Ku"
	g.specs = make(map[*stg.Type]string)
	g.Gen(u)
	_, err := w.Write(g.Bytes())
	return err
}
