package elf64le

import (
	"bytes"
	"encoding/binary"
	"io"
)

func (f *File) Encode(w io.Writer) error {
	return nil
}

func (f *File) EncodeBytes() ([]byte, error) {
	encoder := NewEncoder(f)
	encoder.encode()
	return encoder.Bytes(), nil
}

type Encoder struct {
	Map FileMap

	buf []byte

	file *File
}

func NewEncoder(file *File) *Encoder {
	e := &Encoder{file: file}
	e.mapFile()
	e.prepareProgramHeaders()
	e.buf = make([]byte, 0, e.Map.Size)
	return e
}

func (e *Encoder) Encode(w io.Writer) error {
	e.encode()
	_, err := io.Copy(w, bytes.NewReader(e.Bytes()))
	return err
}

/*
FileMap is a helper struct which contains encoder data about to-be-written ELF.
Calculated before writing the file itself.

For now this comment also serves as documentation storage about ELF structure.
This info is gathered from all kinds of other places to be near encoder
implementation.

=======================
file header
=======================
program headers
=======================
section headers
=======================
data
=======================
*/
type FileMap struct {
	// Components = Programs + Sections

	Programs []FileComponentMap
	Sections []FileComponentMap

	// Total size of the resulting file.
	Size uint64

	ProgramHeadersOffset uint64
	SectionHeadersOffset uint64

	// Number of bytes into the file where last header ends.
	HeadersEndOffset uint64

	// Number of bytes into the file where programs data will be stored.
	ProgramDataOffset uint64

	// Number of bytes into the file where sections data will be stored.
	SectionDataOffset uint64

	// Number of bytes needed to store all programs data. Includes
	// possible paddings between programs.
	TotalProgramDataSize uint64

	// Number of bytes needed to store all sections data. Includes
	// possible paddings between sections.
	TotalSectionDataSize uint64

	TotalProgramHeadersSize uint32
	TotalSectionHeadersSize uint32

	// Number of bytes between end of components headers and start of components data.
	DataAlignPaddingSize uint32

	ProgramHeadersCount uint16
	SectionHeadersCount uint16
}

type FileComponentMap struct {
	HeaderOffset uint64
	Offset       uint64

	FileSize uint32

	// FileSize + possible alignment padding
	AlignedSize uint32

	AlignPaddingSize uint32
}

// performs all calculations which are necessary to determine
// headers and data offsets and sizes within file
func (e *Encoder) mapFile() {
	e.Map.calc(e.file)
}

const (
	SegmentTypeLoad = 1

	SegmentFlagExec  = 1
	SegmentFlagWrite = 2
	SegmentFlagRead  = 4
)

func (e *Encoder) prepareProgramHeaders() {
	offset := align(e.Map.HeadersEndOffset, e.file.Text.Alignment)

	e.file.Programs = append(e.file.Programs, makeProgramHeaderFromText(&e.file.Text, offset))
	e.file.Header.EntrypointVirtualAddress = e.file.Text.VirtualAddress + e.file.Text.EntrypointOffset
}

func (m *FileMap) calc(f *File) {
	m.ProgramHeadersCount = 1
	m.SectionHeadersCount = 0

	m.TotalProgramHeadersSize = ProgramHeadersSize * uint32(m.ProgramHeadersCount)
	m.TotalSectionHeadersSize = SectionHeadersSize * uint32(m.SectionHeadersCount)

	m.ProgramHeadersOffset = FileHeaderSize

	if m.SectionHeadersCount != 0 {
		m.SectionHeadersOffset = m.ProgramHeadersOffset + uint64(m.TotalProgramHeadersSize)
		m.HeadersEndOffset = m.SectionHeadersOffset + uint64(m.TotalSectionHeadersSize)
	} else {
		m.HeadersEndOffset = m.ProgramHeadersOffset
	}

	m.ProgramDataOffset = alignBy16(m.HeadersEndOffset)
	m.DataAlignPaddingSize = uint32(m.ProgramDataOffset - m.HeadersEndOffset)

	var programs []FileComponentMap
	var sections []FileComponentMap

	headerOffset := m.ProgramHeadersOffset
	offset := m.ProgramDataOffset

	if len(f.Programs) != 0 {
		programs = make([]FileComponentMap, 0, len(f.Programs))

		for _, _ = range f.Programs {
			size := uint32(0)
			alignedSize := alignSizeBy16(size)

			programs = append(programs, FileComponentMap{
				FileSize:         size,
				Offset:           offset,
				AlignedSize:      alignedSize,
				HeaderOffset:     headerOffset,
				AlignPaddingSize: alignedSize - size,
			})

			offset += uint64(alignedSize)
			headerOffset += ProgramHeadersSize
		}
	}

	m.TotalProgramDataSize = offset - m.ProgramDataOffset
	m.SectionDataOffset = offset

	headerOffset = m.SectionHeadersOffset

	if len(f.Sections) != 0 {
		sections = make([]FileComponentMap, 0, len(f.Sections))

		for _, _ = range f.Sections {
			size := uint32(0)
			alignedSize := alignSizeBy16(size)

			programs = append(programs, FileComponentMap{
				FileSize:         size,
				Offset:           offset,
				AlignedSize:      alignedSize,
				HeaderOffset:     headerOffset,
				AlignPaddingSize: alignedSize - size,
			})

			offset += uint64(alignedSize)
			headerOffset += SectionHeadersSize
		}
	}

	m.TotalSectionDataSize = offset - m.SectionDataOffset

	m.Programs = programs
	m.Sections = sections
	m.Size = offset
}

