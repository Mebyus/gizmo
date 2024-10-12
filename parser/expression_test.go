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
			Kind: token.DecInteger,
			Val:  v,
		},
	}
}

func flt(l string) ast.BasicLiteral {
	return lit(token.DecFloat, l)
}

func word(name string) ast.Identifier {
	return ast.Identifier{Lit: name}
}

func sym(name string) ast.SymbolExp {
	return ast.SymbolExp{Identifier: word(name)}
}

func par(x ast.Exp) ast.ParenExp {
	return ast.ParenExp{Inner: x}
}

func not(expr ast.Exp) *ast.UnaryExp {
	return uex(uop.Not, expr)
}

func plus(expr ast.Exp) *ast.UnaryExp {
	return uex(uop.Plus, expr)
}

func neg(expr ast.Exp) *ast.UnaryExp {
	return uex(uop.Minus, expr)
}

func uex(kind uop.Kind, expr ast.Exp) *ast.UnaryExp {
	return &ast.UnaryExp{
		Operator: ast.UnaryOperator{Kind: kind},
		Inner:    expr,
	}
}

func sub(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.Sub, left, right)
}

func add(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.Add, left, right)
}

func mul(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.Mul, left, right)
}

func and(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.And, left, right)
}

func eq(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.Equal, left, right)
}

func neq(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.NotEqual, left, right)
}

func gr(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.Greater, left, right)
}

func ls(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.Less, left, right)
}

func leq(left ast.Exp, right ast.Exp) ast.BinExp {
	return bin(bop.LessOrEqual, left, right)
}

func bin(kind bop.Kind, left ast.Exp, right ast.Exp) ast.BinExp {
	return ast.BinExp{
		Operator: ast.BinaryOperator{Kind: kind},
		Left:     left,
		Right:    right,
	}
}

func ch(start string, parts ...ast.ChainPart) ast.ChainOperand {
	return ast.ChainOperand{
		Identifier: word(start),
		Parts:      parts,
	}
}

func sel(name string) ast.SelectPart {
	return ast.SelectPart{Name: word(name)}
}

func idx(index ast.Exp) ast.IndexPart {
	return ast.IndexPart{Index: index}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    ast.Exp
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
			name: "16 select",
			str:  "a.b",
			want: ch("a", sel("b")),
		},
		{
			name: "17 index",
			str:  "a[3]",
			want: ch("a", idx(dec(3))),
		},
		{
			name: "18 index on member",
			str:  "a.b[3]",
			want: ch("a", sel("b"), idx(dec(3))),
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
			name: "20 select on receiver",
			str:  "g.b",
			want: ch("g", sel("b")),
		},
		{
			name:    "21 unfinished binary expression",
			str:     "1 +",
			wantErr: true,
		},
		{
			name: "22 binary expression (3 operands)",
			str:  "a + b - 3",
			want: sub(
				add(sym("a"), sym("b")),
				dec(3),
			),
		},
		{
			name: "23 binary expression (3 operands)",
			str:  "a + b * (3 + 1)",
			want: add(
				sym("a"),
				mul(
					sym("b"),
					par(add(dec(3), dec(1))),
				),
			),
		},
		{
			name: "24 binary expression (3 operands)",
			str:  "a * 9 + b",
			want: add(
				mul(sym("a"), dec(9)),
				sym("b"),
			),
		},
		{
			name: "25 binary expression (4 operands)",
			str:  "a * b + c * d", //  = (a * b) + (c * d)
			want: add(
				mul(sym("a"), sym("b")),
				mul(sym("c"), sym("d")),
			),
		},
		{
			name: "26 binary expression (4 operands)",
			str:  "a < b + c * d", // = a < (b + (c * d))
			want: ls(
				sym("a"),
				add(
					sym("b"),
					mul(sym("c"), sym("d")),
				),
			),
		},
		// TODO: add more test cases
		// a + b + c * d = ((a + b) + (c * d))
		// a + b * c * d = (a + ((b * c) * d))
		// a + b * c + d = ((a + (b * c)) + d)
		// a + b + c + d = (((a + b) + c) + d)
		// a + b + c = ((a + b) + c)

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
