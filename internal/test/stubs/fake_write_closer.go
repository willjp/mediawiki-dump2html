package test

import (
	"strings"
)

// Stub for io.WriteCloser.
//   Provide your own Writer so you can validate written text.
//   Provide a WriteError or CloseError to alter return value and test error handling.
type FakeWriteCloser struct {
	Writer     *strings.Builder
	CloseError error
	WriteError error
}

// Writes to provided Writer, returns WriteError if provided.
func (this *FakeWriteCloser) Write(p []byte) (n int, err error) {
	n, err = this.Writer.Write(p)
	if this.WriteError != nil {
		return n, this.WriteError
	}
	return n, err
}

// Returns CloseError, if provided.
func (this *FakeWriteCloser) Close() error {
	return this.CloseError
}

// Returns written text.
func (this *FakeWriteCloser) String() string {
	return this.Writer.String()
}
