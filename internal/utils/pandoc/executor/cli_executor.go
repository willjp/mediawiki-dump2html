package utils

import (
	"io"

	ipandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc/interfaces"
)

type CliExecutor struct{}

func (this *CliExecutor) Execute(cmd ipandoc.CmdExecutor, stdin io.Reader) (render string, errs []error) {
	return cmd.Execute(stdin)
}
