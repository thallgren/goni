package util

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Indenter struct {
	b *bytes.Buffer
	i int
}

func NewIndenter() *Indenter {
	return &Indenter{b: &bytes.Buffer{}, i: 0}
}

// Reset resets the internal buffer. It does not reset the indent
func (i *Indenter) Reset() {
	i.b.Reset()
}

func (i *Indenter) String() string {
	n := bytes.NewBuffer(make([]byte, 0, i.b.Len()))
	wb := &bytes.Buffer{}
	for {
		r, _, err := i.b.ReadRune()
		if err == io.EOF {
			break
		}
		if r == ' ' || r == '\t' {
			// Defer whitespace output
			wb.WriteByte(byte(r))
			continue
		}
		if r == '\n' {
			// Truncate trailing space
			wb.Reset()
		} else {
			if wb.Len() > 0 {
				n.Write(wb.Bytes())
				wb.Reset()
			}
		}
		n.WriteRune(r)
	}
	return n.String()
}

func (i *Indenter) WriteString(s string) (n int, err error) {
	return i.b.WriteString(s)
}

func (i *Indenter) Write(p []byte) (n int, err error) {
	return i.b.Write(p)
}

func (i *Indenter) AppendRune(r rune) {
	i.b.WriteRune(r)
}

func (i *Indenter) Append(s string) {
	WriteString(i.b, s)
}

// AppendIndented is like Indent but replaces all occurrences of newline with an indented newline
func (i *Indenter) AppendIndented(s string) {
	for ni := strings.IndexByte(s, '\n'); ni >= 0; ni = strings.IndexByte(s, '\n') {
		if ni > 0 {
			WriteString(i.b, s[:ni])
		}
		i.NewLine()
		ni++
		if ni >= len(s) {
			return
		}
		s = s[ni:]
	}
	if len(s) > 0 {
		WriteString(i.b, s)
	}
}

func (i *Indenter) AppendBool(b bool) {
	var s string
	if b {
		s = `true`
	} else {
		s = `false`
	}
	WriteString(i.b, s)
}

func (i *Indenter) AppendInt(b int) {
	WriteString(i.b, strconv.Itoa(b))
}

func (i *Indenter) Indent() *Indenter {
	return &Indenter{b: i.b, i: i.i + 1}
}

func (i *Indenter) Printf(s string, args ...interface{}) {
	Fprintf(i.b, s, args...)
}

// NewLine writes a newline followed by the current indent after trimming trailing whitespaces
func (i *Indenter) NewLine() {
	i.b.WriteByte('\n')
	for n := 0; n < i.i; n++ {
		WriteString(i.b, `  `)
	}
}

func WriteString(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}

func Fprintf(w io.Writer, s string, args ...interface{}) {
	_, err := fmt.Fprintf(w, s, args...)
	if err != nil {
		panic(err)
	}
}
