package highlight

import (
	"fmt"
	"os"

	"github.com/mebyus/gizmo/lexer"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

func Highlight(filename string, tokidx int) error {
	src, err := source.Load(filename)
	if err != nil {
		return err
	}
	lx := lexer.FromSource(src)

	var found token.Token
	i := 0
	for {
		tok := lx.Lex()

		if i == tokidx {
			found = tok
			break
		}
		i++

		if tok.IsEOF() {
			return fmt.Errorf("unable to find token at index %d, stream contains only %d tokens", tokidx, i)
		}
	}

	err = source.Render(os.Stdout, found, source.LineWindow{
		Before: 5,
		After:  5,
	})
	if err != nil {
		return err
	}
	return nil
}
