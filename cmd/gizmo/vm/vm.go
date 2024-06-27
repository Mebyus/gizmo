package vm

import (
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/gizmo/vm/lex"
)

var VM = &butler.Lackey{
	Name: "vm",

	Short: "use kuvm command",
	Usage: "gizmo vm <command>",

	Sub: []*butler.Lackey{
		lex.Lex,
	},
}
