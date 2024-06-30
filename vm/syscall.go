package vm

import "fmt"

const (
	// TODO: change to 1, make value 0 a trap
	SysCallWrite = 0
)

func (m *Machine) syscall() error {
	var err error
	switch m.sn {
	case SysCallWrite:
		err = m.sysCallWrite()
	default:
		return fmt.Errorf("unknown syscall %d", m.sn)
	}
	return err
}

func (m *Machine) sysCallWrite() error {
	// r0 - select file descriptor for output
	// r1 - pointer to start of the data
	// r2 - number of bytes to write, starting from specified pointer
	switch m.r[0] {
	case 0:
		// stdin
		return fmt.Errorf("stdin write attempt")
	case 1:
		// stdout
	case 2:
		// stderr
	default:
		return fmt.Errorf("arbitrary file descriptors are not implemented")
	}
	return nil
}
