package stg

import "strings"

type RenderBuffer struct {
	s strings.Builder
}

func (g *RenderBuffer) puts(s string) {
	g.s.WriteString(s)
}

func (t *Type) Render(g *RenderBuffer) {

}

func (t *Type) String() string {
	var buf RenderBuffer
	t.Render(&buf)
	return buf.s.String()
}
