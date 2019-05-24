package option

import (
	"bytes"
	"github.com/lyraproj/goni/util"
	"io"
)

type Type int

const (
	None             = Type(0)
	IgnoreCase       = Type(1 << 0)
	Extend           = Type(1 << 1)
	MultiLine        = Type(1 << 2)
	SingleLine       = Type(1 << 3)
	FindLongest      = Type(1 << 4)
	FindNotEmpty     = Type(1 << 5)
	NegateSingleLine = Type(1 << 6)
	DontCaptureGroup = Type(1 << 7)
	CaptureGroup     = Type(1 << 8)

	// search time
	NotBOL      = Type(1 << 9)
	NotEOL      = Type(1 << 10)
	PosixRegion = Type(1 << 11)

	// ctype range
	AsciiRange           = Type(1 << 12)
	PosixBracketAllRange = Type(1 << 13)
	WordBoundAllRange    = Type(1 << 14)

	// newline
	NewlineCRLF = Type(1 << 15)
	NotBOS      = Type(1 << 16)
	NotEOS      = Type(1 << 17)
	CR7Bit      = Type(1 << 18)
	MaxBit      = Type(1 << 19)

	Default = None
)

func (v Type) String() string {
	w := &bytes.Buffer{}
	v.AppendString(w)
	return w.String()
}

func (v Type) AppendString(w io.Writer) {
	if v.IsIgnoreCase() {
		util.WriteString(w, `IgnoreCase `)
	}
	if v.IsExtend() {
		util.WriteString(w, `Extend `)
	}
	if v.IsMultiLine() {
		util.WriteString(w, `MultiLine `)
	}
	if v.IsSingleLine() {
		util.WriteString(w, `SingleLine `)
	}
	if v.IsFindLongest() {
		util.WriteString(w, `FindLongest `)
	}
	if v.IsFindNotEmpty() {
		util.WriteString(w, `FindNotEmpty `)
	}
	if v.IsNegateSingleLine() {
		util.WriteString(w, `NegateSingleLine `)
	}
	if v.IsDontCaptureGroup() {
		util.WriteString(w, `DontCaptureGroup `)
	}
	if v.IsCaptureGroup() {
		util.WriteString(w, `CaptureGroup `)
	}
	if v.IsNotBOL() {
		util.WriteString(w, `NotBOL `)
	}
	if v.IsNotEOL() {
		util.WriteString(w, `NotEOL `)
	}
	if v.IsPosixRegion() {
		util.WriteString(w, `PosixRegion `)
	}
	if v.IsAsciiRange() {
		util.WriteString(w, `AsciiRange `)
	}
	if v.IsPosixBracketAllRange() {
		util.WriteString(w, `PosixBracketAllRange `)
	}
	if v.IsWordBoundAllRange() {
		util.WriteString(w, `WordBoundAllRange `)
	}
	if v.IsNewlineCRLF() {
		util.WriteString(w, `NewlineCRLF `)
	}
	if v.IsNotBOS() {
		util.WriteString(w, `NotBOS `)
	}
	if v.IsNotEOS() {
		util.WriteString(w, `NotEOS `)
	}
	if v.IsCR7Bit() {
		util.WriteString(w, `CR7Bit `)
	}
	if v.IsMaxBit() {
		util.WriteString(w, `MaxBit `)
	}
}

func (t Type) IsType(ot Type) bool {
	return (t & ot) != 0
}

func (v Type) IsIgnoreCase() bool {
	return (v & IgnoreCase) != 0
}

func (v Type) IsExtend() bool {
	return (v & Extend) != 0
}

func (v Type) IsMultiLine() bool {
	return (v & MultiLine) != 0
}

func (v Type) IsSingleLine() bool {
	return (v & SingleLine) != 0
}

func (v Type) IsFindLongest() bool {
	return (v & FindLongest) != 0
}

func (v Type) IsFindNotEmpty() bool {
	return (v & FindNotEmpty) != 0
}

func (v Type) IsNegateSingleLine() bool {
	return (v & NegateSingleLine) != 0
}

func (v Type) IsDontCaptureGroup() bool {
	return (v & DontCaptureGroup) != 0
}

func (v Type) IsCaptureGroup() bool {
	return (v & CaptureGroup) != 0
}

func (v Type) IsNotBOL() bool {
	return (v & NotBOL) != 0
}

func (v Type) IsNotEOL() bool {
	return (v & NotEOL) != 0
}

func (v Type) IsPosixRegion() bool {
	return (v & PosixRegion) != 0
}

func (v Type) IsAsciiRange() bool {
	return (v & AsciiRange) != 0
}

func (v Type) IsPosixBracketAllRange() bool {
	return (v & PosixBracketAllRange) != 0
}

func (v Type) IsWordBoundAllRange() bool {
	return (v & WordBoundAllRange) != 0
}

func (v Type) IsNewlineCRLF() bool {
	return (v & NewlineCRLF) != 0
}

func (v Type) IsNotBOS() bool {
	return (v & NotBOS) != 0
}

func (v Type) IsNotEOS() bool {
	return (v & NotEOS) != 0
}

func (v Type) IsCR7Bit() bool {
	return (v & CR7Bit) != 0
}

func (v Type) IsMaxBit() bool {
	return (v & MaxBit) != 0
}
