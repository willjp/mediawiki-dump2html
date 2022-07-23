package cli

import (
	commands "willpittman.net/x/mediawiki-to-sphinxdoc/internal/cli/commands"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

type ArgumentParser struct {
	CliArgs []string
}

type RenderOpts struct {
	XMLDump string
	OutDir  string
}

func New() ArgumentParser {
	return ArgumentParser{}
}

func (this *ArgumentParser) Parse() interfaces.CliCommand {
	shift := 0
	opts := RenderOpts{}
	for i, arg := range this.CliArgs[1:] {
		// skip N iterations, if flag consumes more than one param
		if shift > 0 {
			shift -= 1
			continue
		}

		switch arg {
		case "-h", "--help":
			return &commands.ShowHelp{}
		case "-i", "--input":
			if len(this.CliArgs) <= i+1 {
				return &commands.InvalidArgs{}
			}
			opts.XMLDump = this.CliArgs[i+1]
			shift += 1
		case "-o", "--outdir":
			if len(this.CliArgs) <= i+1 {
				return &commands.InvalidArgs{}
			}
			opts.OutDir = this.CliArgs[i+1]
			shift += 1
		}
	}
	if opts.XMLDump != "" && opts.OutDir != "" {
		return &commands.Build{}
	}
	return &commands.InvalidArgs{}
}
