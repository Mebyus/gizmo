package amd64

import (
	"encoding/binary"
	"fmt"
)

// LabelMap is a container for labels inside the segment.
type LabelMap struct {
	// Maps label index to its mark in segment.
	Marks []LabelMark

	// Maps label name to its index.
	Index map[string]uint64
}

func NewLabelMap() *LabelMap {
	return &LabelMap{Index: make(map[string]uint64)}
}

// Push tries to add label by its name to container.
// If label is already present inside the container then nothing happens.
// Otherwise stores and assigns an index to this label.
//
// Returns label index in both cases.
func (m *LabelMap) Push(name string) uint64 {
	i, ok := m.Index[name]
	if ok {
		return i
	}
	i = m.Size()
	m.Index[name] = i
	m.Marks = append(m.Marks, LabelMark{})
	return i
}

func (m *LabelMap) Size() uint64 {
	return uint64(len(m.Marks))
}

type Encoder struct {
	buf []byte

	Labels *LabelMap

	patches []LabelBackpatch
}

// LabelMark marks label position in machine code.
type LabelMark struct {
	// Offset into segment with code.
	pos uint64

	// True if label position is known inside the segment.
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

func (e *Encoder) putBack32At(pos uint64, x uint32) {
	binary.LittleEndian.PutUint32(e.buf[pos-4:pos], x)
}

func (e *Encoder) placeLabel(label uint64) {
	e.Labels.Marks[label].pos = e.pos()
	e.Labels.Marks[label].ok = true
}

// Accepts label index as argument.
// Returns (offset, true) relative to current encoder position.
// Returns (0, false) if label position is yet unknown.
func (e *Encoder) getRelOffset(label uint64) (int64, bool) {
	mark := e.Labels.Marks[label]
	if !mark.ok {
		return 0, false
	}

	// Resulting offset is always negative in this case
	// because label mark can only be placed after encoder
	// has already passed label position. Thus len(e.buf) > mark.pos
	// is always true for this branch.
	return -int64(e.pos() - mark.pos), true
}

func (e *Encoder) putBackRel32OffsetOrBackpatch(label uint64) {
	e.u32(0) // adjust encoder position
	offset, ok := e.getRelOffset(label)
	if ok {
		e.putTail32(uint32(offset))
	} else {
		e.saveLabelBackpatchRel32(label, e.pos())
	}
}

func (e *Encoder) saveLabelBackpatchRel32(label, pos uint64) {
	e.patches = append(e.patches, LabelBackpatch{
		label: label,
		pos:   pos,
	})
}

func (e *Encoder) backpatch() error {
	for _, p := range e.patches {
		mark := e.Labels.Marks[p.label]
		if !mark.ok {
			return fmt.Errorf("label (%d) has no placement inside the segment", p.label)
		}

		// Resulting offset is always positive in this case
		// because we a doing a backpatch. That means label
		// is placed after it was referenced.
		e.putBack32At(p.pos, uint32(mark.pos-p.pos))
	}

	return nil
}
