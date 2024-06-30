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

	m.set(dr, v)
	return nil
}
