package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func Decode(r io.Reader) (*Prog, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(b) < headerSize {
		return nil, fmt.Errorf("file is too small: %d bytes", len(b))
	}

	if !bytes.Equal(b[:4], magic[:]) {
		return nil, fmt.Errorf("file is not a kuvm bytecode program")
	}

	textOffset := binary.LittleEndian.Uint64(b[6:14])
	textLen := binary.LittleEndian.Uint32(b[14:18])
	if uint64(len(b)) < textOffset+uint64(textLen) {
		return nil, fmt.Errorf("not enough bytes to read program text")
	}
	text := b[textOffset : textOffset+uint64(textLen)]

	dataOffset := binary.LittleEndian.Uint64(b[18:26])
	dataLen := binary.LittleEndian.Uint32(b[26:34])
	if uint64(len(b)) < dataOffset+uint64(dataLen) {
		return nil, fmt.Errorf("not enough bytes to read program data")
	}
	data := b[dataOffset : dataOffset+uint64(dataLen)]

	return &Prog{
		Text: text,
		Data: data,
	}, nil
}
