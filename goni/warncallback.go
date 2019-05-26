package goni

import (
	"fmt"
	"os"
)

type WarnCallback interface {
	Warn(message string)
}

type WarnNone struct {}

func (*WarnNone) Warn(message string) {
}

type WarnDefault struct {}

func (*WarnDefault) Warn(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
}
