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

func (m *Machine) jumpAddrNotZero() error {
	d := m.id(4)
	a := val32(d)

	if m.zeroFlag() {
		// no jump
		return nil
	}
	return m.jmp(a)
}

func (m *Machine) zeroFlag() bool {
	return m.cf&1 == 1
}

func (m *Machine) testRegVal() error {
	d := m.id(9)
	r := d[0]
	v2 := val64(d[1:])

	v1, err := m.get(r)
	if err != nil {
		return err
	}

	// TODO: set flags without branching
	v := v1 - v2
	if v == 0 {
		m.cf = 1
	} else {
		m.cf = 0
	}

	return nil
}
