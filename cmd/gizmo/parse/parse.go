package parse

import (
	"fmt"

	"github.com/mebyus/gizmo/parser"
)

func Parse(filename string) error {
	unit, err := parser.ParseFile(filename)
	if err != nil {
		return err
	}
	fmt.Printf("%+#v\n", unit)
	return nil
}
