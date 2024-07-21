package source

import "strconv"

// Pos represents a starting position of something (token, symbol, code block, etc.) in source code
type Pos struct {
	File *File

	// byte offset in stream
	Ofs uint32

	// order number, starts from zero
	Num uint32

	// starts from zero
	Line uint32

	// column number, starts from zero
	Col uint32
}

// Pin is an interface that represents something with position in source code
type Pin interface {
	Pin() Pos
}

func (p Pos) Pin() Pos {
	return p
}

func (p Pos) Len() uint32 {
	return 0
}

// Erase clears position information
func (p *Pos) Clear() {
	p.File = nil
	p.Ofs = 0
	p.Num = 0
	p.Line = 0
	p.Col = 0
}

// Reset sets position to stream start
func (p *Pos) Reset() {
	p.Ofs = 0
	p.Line = 0
	p.Col = 0
}

func (p Pos) Long() string {
	if p.File == nil {
		return "<unknown>:" + p.Short()
	}
	return p.File.Path + ":" + p.Short()
}

func (p Pos) String() string {
	return p.Long()
}

// Short formats Pos into string <line>:<column>. Both numbers
// are increased by 1, thus starting position appears as 1:1
func (p Pos) Short() string {
	return strconv.FormatUint(uint64(p.Line+1), 10) + ":" + strconv.FormatUint(uint64(p.Col+1), 10)
}

func (p Pos) Full() string {
	return p.Long() + "[0x" + strconv.FormatUint(uint64(p.Ofs), 16) + "]"
}
