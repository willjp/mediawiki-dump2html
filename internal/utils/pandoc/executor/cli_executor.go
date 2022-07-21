package utils

import (
	"io"

	ipandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc/interfaces"
)

type CliExecutor struct {
	cmd ipandoc.CmdExecutor
}

func (this *CliExecutor) Execute(cmd ipandoc.CmdExecutor, stdin io.Reader) (render string, errs []error) {
	this.cmd = cmd
	return cmd.Execute(stdin)
}

func (this *CliExecutor) Args() []string {
	return this.cmd.Args()
}
