package kbs

import (
	"fmt"

	"github.com/mebyus/gizmo/ast/bop"
)

type Evaluator struct {
	out ScriptOutput

	// maps env name to its value
	env map[ /* env name */ string]string
}

func (e *Evaluator) Reset() {
	e.out = ScriptOutput{}
}

func (e *Evaluator) Eval(dirs []Directive) (*ScriptOutput, error) {
	e.Reset()
	err := e.eval(dirs)
	if err != nil {
		return nil, err
	}
	return &e.out, nil
}

func (e *Evaluator) eval(dirs []Directive) error {
	for _, dir := range dirs {
		switch d := dir.(type) {
		case Include:
			e.out.Includes = append(e.out.Includes, d.String)
		case Link:
			e.out.Links = append(e.out.Links, d.String)
		case Name:
			if e.out.Name != "" {
				return fmt.Errorf("%s: duplicate name directive", d.Pos)
			}
			e.out.Name = d.String
		case If:
			exp, err := e.evalExp(d.Clause.Exp)
			if err != nil {
				return err
			}
			switch exp.(type) {
			case True:
				err = e.eval(d.Clause.Body.Directives)
				if err != nil {
					return err
				}
			case False:
				// skip body and continue
			default:
				return fmt.Errorf("%s: condition expression must evaluate to boolean value", d.Clause.Pos)
			}
		default:
			panic(fmt.Errorf("unexpected %v (%T) directive", d, d))
		}
	}
	return nil
}

// evaluates the expression down to terminal value
func (e *Evaluator) evalExp(exp Exp) (Exp, error) {
	switch exp := exp.(type) {
	case True:
		return exp, nil
	case False:
		return exp, nil
	case String:
		return exp, nil
	case EnvSymbol:
		v, ok := e.env[exp.Name]
		if !ok {
			return nil, fmt.Errorf("%s: undefined env (%s)", exp.Pos, exp.Name)
		}
		return String{Pos: exp.Pos, Value: v}, nil
	case BinExp:
		a, err := e.evalExp(exp.Left)
		if err != nil {
			return nil, err
		}
		b, err := e.evalExp(exp.Right)
		if err != nil {
			return nil, err
		}
		return evalSimpleBinExp(exp.Op, a, b)
	default:
		panic(fmt.Errorf("unexpected %s expression %T", exp, exp))
	}
}

func evalSimpleBinExp(op bop.Kind, a, b Exp) (Exp, error) {
	switch op {
	case bop.Equal:
		switch a := a.(type) {
		case String:
			b, ok := b.(String)
			if ok && a.Value == b.Value {
				return True{}, nil
			}
			return False{}, nil
		default:
			panic(fmt.Errorf("unexpected %s expression %T", a, a))
		}
	default:
		panic(fmt.Errorf("unexpected %s operator", op))
	}
}
