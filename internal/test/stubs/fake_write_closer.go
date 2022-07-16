package test

import "strings"

type FakeWriteCloser struct {
	Writer     strings.Builder
	CloseError error
}

func (this FakeWriteCloser) Write(p []byte) (n int, err error) {
	return this.Writer.Write(p)
}

func (this FakeWriteCloser) Close() error {
	return this.CloseError
}
