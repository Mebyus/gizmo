package amd64

import "encoding/binary"

type Encoder struct {
	buf []byte

	labels []LabelMark

	patches []LabelBackpatch
}

// LabelMark marks label position in machine code.
type LabelMark struct {
	// Offset into segment with code.
	pos uint64

	// True if label position is known.
	ok bool
}

type LabelBackpatch struct {
	label uint64
	pos   uint64
}

func (e *Encoder) Bytes() []byte {
	return e.buf
}

func (e *Encoder) pos() uint64 {
	return uint64(len(e.buf))
}

func (e *Encoder) u8s(x ...byte) {
	e.buf = append(e.buf, x...)
}

func (e *Encoder) u8(x uint8) {
	e.buf = append(e.buf, x)
}

func (e *Encoder) u16(x uint16) {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], x)
	e.buf = append(e.buf, buf[:]...)
}

func (e *Encoder) u32(x uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], x)
	e.buf = append(e.buf, buf[:]...)
}

func (e *Encoder) u64(x uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	e.buf = append(e.buf, buf[:]...)
}

func (e *Encoder) putTail32(x uint32) {
	binary.LittleEndian.PutUint32(e.buf[len(e.buf)-4:], x)
}

// Accepts label index as argument.
// Returns (offset, true) relative to current encoder position.
// Returns (0, false) if label position is yet unknown.
func (e *Encoder) getRelOffset(label uint64) (int64, bool) {
	mark := e.labels[label]
	if !mark.ok {
		return 0, false
	}

	// Resulting offset is always negative in this case
	// because label mark can only be placed after encoder
	// has already passed label position. Thus len(e.buf) > mark.pos
	// is always true for this branch.
	return -int64(uint64(len(e.buf)) - mark.pos), true
}

func (e *Encoder) saveLabelBackpatchRel32(label, pos uint64) {
	e.patches = append(e.patches, LabelBackpatch{
		label: label,
		pos:   pos,
	})
}
