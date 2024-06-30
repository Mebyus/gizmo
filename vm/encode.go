package vm

import (
	"io"
)

var magic = [...]byte{'K', 'U', 'B', 0xE9}

func Encode(w io.Writer, prog *Prog) error {
	e := Encoder{}
	e.encode(prog)
	_, err := w.Write(e.Bytes())
	return err
}

type Encoder struct {
	Buffer

	prog *Prog
}

func (e *Encoder) encode(prog *Prog) {
	e.prog = prog
	e.alloc()
	e.magic()
	e.version()
	e.textHeader()
	e.dataHeader()
	e.globalHeader()
	e.Append(prog.Text)
	e.Append(prog.Data)
	e.Append(prog.Global)
}

// header size = magic + version + text header + data header + global header +
const headerSize = 4 + 2 + (8 + 4) + (8 + 4) + (8 + 4)

func (e *Encoder) alloc() {
	size := headerSize + len(e.prog.Text) + len(e.prog.Data) + len(e.prog.Global)
	e.Init(size)
}

func (e *Encoder) magic() {
	e.Append(magic[:])
}

func (e *Encoder) version() {
	e.Val16(0)
}

func (e *Encoder) textHeader() {
	if len(e.prog.Text) == 0 {
		e.header(0, 0)
		return
	}
	offset := headerSize
	e.header(uint64(offset), uint32(len(e.prog.Text)))
}

func (e *Encoder) dataHeader() {
	if len(e.prog.Data) == 0 {
		e.header(0, 0)
		return
	}
	offset := headerSize + len(e.prog.Text)
	e.header(uint64(offset), uint32(len(e.prog.Data)))
}

func (e *Encoder) globalHeader() {
	if len(e.prog.Global) == 0 {
		e.header(0, 0)
		return
	}
	offset := headerSize + len(e.prog.Text) + len(e.prog.Data)
	e.header(uint64(offset), uint32(len(e.prog.Global)))
}

func (e *Encoder) header(offset uint64, size uint32) {
	e.Val64(offset)
	e.Val32(size)
}
