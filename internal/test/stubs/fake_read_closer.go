package test

import "io"

type FakeReadCloser struct {
	Reader     io.Reader
	CloseError error
}

func (this FakeReadCloser) Read(p []byte) (int, error) {
	return this.Reader.Read(p)
}

func (this FakeReadCloser) Close() error {
	return this.CloseError
}
