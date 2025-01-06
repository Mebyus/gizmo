package kbs

import (
	"fmt"
	"testing"
)

func str(s string) String {
	return String{Value: s}
}

func checkEqualEvalResult(got, want Exp) bool {
	if want == nil {
		return got == nil
	}

	switch w := want.(type) {
	case True:
		_, ok := got.(True)
		return ok
	case False:
		_, ok := got.(False)
		return ok
	case String:
		g, ok := got.(String)
		return ok && g.Value == w.Value
	default:
		panic(fmt.Sprintf("unexpected want %v type %T", w, w))
	}
}

func TestEvaluator_evalExp(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		exp     string
		want    Exp
		wantErr bool
	}{
		{
			name: "1 string",
			exp:  `"hello"`,
			want: str("hello"),
		},
		{
			name: "2 defined env",
			exp:  "#:HELLO",
			env:  map[string]string{"HELLO": "test"},
			want: str("test"),
		},
		{
			name:    "3 undefined env",
			exp:     "#:HELLO",
			wantErr: true,
		},
		{
			name: "4 env equal string",
			exp:  `#:HELLO == "test"`,
			env:  map[string]string{"HELLO": "test"},
			want: True{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := ParseExp(tt.exp)
			if err != nil {
				t.Errorf("ParseExp() error = %v", err)
				return
			}
			e := &Evaluator{env: tt.env}
			got, err := e.evalExp(exp)
			if (err != nil) != tt.wantErr {
				t.Errorf("evalExp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !checkEqualEvalResult(got, tt.want) {
				t.Errorf("evalExp() = %s, want %s", got, tt.want)
			}
		})
	}
}
