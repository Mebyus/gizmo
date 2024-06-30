package vm

import "fmt"

func (m *Machine) jmp(a uint32) error {
	if a >= uint32(len(m.text)) {
		return fmt.Errorf("jump to 0x%x is out of program text range", a)
	}
	m.ip = uint64(a)
	m.jump = true
	return nil
}

func (m *Machine) jumpAddr() error {
	d := m.id(4)
	a := val32(d)
	return m.jmp(a)
}
