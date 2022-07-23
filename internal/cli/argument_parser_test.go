package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	commands "willpittman.net/x/mediawiki-to-sphinxdoc/internal/cli/commands"
	html "willpittman.net/x/mediawiki-to-sphinxdoc/internal/renderers/html"
)

func TestArgumentParser(t *testing.T) {
	tcases := []struct {
		test    string
		args    []string
		expects string
	}{
		{
			test:    "-h returns ShowHelp command",
			args:    []string{"mw2html", "-h"},
			expects: fmt.Sprintf("%T", &commands.ShowHelp{}),
		},
		{
			test:    "--help returns ShowHelp command",
			args:    []string{"mw2html", "--help"},
			expects: fmt.Sprintf("%T", &commands.ShowHelp{}),
		},
		{
			test:    "-i, -o together returns Build command",
			args:    []string{"mw2html", "-i", "dump.xml", "-o", "/var/tmp/out"},
			expects: fmt.Sprintf("%T", &commands.Build{}),
		},
		{
			test:    "-i without -o returns InvalidArgs command",
			args:    []string{"mw2html", "-i", "dump.xml"},
			expects: fmt.Sprintf("%T", &commands.InvalidArgs{}),
		},
		{
			test:    "-i, -o together but missing values returns InvalidArgs command",
			args:    []string{"mw2html", "-i", "-o"},
			expects: fmt.Sprintf("%T", &commands.InvalidArgs{}),
		},
		{
			test:    "-o, -i together but missing values returns InvalidArgs command",
			args:    []string{"mw2html", "-o", "-i"},
			expects: fmt.Sprintf("%T", &commands.InvalidArgs{}),
		},
		{
			test:    "-o without -i returns InvalidArgs command",
			args:    []string{"mw2html", "-o", "/var/tmp/out"},
			expects: fmt.Sprintf("%T", &commands.InvalidArgs{}),
		},
		{
			test:    "--invalid-param returns InvalidArgs command",
			args:    []string{"mw2html", "--invalid-param"},
			expects: fmt.Sprintf("%T", &commands.InvalidArgs{}),
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.test, func(t *testing.T) {
			parser := ArgumentParser{CliArgs: tcase.args}
			command := parser.Parse()
			assert.Equal(t, tcase.expects, fmt.Sprintf("%T", command))

		})
	}

	t.Run("Build Command Receives Opts", func(t *testing.T) {
		parser := ArgumentParser{CliArgs: []string{"mw2html", "-i", "dump.xml", "-o", "/var/tmp/out"}}
		cmd := parser.Parse()
		cast := cmd.(*commands.Build)
		assert.Equal(t, "dump.xml", cast.Opts.XMLDump)
		assert.Equal(t, "/var/tmp/out", cast.Opts.OutDir)
	})

	t.Run("Build Command Receives Renderer", func(t *testing.T) {
		parser := ArgumentParser{CliArgs: []string{"mw2html", "-i", "dump.xml", "-o", "/var/tmp/out"}}
		cmd := parser.Parse()
		cast := cmd.(*commands.Build)
		expects := fmt.Sprintf("%T", &html.HTML{})
		assert.Equal(t, expects, fmt.Sprintf("%T", cast.Renderer))
	})
}
