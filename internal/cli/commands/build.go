package cli

import "fmt"

type Build struct{}

func (this *Build) Call() error {
	fmt.Println("build")
type BuildOpts struct {
	XMLDump string // path to xmldump
	OutDir  string // directory to write statichtml files
}

type Build struct {
	Opts BuildOpts
}
	return nil
}
