package atom

import (
	"github.com/mebyus/gizmo/butler"
	"github.com/mebyus/gizmo/cmd/ku/atom/header"
)

var Atom = &butler.Lackey{
	Name: "atom",

	Short: "output info about given atom (gizmo source file)",
	Usage: "gizmo atom <command>",

	Sub: []*butler.Lackey{
		header.Header,
	},
}
