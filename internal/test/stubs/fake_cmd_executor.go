package test

import (
	"bytes"
	"io"
)

// implements interfaces/CmdExecutor
type FakeCmd struct {
	Errors  []error
	Render  string
	CliArgs []string
	Stdin   *bytes.Buffer
}

func (this *FakeCmd) Execute(stdin io.Reader) (render string, errs []error) {
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

func (this *FakeCmd) Args() []string {
	return this.CliArgs
}
