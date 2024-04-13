package gencpp

import (
	"io"
	"strings"

	"github.com/mebyus/gizmo/ast"
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

	builder := NewBuilder(cfg)

	_, err := io.Copy(w, builder)
	return err
}
