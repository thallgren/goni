package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/coderange"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

const flagNcCClassNot = 1 << 0

type CClassNode struct {
	abstractNode
	flags int
	bs    *goni.BitSet
	mbuf  *coderange.Buffer
}

func (c *CClassNode) String() string {
	goni.String(c)
}

func (c *CClassNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`flags: `)
	c.appendFlags(w)
	w.NewLine()
	w.Append(`bs:`)
	c.bs.AppendTo(w.Indent())
	if c.mbuf != nil {
		w.NewLine()
		w.Append(`mbuf: `)
		c.mbuf.AppendTo(w.Indent())
	}
}

func (c *CClassNode) Name() string {
	return `Character Class`
}

func NewCClass() *CClassNode {
	return &CClassNode{abstractNode: abstractNode{nodeType: node.CClass}, bs: goni.NewBitSet()}
}

func (c *CClassNode) Clear() {
	c.bs.ClearAll()
	c.flags = 0
	c.mbuf = nil
}

func (c *CClassNode) isEmpty() bool {
	return c.mbuf == nil && c.bs.IsEmpty()
}

func (c *CClassNode) isNot() bool {
	return (c.flags & flagNcCClassNot) != 0
}

func (c *CClassNode) clearNot() {
	c.flags &= ^flagNcCClassNot
}

func (c *CClassNode) appendFlags(w *util.Indenter) {
	if c.isNot() {
		w.Append(`NOT`)
	}
}

func (c *CClassNode) addCodeRangeToBuf(env goni.ScanEnvironment, from, to int, checkDup bool) {
	c.mbuf = coderange.AddToBuffer(c.mbuf, env, from, to, checkDup)
}

func (c *CClassNode) addCodeRange(env goni.ScanEnvironment, from, to int, checkDup bool) {
	c.mbuf = coderange.Add(c.mbuf, env, from, to, checkDup)
}

func (c *CClassNode) addAllMultiByteRange(env goni.ScanEnvironment) {
	c.mbuf = coderange.AddAllMultiByte(env, c.mbuf)
}

func (c *CClassNode) clearNotFlag(env goni.ScanEnvironment) {
	if c.isNot() {
		c.bs.InvertAll()
		if !env.Encoding().IsSingleByte() {
			c.mbuf = coderange.NotBuffer(env, c.mbuf)
		}
		c.clearNot()
	}
}

func (cn *CClassNode) isOneChar() int {
	if cn.isNot() {
		return -1
	}
	c := -1
	if cn.mbuf != nil {
		rng := cn.mbuf.Range()
		c := (uint)(rng[1])
		if rng[0] == 1 && c == (uint)(rng[2]) {
			if c < goni.SingleByteSize && cn.bs.At(c) {
				c = -1
			}
		} else {
			return -1
		}
	}

	for i := 0; i < goni.BitSetSize; i++ {
		b1 := cn.bs.RoomAt(i)
		if b1 != 0 {
			if (b1&(b1-1)) == 0 && c == -1 {
				c = goni.BitsInRoom*i + goni.BitCount(b1-1)
			} else {
				return -1
			}
		}
	}
	return c
}

