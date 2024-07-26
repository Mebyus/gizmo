package stg

type MethodDef struct {
	nodeSymDef

	Signature

	Body Block

	// List of top-level unit symbols which are used in method.
	Refs []*Symbol
}
