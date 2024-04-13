package gencpp

import (
	"io"
	"strings"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
)

func Gen(w io.Writer, cfg *Config, atom ast.Atom) error {

	if cfg == nil {
		cfg = &Config{
			DefaultNamespace:       "<default>",
			SourceLocationComments: true,
		}
	}
	if cfg.GlobalNamespacePrefix != "" && !strings.HasSuffix(cfg.GlobalNamespacePrefix, "::") {
		cfg.GlobalNamespacePrefix = cfg.GlobalNamespacePrefix + "::"
	}

	builder := source.NewBuilder(cfg.Size)

	_, err := io.Copy(w, builder)
	return err
}
