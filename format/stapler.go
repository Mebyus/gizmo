package format

import (
	"fmt"

	"github.com/mebyus/gizmo/token"
)

type Stapler struct {
	// read-only input field
	tokens []token.Token

	// accumulator output field
	nodes []Node

	// ordering number of previous token
	num int
}

func NewStapler(tokens []token.Token) *Stapler {
	return &Stapler{
		tokens: tokens,
		num:    -1,
	}
}

func (s *Stapler) Nodes() []Node {
	return s.nodes
}

func Staple(tokens []token.Token, nodes []Node) []Node {
	s := NewStapler(tokens)

	for _, node := range nodes {
		switch node.Kind {
		case GenNode:
			s.gen(node)
		case TokNode:
			s.tok(node)
		case StartBlockNode:
			s.sb()
		case EndBlockNode:
			s.eb()
		default:
			s.add(node)
		}
	}

	return s.Nodes()
}

func (s *Stapler) add(node Node) {
	s.nodes = append(s.nodes, node)
}

func (s *Stapler) node(kind NodeKind) {
	s.add(Node{Kind: kind})
}

func (s *Stapler) sb() {
	s.advance(token.LeftCurly)
	s.node(StartBlockNode)
}

func (s *Stapler) eb() {
	s.advance(token.RightCurly)
	s.node(EndBlockNode)
}

func (s *Stapler) advance(kind token.Kind) {
	i := s.num + 1
	for ; i < len(s.tokens); i += 1 {
		tok := s.tokens[i]

		if tok.Kind == token.LineComment {
			s.add(Node{Kind: TokNode, Val: uint32(i)})
		} else if tok.Kind == kind {
			break
		} else {
			panic(fmt.Sprintf("unexpected token %s", tok.Kind.String()))
		}
	}
	s.num = i
}

func (s *Stapler) addCursorToken() {
	s.add(Node{Kind: TokNode, Val: uint32(s.num)})
}

func (s *Stapler) gen(node Node) {
	kind := token.Kind(node.Val)
	s.advance(kind)
	s.addCursorToken()
}

func (s *Stapler) tok(node Node) {
	num := node.Val
	prev := s.num
	if int(num) <= prev {
		panic(fmt.Sprintf("abnormal token ordering %d => %d", prev, num))
	}
	s.num = int(num)

	if int(num) == prev+1 {
		s.add(node)
		return
	}

	for i := prev + 1; i < int(num); i += 1 {
		tok := s.tokens[i]
		if tok.Kind == token.LineComment {
			s.add(Node{Kind: TokNode, Val: uint32(i)})
		}
	}

	s.add(node)
}
