package genstf

import "github.com/mebyus/gizmo/stg"

func (g *Builder) imports(units []*stg.Unit) {
	g.puts("unit main")
	g.nl()
	g.nl()

	g.puts("import std {")
	g.nl()

	g.inc()
	g.indent()
	g.puts("stf => \"stf\"")
	g.dec()

	g.nl()
	g.puts("}")
	g.nl()
	g.nl()

	g.puts("import {")
	g.nl()

	g.inc()
	for _, u := range units {
		if len(u.Tests) == 0 {
			continue
		}

		g.indent()
		g.puts(u.Name)
		g.puts(" => ")
		g.putb('"')
		g.puts(u.Path.ImpStr)
		g.putb('"')
		g.nl()
	}
	g.dec()

	g.puts("}")
	g.nl()
}
