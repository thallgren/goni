package posix

import "github.com/lyraproj/goni/goni/character"

var PBSNameStrings = []string {
	`alnum`,
	`alpha`,
	`blank`,
	`cntrl`,
	`digit`,
	`graph`,
	`lower`,
	`print`,
	`punct`,
	`space`,
	`upper`,
	`xdigit`,
	`ascii`,
	`word`}

var PBSValues  = []character.Type{
	character.Alnum,
	character.Alpha,
	character.Blank,
	character.Cntrl,
	character.Digit,
	character.Graph,
	character.Lower,
	character.Print,
	character.Punct,
	character.Space,
	character.Upper,
	character.XDigit,
	character.Ascii,
	character.Word,
}

var PBSNamesLower [][]byte

var PBSTableUpper map[string]character.Type

func init() {
	PBSNamesLower = make([][]byte, len(PBSNameStrings))
	PBSTableUpper = make(map[string]character.Type, len(PBSNameStrings))
	for i, n := range PBSNameStrings {
		PBSNamesLower[i] = []byte(n)
		PBSTableUpper[n] = PBSValues[i]
	}
}

