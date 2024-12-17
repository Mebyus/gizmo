package elf64le

type File struct {
	Header   Header
	Programs []ProgramHeader
	Sections []SectionHeader

	Text Text
}

type Header struct {
	EntrypointVirtualAddress uint64

	// Index into Sections slice. Points to section which contains
	// table of names (strings) for other sections.
	SectionNameTableIndex uint16
}

type SectionHeader struct {
}

type ProgramHeader struct {
	Type  uint32
	Flags uint32

	ImageOffset     uint64
	VirtualAddress  uint64
	PhysicalAddress uint64

	ImageSize  uint64
	MemorySize uint64
	Alignment  uint64
}
