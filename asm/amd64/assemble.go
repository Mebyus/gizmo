package amd64

import (
	"os"

	"github.com/mebyus/gizmo/elf64le"
)

func AssembleExecutableFromSourceFile(source string, out string) error {
	p, err := NewParserFromFile(source)
	if err != nil {
		return err
	}

	var e Encoder
	e.Labels = p.Labels

	var o Operation
	for p.Next(&o) {
		c := translate(&o)
		c.do(&e)
	}
	err = p.Error
	if err != nil {
		return err
	}
	err = e.backpatch()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o775)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := elf64le.NewEncoder(&elf64le.File{
		Text: elf64le.Text{
			// Data: []byte{
			// 	0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, // mov    $0x3c,%rax
			// 	0x48, 0xc7, 0xc7, 0x05, 0x00, 0x00, 0x00, // mov    $0x5,%rdi
			// 	0x0f, 0x05, // syscall
			// },
			Data: e.Bytes(),

			EntrypointOffset: 0,
			VirtualAddress:   0x400000,
			Alignment:        0x1000,
		},
	})
	return encoder.Encode(file)
}
