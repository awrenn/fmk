package fmk

import (
	"os"
)

func init() {
	LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	return
}
