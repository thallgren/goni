package internal

const RegionNotpos = -1

type Region struct {
	numRegs int
	beg []int;
	end []int;
	historyRoot *captureTreeNode
}
