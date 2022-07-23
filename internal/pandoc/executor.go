package pandoc

import (
	"io"

	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

// Executes a pandoc.Cmd
// Implements interfaces.CmdExecutor
type Executor struct {
	cmd interfaces.CmdExecutor
}

func (this *Executor) Execute(cmd interfaces.CmdExecutor, stdin io.Reader) (render string, errs []error) {
	this.cmd = cmd
	return cmd.Execute(stdin)
}

func (this *Executor) Args() []string {
	return this.cmd.Args()
}
