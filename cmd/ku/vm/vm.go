package vm

import (
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/ku/vm/asm"
	"github.com/mebyus/gizmo/cmd/ku/vm/exec"
	"github.com/mebyus/gizmo/cmd/ku/vm/gui"
	"github.com/mebyus/gizmo/cmd/ku/vm/lex"
)

var VM = &butler.Lackey{
	Name: "vm",

	Short: "use kuvm command",
	Usage: "ku vm <command>",

	Sub: []*butler.Lackey{
		lex.Lex,
		asm.Asm,
		exec.Exec,
		gui.Serve,
	},
}
