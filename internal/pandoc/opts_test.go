package pandoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpts(t *testing.T) {
	tcases := []struct {
		test    string
		opts    Opts
		expects []string
	}{
		{
			test:    "Command() parses 'To,From' as params '-f FMT -t FMT'",
			opts:    Opts{From: "mediawiki", To: "html"},
			expects: []string{"pandoc", "-f", "mediawiki", "-t", "html"},
		},
		{
			test:    "Command() parses 'Standalone' as params '--standalone'",
			opts:    Opts{From: "mediawiki", To: "html", Standalone: true},
			expects: []string{"pandoc", "-f", "mediawiki", "-t", "html", "--standalone"},
		},
		{
			test:    "Command() parses 'Template' as params '--template=TEMPLATE'",
			opts:    Opts{Template: "/foo/bar/baz.template"},
			expects: []string{"pandoc", "--template=/foo/bar/baz.template"},
		},
	}
	for _, tcase := range tcases {
		cmd := tcase.opts.Command()
		t.Run(tcase.test, func(t *testing.T) {
			assert.Equal(t, tcase.expects, cmd.Args())
		})
	}
}
