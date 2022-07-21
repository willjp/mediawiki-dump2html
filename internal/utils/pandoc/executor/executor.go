package utils

import (
	"io"

	pandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc"
)

type Executor interface {
	Execute(cmd *pandoc.Cmd, stdin io.Reader) (render string, errs []error)
}
