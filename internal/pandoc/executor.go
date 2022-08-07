package pandoc

import (
	"io"

	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
)

// Executes a pandoc.Cmd
// Implements interfaces.CmdExecutor
type Executor struct {
	cmd interfaces.ExecCmdExecutor
}

func (this *Executor) Execute(cmd interfaces.ExecCmdExecutor, stdin io.Reader) (render string, errs []error) {
	this.cmd = cmd
	return cmd.Execute(stdin)
}

func (this *Executor) Args() []string {
	return this.cmd.Args()
}
