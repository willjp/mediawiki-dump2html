package pandoc

import "os/exec"

// Struct representing a pandoc command.
// (primarily intended for format conversions).
type Opts struct {
	From       string
	To         string
	Standalone bool
}

// Commandline Arguments for Pandoc
func (this *Opts) args() []string {
	args := []string{"-f", this.From, "-t", this.To}
	if this.Standalone == true {
		args = append(args, "--standalone")
	}
	return args
}

func (this *Opts) Command() Cmd {
	args := this.args()
	cmd := exec.Command("pandoc", args...)
	return Cmd{ExecCmd: cmd, args: cmd.Args}
}
