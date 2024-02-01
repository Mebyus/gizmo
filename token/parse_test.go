package token

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    Token
		wantErr bool
	}{
		{
			name:    "1 empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "2 pos spec only",
			line:    "3:7",
			wantErr: true,
		},
		{
			name: "3 EOF",
			line: "3:7 EOF",
			want: Token{
				Kind: EOF,
				Pos:  pos(3, 7),
			},
		},
		{
			name: "4 identifier",
			line: "1:1 IDENT abc",
			want: Token{
				Kind: Identifier,
				Pos:  pos(1, 1),
				Lit:  "abc",
			},
		},
		{
			name: "5 decimal integer",
			line: "101:22   INT.DEC     4510",
			want: Token{
				Kind: DecimalInteger,
				Pos:  pos(101, 22),
				Val:  4510,
			},
		},
		{
			name: "6 plus",
			line: "30:71 \t+",
			want: Token{
				Kind: Plus,
				Pos:  pos(30, 71),
			},
		},
		{
			name: "7 string",
			line: "716:9\t\tSTR   \"Hello, world!\"",
			want: Token{
				Kind: String,
				Pos:  pos(716, 9),
				Lit:  "Hello, world!",
			},
		},
		{
			name: "8 float",
			line: "101:22   FLT.DEC     4.51",
			want: Token{
				Kind: DecimalFloat,
				Pos:  pos(101, 22),
				Lit:  "4.51",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				got, err := Parse(tt.line)
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					var msg string
					if got.Lit != tt.want.Lit {
						msg = "literals are not equal"
					} else if got.Val != tt.want.Val {
						msg = "values are not equal"
					}
					if msg != "" {
						msg = "\n" + msg
					}
					t.Errorf("%s\nParse() = %v\nwant      %v", msg, got, tt.want)
				}
			},
		)
	}
}
