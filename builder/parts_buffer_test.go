package builder

import (
	"io"
	"strings"
	"testing"
)

func clipString(s string, i int) string {
	if len(s) <= i {
		return s
	}
	return s[:i]
}

func TestPartsBuffer(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
	}{
		{
			name:  "nil slice",
			parts: nil,
		},
		{
			name:  "empty slice",
			parts: []string{},
		},
		{
			name:  "single empty string",
			parts: []string{""},
		},
		{
			name:  "multiple empty strings",
			parts: []string{"", "", "", ""},
		},
		{
			name:  "single string",
			parts: []string{"abc_ 123"},
		},
		{
			name:  "two strings",
			parts: []string{"abc_ 123", "  kk()"},
		},
		{
			name:  "leading empty string",
			parts: []string{"", "abc_ 123", "  kk()"},
		},
		{
			name:  "unicode characters",
			parts: []string{"abc_ 123", "[", "ыййц", "", "]"},
		},
		{
			name:  "multiple empty strings in the middle",
			parts: []string{"", "42 gf", "", "", "", "__+!", ""},
		},
	}

	for _, tt := range tests {
		t.Run("(Read) "+tt.name, func(t *testing.T) {
			buf := PartsBuffer{}
			buf.AddStr(tt.parts...)

			const maxLen = 1 << 16
			var readBuf [maxLen]byte

			want := clipString(strings.Join(tt.parts, ""), maxLen)

			n, err := buf.Read(readBuf[:])
			if err != nil && err != io.EOF {
				t.Errorf("Read() error = %v, want <nil>", err)
				return
			}
			if n != len(want) {
				t.Errorf("Read() = %d, want %d", n, len(want))
				return
			}

			got := string(readBuf[:n])
			if got != want {
				t.Errorf("got = \"%s\", want \"%s\"", got, want)
				return
			}
		})

		t.Run("(WriteTo) "+tt.name, func(t *testing.T) {
			buf := PartsBuffer{}
			buf.AddStr(tt.parts...)

			want := strings.Join(tt.parts, "")
			gotLen := buf.Len()
			if gotLen != len(want) {
				t.Errorf("Len() = %d, want %d", gotLen, len(want))
				return
			}

			gotBuf := strings.Builder{}
			n, err := buf.WriteTo(&gotBuf)
			if err != nil {
				t.Errorf("WriteTo() error = %v, want <nil>", err)
				return
			}
			if n != int64(len(want)) {
				t.Errorf("WriteTo() = %d, want %d", n, len(want))
				return
			}

			got := gotBuf.String()
			if got != want {
				t.Errorf("got = \"%s\", want \"%s\"", got, want)
				return
			}
		})

		t.Run("(Copy) "+tt.name, func(t *testing.T) {
			buf := PartsBuffer{}
			buf.AddStr(tt.parts...)

			want := strings.Join(tt.parts, "")

			gotBuf := strings.Builder{}
			n, err := io.Copy(&gotBuf, &buf)
			if err != nil {
				t.Errorf("Copy() error = %v, want <nil>", err)
				return
			}
			if n != int64(len(want)) {
				t.Errorf("Copy() = %d, want %d", n, len(want))
				return
			}

			got := gotBuf.String()
			if got != want {
				t.Errorf("got = \"%s\", want \"%s\"", got, want)
				return
			}
		})
	}
}
