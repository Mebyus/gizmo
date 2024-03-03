package butler

import (
	"fmt"
	"strings"
)

type ParamKind uint32

const (
	empty ParamKind = iota

	Bool
	String
	List
	CommaList
)

var kindText = [...]string{
	empty: "<nil>",

	Bool:      "bool",
	String:    "string",
	List:      "list",
	CommaList: "comma list",
}

func (k ParamKind) String() string {
	return kindText[k]
}

type Param struct {
	Desc string

	ValidValues []string

	Def any

	// stored bound param values
	val any

	Name string

	Kind ParamKind
}

func (p *Param) Bind(val string) error {
	switch p.Kind {
	case empty:
		panic("unspecified kind")
	case String:
		p.val = val
	default:
		panic(fmt.Sprintf("unxpected param kind: {%d}", p.Kind))
	}
	return nil
}

type Params interface {
	Apply(*Param) error
	Recipe() []Param
}

// Parse parses supplied arguments according to params container
//
// If parsing was successful returns unbound args and error otherwise.
// Param values which were found in arguments are applied to params container
func Parse(params Params, args []string) ([]string, error) {
	pms := params.Recipe()
	if len(pms) == 0 {
		panic("params container must be nil or have at least one param in recipe")
	}

	// maps param name to its description
	var m map[string]*Param

	if len(pms) != 0 {
		m = make(map[string]*Param, len(pms))
		for _, p := range pms {
			if p.Name == "" {
				panic("param name cannot be empty")
			}
			if p.Kind == empty {
				panic(fmt.Sprintf("param {%s} is of unspecified kind", p.Name))
			}
			_, ok := m[p.Name]
			if ok {
				panic(fmt.Sprintf("param {%s} is not unique", p.Name))
			}

			allocParam := p
			m[p.Name] = &allocParam
		}
	}

	unbound, err := bindArgsToParams(m, args)
	if err != nil {
		return nil, err
	}

	for _, p := range m {
		err := params.Apply(p)
		if err != nil {
			return nil, fmt.Errorf("apply param {%s} value {%v}: %w", p.Name, p.val, err)
		}
	}

	return unbound, nil
}

func bindArgsToParams(m map[string]*Param, args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, nil
	}

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "" {
			panic("empty arg")
		}

		if !strings.HasPrefix(arg, "--") {
			// end parsing because we encountered first unbound argument
			return args[i:], nil
		}

		arg = strings.TrimPrefix(arg, "--")
		j := strings.Index(arg, "=")
		if j < 0 {
			// probably we need to check next arg for param value
			panic("not implemented")
		} else {
			name := arg[:j]
			value := arg[j+1:]

			if name == "" || value == "" {
				return nil, fmt.Errorf("invalid parameter syntax: {%s}", arg)
			}

			p, ok := m[name]
			if !ok {
				return nil, fmt.Errorf("unknown param: {%s}", name)
			}
			err := p.Bind(value)
			if err != nil {
				return nil, fmt.Errorf("bind param {%s} value {%s}: %w", name, value, err)
			}
		}

		i += 1
	}

	return nil, nil
}

func (p *Param) Bool() bool {
	if p.val == nil {
		if p.Def == nil {
			return false
		}
		return p.Def.(bool)
	}
	return p.val.(bool)
}

func (p *Param) List() []string {
	if p.val == nil {
		if p.Def == nil {
			return nil
		}
		return p.Def.([]string)
	}
	return p.val.([]string)
}

func (p *Param) Str() string {
	if p.val == nil {
		if p.Def == nil {
			return ""
		}
		return p.Def.(string)
	}
	return p.val.(string)
}
