package interfaces

import "io"

type CmdExecutor interface {
	Execute(stdin io.Reader) (render string, errs []error)
	Args() []string
}
