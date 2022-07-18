package test

type FakeFile struct {
	Path string
}

func (this FakeFile) Name() string {
	return this.Path
}
