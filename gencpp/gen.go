package gencpp

import (
	"io"
	"strings"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ir"
)

func Gen(w io.Writer, cfg *Config, atom ast.UnitAtom) error {
	meta := ir.Index(atom)

	if cfg == nil {
		cfg = &Config{
			DefaultNamespace:       "<default>",
			SourceLocationComments: true,
		}
	}
	if cfg.GlobalNamespacePrefix != "" && !strings.HasSuffix(cfg.GlobalNamespacePrefix, "::") {
		cfg.GlobalNamespacePrefix = cfg.GlobalNamespacePrefix + "::"
	}

	builder := NewBuilder(cfg)
	builder.meta = meta
	builder.UnitAtom(atom)

	_, err := io.Copy(w, builder)
	return err
}
