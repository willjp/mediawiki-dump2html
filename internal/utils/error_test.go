package utils

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
)

func TestRmFileOn(t *testing.T) {
	t.Run("Deletes file on error", func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		file, err := appfs.AppFs.Create("/var/tmp/foo.txt")
		assert.Nil(t, err)

		RmFileOn(file, fmt.Errorf("error"))
		exists, err := Os.Exists("/var/tmp/foo.txt")
		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("Does not delete file when no error", func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		file, err := appfs.AppFs.Create("/var/tmp/foo.txt")
		assert.Nil(t, err)

		RmFileOn(file, nil)
		exists, err := Os.Exists("/var/tmp/foo.txt")
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}
