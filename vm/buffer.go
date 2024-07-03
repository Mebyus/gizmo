package vm

// Buffer is a helper object for convenient binary
// data encoding in little endian format.
type Buffer struct {
	buf []byte
	pos int
}

func (g *Buffer) Init(size int) {
	g.buf = make([]byte, size)
}

func (g *Buffer) Bytes() []byte {
	return g.buf[:g.pos]
}

func (g *Buffer) Add(b []byte) {
	n := copy(g.buf[g.pos:], b)
	if n != len(b) {
		panic("buffer size is too small")
	}
	g.pos += n
}

func (g *Buffer) AddStr(s string) {
	n := copy(g.buf[g.pos:], s)
	if n != len(s) {
		panic("buffer size is too small")
	}
	g.pos += n
}

func (g *Buffer) Pad(n int) {
	g.pos += n
}

func (g *Buffer) Align8() {
	n := AlignBy8(uint64(g.pos)) - uint64(g.pos)
	g.Pad(int(n))
}

func (g *Buffer) Align16() {
	n := AlignBy16(uint64(g.pos)) - uint64(g.pos)
	g.Pad(int(n))
}

func (g *Buffer) Val8(v uint8) {
	g.buf[g.pos] = v
	g.pos += 1
}

func (g *Buffer) Val16(v uint16) {
	b := g.buf[g.pos : g.pos+2]
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	g.pos += 2
}

func (g *Buffer) Val32(v uint32) {
	b := g.buf[g.pos : g.pos+4]
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	g.pos += 4
}

func (g *Buffer) Val64(v uint64) {
	b := g.buf[g.pos : g.pos+8]
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	g.pos += 8
}

func AlignBy16(v uint64) uint64 {
	a := v & 0xf
	a = ((^a) + 1) & 0xf
	return v + a
}

func AlignSizeBy16(v uint32) uint32 {
	return uint32(AlignBy16(uint64(v)))
}

func AlignBy8(v uint64) uint64 {
	a := v & 0x7
	a = ((^a) + 1) & 0x7
	return v + a
}
