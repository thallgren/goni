package internal

import (
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/issue/issue"
	"math"
)

type scannerSupport struct {
	value       int
	enc         goni.Encoding // fast access to encoding
	bytes       []byte        // pattern
	p           int           // current scanner position
	stop        int           // pattern end (mutable)
	lastFetched int           // last fetched value for unfetch support
	c           int           // current code point

	begin int // pattern begin position for reset() support
	end   int // pattern end position for reset() support
	_p    int // used by mark()/restore() to mark positions
}

func (ss *scannerSupport) init(enc goni.Encoding, bytes []byte, p, end int) {
	ss.enc = enc
	ss.bytes = bytes
	ss.begin = p
	ss.end = end
}

const intSignBit = 1 << 31

func digitVal(code int) int {
	return code - '0'
}

func xdigitVal(enc goni.Encoding, code int) int {
		if (enc.IsDigit(code)) {
			return digitVal(code);
		} else {
			if enc.IsUpper(code) {
				return code - 'A' + 10;
			}
			return code - 'a' + 10;
		}
	}


func (ss *scannerSupport) scanUnsignedNumber() int {
	enc := ss.enc
	last := ss.c
	num := 0
	for ss.left() {
		ss.fetch()
		if enc.IsDigit(ss.c) {
			onum := num
			num = num*10 + digitVal(ss.c)
			if ((onum ^ num) & intSignBit) != 0 {
				return -1
			}
		} else {
			ss.unfetch()
			break
		}
	}
	ss.c = last
	return num
}

func (ss *scannerSupport) scanUnsignedHexadecimalNumber(minLength, maxLength int) int {
	enc := ss.enc
	last := ss.c
	num := 0
	restLen := maxLength - minLength;
	for ss.left() && maxLength > 0 {
		maxLength--
		ss.fetch()
		if enc.IsXDigit(ss.c) {
			val := xdigitVal(enc, ss.c)
			if ((math.MaxInt64 - val) / 16 < num) {
				return -1
			}
			num = (num << 5) + val
		} else {
			ss.unfetch()
			maxLength++
			break
		}
	}
	if maxLength > restLen {
		return -2
	}
	ss.c = last
	return num
}

func (ss *scannerSupport) scanUnsignedOctalNumber(maxLength int) int {
	enc := ss.enc
	last := ss.c
	num := 0
	for ss.left() && maxLength > 0 {
		maxLength--
		ss.fetch()
		if enc.IsDigit(ss.c) && ss.c < '8' {
			onum := num
			val := digitVal(ss.c)
			num = (num << 3) + val
			if ((onum ^ num) & intSignBit) != 0 {
				return -1
			}
		} else {
			ss.unfetch()
			break
		}
	}
	ss.c = last
	return num
}

func (ss *scannerSupport) reset() {
	ss.p = ss.begin
	ss.stop = ss.end
}

func (ss *scannerSupport) mark() {
	ss._p = ss.p
}

func (ss *scannerSupport) restore() {
	ss.p = ss._p
}

func (ss *scannerSupport) inc() {
	ss.lastFetched = ss.p
	ss.p += ss.enc.Length(ss.bytes, ss.p, ss.stop)
}

func (ss *scannerSupport) fetch() {
	var ln int
	ss.c, ln = ss.enc.MbcToCode(ss.bytes, ss.p, ss.stop)
	ss.lastFetched = ss.p
	ss.p += ln
}

func (ss *scannerSupport) fetchTo() int {
	to, ln := ss.enc.MbcToCode(ss.bytes, ss.p, ss.stop)
	ss.lastFetched = ss.p
	ss.p += ln
	return to
}

func (ss *scannerSupport) unfetch() {
	ss.p = ss.lastFetched
}

func (ss *scannerSupport) peek() int {
	if ss.p < ss.stop {
		c, _ := ss.enc.MbcToCode(ss.bytes, ss.p, ss.stop)
		return c
	}
	return 0
}

func (ss *scannerSupport) peekIs(c int) bool {
	return ss.peek() == c
}

func (ss *scannerSupport) left() bool {
	return ss.p < ss.stop
}

func (ss* scannerSupport) newValueException(code issue.Code, p, end int) issue.Reported {
	return err.WithArgs(code, issue.H{`n`: string(ss.bytes[p:end]) })
}

func newSyntaxException(code issue.Code) issue.Reported {
	return err.NoArgs(code)
}
