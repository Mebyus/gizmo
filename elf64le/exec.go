package elf64le

// Text represents contents of text section.
type Text struct {
	// Raw binary contents of this section.
	Data []byte

	// Position of entrypoint inside Data slice.
	EntrypointOffset uint64

	// Where this section Data will be placed in virtual memory.
	VirtualAddress uint64

	// Power of 2. Typically equals page size 0x1000 (4 KB).
	Alignment uint64
}

func makeProgramHeaderFromText(text *Text, offset uint64) ProgramHeader {
	return ProgramHeader{
		Type:  SegmentTypeLoad,
		Flags: SegmentFlagExec | SegmentFlagRead,

		ImageOffset:    offset,
		VirtualAddress: text.VirtualAddress,

		ImageSize:  uint64(len(text.Data)),
		MemorySize: uint64(len(text.Data)),
		Alignment:  text.Alignment,
	}
}
