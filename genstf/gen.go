package genstf

import (
	"io"

	"github.com/mebyus/gizmo/uwalk"
)

// Generate ku driver code for running tests from the given program.
func Gen(w io.Writer, p *uwalk.Program) error {
	var g Builder

	g.imports(p.Units)
	g.logbuf()
	g.main(p.Units)

	/*

		TODO: generate code like the following

		import {
			example => "example"
		}

		var t stf.Test

		example.test.add(&t) // note special "test" chain part

		// check t.err, print logs and test result

	*/

	_, err := w.Write(g.Bytes())
	return err
}
