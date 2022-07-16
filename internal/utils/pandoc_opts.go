package utils

import "os/exec"

// Struct representing a pandoc command.
// (primarily intended for format conversions).
type PandocOpts struct {
	From       string
	To         string
	Standalone bool
}

// Commandline Arguments for Pandoc
func (this *PandocOpts) args() []string {
	args := []string{"-f", this.From, "-t", this.To}
	if this.Standalone == true {
		args = append(args, "--standalone")
	}
	return args
}

func (this *PandocOpts) Command() *PandocCmd {
	args := this.args()
	cmd := exec.Command("pandoc", args...)
	return &PandocCmd{Cmd: cmd, Args: cmd.Args}
}
