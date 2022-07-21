package test

import (
	"bytes"
	"io"

	pandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc"
)

// A FakePandocExecutor that stubs a pandoc execution, and records the provided stdin.
type FakePandocExecutor struct {
	Errors []error
	Render string
	Stdin  *bytes.Buffer
	Args   []string
}

//
// Reads provided `stdin` into `this.Stdin`, exposing a seam to test the text that was written to it.
func (this *FakePandocExecutor) Execute(cmd *pandoc.Cmd, stdin io.Reader) (render string, errs []error) {
	this.Args = cmd.Args
	this.Stdin = bytes.NewBuffer(nil)
	conts, err := io.ReadAll(stdin)
	if err != nil {
		panic(err)
	}
	_, err = this.Stdin.Write(conts)
	if err != nil {
		panic(err)
	}
	return this.Render, this.Errors
}
