package internal

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/option"
	"github.com/lyraproj/goni/goni/search"
)

type regex struct {
	code []int             /* compiled pattern */
	codeLength int
	requireStack bool

	numMem int             /* used memory(...) num counted from 1 */
	numRepeat int          /* OP_REPEAT/OP_REPEAT_NG id-counter */
	numNullCheck int       /* OP_NULL_CHECK_START/END id counter */
	numCombExpCheck int    /* combination explosion check */
	numCall int            /* number of subexp call */
	captureHistory int     /* (?@...) flag (1-31) */
	btMemStart int         /* need backtrack flag */
	btMemEnd int           /* need backtrack flag */

	stackPopLevel int

	o []int
	i []int

	factory goni.MatcherFactory

	enc goni.Encoding
	options option.Type
	userOptions int
	userObject interface{}
	caseFoldFlag int

	nameTable map[string]int // named entries

	/* optimization info (string search, char-map and anchors) */
	forward search.Forward ;                 /* optimize flag */
	backward search.Backward ;
	thresholdLength int                    /* search str-length for apply optimize */
	anchor int                             /* BEGIN_BUF, BEGIN_POS, (SEMI_)END_BUF */
	anchorDmin int                         /* (SEMI_)END_BUF anchor distance */
	anchorDmax int                         /* (SEMI_)END_BUF anchor distance */
	subAnchor int                          /* start-anchor for exact or map */

	exact []byte
	exactP int
	exactEnd int

	mp []byte;                              /* used as BM skip or char-map */
	p []int                            /* BM skip for exact_len > 255 */
	d []int                    /* BM skip for backward search */
	dMin int                               /* min-distance of exact or map */
	dMax int                               /* max-distance of exact or map */

	templates [][]byte                      /* fixed pattern strings not embedded in bytecode */
	templateNum int
}

func (rx *regex) nameToGroupNumbers([]byte, int, int) *goni.NameEntry {
	// TODO
	return nil
}

func (rx *regex) nameAdd(bytes []byte, i int, i2 int, i3 int, syntax *goni.Syntax) {

}


