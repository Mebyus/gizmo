package tt

import (
	"fmt"

	"github.com/mebyus/gizmo/source"
)

type Warn struct {
	Pos source.Pos
	Msg string
}

func (w Warn) Error() string {
	return fmt.Sprintf("[warn] %s: %s", w.Pos.String(), w.Msg)
}

// adds warning to merger output
func (m *Merger) warn(pos source.Pos, msg string) {
	m.Warns = append(m.Warns, Warn{Pos: pos, Msg: msg})
}
