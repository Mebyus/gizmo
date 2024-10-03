package genstf

import "github.com/mebyus/gizmo/stg"

func (g *Builder) main(units []*stg.Unit) {
	g.puts("fun main() {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("var failed: uint = 0;")
	g.nl()

	g.indent()
	g.puts("var t: stf.Test;")
	g.nl()

	g.indent()
	g.puts("t.init(logbuf[:]);")
	g.nl()

	for _, u := range units {
		for _, t := range u.Tests {
			g.nl()
			g.test(u, t)
		}
	}

	g.nl()
	g.indent()
	g.puts("if failed != 0 {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("os.exit(1);")
	g.nl()

	g.dec()
	g.indent()
	g.puts("}")
	g.nl()

	g.dec()
	g.puts("}")
	g.nl()
}

func (g *Builder) test(unit *stg.Unit, t *stg.Symbol) {
	g.indent()
	g.puts("t.reset();")
	g.nl()

	g.indent()
	g.puts(unit.Name)
	g.puts(".test.")
	g.puts(t.Name)
	g.puts("(t.&);")
	g.nl()

	g.indent()
	g.puts("if !t.ok() {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("failed += 1;")
	g.nl()

	g.indent()
	g.puts("print(\"")
	g.puts("[FAIL] ")
	g.puts(t.Name)
	g.puts("\");")
	g.nl()

	g.dec()
	g.indent()
	g.puts("}")

	g.nl()
}
