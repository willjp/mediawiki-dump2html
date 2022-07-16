package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPandocOpts(t *testing.T) {
	tcases := []struct {
		test    string
		opts    PandocOpts
		expects []string
	}{
		{
			test:    "Command() parses 'To,From' as params '-f FMT -t FMT'",
			opts:    PandocOpts{From: "mediawiki", To: "html"},
			expects: []string{"pandoc", "-f", "mediawiki", "-t", "html"},
		},
		{
			test:    "Command() parses 'Standalone' as params '--standalone'",
			opts:    PandocOpts{From: "mediawiki", To: "html", Standalone: true},
			expects: []string{"pandoc", "-f", "mediawiki", "-t", "html", "--standalone"},
		},
	}
	for _, tcase := range tcases {
		cmd := tcase.opts.Command()
		t.Run(tcase.test, func(t *testing.T) {
			assert.Equal(t, tcase.expects, cmd.Args)
		})
	}
}