func (e *Encoder) encode() {
	e.header()
	e.programHeaders()
	e.sectionHeaders()
	// e.addDataAlignPadding()
	e.fileData()
}

func (e *Encoder) Bytes() []byte {
	return e.buf[:]
}

func alignBy16(v uint64) uint64 {
	a := v & 0xf
	a = ((^a) + 1) & 0xf
	return v + a
}

func align(v, a uint64) uint64 {
	mask := a - 1
	x := v & mask
	x = ((^x) + 1) & mask
	return v + x
}

func alignSizeBy16(v uint32) uint32 {
	a := v & 0xf
	a = ((^a) + 1) & 0xf
	return v + a
}

func (e *Encoder) programHeaders() {
	for i := 0; i < len(e.file.Programs); i += 1 {
		e.programHeader(i)
	}
}

func (e *Encoder) sectionHeaders() {
	for i := 0; i < len(e.file.Sections); i += 1 {
		e.sectionHeader(i)
	}
}

func (e *Encoder) programHeader(i int) {
	h := &e.file.Programs[i]

	e.u32(h.Type)
	e.u32(h.Flags)

	e.u64(h.ImageOffset)
	e.u64(h.VirtualAddress)
	e.u64(h.PhysicalAddress)

	e.u64(h.ImageSize)
	e.u64(h.MemorySize)
	e.u64(h.Alignment)
}

func (e *Encoder) sectionHeader(i int) {
	e.pad(64)
}

func (e *Encoder) addDataAlignPadding() {
	e.pad(e.Map.DataAlignPaddingSize)
}

func (e *Encoder) fileData() {
	e.pad(uint32(e.file.Programs[0].ImageOffset - uint64(len(e.buf))))
	e.data(e.file.Text.Data)
}

// elf file header
//
// 64 bytes
func (e *Encoder) header() {
	e.ident()
	e.fileType()
	e.machine()
	e.version()
	e.entrypoint()
	e.programHeadersOffset()
	e.sectionHeadersOffset()
	e.fileFlags()
	e.fileHeaderSize()
	e.programHeadersSize()
	e.programHeadersNumber()
	e.sectionHeadersSize()
	e.sectionHeadersNumber()
	e.sectionNamesHeaderIndex()
}

const (
	FileHeaderSize     = 64
	ProgramHeadersSize = 56
	SectionHeadersSize = 64
)

func (e *Encoder) fileHeaderSize() {
	e.u16(FileHeaderSize)
}

// size of one program header
func (e *Encoder) programHeadersSize() {
	e.u16(ProgramHeadersSize)
}

// size of one section header
func (e *Encoder) sectionHeadersSize() {
	e.u16(SectionHeadersSize)
}

func (e *Encoder) programHeadersNumber() {
	n := uint16(len(e.file.Programs))
	e.u16(n)
}

func (e *Encoder) sectionHeadersNumber() {
	n := uint16(len(e.file.Sections))
	e.u16(n)
}

func (e *Encoder) sectionNamesHeaderIndex() {
	e.u16(0)
}

// elf identification
//
// 16 bytes
func (e *Encoder) ident() {
	e.magic()
	e.class()
	e.endianness()
	e.identVersion()
	e.osABI()
	e.abiVersion()
	e.pad(7)
}

func (e *Encoder) magic() {
	e.u8(0x7f)
	e.u8('E')
	e.u8('L')
	e.u8('F')
}

func (e *Encoder) class() {
	const class64bit = 0x02
	e.u8(class64bit)
}

func (e *Encoder) endianness() {
	const lsb = 0x01 // Least Significant Byte First (Little Endian)
	e.u8(lsb)
}

// elf identification version
func (e *Encoder) identVersion() {
	e.u8(0x01)
}

func (e *Encoder) osABI() {
	const systemV = 0x00
	e.u8(systemV)
}

func (e *Encoder) abiVersion() {
	e.u8(0x00)
}

func (e *Encoder) fileType() {
	const executable = 2
	e.u16(executable)
}

func (e *Encoder) machine() {
	const amd64 = 0x3E
	e.u16(amd64)
}

func (e *Encoder) version() {
	e.u32(1)
}

func (e *Encoder) entrypoint() {
	e.u64(e.file.Header.EntrypointVirtualAddress)
}

func (e *Encoder) programHeadersOffset() {
	e.u64(e.Map.ProgramHeadersOffset)
}

func (e *Encoder) sectionHeadersOffset() {
	e.u64(e.Map.SectionHeadersOffset)
}

func (e *Encoder) fileFlags() {
	e.u32(0)
}

func (e *Encoder) data(b []byte) {
	e.buf = append(e.buf, b...)
}

func (e *Encoder) pad(n uint32) {
	for range n {
		e.buf = append(e.buf, 0)
	}
}

func (e *Encoder) u8(v uint8) {
	e.buf = append(e.buf, v)
}

func (e *Encoder) u16(v uint16) {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	e.buf = append(e.buf, buf[:]...)
}

func (e *Encoder) u32(v uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	e.buf = append(e.buf, buf[:]...)
}

func (e *Encoder) u64(v uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	e.buf = append(e.buf, buf[:]...)
}
