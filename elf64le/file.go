package elf64le

type File struct {
	Header   Header
	Programs []*Program
	Sections []*Section
}

type Header struct {
	// Index into Sections slice. Points to section which contains
	// table of names (strings) for other sections.
	SectionNameTableIndex uint16
}

type Section struct {
	Header SectionHeader

	Data []byte
}

type Program struct {
	Header ProgramHeader

	Data []byte
}

type SectionHeader struct {
}

type ProgramHeader struct {
}
