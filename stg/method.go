package stg

type MethodDef struct {
	nodeSymDef

	Signature

	Body Block

	// List of top-level unit symbols which are used in method.
	Refs []*Symbol

	// Always not nil.
	Receiver *Type

	// Temporary field for transferring context between inspect
	// and full scan phases.
	ctx *Context
}
