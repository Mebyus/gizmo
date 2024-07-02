package vm

import (
	"fmt"
	"io"
	"os"
)

const (
	// TODO: change to 1, make value 0 a trap
	SysCallWrite = 1
)

func (m *Machine) syscall() error {
	var err error
	switch m.sc {
	case SysCallWrite:
		err = m.sysCallWrite()
	default:
		return fmt.Errorf("unknown syscall %d", m.sc)
	}
	return err
}

func (m *Machine) sysCallWrite() error {
	// r0 - select file descriptor for output
	// r1 - pointer to start of the data
	// r2 - number of bytes to write, starting from specified pointer
	var w io.Writer
	switch m.r[0] {
	case 0:
		// stdin
		return fmt.Errorf("stdin write attempt")
	case 1:
		// stdout
		w = os.Stdout
	case 2:
		// stderr
		w = os.Stderr
	default:
		return fmt.Errorf("arbitrary file descriptors are not implemented")
	}

	ptr := m.r[1]
	n := m.r[2]
	b, err := m.memslice(ptr, n)
	if err != nil {
		return err
	}

	var result uint64
	nw, err := w.Write(b)
	if err != nil {
		result = 1
	} else {
		result = 0
	}

	// sc - store result code of syscall
	// r0 - number of bytes written
	m.sc = result
	m.r[0] = uint64(nw)

	return nil
}
