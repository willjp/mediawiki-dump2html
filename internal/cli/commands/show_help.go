package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/lithammer/dedent"
)

// Prints the help menu on the commandline
// Concretion of interfaces.CliCommand
type ShowHelp struct{}

func (this *ShowHelp) Call() error {
	executable := path.Base(os.Args[0])
	fmt.Printf(dedent.Dedent(`
	%s [-i INPUT -o OUTDIR] [-h]

	DESCRIPTION:
	    Converts a mediawiki XML dump to a statichtml website

	PARAMS:
	  -h --help
	      show this help menu

	  -i --input INPUT
	      the xml dump you'd like to convert

	  -o --outdir OUTDIR
	      the directory you'd like to write html files to

	EXAMPLES:
	    %s -i dump.xml -o /var/tmp/website/  # convert dump.xml to statichtml in /var/tmp/website
	`), executable, executable)
	return nil
}
