package interfaces

import "io"

// Executes an exec.Cmd
type CmdExecutor interface {
	Execute(stdin io.Reader) (render string, errs []error)
	Args() []string
}
