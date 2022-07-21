package test

import "io"

// Stub for interface of exec.Cmd that does not start a subprocess.
//  Provide stubs for Stdin, Stdout, Stderr to fake process inputs and behaviour.
type FakeExecCmd struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Args   []string
}

func (this FakeExecCmd) StdinPipe() (io.WriteCloser, error) {
	return this.Stdin, nil
}

func (this FakeExecCmd) StdoutPipe() (io.ReadCloser, error) {
	return this.Stdout, nil
}

func (this FakeExecCmd) StderrPipe() (io.ReadCloser, error) {
	return this.Stderr, nil
}

func (this FakeExecCmd) Start() error {
	return nil
}

func (this FakeExecCmd) Wait() error {
	return nil
}
