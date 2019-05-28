package internal

const stackMachineInvalidIndex = -1

type stackMachine struct {
	matcher

	stack []*stackEntry
	stk int  // stkEnd
	repeatStk []int
	memStartStk int
	memEndStk int
	stateCheckBuff []byte // CEC, move to int[] ?
	stateCheckBuffSize int
}
