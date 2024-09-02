package source

import (
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Kind uint8

const (
	NoExt Kind = iota
	Unknown

	KU
	CPP
	ASM
)

var kindText = [...]string{
	NoExt:   "<nil>",
	Unknown: "<unknown>",

	KU:  "ku",
	CPP: "c",
	ASM: "asm",
}

func (k Kind) String() string {
	return kindText[k]
}

func ParseKindFromExtension(ext string) Kind {
	switch ext {
	case "":
		return NoExt
	case ".ku":
		return KU
	case ".c":
		return CPP
	case ".asm":
		return ASM
	default:
		return Unknown
	}
}

type File struct {
	// Raw binary data from file
	//
	// May be nil if file reading was deferred or this struct
	// was loaded from build cache
	Bytes []byte

	// Time of last modification, as reported by OS info call
	ModTime time.Time

	// Path to this file
	Path string

	// Base file name (without directories). Includes extension
	Name string

	// File extension, including dot character. Examples:
	//
	//	- ".gm"
	//	- ".cc"
	//	- ".cpp"
	//	- ".asm"
	//
	// This field will be empty if file does not have extension in its name
	Ext string

	// In bytes, as reported by OS info call
	//
	// Check this field first in case bytes were not yet read and stored
	// inside this struct
	Size uint64

	// For quick equality comparison
	Hash uint64

	// Ordering information given to a group of source files for consistent
	// processing results between different compiler runs. Such group is formed
	// when processing files inside a unit directory. Ordering starts from 0.
	Num uint32

	Kind Kind
}

func Hash(b []byte) uint64 {
	h := fnv.New64a()

	// implementation always writes all given bytes
	// and never returns error
	h.Write(b)

	return h.Sum64()
}

// Info returns basic information about a file, without reading it
//
// Bytes and Hash fields in the resulting File struct will be nil/zero
// regardless of actual file contents
func Info(path string) (*File, error) {
	if path == "." || path == "" {
		panic("empty path")
	}

	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path \"%s\" leads to a directory", path)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("path \"%s\" does not lead to a regular file", path)
	}

	name := filepath.Base(path)
	ext := filepath.Ext(name)
	return &File{
		Name:    name,
		Ext:     ext,
		Kind:    ParseKindFromExtension(ext),
		Path:    path,
		ModTime: info.ModTime(),
		Size:    uint64(info.Size()),
	}, nil
}

// Load reads file contents and calculates its hash
//
// If file size has changed this method updates it, making
// len(f.Bytes) == f.Size true after the call
func (f *File) Load() error {
	r, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer r.Close()

	h := fnv.New64a()
	b := make([]byte, 0, f.Size)
	for {
		if len(b) >= cap(b) {
			d := append(b[:cap(b)], 0)
			b = d[:len(b)]
		}
		buf := b[len(b):cap(b)]
		n, err := r.Read(buf)
		b = b[:len(b)+n]
		h.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				f.Bytes = b
				f.Size = uint64(len(b))
				f.Hash = h.Sum64()
				return nil
			}
			return err
		}
	}
}

// Load reads full information about a file including its hash and contents
func Load(path string) (*File, error) {
	if path == "." || path == "" {
		panic("empty path")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := uint64(info.Size())
	size++ // one byte for final read at EOF

	h := fnv.New64a()
	b := make([]byte, 0, size)
	for {
		if len(b) >= cap(b) {
			d := append(b[:cap(b)], 0)
			b = d[:len(b)]
		}
		buf := b[len(b):cap(b)]
		n, err := f.Read(buf)
		b = b[:len(b)+n]
		h.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				name := filepath.Base(path)
				ext := filepath.Ext(name)
				return &File{
					Name:    name,
					Ext:     ext,
					Kind:    ParseKindFromExtension(ext),
					Path:    path,
					Bytes:   b,
					ModTime: info.ModTime(),
					Size:    uint64(len(b)),
					Hash:    h.Sum64(),
				}, nil
			}
			return nil, err
		}
	}
}

// SortAndOrder sorts given slice of files and assigns File.Num order number
// to each element according to resulting order of elements after sorting is done.
func SortAndOrder(files []*File) {
	if len(files) == 0 {
		panic("invalid argument: <nil>")
	}

	if len(files) == 1 {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for i := 0; i < len(files); i += 1 {
		files[i].Num = uint32(i)
	}
}

const MaxFileSize = 1 << 26

func LoadUnitFiles(dir string, includeTestFiles bool) ([]*File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("directory \"%s\" is empty", dir)
	}

	var files []*File
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".ku" {
			continue
		}
		if !includeTestFiles && strings.HasSuffix(name, ".test.ku") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if !info.Mode().IsRegular() {
			continue
		}

		size := info.Size()
		if size > MaxFileSize {
			return nil, fmt.Errorf("file \"%s\" is larger than max allowed size (64mb)", name)
		}

		path := filepath.Join(dir, name)
		file, err := Load(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("directory \"%s\" does not contain gizmo source files", dir)
	}
	SortAndOrder(files)

	return files, nil
}
