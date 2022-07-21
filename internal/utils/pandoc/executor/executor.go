package utils

import (
	"io"

	ipandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc/interfaces"
)

type Executor interface {
	Execute(cmd ipandoc.CmdExecutor, stdin io.Reader) (render string, errs []error)
	Args() []string
}
