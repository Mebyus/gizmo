package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/enums/exk"
)

func (s *Scope) evalStaticExp(exp Exp) (Exp, error) {
	if !exp.Type().Static() {
		return nil, fmt.Errorf("%s: expression is not static", exp.Pin())
	}
	return s.evalExp(exp)
}

func (s *Scope) evalExp(exp Exp) (Exp, error) {
	switch exp.Kind() {
	case exk.Integer:
		return exp, nil
	case exk.Symbol:
		s := exp.(*SymbolExp).Sym
		c := s.Def.(*ConstDef)
		return c.Exp, nil
	case exk.Paren:
		return exp.(*ParenExp).Inner, nil
	case exk.Unary:
		return s.evalUnaryExp(exp.(*UnaryExp))
	case exk.Binary:
		return s.evalBinExp(exp.(*BinExp))
	default:
		panic(fmt.Sprintf("unexpected %s expression", exp.Kind()))
	}
}

func (s *Scope) evalBinExp(exp *BinExp) (Exp, error) {
	l, err := s.evalExp(exp.Left)
	if err != nil {
		return nil, err
	}
	r, err := s.evalExp(exp.Right)
	if err != nil {
		return nil, err
	}
	if l.Kind() == exk.Integer && r.Kind() == exk.Integer {
		a := l.(Integer)
		b := r.(Integer)
		o := exp.Operator

		// TODO: make separate function for code below
		switch o.Kind {
		case bop.Add:
			return addIntegers(a, b)
		case bop.Sub:
			return subIntegers(a, b)
		case bop.LeftShift:
			return leftShiftInteger(a, b)
		default:
			panic(fmt.Sprintf("not implemented for %s operator", o.Kind))
		}
	}

	panic(fmt.Sprintf("not implemented for %s and %s expressions", l.Kind(), r.Kind()))
}

func (s *Scope) evalUnaryExp(exp *UnaryExp) (Exp, error) {
	inner, err := s.evalExp(exp.Inner)
	if err != nil {
		return nil, err
	}
	switch inner.Kind() {
	case exk.Integer:
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

func leftShiftInteger(a, b Integer) (Integer, error) {
	if b.Neg {
		return Integer{}, fmt.Errorf("%s: negative left shift", b.Pos)
	}
	if a.Neg {
		return Integer{}, fmt.Errorf("%s: left shift on negative integer", a.Pos)
	}
	if b.Val == 0 {
		return a, nil
	}
	if b.Val >= 64 {
		panic("not implemented")
	}

	i := Integer{
		typ: a.typ,
		Pos: a.Pos,
	}

	i.Val = a.Val << b.Val
	return i, nil
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
