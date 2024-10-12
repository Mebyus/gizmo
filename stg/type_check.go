package stg

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/bop"
)

// typeCheckDesugarBinExp performs type checking and desugaring of a given
// binary expression and returns the resulting expression.
func typeCheckDesugarBinExp(exp *BinExp) (Exp, error) {
	switch exp.Operator.Kind {
	case bop.Equal:
		return typeCheckDesugarEqualExp(exp)
	case bop.NotEqual:
		return typeCheckDesugarNotEqualExp(exp)
	case bop.Less, bop.LessOrEqual, bop.Greater, bop.GreaterOrEqual:
		return typeCheckIntegersComparison(exp)
	default:
	}

	a := exp.Left.Type()
	b := exp.Right.Type()
	if a == b {
		// TODO: add operator compatibility checks
		exp.typ = a
		return exp, nil
	}
	if a.IsStaticInteger() {
		exp.typ = b
		return exp, nil
	}
	if b.IsStaticInteger() {
		exp.typ = a
		return exp, nil
	}

	panic(fmt.Sprintf("not implemented for %s and %s types", a.Kind, b.Kind))
}

func typeCheckDesugarEqualExp(exp *BinExp) (Exp, error) {
	if exp.Left.Type().IsStringType() {
		return typeCheckDesugarStringsEqualExp(exp)
	}

	if exp.Left.Type().IsIntegerType() {
		return typeCheckIntegersComparison(exp)
	}

	panic("not implemented")
}

func typeCheckDesugarNotEqualExp(exp *BinExp) (Exp, error) {
	if exp.Left.Type().IsStringType() {
		return typeCheckDesugarStringsNotEqualExp(exp)
	}

	if exp.Left.Type().IsIntegerType() {
		return typeCheckIntegersComparison(exp)
	}

	panic("not implemented")
}

func typeCheckIntegersComparison(exp *BinExp) (Exp, error) {
	a := exp.Left.Type()
	b := exp.Right.Type()

	if !b.IsIntegerType() {
		return nil, fmt.Errorf("%s: mismatched types %s and %s in comparison", exp.Left.Pin(), a.Kind, b.Kind)
	}

	if a.IsStaticInteger() && b.IsStaticInteger() {
		exp.typ = StaticBooleanType
		return exp, nil
	}

	if a.Custom() && b.Custom() {
		if a != b {
			return nil, fmt.Errorf("%s: mismatched custom types %s and %s in comparison",
				exp.Left.Pin(), a.Kind, b.Kind)
		}

		exp.typ = BooleanType
		return exp, nil
	}

	if a.IsStaticInteger() || b.IsStaticInteger() {
		exp.typ = BooleanType
		return exp, nil
	}

	if a != b {
		return nil, fmt.Errorf("%s: mismatched types %s and %s in comparison", exp.Left.Pin(), a.Kind, b.Kind)
	}

	exp.typ = BooleanType
	return exp, nil
}

func typeCheckDesugarStringsEqualExp(exp *BinExp) (Exp, error) {
	a := exp.Left.Type()
	b := exp.Right.Type()

	if !b.IsStringType() {
		return nil, fmt.Errorf("%s: mismatched types %s and %s in equal comparison", exp.Left.Pin(), a.Kind, b.Kind)
	}

	if a.Custom() && b.Custom() {
		if a != b {
			return nil, fmt.Errorf("%s: mismatched custom string types %s and %s in equal comparison",
				exp.Left.Pin(), a.Kind, b.Kind)
		}

		return &StringsEqual{
			Left:  exp.Left,
			Right: exp.Right,
		}, nil
	}

	// for code below one of the types a or b is builtin type

	if a.IsStaticString() {
		if isEmptyString(exp.Left) {
			return &StringEmpty{Exp: exp.Right}, nil
		}

		return &StringsEqual{
			Left:  exp.Left,
			Right: exp.Right,
		}, nil
	}
	if b.IsStaticString() {
		if isEmptyString(exp.Right) {
			return &StringEmpty{Exp: exp.Left}, nil
		}

		return &StringsEqual{
			Left:  exp.Left,
			Right: exp.Right,
		}, nil
	}

	if a != b {
		return nil, fmt.Errorf("%s: mismatched types %s and %s in equal comparison", exp.Left.Pin(), a.Kind, b.Kind)
	}

	return &StringsEqual{
		Left:  exp.Left,
		Right: exp.Right,
	}, nil
}

func typeCheckDesugarStringsNotEqualExp(exp *BinExp) (Exp, error) {
	a := exp.Left.Type()
	b := exp.Right.Type()

	if !b.IsStringType() {
		return nil, fmt.Errorf("%s: mismatched types %s and %s in not equal comparison", exp.Left.Pin(), a.Kind, b.Kind)
	}

	if a.Custom() && b.Custom() {
		if a != b {
			return nil, fmt.Errorf("%s: mismatched custom string types %s and %s in not equal comparison",
				exp.Left.Pin(), a.Kind, b.Kind)
		}

		return &StringsNotEqual{
			Left:  exp.Left,
			Right: exp.Right,
		}, nil
	}

	// for code below one of the types a or b is builtin type

	if a.IsStaticString() {
		if isEmptyString(exp.Left) {
			return &StringNotEmpty{Exp: exp.Right}, nil
		}

		return &StringsNotEqual{
			Left:  exp.Left,
			Right: exp.Right,
		}, nil
	}
	if b.IsStaticString() {
		if isEmptyString(exp.Right) {
			return &StringNotEmpty{Exp: exp.Left}, nil
		}

		return &StringsNotEqual{
			Left:  exp.Left,
			Right: exp.Right,
		}, nil
	}

	if a != b {
		return nil, fmt.Errorf("%s: mismatched types %s and %s in not equal comparison", exp.Left.Pin(), a.Kind, b.Kind)
	}

	return &StringsNotEqual{
		Left:  exp.Left,
		Right: exp.Right,
	}, nil
}

// returns true if a given expression is a static empty string.
func isEmptyString(exp Exp) bool {
	return exp.(String).Size == 0
}
