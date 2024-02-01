package source

import (
	"hash/fnv"
	"io"
	"os"
	"time"
)

type File struct {
	// Raw binary data from file
	//
	// May be nil if file reading was deferred or this struct
	// was loaded from build cache
	Bytes []byte

	// Time of last modification, as reported by OS info call
	ModTime time.Time

	// name/path to a file
	Name string

	// In bytes, as reported by OS info call
	//
	// Check this field first in case bytes were not yet read and stored
	// inside this struct
	Size uint64

	// For quick equality comparison
	Hash uint64
}

func Hash(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}

// Check returns basic information about a file, without reading it
//
// Bytes and Hash fields in the resulting File struct will be nil/zero
// regardless of actual file contents
func Check(name string) (*File, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return &File{
		Name:    name,
		ModTime: info.ModTime(),
		Size:    uint64(info.Size()),
	}, nil
}

// Load reads file contents and calculates its hash
//
// If file size has changed this method updates it, making
// len(f.Bytes) == f.Size true after the call
func (f *File) Load() error {
	r, err := os.Open(f.Name)
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
func Load(name string) (*File, error) {
	f, err := os.Open(name)
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
				return &File{
					Name:    name,
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
