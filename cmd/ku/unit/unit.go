package unit

import (
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/ku/unit/gen"
)

var Unit = &butler.Lackey{
	Name: "unit",

	Short: "use command on gizmo unit",
	Usage: "gizmo unit <command>",

	Sub: []*butler.Lackey{
		gen.Gen,
	},
}
