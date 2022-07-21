package interfaces

import (
	"io"
)

type PandocExecutor interface {
	Execute(cmd CmdExecutor, stdin io.Reader) (render string, errs []error)
	Args() []string
}
