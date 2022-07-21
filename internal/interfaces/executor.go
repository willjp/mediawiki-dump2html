package interfaces

import (
	"io"
)

type Executor interface {
	Execute(cmd CmdExecutor, stdin io.Reader) (render string, errs []error)
	Args() []string
}
