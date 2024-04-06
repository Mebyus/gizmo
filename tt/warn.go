package tt

// adds warning to merger output
func (m *Merger) warn(err error) {
	m.Warns = append(m.Warns, err)
}
