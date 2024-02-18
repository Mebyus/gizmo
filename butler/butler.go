package butler

import (
	"fmt"
	"io"
	"os"
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
		r.help()
		return nil
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

func (r *Lackey) help() {
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

	for _, s := range r.Sub {
		r.indent()
		r.write(s.Name)
		r.write("  ")
		r.write(s.Short)
		r.nl()
	}

	r.nl()
	r.write(fmt.Sprintf(`Use "%s help <command>" for more information about a specific command.`, r.Name))
	r.nl()
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
