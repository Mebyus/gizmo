package parser

import (
	"reflect"
	"testing"

	"github.com/mebyus/gizmo/ast"
	"github.com/mebyus/gizmo/ast/bop"
	"github.com/mebyus/gizmo/ast/uop"
	"github.com/mebyus/gizmo/token"
)

func lit(kind token.Kind, lit string) ast.BasicLiteral {
	return ast.BasicLiteral{
		Token: token.Token{
			Kind: kind,
			Lit:  lit,
		},
	}
}

func dec(v uint64) ast.BasicLiteral {
	return ast.BasicLiteral{
		Token: token.Token{
			Kind: token.DecimalInteger,
			Val:  v,
		},
	}
}

func flt(l string) ast.BasicLiteral {
	return lit(token.DecimalFloat, l)
}

func idn(name string) ast.Identifier {
	return ast.Identifier{Lit: name}
}

func sym(name string) ast.SymbolExpression {
	return ast.SymbolExpression{Identifier: idn(name)}
}

func par(x ast.Expression) ast.ParenthesizedExpression {
	return ast.ParenthesizedExpression{Inner: x}
}

func not(expr ast.Expression) *ast.UnaryExpression {
	return uex(uop.Not, expr)
}

func plus(expr ast.Expression) *ast.UnaryExpression {
	return uex(uop.Plus, expr)
}

func neg(expr ast.Expression) *ast.UnaryExpression {
	return uex(uop.Minus, expr)
}

func uex(kind uop.Kind, expr ast.Expression) *ast.UnaryExpression {
	return &ast.UnaryExpression{
		Operator: ast.UnaryOperator{Kind: kind},
		Inner:    expr,
	}
}

func sub(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.Sub, left, right)
}

func add(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.Add, left, right)
}

func mul(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.Mul, left, right)
}

func and(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.And, left, right)
}

func eq(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.Equal, left, right)
}

func neq(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.NotEqual, left, right)
}

func gr(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.Greater, left, right)
}

func leq(left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return bin(bop.LessOrEqual, left, right)
}

func bin(kind bop.Kind, left ast.Expression, right ast.Expression) ast.BinaryExpression {
	return ast.BinaryExpression{
		Operator: ast.BinaryOperator{Kind: kind},
		Left:     left,
		Right:    right,
	}
}

func ch(start string, parts ...ast.ChainPart) ast.ChainOperand {
	return ast.ChainOperand{
		Identifier: idn(start),
		Parts:      parts,
	}
}

func chg(parts ...ast.ChainPart) ast.ChainOperand {
	return ast.ChainOperand{
		Identifier: idn(""),
		Parts:      parts,
	}
}

func mbr(name string) ast.MemberPart {
	return ast.MemberPart{Member: idn(name)}
}

func idx(index ast.Expression) ast.IndexPart {
	return ast.IndexPart{Index: index}
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
			want: dec(42),
		},
		{
			name: "3 decimal float literal",
			str:  "42.042",
			want: flt("42.042"),
		},
		{
			name: "4 identifier",
			str:  "abc",
			want: sym("abc"),
		},
		{
			name: "6 integer in parentheses",
			str:  "(3)",
			want: par(dec(3)),
		},
		{
			name: "7 unary expression on integer",
			str:  "+49",
			want: plus(dec(49)),
		},
		{
			name: "8 unary expression on identifier",
			str:  "!is_good",
			want: not(sym("is_good")),
		},
		{
			name: "9 binary expression on integers",
			str:  "49 - 90",
			want: sub(dec(49), dec(90)),
		},
		{
			name: "10 binary expression in double parentheses",
			str:  "((49 - 90))",
			want: par(par(sub(dec(49), dec(90)))),
		},
		{
			name: "11 nil literal",
			str:  "nil",
			want: lit(token.Nil, ""),
		},
		{
			name: "12 four binary plus operators",
			str:  "0 + 1 + 2 + 3 + 4",
			want: add(
				add(
					add(
						add(dec(0), dec(1)),
						dec(2),
					),
					dec(3),
				),
				dec(4),
			),
		},
		{
			name: "13 four binary operators",
			str:  "0 + 1 + 2 + 3 * 4",
			want: add(
				add(
					add(dec(0), dec(1)),
					dec(2),
				),
				mul(
					dec(3),
					dec(4),
				),
			),
		},
		{
			name: "14 comparison and logic",
			str:  "a == 3.13 && b != 5.4",
			want: and(
				eq(sym("a"), flt("3.13")),
				neq(sym("b"), flt("5.4")),
			),
		},
		{
			name: "15 comparison and logic",
			str:  "a > 3.13 && b <= c",
			want: and(
				gr(
					sym("a"),
					flt("3.13"),
				),
				leq(
					sym("b"),
					sym("c"),
				),
			),
		},
		{
			name: "16 member",
			str:  "a.b",
			want: ch("a", mbr("b")),
		},
		{
			name: "17 index",
			str:  "a[3]",
			want: ch("a", idx(dec(3))),
		},
		{
			name: "18 index on member",
			str:  "a.b[3]",
			want: ch("a", mbr("b"), idx(dec(3))),
		},
		{
			name: "19 binary expression with parentheses",
			str:  "(1 + 1) - 0 * 2 * (-1)",
			want: sub(
				par(add(dec(1), dec(1))),
				mul(
					mul(dec(0), dec(2)),
					par(neg(dec(1))),
				),
			),
		},
		{
			name: "20 member on receiver",
			str:  "g.b",
			want: chg(mbr("b")),
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
