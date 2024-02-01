package token

import "fmt"

// Compare returns a non-nil error if given tokens are not equal
func Compare(a, b Token) error {
	if a.Kind != b.Kind {
		return fmt.Errorf("kind: got %s, want %s", a.Kind.String(), b.Kind.String())
	}
	if a.Pos != b.Pos {
		if a.Pos.Ofs != b.Pos.Ofs {
			return fmt.Errorf("pos.ofs: got %d, want %d", a.Pos.Ofs, b.Pos.Ofs)
		}
		return fmt.Errorf("pos: got %s, want %s", a.Pos.String(), b.Pos.String())
	}
	if a.Kind.hasStaticLiteral() {
		return nil
	}
	if a.Lit == "" && b.Lit == "" {
		if a.Val != b.Val {
			return fmt.Errorf("val: got %d, want %d", a.Val, b.Val)
		}
		return nil
	}
	if a.Lit != b.Lit {
		return fmt.Errorf("lit: got %s, want %s", a.Lit, b.Lit)
	}
	return nil
}
