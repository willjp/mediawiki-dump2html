package test

import "io"

// implements interfaces/CmdExecutor
type FakeCmd struct {
	CliArgs []string
}

func (this *FakeCmd) Execute(stdin io.Reader) (render string, errs []error) {
	return "", nil
}

func (this *FakeCmd) Args() []string {
	return this.CliArgs
}
