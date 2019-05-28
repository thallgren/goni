package internal

type byteCodeMachine struct {
	stackMachine

	interruptCheckCounter int // we modulos this to occasionally check for interrupts
	interrupted bool

	bestLen int          // return value
	s int         // current char

	rng int            // right range
	sprev int
	sstart int
	sbegin int
	pkeep int

	code []int        // byte code
	ip int                 // instruction pointer
}
