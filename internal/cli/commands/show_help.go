package cli

import (
	"fmt"
	"os"
	"regexp"
)

type ShowHelp struct{}

func (this *ShowHelp) Call() error {
	leadingWhitespaceRx := regexp.MustCompile(`(?m)(^\w+)`)
	fmt.Printf(
		leadingWhitespaceRx.ReplaceAllString(`
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
	`, ""), os.Args[0], os.Args[0])
	return nil
}
