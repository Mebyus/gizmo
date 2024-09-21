package stg

import "fmt"

// typeCheckBinExp performs type checking of binary expression and
// returns the resulting type of expression.
func typeCheckBinExp(exp *BinExp) (*Type, error) {
	a := exp.Left.Type()
	b := exp.Right.Type()
	if a == b {
		// TODO: add operator compatibility checks
		return a, nil
	}
	if a.PerfectInteger() {
		return b, nil
	}
	if b.PerfectInteger() {
		return a, nil
	}

	panic(fmt.Sprintf("not implemented for %s and %s types", a.Kind, b.Kind))
	return nil, nil
}
