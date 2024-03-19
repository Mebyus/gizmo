package builder

import (
	"io"
	"os"
)

// PartsBuffer stores multiple parts of data (byte slices) in sequential manner
// and implements Read and WriteTo methods which use stored parts to stream
// data. In other words this buffer provides abstraction analogous to bytes.Buffer
// for working with multiple byte slices as if they were one continuous source of data
//
// Zero value of PartsBuffer is ready for usage. Although reading from an empty buffer
// will always result in EOF. This implementation is not safe for concurrent usage. Users
// must use a synchronization primitive with this buffer if they want to use it from
// mulptiple goroutines
type PartsBuffer struct {
	// stored underlying parts
	parts [][]byte

	// reading index inside current part
	i int

	// index of current part (for reading)
	j int

	// total length of all stored parts
	len int
}

// NewPartsBuffer make a new buffer with preallocated space to store at least
// specified number of parts
func NewPartsBuffer(size int) *PartsBuffer {
	return &PartsBuffer{parts: make([][]byte, 0, size)}
}

// Add appends data parts to the end of buffer. Takes ownership of each given slice.
// Parts with no data (len = 0) are discarded
func (p *PartsBuffer) Add(parts ...[]byte) {
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		p.parts = append(p.parts, part)
		p.len += len(part)
	}
}

// AddStr does the same operation as Add, but works with strings instead of byte slices
func (p *PartsBuffer) AddStr(parts ...string) {
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		p.parts = append(p.parts, []byte(part))
		p.len += len(part)
	}
}

// Len returns total length of all stored parts
func (p *PartsBuffer) Len() int {
	return p.len
}

// Read implements io.Reader. Never returns error other than io.EOF. Read and Add
// calls can be intermixed if they are performed from the same goroutine
//
// After reaching EOF Add method can still be used to add new parts. In that case
// Read will pick up from where it stopped previously and read given data until new
// EOF is reached
func (p *PartsBuffer) Read(b []byte) (n int, err error) {
	if p.j >= len(p.parts) {
		return 0, io.EOF
	}

	for {
		nw := copy(b[n:], p.parts[p.j][p.i:])
		p.i += nw
		n += nw

		if p.i >= len(p.parts[p.j]) {
			p.j += 1
			p.i = 0
		}

		if nw == 0 || p.j >= len(p.parts) {
			return
		}
	}
}

// WriteTo implements io.WriterTo
func (p *PartsBuffer) WriteTo(w io.Writer) (n int64, err error) {
	for p.j < len(p.parts) {
		nw, err := w.Write(p.parts[p.j][p.i:])
		p.i += nw
		n += int64(nw)

		if p.i >= len(p.parts[p.j]) {
			p.j += 1
			p.i = 0
		}

		if err != nil {
			return n, err
		}
	}
	return
}

func SavePartsBuffer(filename string, buf *PartsBuffer) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = buf.WriteTo(file)
	return err
}
