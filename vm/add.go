package vm

func (m *Machine) incReg() error {
	d := m.id(1)
	r := d[0]

	v, err := m.get(r)
	if err != nil {
		return err
	}

	m.unsafeSet(r, v+1)
	return nil
}

func (m *Machine) addRegReg() error {
	d := m.id(2)
	dr := d[0]
	sr := d[1]

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
