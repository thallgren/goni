package ast

const FlagNcCClassNot = 1<<0

type CClass struct {
	abstractNode
	flags int
}


