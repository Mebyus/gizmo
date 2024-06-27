package vm

func (m *Machine) loadValReg() error {
	d := m.instd(9)
	r := d[0]
	v := val64(d[1:])
	return m.set(r, v)
}

func (m *Machine) loadRegReg() error {
	d := m.instd(2)
	sr := d[0]
	dr := d[1]

	v1, err := m.get(sr)
	if err != nil {
		return err
	}
	v2, err := m.get(dr)
	if err != nil {
		return err
	}

	m.unsafeSet(dr, v1+v2)
	return nil
}
