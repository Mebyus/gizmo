package parser

import (
	"reflect"
	"testing"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/oper"
	"github.com/mebyus/gizmo/token"
)

func tok(k token.Kind) token.Token {
	return token.Token{
		Kind: k,
	}
}

func lit(kind token.Kind, lit string) ast.BasicLiteral {
	return ast.BasicLiteral{
		Token: token.Token{
			Kind: kind,
			Lit:  lit,
		},
	}
}

func dint(v uint64) ast.BasicLiteral {
	return ast.BasicLiteral{
		Token: token.Token{
			Kind: token.DecimalInteger,
			Val:  v,
		},
	}
}

func dflt(l string) ast.BasicLiteral {
	return lit(token.DecimalFloat, l)
}

func idn(lit string) ast.Identifier {
	return ast.Identifier{Lit: lit}
}

func cst(lit ...string) ast.ChainStart {
	return ast.ChainStart{Identifier: sidn(lit...)}
}

func sidn(lits ...string) ast.ScopedIdentifier {
	var scopes []ast.Identifier
	if len(lits) > 1 {
		scopes = make([]ast.Identifier, 0, len(lits)-1)
		for _, lit := range lits[:len(lits)-1] {
			scopes = append(scopes, idn(lit))
		}
	}

	return ast.ScopedIdentifier{
		Scopes: scopes,
		Name:   idn(lits[len(lits)-1]),
	}
}

func subs(names ...string) ast.SubsExpression {
	return ast.SubsExpression{Identifier: sidn(names...)}
}

func par(x ast.Expression) ast.ParenthesizedExpression {
	return ast.ParenthesizedExpression{Inner: x}
}

func uex(kind token.Kind, inner ast.Expression) *ast.UnaryExpression {
	return &ast.UnaryExpression{
		Operator: oper.NewUnary(tok(kind)),
		Inner:    inner,
	}
}

func bin(kind token.Kind, left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return ast.BinaryExpression{
		Operator: oper.NewBinary(tok(kind)),
		Left:     left,
		Right:    right,
	}
}

func sel(target ast.ChainOperand, selected string) ast.SelectorExpression {
	return ast.SelectorExpression{
		Target:   target,
		Selected: idn(selected),
	}
}

func idx(target ast.ChainOperand, index ast.Expression) ast.IndexExpression {
	return ast.IndexExpression{
		Target: target,
		Index:  index,
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    ast.Expression
		wantErr bool
	}{
		{
			name:    "1 empty string",
			str:     "",
			wantErr: true,
		},
		{
			name: "2 decimal integer literal",
			str:  "42",
			want: dint(42),
		},
		{
			name: "3 decimal float literal",
			str:  "42.042",
			want: dflt("42.042"),
		},
		{
			name: "4 identifier",
			str:  "abc",
			want: subs("abc"),
		},
		{
			name: "6 integer in parentheses",
			str:  "(3)",
			want: par(dint(3)),
		},
		{
			name: "7 unary expression on integer",
			str:  "+49",
			want: uex(token.Plus, dint(49)),
		},
		{
			name: "8 unary expression on identifier",
			str:  "!is_good",
			want: uex(token.Not, subs("is_good")),
		},
		{
			name: "9 binary expression on integers",
			str:  "49 - 90",
			want: bin(token.Minus, dint(49), dint(90)),
		},
		{
			name: "10 binary expression in double parentheses",
			str:  "((49 - 90))",
			want: par(par(bin(token.Minus, dint(49), dint(90)))),
		},
		{
			name: "11 nil literal",
			str:  "nil",
			want: lit(token.Nil, ""),
		},
		{
			name: "12 four binary plus operators",
			str:  "0 + 1 + 2 + 3 + 4",
			want: bin(token.Plus, bin(token.Plus, bin(token.Plus, bin(token.Plus, dint(0), dint(1)), dint(2)), dint(3)), dint(4)),
		},
		{
			name: "13 four binary operators",
			str:  "0 + 1 + 2 + 3 * 4",
			want: bin(token.Plus,
				bin(token.Plus,
					bin(token.Plus,
						dint(0),
						dint(1),
					),
					dint(2),
				),
				bin(token.Asterisk, dint(3), dint(4)),
			),
		},
		{
			name: "14 comparison and logic",
			str:  "a == 3.13 && b != 5.4",
			want: bin(token.LogicalAnd,
				bin(token.Equal, subs("a"), dflt("3.13")),
				bin(token.NotEqual, subs("b"), dflt("5.4")),
			),
		},
		{
			name: "15 comparison and logic",
			str:  "a > 3.13 && b <= c",
			want: bin(token.LogicalAnd,
				bin(token.Greater, subs("a"), dflt("3.13")),
				bin(token.LessOrEqual, subs("b"), subs("c")),
			),
		},
		{
			name: "16 selector",
			str:  "a.b",
			want: sel(cst("a"), "b"),
		},
		{
			name: "17 index",
			str:  "a[3]",
			want: idx(cst("a"), dint(3)),
		},
		{
			name: "18 index on selector",
			str:  "a.b[3]",
			want: idx(sel(cst("a"), "b"), dint(3)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseExpression(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("\nParseExpression() error = %v\nwantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nParseExpression() = %+v\nwant %+v", got, tt.want)
			}
		})
	}
}
