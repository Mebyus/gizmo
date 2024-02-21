package treeview

import (
	"fmt"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/tps"
)

func ConvertFieldDefinition(field ast.FieldDefinition) []Node {
	nameTitle := "name: "
	if len(field.Name.Lit) == 0 {
		nameTitle += "<nil>"
	} else {
		nameTitle += field.Name.Lit
	}

	return []Node{
		{
			Text: nameTitle,
		},
		ConvertTypeSpecifier(field.Type),
	}
}

func ConvertTypeSpecifier(spec ast.TypeSpecifier) Node {
	typeTitle := "type: "
	if spec.Kind() == tps.Name {
		typeTitle += spec.(ast.TypeName).Name.String()
	} else {
		typeTitle += fmt.Sprintf("<%s not implemented>", spec.Kind().String())
	}

	return Node{
		Text: typeTitle,
	}
}
