package gencpp

// Builder is used to generate formatted output (C++ code)
//
// Output is stored in memory
type Builder struct {
	// Stored output
	out []byte

	// Identation buffer
	//
	// Stores sequence of bytes which is used for indenting current line
	// in output. When a new line starts this buffer is used to add identation
	ibuf []byte
}

// NewBuilder create a ready to use builder with specified initial
// preallocated memory capacity
func NewBuilder(size int) *Builder {
	return &Builder{
		out: make([]byte, 0, size),
	}
}

// Bytes returns stored output in form of raw bytes
func (g *Builder) Bytes() []byte {
	return g.out
}

// Reset builder to a clean state, but keep underlying memory
// for further usage
func (g *Builder) Reset() {
	g.out = g.out[:0]
	g.ibuf = g.ibuf[:0]
}

func (g *Builder) write(s string) {
	g.out = append(g.out, []byte(s)...)
}

// writes a single byte to output
func (g *Builder) wb(b byte) {
	g.out = append(g.out, b)
}

func (g *Builder) inc() {
	g.ibuf = append(g.ibuf, '\t')
}

func (g *Builder) dec() {
	g.ibuf = g.ibuf[:len(g.ibuf)-1]
}

func (g *Builder) comment(s string) {
	g.write("// ")
	g.write(s)
	g.nl()
}

func (g *Builder) space() {
	g.wb(' ')
}

// start new line
func (g *Builder) nl() {
	g.wb('\n')
}

// add indentation to current line
func (g *Builder) indent() {
	g.out = append(g.out, g.ibuf...)
}
