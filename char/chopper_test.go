package char

import "testing"

func TestChopper_Advance(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "1 empty input",
			input: "",
		},
		{
			name:  "2 one byte input",
			input: "h",
		},
		{
			name:  "3 two byte input",
			input: "he",
		},
		{
			name:  "4 three byte input",
			input: "hel",
		},
		{
			name:  "5 two line input",
			input: "hello,\nworld!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Chopper
			c.Init([]byte(tt.input))

			for i := range len(tt.input) {
				if c.C != tt.input[i] {
					t.Errorf("wrong byte at offset %d: got %c (0x%02x), want %c (0x%02x)",
						i, c.C, c.C, tt.input[i], tt.input[i])
					return
				}
				c.Advance()
			}

			if !c.EOF {
				t.Errorf("expected to be at EOF")
				return
			}
		})
	}
}
