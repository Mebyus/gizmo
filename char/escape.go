package char

func Unescape(s string) string {
	// first scan the string to check if it has at least one escape sequence
	i := 0
	for ; i < len(s); i++ {
		if s[i] == '\\' {
			// escape sequence found
			break
		}
	}
	if i >= len(s) {
		// return the string "as is", because it does not have escape sequences
		return s
	}

	buf := make([]byte, i, len(s))
	copy(buf[:i], s[:i])

	for ; i < len(s); i++ {
		if s[i] == '\\' {
			i += 1
			var b byte
			switch s[i] {
			case '\\':
				b = '\\'
			case 'n':
				b = '\n'
			case 't':
				b = '\t'
			case 'r':
				b = '\r'
			case '"':
				b = '"'
			default:
				panic("bad escape sequence")
			}
			buf = append(buf, b)
		} else {
			buf = append(buf, s[i])
		}
	}

	return string(buf)
}
