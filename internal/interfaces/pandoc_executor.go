package interfaces

import (
	"io"
)

// Executes a pandoc.Cmd
type PandocExecutor interface {
	Execute(cmd ExecCmdExecutor, stdin io.Reader) (render string, errs []error)
	Args() []string
}
