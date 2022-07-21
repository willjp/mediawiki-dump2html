package test

import (
	"bytes"
	"io"

	ipandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc/interfaces"
)

// A FakePandocExecutor that stubs a pandoc execution, and records the provided stdin.
type FakePandocExecutor struct {
	Errors []error
	Render string
	Stdin  *bytes.Buffer
	args   []string
}

//
// Reads provided `stdin` into `this.Stdin`, exposing a seam to test the text that was written to it.
func (this *FakePandocExecutor) Execute(cmd ipandoc.CmdExecutor, stdin io.Reader) (render string, errs []error) {
	this.args = cmd.Args()
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

func (this *FakePandocExecutor) Args() []string {
	return this.args
}
