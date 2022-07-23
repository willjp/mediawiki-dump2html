package cli

import (
	"os"

	commands "willpittman.net/x/mediawiki-to-sphinxdoc/internal/cli/commands"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	html "willpittman.net/x/mediawiki-to-sphinxdoc/internal/renderers/html"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/writers"
)

type ArgumentParser struct {
	CliArgs []string
}

func New() ArgumentParser {
	return ArgumentParser{CliArgs: os.Args}
}

func (this *ArgumentParser) Parse() interfaces.CliCommand {
	shift := 0
	opts := commands.BuildOpts{}
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
			if len(this.CliArgs) <= i+2 {
				return &commands.InvalidArgs{}
			}
			opts.XMLDump = this.CliArgs[i+2]
			shift += 1
		case "-o", "--outdir":
			if len(this.CliArgs) <= i+2 {
				return &commands.InvalidArgs{}
			}
			opts.OutDir = this.CliArgs[i+2]
			shift += 1
		}
	}
	if opts.XMLDump != "" && opts.OutDir != "" {
		renderer := html.New()
		writer := writers.RenderWriter{}
		return &commands.Build{Opts: opts, Renderer: &renderer, RenderWriter: &writer}
	}
	return &commands.InvalidArgs{}
}
