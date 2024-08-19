package genc

import (
	"strings"

	"github.com/mebyus/gizmo/enums/smk"
	"github.com/mebyus/gizmo/stg"
	"github.com/mebyus/gizmo/stg/scp"
)

func (g *Builder) getSymbolName(s *stg.Symbol) string {
	if s.Scope.Kind == scp.Global {
		return s.Name
	}
	if s.Scope.Kind == scp.Unit {
		if s.Kind == smk.Type {
			return g.tprefix + s.Name
		}
		if s.Kind == smk.Method {
			return g.prefix + strings.Replace(s.Name, ".", "_", 1)
		}
		return g.prefix + s.Name
	}
	return s.Name
}

func (g *Builder) SymbolName(s *stg.Symbol) {
	g.puts(g.getSymbolName(s))
}

func (g *Builder) SymbolSourceComment(s *stg.Symbol) {
	g.LineComment(s.Pos.String())
}
