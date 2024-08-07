package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/exn"
	"github.com/mebyus/gizmo/ast/uop"
)

func (s *Scope) evalStaticExp(exp Expression) (Expression, error) {
	if !exp.Type().Static() {
		return nil, fmt.Errorf("%s: expression is not static", exp.Pin())
	}
	return s.evalExp(exp)
}

func (s *Scope) evalExp(exp Expression) (Expression, error) {
	switch exp.Kind() {
	case exn.Integer:
		return exp, nil
	case exn.Symbol:
		s := exp.(*SymbolExpression).Sym
		c := s.Def.(*ConstDef)
		return c.Exp, nil
	case exn.Paren:
		return exp.(*ParenthesizedExpression).Inner, nil
	case exn.Unary:
		return s.evalUnaryExp(exp.(*UnaryExpression))
	case exn.Binary:
		return s.evalBinaryExp(exp.(*BinaryExpression))
	default:
		panic(fmt.Sprintf("unexpected %s expression", exp.Kind()))
	}
}

func (s *Scope) evalBinaryExp(exp *BinaryExpression) (Expression, error) {
	l, err := s.evalExp(exp.Left)
	if err != nil {
		return nil, err
	}
	r, err := s.evalExp(exp.Right)
	if err != nil {
		return nil, err
	}
	if l.Kind() == exn.Integer && r.Kind() == exn.Integer {
		a := l.(Integer)
		b := r.(Integer)
		o := exp.Operator

		// TODO: make separate function for code below
		switch o.Kind {
		case bop.Add:
			return addIntegers(a, b)
		case bop.Sub:
			return subIntegers(a, b)
		default:
			panic(fmt.Sprintf("not implemented for %s operator", o.Kind))
		}
	}

	panic(fmt.Sprintf("not implemented for %s and %s expressions", l.Kind(), r.Kind()))
}

func (s *Scope) evalUnaryExp(exp *UnaryExpression) (Expression, error) {
	inner, err := s.evalExp(exp.Inner)
	if err != nil {
		return nil, err
	}
	switch inner.Kind() {
	case exn.Integer:
		i := inner.(Integer)
		i.Pos = exp.Pin()

		o := exp.Operator
		switch o.Kind {
		case uop.Plus:
			return i, nil
		case uop.Minus:
			i.Neg = !i.Neg
			return i, nil
		default:
			panic(fmt.Sprintf("unexpected %s unary operator on integer", o.Kind))
		}
	default:
		panic(fmt.Sprintf("unexpected %s expression", inner.Kind()))
	}
}

func subIntegers(a, b Integer) (Integer, error) {
	b.Neg = !b.Neg
	return addIntegers(a, b)
}

func addIntegers(a, b Integer) (Integer, error) {
	i := Integer{
		typ: a.typ,
		Pos: a.Pos,
	}

	var val uint64
	var neg bool

	switch {
	case a.Neg && b.Neg:
		neg = true
		val = a.Val + b.Val
	case a.Neg && !b.Neg:
		if a.Val > b.Val {
			neg = true
			val = a.Val - b.Val
		} else {
			val = b.Val - a.Val
		}
	case !a.Neg && b.Neg:
		if a.Val >= b.Val {
			val = a.Val - b.Val
		} else {
			neg = true
			val = b.Val - a.Val
		}
	case !a.Neg && !b.Neg:
		val = a.Val + b.Val
	default:
		panic("impossible condition")
	}
	if val == 0 {
		neg = false
	}
	i.Val = val
	i.Neg = neg

	return i, nil
}
