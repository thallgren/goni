package node

type AnyChar struct {
	node
}

func (a *AnyChar) String() string {
	NodeString(a)
}

func (a *AnyChar) Name() string {
	return `Any Char`
}

func (a *AnyChar) levelString(level int) string {
	return ``
}

func NewAnyChar() Node {
	return &AnyChar{node: node{nodeType: TypeCAny}}
}
