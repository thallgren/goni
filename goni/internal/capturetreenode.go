package internal

type captureTreeNode struct {
	group int
	beg int
	end int

	children []*captureTreeNode
}

const historyTreeInitAllocSize = 8

func newCaptureTreeNode() *captureTreeNode {
	return &captureTreeNode{beg: RegionNotpos, end: RegionNotpos, group: -1 }
}

func (cn *captureTreeNode) addChild(child *captureTreeNode) {
	if cn.children == nil {
		cn.children = make([]*captureTreeNode, 0, historyTreeInitAllocSize)
	}
	cn.children = append(cn.children, child)
}

func (cn *captureTreeNode) clear() {
	cn.children = cn.children[:0]
	cn.beg = RegionNotpos
	cn.end = RegionNotpos
	cn.group = -1
}

func (cn *captureTreeNode) cloneTree() *captureTreeNode {
	clone := newCaptureTreeNode()
	clone.beg = cn.beg
	clone.end = cn.end
	if cn.children != nil {
		clone.children = make([]*captureTreeNode, len(cn.children), cap(cn.children))
		for i, c := range cn.children {
			clone.children[i] = c.cloneTree()
		}
	}
	return clone
}
