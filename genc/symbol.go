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
		prefix := g.prefix + g.getMangledUnitName(s.Scope.Unit) + "_"
		if s.Kind == smk.Method {
			return prefix + "g" + "_" + strings.Replace(s.Name, ".", "_", 1)
		}
		return prefix + s.Name
	}
	switch s.Name {
	case "long", "unsigned", "signed", "static", "inline":
		return g.prefix + s.Name
	default:
		return s.Name
	}
}

func (g *Builder) SymbolName(s *stg.Symbol) {
	g.puts(g.getSymbolName(s))
}

func (g *Builder) SymbolSourceComment(s *stg.Symbol) {
	g.LineComment(s.Pos.String())
}
