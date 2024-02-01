package source

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "001_file",
			data: nil,
		},
		{
			name: "002_file",
			data: []byte{10},
		},
		{
			name: "003_file",
			data: bytes.Repeat([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 1<<11),
		},
		{
			name: "004_file",
			data: []byte{10, 0, 255},
		},
	}

	dir := t.TempDir()

	// test package level Load function
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "" {
				panic("empty test/file name")
			}
			path := filepath.Join(dir, tt.name)
			err := os.WriteFile(path, tt.data, 0o664)
			if err != nil {
				t.Errorf("Write temp file error = %v", err)
				return
			}

			got, err := Load(path)
			if err != nil {
				t.Errorf("Load() error = %v", err)
				return
			}
			if got.Size != uint64(len(tt.data)) {
				t.Errorf("Load() size = %d, want %d", got.Size, uint64(len(tt.data)))
			}
			if got.Name != path {
				t.Errorf("Load() name = \"%s\", want \"%s\"", got.Name, path)
			}
			wantHash := Hash(tt.data)
			if got.Hash != wantHash {
				t.Errorf("Load() hash = %016x, want %016x", got.Hash, wantHash)
			}
			if !bytes.Equal(got.Bytes, tt.data) {
				t.Errorf("Load() bytes are not equal")
			}
		})
	}

	// test package level Check function + Load method
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name)
			got, err := Check(path)
			if err != nil {
				t.Errorf("Check() error = %v", err)
				return
			}
			err = got.Load()
			if err != nil {
				t.Errorf("File.Load() error = %v", err)
				return
			}
			if got.Size != uint64(len(tt.data)) {
				t.Errorf("File.Load() size = %d, want %d", got.Size, uint64(len(tt.data)))
			}
			if got.Name != path {
				t.Errorf("File.Load() name = \"%s\", want \"%s\"", got.Name, path)
			}
			wantHash := Hash(tt.data)
			if got.Hash != wantHash {
				t.Errorf("File.Load() hash = %016x, want %016x", got.Hash, wantHash)
			}
			if !bytes.Equal(got.Bytes, tt.data) {
				t.Errorf("File.Load() bytes are not equal")
			}
		})
	}
}
