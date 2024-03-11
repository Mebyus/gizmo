package ast

import (
	"github.com/mebyus/gizmo/ast/prv"
	"github.com/mebyus/gizmo/source"
)

// <Prop> = <KeyValueProp> | <TagsProp>
//
// <KeyValueProp> = "#[" <Key> = <Value> "]"
//
// <Key> = <Identifier> { "." <Identifier> }
// <Value> = <String> | <Integer> | <True> | <False>
//
// <TagsProp> = "#[" <Tag> { "," <Tag> } [ "," ] "]"
// <Tag> = <Identifier>
type Prop struct {
	Pos source.Pos

	Value PropValue

	// Always empty for TagsProp
	Key string
}

type PropValue interface {
	Node

	// dummy discriminator method
	PropValue()

	Kind() prv.Kind
}

// Dummy prop value node provides quick, easy to use implementation of discriminator PropValue() method
//
// Used for embedding into other (non-dummy) prop value nodes
type nodePropValue struct{ uidHolder }

func (nodePropValue) PropValue() {}

type PropValueInteger struct {
	nodePropValue

	Pos source.Pos

	Val uint64
}

var _ PropValue = PropValueInteger{}

func (PropValueInteger) Kind() prv.Kind {
	return prv.Integer
}

func (v PropValueInteger) Pin() source.Pos {
	return v.Pos
}

type PropValueString struct {
	nodePropValue

	Pos source.Pos

	Val string
}

var _ PropValue = PropValueString{}

func (PropValueString) Kind() prv.Kind {
	return prv.String
}

func (v PropValueString) Pin() source.Pos {
	return v.Pos
}

type PropValueBool struct {
	nodePropValue

	Pos source.Pos

	Val bool
}

var _ PropValue = PropValueBool{}

func (PropValueBool) Kind() prv.Kind {
	return prv.Bool
}

func (v PropValueBool) Pin() source.Pos {
	return v.Pos
}

type PropValueTags struct {
	nodePropValue

	Pos source.Pos

	// Contains at least one element
	Tags []Identifier
}

var _ PropValue = PropValueTags{}

func (PropValueTags) Kind() prv.Kind {
	return prv.Tag
}

func (v PropValueTags) Pin() source.Pos {
	return v.Pos
}
