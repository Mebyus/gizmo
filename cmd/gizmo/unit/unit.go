package unit

import (
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/gizmo/unit/gen"
	"github.com/mebyus/gizmo/cmd/gizmo/unit/utyp"
)

var Unit = &butler.Lackey{
	Name: "unit",

	Short: "use command on gizmo unit",
	Usage: "gizmo unit <command>",

	Sub: []*butler.Lackey{
		utyp.Utyp,
		gen.Gen,
	},
}
