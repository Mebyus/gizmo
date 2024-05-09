package format

import (
	"fmt"

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

	// Mandatory point of whitespace separation. By default it results
	// in a single space character in output. Depending on surrounding
	// tokens it may result in:
	//
	//	- space character (default)
	//	- newline
	//	- newline + indent
	SpaceNode

	IndentNode

	// Mandatory single space character.
	StrictSpaceNode

	// Possible point of whitespace separation. By default is nop.
	// Depending on surrounding tokens it may result in:
	//
	//	- nop (default)
	//	- space
	//	- newline
	//	- newline + indent
	SepNode

	// Mandatory newline.
	NewlineNode

	// Mandatory newline with indent.
	NewlineIndentNode

	// Mandatory blank line.
	BlankNode

	TrailCommaNode
	StartBlockNode
	EndBlockNode
	StartNode
)

var kindText = [...]string{
	emptyNode: "<nil>",

	GenNode:           "gen",
	TokNode:           "tok",
	IncNode:           "inc",
	DecNode:           "dec",
	SpaceNode:         "space",
	IndentNode:        "indent",
	StrictSpaceNode:   "ss",
	SepNode:           "sep",
	NewlineNode:       "nl",
	NewlineIndentNode: "nli",
	TrailCommaNode:    "trail_comma",
	StartBlockNode:    "start_block",
	EndBlockNode:      "end_block",
	StartNode:         "start",
	BlankNode:         "blank",
}

func (k NodeKind) String() string {
	return kindText[k]
}

type Node struct {
	// Meaning of this field depends on Node.Kind:
	//
	//	- GenNode           - kind (Token.Kind) of generated token
	//	- TokNode           - ordering number of token position (Token.Pos.Num)
	//	- IncNode           - (reserved) always 0
	//	- DecNode           - (reserved) always 0
	//	- SpaceNode         - (reserved) always 0
	//	- IndentNode        - (reserved) always 0
	//	- StrictSpaceNode   - (reserved) always 0
	//	- SepNode           - (reserved) always 0
	//	- NewlineNode       - (reserved) always 0
	//	- NewlineIndentNode - (reserved) always 0
	//	- BlankNode         - (reserved) always 0
	//	- TrailCommaNode    - (reserved) always 0
	//	- StartNode         - (reserved) always 0
	//
	Val uint32

	Kind NodeKind
}

func (g *Noder) verify(kind token.Kind, num uint32) {
	exp := g.tokens[num].Kind
	if kind != exp {
		panic(fmt.Sprintf("unexpected generated token %s instead of %s", kind.String(), exp.String()))
	}
}

func (g *Noder) add(kind NodeKind, val uint32) {
	g.nodes = append(g.nodes, Node{Kind: kind, Val: val})
}

// place original token with source position information into output
func (g *Noder) tok(tok token.Token) {
	g.add(TokNode, tok.Pos.Num)
}

// place generated token into output
func (g *Noder) gen(kind token.Kind) {
	g.add(GenNode, uint32(kind))
}

// place identifier token into output
func (g *Noder) idn(idn ast.Identifier) {
	g.verify(token.Identifier, idn.Pos.Num)
	g.add(TokNode, idn.Pos.Num)
}

func (g *Noder) bop(op ast.BinaryOperator) {
	g.add(TokNode, op.Pos.Num)
}

// place generated token with source position information into output
func (g *Noder) genpos(kind token.Kind, pos source.Pos) {
	g.verify(kind, pos.Num)
	g.add(TokNode, pos.Num)
}

// place generated semicolon token into output
func (g *Noder) semi() {
	g.gen(token.Semicolon)
}

// increment indentation buffer by one level.
func (g *Noder) inc() {
	g.add(IncNode, 0)
}

// decrement indentation buffer by one level
func (g *Noder) dec() {
	g.add(DecNode, 0)
}

// add a space between tokens into output
func (g *Noder) space() {
	g.add(SpaceNode, 0)
}

// place a "strict space" which cannot be substituted with newline break,
// such space between two tokens always leads to a space character in output
func (g *Noder) ss() {
	g.add(StrictSpaceNode, 0)
}

// place a blank line into output
func (g *Noder) blank() {
	g.add(BlankNode, 0)
}

// add potential separator into output
func (g *Noder) sep() {
	g.add(SepNode, 0)
}

// place a copy of indentation buffer into output
func (g *Noder) indent() {

}

// start a new statement in output
func (g *Noder) start() {
	g.add(StartNode, 0)
}

func (g *Noder) startBlock(pos source.Pos) {
	g.verify(token.LeftCurly, pos.Num)
	g.add(StartBlockNode, 0)
}

func (g *Noder) endBlock() {
	g.add(EndBlockNode, 0)
}

// place on optional trailing comma, if next brace token is on the same line
// it will be skipped by stapler
func (g *Noder) trailComma() {
	g.add(TrailCommaNode, 0)
}

// place "pub" keyword and start new line
func (g *Noder) pub() {
	g.gen(token.Pub)
	g.nl()
}

// start new line in generated output
func (g *Noder) nl() {
	g.add(NewlineNode, 0)
}

// start new line and place indentation into output
func (g *Noder) nli() {
	g.nl()
	g.indent()
}
