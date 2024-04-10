package tt

import "github.com/mebyus/gizmo/tt/sfp"

type FlowPoint struct {
	// Index (inside a block) of statement with special flow of execution.
	Index uint32

	Kind sfp.Kind
}
