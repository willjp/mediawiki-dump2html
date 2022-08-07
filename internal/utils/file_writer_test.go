package utils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
)

func TestFileWriter(t *testing.T) {
	t.Run("Implements Interface FileWriter", func(t *testing.T) {
		implementsInterface := func(interfaces.FileWriter) bool {
			return true
		}
		assert.True(t, implementsInterface(&FileWriter{}))
	})

	t.Run("Writes File", func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		filename := "out.txt"
		file, err := Os.Create(filename)
		assert.Nil(t, err)

		writer := FileWriter{}
		written, err := writer.Write(file, []byte("abc"))
		assert.Nil(t, err)
		assert.Equal(t, 3, written)

		success, err := Os.FileContainsBytes(filename, []byte("abc"))
		assert.Nil(t, err)
		assert.True(t, success)
	})
}
