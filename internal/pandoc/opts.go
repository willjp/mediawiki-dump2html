package pandoc

import "os/exec"

// Struct representing a pandoc command.
// (primarily intended for format conversions).
type Opts struct {
	From       string // -f mediawiki
	To         string // -t html
	Template   string // --template=highlighting-css.template
	Standalone bool   // --standalone
}

// Commandline Arguments for Pandoc
func (this *Opts) args() []string {
	args := []string{}
	if this.From != "" {
		args = append(args, "-f", this.From)
	}
	if this.To != "" {
		args = append(args, "-t", this.To)
	}
	if this.Template != "" {
		args = append(args, "--template="+this.Template)
	}
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
