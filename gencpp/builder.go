package gencpp

type Config struct {
	DefaultNamespace string

	GlobalNamespacePrefix string

	// Initial preallocated memory capacity for storing output
	Size int

	SourceLocationComments bool
}
