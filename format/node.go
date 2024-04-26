package format

import (
	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/source"
	"github.com/mebyus/gizmo/token"
)

type NodeKind uint8

const (
	emptyNode NodeKind = iota

	GenNode
	TokNode
	IncNode
	DecNode
	SpaceNode
	IndentNode
	StrictSpaceNode
	SepNode
	NewlineNode
	TrailCommaNode
	StartBlockNode
	EndBlockNode
	StartNode
)

var kindText = [...]string{
	emptyNode: "<nil>",

	GenNode:         "gen",
	TokNode:         "tok",
	IncNode:         "inc",
	DecNode:         "dec",
	SpaceNode:       "space",
	IndentNode:      "indent",
	StrictSpaceNode: "ss",
	SepNode:         "sep",
	NewlineNode:     "nl",
	TrailCommaNode:  "tc",
	StartBlockNode:  "sb",
	EndBlockNode:    "eb",
	StartNode:       "start",
}

func (k NodeKind) String() string {
	return kindText[k]
}

type Node interface {
	Node()
	Kind() NodeKind
}

type node struct{}

func (node) Node() {}

type Gen struct {
	node

	TokenKind token.Kind
}

func (Gen) Kind() NodeKind {
	return GenNode
}

type Tok struct {
	node

	Token token.Token
}

func (Tok) Kind() NodeKind {
	return TokNode
}

type Start struct {
	node
}

func (Start) Kind() NodeKind {
	return StartNode
}

type StrictSpace struct {
	node
}

func (StrictSpace) Kind() NodeKind {
	return StrictSpaceNode
}

type Inc struct {
	node
}

func (Inc) Kind() NodeKind {
	return IncNode
}

type Dec struct {
	node
}

func (Dec) Kind() NodeKind {
	return DecNode
}

type StartBlock struct {
	node
}

func (StartBlock) Kind() NodeKind {
	return StartBlockNode
}

type EndBlock struct {
	node
}

func (EndBlock) Kind() NodeKind {
	return EndBlockNode
}

// place original token with source position information into output
func (g *Builder) tok(tok token.Token) {
	g.s.add(Tok{Token: tok})
}

// place generated token into output
func (g *Builder) gen(kind token.Kind) {
	g.s.add(Gen{TokenKind: kind})
}

// place identifier token into output
func (g *Builder) idn(idn ast.Identifier) {
	g.s.add(Tok{Token: idn.Token()})
}

func (g *Builder) bop(op ast.BinaryOperator) {

}

// place generated token with source position information into output
func (g *Builder) genpos(kind token.Kind, pos source.Pos) {
	g.s.add(Tok{
		Token: token.Token{
			Kind: kind,
			Pos:  pos,
		},
	})
}

// place generated semicolon token into output
func (g *Builder) semi() {
	g.gen(token.Semicolon)
}

// increment indentation buffer by one level.
func (g *Builder) inc() {
	g.s.add(Inc{})
}

// decrement indentation buffer by one level
func (g *Builder) dec() {
	g.s.add(Dec{})
}

// add a space between tokens into output
func (g *Builder) space() {

}

// place a "strict space" which cannot be substituted with newline break,
// such space between two tokens always leads to a space character in output
func (g *Builder) ss() {
	g.s.add(StrictSpace{})
}

// add potential separator into output
func (g *Builder) sep() {

}

// place a copy of indentation buffer into output
func (g *Builder) indent() {

}

// start a new statement in output
func (g *Builder) start() {
	g.s.add(Start{})
}

func (g *Builder) startBlock(pos source.Pos) {
	g.s.add(StartBlock{})
}

func (g *Builder) endBlock() {
	g.s.add(EndBlock{})
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
