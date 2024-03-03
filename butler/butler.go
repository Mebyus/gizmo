package butler

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Lackey struct {
	// Fields for supplying user help

	Short string
	Usage string
	Full  string

	// Dispatching fields

	Name string
	Sub  []*Lackey
	Def  *Lackey

	// Execution fields

	// This function will be checked first before dispatching
	// to sublackeys. Must be nil if lackey is a dispatcher
	Exec func(r *Lackey, args []string) error

	// Template for parameters (flags) parsing
	Params Params
}

func (r *Lackey) Run(args []string) error {
	if r.Exec != nil {
		if r.Params == nil {
			return r.Exec(r, args)
		}
		unboundArgs, err := Parse(r.Params, args)
		if err != nil {
			return err
		}
		return r.Exec(r, unboundArgs)
	}

	if len(args) == 0 {
		return r.def()
	}
	name := args[0]
	if name == "help" {
		if len(args) == 1 {
			r.help()
			return nil
		}
		subName := args[1]
		return r.subHelp(subName)
	}
	for _, s := range r.Sub {
		if name == s.Name {
			return s.Run(args[1:])
		}
	}
	return r.unknown(name)
}

// Run default lackey
//
// If default lackey is not provided then help lackey is run
func (r *Lackey) def() error {
	if r.Def != nil {
		return r.Def.Run(nil)
	}
	r.help()
	return nil
}

func (r *Lackey) subHelp(name string) error {
	if name == "" {
		r.help()
		return nil
	}

	for _, s := range r.Sub {
		if name == s.Name {
			s.help()
			return nil
		}
	}
	return fmt.Errorf("unknown help topic \"%s\"", name)
}

func (r *Lackey) execHelp() {
	r.write(r.Short)
	r.nl()
	r.nl()

	r.write("Usage:")
	r.nl()
	r.nl()
	r.indent()
	r.write(r.Usage)
	r.nl()
	r.nl()

	if r.Full != "" {
		r.write(r.Full)
		r.nl()
		r.nl()
	}

	if r.Params == nil {
		return
	}

	r.write("Available options:")
	r.nl()
	r.nl()
	params := r.Params.Recipe()
	for _, param := range params {
		r.indent()
		r.write("--")
		r.write(param.Name)
		r.write(" (")
		r.write(param.Kind.String())
		r.write(")")
		r.write("    ")
		r.write(param.Desc)
		r.nl()
	}
	r.nl()
}

func (r *Lackey) help() {
	if r.Exec != nil {
		r.execHelp()
		return
	}

	r.write(r.Short)
	r.nl()
	r.nl()

	r.write("Usage:")
	r.nl()
	r.nl()
	r.indent()
	r.write(r.Usage)
	r.nl()
	r.nl()

	r.write("Available commands:")
	r.nl()
	r.nl()

	subNamesColumnWidth := r.maxSubNameWidth() + 4
	for _, s := range r.Sub {
		r.indent()
		r.write(s.Name)
		r.write(strings.Repeat(" ", subNamesColumnWidth-len(s.Name)))
		r.write(s.Short)
		r.nl()
	}

	r.nl()
	r.write(fmt.Sprintf(`Use "%s help <command>" for more information about a specific command.`, r.Name))
	r.nl()
}

func (r *Lackey) maxSubNameWidth() int {
	m := 0
	for _, s := range r.Sub {
		if m < len(s.Name) {
			m = len(s.Name)
		}
	}
	return m
}

func (r *Lackey) nl() {
	r.write("\n")
}

func (r *Lackey) indent() {
	r.write("\t")
}

func (r *Lackey) write(s string) {
	io.WriteString(os.Stderr, s)
}

func (r *Lackey) unknown(name string) error {
	return fmt.Errorf("%s does not recognize '%s' command", r.Name, name)
}