var x = `
    // and_cclass
    public void and(CClassNode other, ScanEnvironment env) {
        boolean not1 = isNot();
        BitSet bsr1 = bs;
        Buffer buf1 = mbuf;
        boolean not2 = other.isNot();
        BitSet bsr2 = other.bs;
        Buffer buf2 = other.mbuf;

        if (not1) {
            BitSet bs1 = new BitSet();
            bsr1.invertTo(bs1);
            bsr1 = bs1;
        }

        if (not2) {
            BitSet bs2 = new BitSet();
            bsr2.invertTo(bs2);
            bsr2 = bs2;
        }

        bsr1.and(bsr2);

        if (bsr1 != bs) {
            bs.copy(bsr1);
            bsr1 = bs;
        }

        if (not1) {
            bs.invert();
        }

        Buffer pbuf = null;

        if (!env.enc.isSingleByte()) {
            if (not1 && not2) {
                pbuf = Buffer.orCodeRangeBuff(env, buf1, false, buf2, false);
            } else {
                pbuf = Buffer.andCodeRangeBuff(buf1, not1, buf2, not2, env);

                if (not1) {
                    pbuf = Buffer.notCodeRangeBuff(env, pbuf);
                }
            }
            mbuf = pbuf;
        }

    }

    // or_cclass
    public void or(CClassNode other, ScanEnvironment env) {
        boolean not1 = isNot();
        BitSet bsr1 = bs;
        Buffer buf1 = mbuf;
        boolean not2 = other.isNot();
        BitSet bsr2 = other.bs;
        Buffer buf2 = other.mbuf;

        if (not1) {
            BitSet bs1 = new BitSet();
            bsr1.invertTo(bs1);
            bsr1 = bs1;
        }

        if (not2) {
            BitSet bs2 = new BitSet();
            bsr2.invertTo(bs2);
            bsr2 = bs2;
        }

        bsr1.or(bsr2);

        if (bsr1 != bs) {
            bs.copy(bsr1);
            bsr1 = bs;
        }

        if (not1) {
            bs.invert();
        }

        if (!env.enc.isSingleByte()) {
            Buffer pbuf = null;
            if (not1 && not2) {
                pbuf = Buffer.andCodeRangeBuff(buf1, false, buf2, false, env);
            } else {
                pbuf = Buffer.orCodeRangeBuff(env, buf1, not1, buf2, not2);
                if (not1) {
                    pbuf = Buffer.notCodeRangeBuff(env, pbuf);
                }
            }
            mbuf = pbuf;
        }
    }

    // add_ctype_to_cc_by_range // Encoding out!
    public void addCTypeByRange(int ctype, boolean not, ScanEnvironment env, int sbOut, int mbr[]) {
        int n = mbr[0];
        int i;

        if (!not) {
            for (i=0; i<n; i++) {
                for (int j=CR_FROM(mbr, i); j<=CR_TO(mbr, i); j++) {
                    if (j >= sbOut) {
                        if (j > CR_FROM(mbr, i)) {
                            addCodeRangeToBuf(env, j, CR_TO(mbr, i));
                            i++;
                        }
                        // !goto sb_end!, remove duplication!
                        for (; i<n; i++) {
                            addCodeRangeToBuf(env, CR_FROM(mbr, i), CR_TO(mbr, i));
                        }
                        return;
                    }
                    bs.set(env, j);
                }
            }
            // !sb_end:!
            for (; i<n; i++) {
                addCodeRangeToBuf(env, CR_FROM(mbr, i), CR_TO(mbr, i));
            }

        } else {
            int prev = 0;

            for (i=0; i<n; i++) {
                for (int j=prev; j < CR_FROM(mbr, i); j++) {
                    if (j >= sbOut) {
                        // !goto sb_end2!, remove duplication
                        prev = sbOut;
                        for (i=0; i<n; i++) {
                            if (prev < CR_FROM(mbr, i)) addCodeRangeToBuf(env, prev, CR_FROM(mbr, i) - 1);
                            prev = CR_TO(mbr, i) + 1;
                        }
                        if (prev < 0x7fffffff/*!!!*/) addCodeRangeToBuf(env, prev, 0x7fffffff);
                        return;
                    }
                    bs.set(env, j);
                }
                prev = CR_TO(mbr, i) + 1;
            }

            for (int j=prev; j<sbOut; j++) {
                bs.set(env, j);
            }

            // !sb_end2:!
            prev = sbOut;
            for (i=0; i<n; i++) {
                if (prev < CR_FROM(mbr, i)) addCodeRangeToBuf(env, prev, CR_FROM(mbr, i) - 1);
                prev = CR_TO(mbr, i) + 1;
            }
            if (prev < 0x7fffffff/*!!!*/) addCodeRangeToBuf(env, prev, 0x7fffffff);
        }
    }

    private static int CR_FROM(int[] range, int i) {
        return range[(i * 2) + 1];
    }

    private static int CR_TO(int[] range, int i) {
        return range[(i * 2) + 2];
    }

    // add_ctype_to_cc
    public void addCType(int ctype, boolean not, boolean asciiRange, ScanEnvironment env, IntHolder sbOut) {
        Encoding enc = env.enc;
        int[]ranges = enc.ctypeCodeRange(ctype, sbOut);
        if (ranges != null) {
            if (asciiRange) {
                CClassNode ccWork = new CClassNode();
                ccWork.addCTypeByRange(ctype, not, env, sbOut.value, ranges);
                if (not) {
                    ccWork.addCodeRangeToBuf(env, 0x80, Buffer.LAST_CODE_POINT, false);
                } else {
                    CClassNode ccAscii = new CClassNode();
                    if (enc.minLength() > 1) {
                        ccAscii.addCodeRangeToBuf(env, 0x00, 0x7F);
                    } else {
                        ccAscii.bs.setRange(env, 0x00, 0x7F);
                    }
                    ccWork.and(ccAscii, env);
                }
                or(ccWork, env);
            } else {
                addCTypeByRange(ctype, not, env, sbOut.value, ranges);
            }
            return;
        }

        int maxCode = asciiRange ? 0x80 : BitSet.SINGLE_BYTE_SIZE;
        switch(ctype) {
        case CharacterType.ALPHA:
        case CharacterType.BLANK:
        case CharacterType.CNTRL:
        case CharacterType.DIGIT:
        case CharacterType.LOWER:
        case CharacterType.PUNCT:
        case CharacterType.SPACE:
        case CharacterType.UPPER:
        case CharacterType.XDIGIT:
        case CharacterType.ASCII:
        case CharacterType.ALNUM:
            if (not) {
                for (int c=0; c<BitSet.SINGLE_BYTE_SIZE; c++) {
                    if (!enc.isCodeCType(c, ctype)) bs.set(env, c);
                }
                addAllMultiByteRange(env);
            } else {
                for (int c=0; c<BitSet.SINGLE_BYTE_SIZE; c++) {
                    if (enc.isCodeCType(c, ctype)) bs.set(env, c);
                }
            }
            break;

        case CharacterType.GRAPH:
        case CharacterType.PRINT:
            if (not) {
                for (int c=0; c<BitSet.SINGLE_BYTE_SIZE; c++) {
                    if (!enc.isCodeCType(c, ctype) || c >= maxCode) bs.set(env, c);
                }
                if (asciiRange) addAllMultiByteRange(env);
            } else {
                for (int c=0; c<maxCode; c++) {
                    if (enc.isCodeCType(c, ctype)) bs.set(env, c);
                }
                if (!asciiRange) addAllMultiByteRange(env);
            }
            break;

        case CharacterType.WORD:
            if (!not) {
                for (int c=0; c<maxCode; c++) {
                    if (enc.isSbWord(c)) bs.set(env, c);
                }
                if (!asciiRange) addAllMultiByteRange(env);
            } else {
                for (int c=0; c<BitSet.SINGLE_BYTE_SIZE; c++) {
                    if (enc.codeToMbcLength(c) > 0 && /* check invalid code point */
                            !(enc.isWord(c) || c >= maxCode)) bs.set(env, c);
                }
                if (asciiRange) addAllMultiByteRange(env);
            }
            break;

        default:
            throw new InternalException(ErrorMessages.PARSER_BUG);
        } // switch
    }

    public static enum CCVALTYPE {
        SB,
        CODE_POINT,
        CLASS
    }

    public static enum CCSTATE {
        VALUE,
        RANGE,
        COMPLETE,
        START
    }

    public static final class CCStateArg {
        public int from;
        public int to;
        public boolean fromIsRaw;
        public boolean toIsRaw;
        public CCVALTYPE inType;
        public CCVALTYPE type;
        public CCSTATE state;
    }

    public void nextStateClass(CCStateArg arg, CClassNode ascCC, ScanEnvironment env) {
        if (arg.state == CCSTATE.RANGE) throw new SyntaxException(ErrorMessages.CHAR_CLASS_VALUE_AT_END_OF_RANGE);

        if (arg.state == CCSTATE.VALUE && arg.type != CCVALTYPE.CLASS) {
            if (arg.type == CCVALTYPE.SB) {
                bs.set(env, arg.from);
                if (ascCC != null) ascCC.bs.set(arg.from);
            } else if (arg.type == CCVALTYPE.CODE_POINT) {
                addCodeRange(env, arg.from, arg.from);
                if (ascCC != null) ascCC.addCodeRange(env, arg.from, arg.from, false);
            }
        }
        arg.state = CCSTATE.VALUE;
        arg.type = CCVALTYPE.CLASS;
    }

    public void nextStateValue(CCStateArg arg, CClassNode ascCc, ScanEnvironment env) {
        switch(arg.state) {
        case VALUE:
            if (arg.type == CCVALTYPE.SB) {
                bs.set(env, arg.from);
                if (ascCc != null) ascCc.bs.set(arg.from);
            } else if (arg.type == CCVALTYPE.CODE_POINT) {
                addCodeRange(env, arg.from, arg.from);
                if (ascCc != null) ascCc.addCodeRange(env, arg.from, arg.from, false);
            }
            break;

        case RANGE:
            if (arg.inType == arg.type) {
                if (arg.inType == CCVALTYPE.SB) {
                    if (arg.from > 0xff || arg.to > 0xff) throw new ValueException(ErrorMessages.ERR_INVALID_CODE_POINT_VALUE);

                    if (arg.from > arg.to) {
                        if (env.syntax.allowEmptyRangeInCC()) {
                            // goto ccs_range_end
                            arg.state = CCSTATE.COMPLETE;
                            break;
                        } else {
                            throw new ValueException(ErrorMessages.EMPTY_RANGE_IN_CHAR_CLASS);
                        }
                    }
                    bs.setRange(env, arg.from, arg.to);
                    if (ascCc != null) ascCc.bs.setRange(null, arg.from, arg.to);
                } else {
                    addCodeRange(env, arg.from, arg.to);
                    if (ascCc != null) ascCc.addCodeRange(env, arg.from, arg.to, false);
                }
            } else {
                if (arg.from > arg.to) {
                    if (env.syntax.allowEmptyRangeInCC()) {
                        // goto ccs_range_end
                        arg.state = CCSTATE.COMPLETE;
                        break;
                    } else {
                        throw new ValueException(ErrorMessages.EMPTY_RANGE_IN_CHAR_CLASS);
                    }
                }
                bs.setRange(env, arg.from, arg.to < 0xff ? arg.to : 0xff);
                addCodeRange(env, arg.from, arg.to);
                if (ascCc != null) {
                    ascCc.bs.setRange(null, arg.from, arg.to < 0xff ? arg.to : 0xff);
                    ascCc.addCodeRange(env, arg.from, arg.to, false);
                }
            }
            // ccs_range_end:
            arg.state = CCSTATE.COMPLETE;
            break;

        case COMPLETE:
        case START:
            arg.state = CCSTATE.VALUE;
            break;

        default:
            break;

        } // switch

        arg.fromIsRaw = arg.toIsRaw;
        arg.from = arg.to;
        arg.type = arg.inType;
    }

    // onig_is_code_in_cc_len
    boolean isCodeInCCLength(int encLength, int code) {
        boolean found;

        if (encLength > 1 || code >= BitSet.SINGLE_BYTE_SIZE) {
            if (mbuf == null) {
                found = false;
            } else {
                found = CodeRange.isInCodeRange(mbuf.getCodeRange(), code);
            }
        } else {
            found = bs.at(code);
        }

        if (isNot()) {
            return !found;
        } else {
            return found;
        }
    }

    // onig_is_code_in_cc
    public boolean isCodeInCC(Encoding enc, int code) {
        int len;
        if (enc.minLength() > 1) {
            len = 2;
        } else {
            len = enc.codeToMbcLength(code);
        }
        return isCodeInCCLength(len, code);
    }

    public void setNot() {
        flags |= FLAG_NCCLASS_NOT;
    }

    public boolean isNot() {
        return (flags & FLAG_NCCLASS_NOT) != 0;
    }
}
`
