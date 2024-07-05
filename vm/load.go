package vm

func (m *Machine) loadValReg() error {
	d := m.id(9)
	r := d[0]
	v := val64(d[1:])
	return m.set(r, v)
}

func (m *Machine) loadRegReg() error {
	d := m.id(2)
	dr := d[0]
	sr := d[1]

	v, err := m.get(sr)
	if err != nil {
		return err
	}

	return m.set(dr, v)
}

func (m *Machine) loadValSysReg() {
	d := m.id(4)
	v := val32(d)
	m.sc = uint64(v)
}

func (m *Machine) clearReg() error {
	d := m.id(1)
	r := d[0]
	return m.set(r, 0)
}
