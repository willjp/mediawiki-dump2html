package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

func TestArgs(t *testing.T) {
	tcases := []struct {
		test    string
		pandoc  utils.Pandoc
		expects []string
	}{
		{
			test:    "To/From as -f FMT -t FMT ",
			pandoc:  utils.Pandoc{From: "mediawiki", To: "html"},
			expects: []string{"-f", "mediawiki", "-t", "html"},
		},
		{
			test:    "Standalone as --standalone",
			pandoc:  utils.Pandoc{From: "mediawiki", To: "html", Standalone: true},
			expects: []string{"-f", "mediawiki", "-t", "html", "--standalone"},
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.test, func(t *testing.T) {
			args := tcase.pandoc.Args()
			assert.Equal(t, args, tcase.expects)
		})
	}
}
