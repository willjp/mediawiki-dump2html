package interfaces

import (
	"io"
)

// Executes a pandoc.Cmd
type PandocExecutor interface {
	Execute(cmd CmdExecutor, stdin io.Reader) (render string, errs []error)
	Args() []string
}
