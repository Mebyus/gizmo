package vm

import (
	"encoding/binary"
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
	buf []byte

	prog *Prog

	pos int
}

func (e *Encoder) encode(prog *Prog) {
	e.prog = prog
	e.alloc()
	e.magic()
	e.version()
	e.textHeader()
	e.dataHeader()
	e.globalHeader()
	e.bytes(prog.Text)
	e.bytes(prog.Data)
	e.bytes(prog.Global)
}

// header size = magic + version + text header + data header + global header +
const headerSize = 4 + 2 + (8 + 4) + (8 + 4) + (8 + 4)

func (e *Encoder) alloc() {
	size := headerSize + len(e.prog.Text) + len(e.prog.Data) + len(e.prog.Global)
	e.buf = make([]byte, size)
}

func (e *Encoder) magic() {
	e.pos += copy(e.buf, magic[:])
}

func (e *Encoder) version() {
	e.u16(0)
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
	e.u64(offset)
	e.u32(size)
}

func (e *Encoder) bytes(b []byte) {
	e.pos += copy(e.buf[e.pos:], b)
}

func (e *Encoder) u8(v uint8) {
	e.buf[e.pos] = v
	e.pos += 1
}

func (e *Encoder) u16(v uint16) {
	b := e.buf[e.pos : e.pos+2]
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	e.pos += 2
}

func (e *Encoder) u32(v uint32) {
	b := e.buf[e.pos : e.pos+4]
	binary.LittleEndian.PutUint32(b, v)
	e.pos += 4
}

func (e *Encoder) u64(v uint64) {
	b := e.buf[e.pos : e.pos+8]
	binary.LittleEndian.PutUint64(b, v)
	e.pos += 8
}

func (e *Encoder) Bytes() []byte {
	return e.buf
}
