package stg

import (
	"github.com/mebyus/gizmo/enums/exk"
	"github.com/mebyus/gizmo/source"
)

// Represents equality comparison of two strings.
type StringsEqual struct {
	NodeE

	Left  Exp
	Right Exp
}

// Explicit interface implementation check
var _ Exp = &StringsEqual{}

func (*StringsEqual) Kind() exk.Kind {
	return exk.StringsEqual
}

func (e *StringsEqual) Pin() source.Pos {
	return e.Left.Pin()
}

func (e *StringsEqual) Type() *Type {
	return BooleanType
}

// Represents inequality comparison of two strings.
type StringsNotEqual struct {
	NodeE

	Left  Exp
	Right Exp
}

// Explicit interface implementation check
var _ Exp = &StringsNotEqual{}

func (*StringsNotEqual) Kind() exk.Kind {
	return exk.StringsNotEqual
}

func (e *StringsNotEqual) Pin() source.Pos {
	return e.Left.Pin()
}

func (e *StringsNotEqual) Type() *Type {
	return BooleanType
}

// Represents empty string check.
type StringEmpty struct {
	NodeE

	Exp Exp
}

// Explicit interface implementation check
var _ Exp = &StringEmpty{}

func (*StringEmpty) Kind() exk.Kind {
	return exk.StringEmpty
}

func (e *StringEmpty) Pin() source.Pos {
	return e.Exp.Pin()
}

func (e *StringEmpty) Type() *Type {
	return BooleanType
}

// Represents not empty string check.
type StringNotEmpty struct {
	NodeE

	Exp Exp
}

// Explicit interface implementation check
var _ Exp = &StringNotEmpty{}

func (*StringNotEmpty) Kind() exk.Kind {
	return exk.StringNotEmpty
}

func (e *StringNotEmpty) Pin() source.Pos {
	return e.Exp.Pin()
}

func (e *StringNotEmpty) Type() *Type {
	return BooleanType
}
