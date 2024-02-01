package source

import (
	"io"
	"strconv"
	"strings"
)

// Span is a piece of source code with known starting position and length (in bytes)
type Span interface {
	Pin

	// continuous length of code piece
	Len() uint32
}

func Render(w io.Writer, s Span, window LineWindow) error {
	lines, index := ExtractLines(s, window)
	lastLine := lines[len(lines)-1]
	maxLineNum := lastLine.Num + 1 // maximum line number in list, one-based

	// how much characters needed to write line number column
	lineCaptionLen := len(strconv.FormatUint(uint64(maxLineNum), 10))

	var lineMaxLen int
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if len(line.Bytes) > lineMaxLen {
			lineMaxLen = len(line.Bytes)
		}
	}

	for i := uint32(0); i <= index; i++ {
		line := lines[i]

		err := formatLine(w, line, lineCaptionLen)
		if err != nil {
			return err
		}
	}

	pos := s.Pin()
	width := s.Len()
	if width == 0 {
		width = 1
	}
	_, err := io.WriteString(w, formatHighlightLine(pos.Col, width, lineCaptionLen))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, formatSeparatorLine(lineMaxLen, lineCaptionLen))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	for i := index + 1; i < uint32(len(lines)); i++ {
		line := lines[i]

		err := formatLine(w, line, lineCaptionLen)
		if err != nil {
			return err
		}
	}

	return nil
}

func formatLine(w io.Writer, line Line, captionWidth int) error {
	_, err := io.WriteString(w, formatLineNumber(line.Num+1, captionWidth)+" | ")
	if err != nil {
		return err
	}
	_, err = w.Write(line.Bytes)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}

func formatHighlightLine(start, width uint32, captionWidth int) string {
	s := strings.Repeat(" ", captionWidth) + " | "
	s += strings.Repeat(" ", int(start))
	s += strings.Repeat("^", int(width))
	return s
}

func formatSeparatorLine(width, captionWidth int) string {
	s := strings.Repeat(" ", captionWidth) + " | "
	s += strings.Repeat("-", width)
	return s
}

func formatLineNumber(n uint32, width int) string {
	num := strconv.FormatUint(uint64(n), 10)
	s := strings.Repeat(" ", width-len(num)) + num
	return s
}

type LineWindow struct {
	Before uint32
	After  uint32
}

// ExtractLines returns Lines (with burrowed bytes) around given
// Span in File
//
// second returned value (index) is index of line which contains Span
// starting Pos
func ExtractLines(s Span, w LineWindow) (lines []Line, index uint32) {
	l := s.Len()
	pos := s.Pin()
	start := pos.Ofs   // start index of Span
	end := pos.Ofs + l // end index of Span
	b := pos.File.Bytes
	// TODO: fix of handle EOF rendering panic
	if end > uint32(len(b)) {
		panic("span is out of bounds")
	}

	var i uint32 // byte index in source bytes
	// search for line start of the line which contains Span start
	for i = start; i != 0 && b[i] != '\n'; i-- {
	}

	var lineStart uint32
	var lineEnd uint32
	if i == 0 {
		lineStart = i
	} else {
		lineStart = i + 1
		lineEnd = i
		i--
	}
	spanLineStart := lineStart

	var j uint32            // difference between starting line and current line
	lineNum := pos.Line - 1 // line number at current index

	// gather lines which come before line with Span start
	for ; i != 0; i-- {
		if b[i] == '\n' {
			lineStart = i + 1
			lines = append(lines, Line{
				Num:   lineNum,
				Bytes: b[lineStart:lineEnd],
			})

			lineEnd = i
			j++
			lineNum--
			if j >= w.Before {
				break
			}
		}
	}

	if i == 0 {
		lines = append(lines, Line{
			Num:   0,
			Bytes: b[:lineEnd],
		})
	}

	// invert lines order in place
	k := 0
	p := len(lines) - 1
	for k < p {
		lines[k], lines[p] = lines[p], lines[k]
		k++
		p--
	}

	// search for line end of the line which contains Span start
	for i = start; i < uint32(len(b)) && b[i] != '\n'; i++ {
	}
	lineEnd = i

	// add line with Span start
	index = uint32(len(lines))
	lines = append(lines, Line{
		Num:   pos.Line,
		Bytes: b[spanLineStart:lineEnd],
	})

	i++
	j = 0
	lineStart = i
	lineNum = pos.Line + 1
	// gather lines which come after line with Span start
	for ; i < uint32(len(b)); i++ {
		if b[i] == '\n' {
			lineEnd = i
			lines = append(lines, Line{
				Num:   lineNum,
				Bytes: b[lineStart:lineEnd],
			})

			lineStart = i + 1
			j++
			lineNum++
			if j >= w.After {
				break
			}
		}
	}

	return lines, index
}

type Line struct {
	// raw bytes which consitute the line without eol characters
	Bytes []byte

	// line number, zero-based
	Num uint32
}
