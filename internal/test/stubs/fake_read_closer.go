package test

import (
	"io"
)

// Stub for io.ReadCloser.
//   Provide your own Reader or ReadError/CloseError and test inputs/error-handling.
type FakeReadCloser struct {
	Reader     io.Reader
	CloseError error
	ReadError  error
}

// Reads from provided Reader, returns ReadError if provided.
func (this *FakeReadCloser) Read(p []byte) (int, error) {
	written, err := this.Reader.Read(p)
	if this.ReadError != nil {
		return written, this.ReadError
	}
	return written, err
}

// Returns CloseError, if provided.
func (this *FakeReadCloser) Close() error {
	return this.CloseError
}
