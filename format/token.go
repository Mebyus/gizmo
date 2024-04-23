package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

// place original token with source position information into output
func (g *Builder) tok(tok token.Token) {

}

// place generated token into output
func (g *Builder) gen(kind token.Kind) {

}

// place identifier token into output
func (g *Builder) idn(idn ast.Identifier) {

}

func (g *Builder) bop(op ast.BinaryOperator) {

}

// place generated token with source position information into output
func (g *Builder) genpos(kind token.Kind, pos source.Pos) {

}

// place generated semicolon token into output
func (g *Builder) semi() {
	g.gen(token.Semicolon)
}

// place generated right curly brace into output
func (g *Builder) rcurl() {
	g.gen(token.RightCurly)
}

// increment indentation buffer by one level.
func (g *Builder) inc() {

}

// decrement indentation buffer by one level
func (g *Builder) dec() {

}

// add a space between tokens into output
func (g *Builder) space() {

}

// place a "strict space" which cannot be substituted with newline break,
// such space between two tokens always leads to a space character in output
func (g *Builder) ss() {

}

// add potential separator into output
func (g *Builder) sep() {

}

// place a copy of indentation buffer into output
func (g *Builder) indent() {

}

// start a new statement in output
func (g *Builder) start() {

}

// place on optional trailing comma, if next brace token is on the same line
// it will be skipped by stapler
func (g *Builder) trailComma() {

}

// place "pub" keyword and start new line
func (g *Builder) pub() {
	g.gen(token.Pub)
	g.nl()
}

// start new line in generated output
func (g *Builder) nl() {

}

// start new line and place indentation into output
func (g *Builder) nli() {
	g.nl()
	g.indent()
}

// place a blank line into output
func (g *Builder) blank() {
	g.nl()
	g.nl()
}

// place a string into output, without any additional logic
func (g *Builder) str(s string) {

}
