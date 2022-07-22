package utils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
)

func TestFileReplace(t *testing.T) {
	t.Run("Writes File", func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}

		errs := FileReplace("abc", "/var/tmp/foo.txt")
		assert.Nil(t, errs)

		res, err := Os.ReadFile("/var/tmp/foo.txt")
		assert.Nil(t, err)
		assert.Equal(t, []byte("abc"), res)
	})
}
