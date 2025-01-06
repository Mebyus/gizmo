package kbs

import "fmt"

type Evaluator struct {
	out ScriptOutput

	// maps env name to its value
	env map[ /* env name */ string]string
}

func (e *Evaluator) Eval(dirs []Directive) (*ScriptOutput, error) {
	err := e.eval(dirs)
	if err != nil {
		return nil, err
	}
	return &e.out, nil
}

func (e *Evaluator) eval(dirs []Directive) error {
	for _, dir := range dirs {
		var err error
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
			err = e.eval(d.Clause.Body.Directives)
		default:
			panic(fmt.Errorf("unexpected %v (%T) directive", d, d))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
