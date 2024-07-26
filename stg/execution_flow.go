package stg

import "github.com/mebyus/gizmo/stg/sfp"

type FlowPoint struct {
	// Index (inside a block) of statement with special flow of execution.
	Index uint32

	Kind sfp.Kind
}
