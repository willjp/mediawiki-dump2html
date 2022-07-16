package test

import "io"

type FakeCmd struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Args   []string
}

func (this FakeCmd) StdinPipe() (io.WriteCloser, error) {
	return this.Stdin, nil
}

func (this FakeCmd) StdoutPipe() (io.ReadCloser, error) {
	return this.Stdout, nil
}

func (this FakeCmd) StderrPipe() (io.ReadCloser, error) {
	return this.Stderr, nil
}

func (this FakeCmd) Start() error {
	return nil
}

func (this FakeCmd) Wait() error {
	return nil
}
