package test

import "bytes"

type FakeFile struct {
	Path       string
	Buffer     *bytes.Buffer
	CloseError error
}

func NewFakeFile() FakeFile {
	buf := bytes.NewBuffer(nil)
	return FakeFile{Buffer: buf}
}

func (this *FakeFile) Name() string {
	return this.Path
}

func (this *FakeFile) WriteString(v string) (n int, err error) {
	return this.Buffer.WriteString(v)
}

func (this *FakeFile) Close() error {
	return this.CloseError
}
