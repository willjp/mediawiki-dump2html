package test

import (
	"bytes"

	"github.com/spf13/afero"
)

type FakeFileWriter struct {
	Written *bytes.Buffer
	Error   error
}

func (this *FakeFileWriter) Write(file afero.File, b []byte) (n int, err error) {
	if this.Error != nil {
		return 0, this.Error
	}
	return this.Written.Write(b)
}
