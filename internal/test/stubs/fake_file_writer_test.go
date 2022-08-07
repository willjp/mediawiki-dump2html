package test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
)

var ExpectedError = errors.New("Expected")

func TestFakeFileWriter(t *testing.T) {
	t.Run("Implements Interface FileWriter", func(t *testing.T) {
		implementsInterface := func(interfaces.FileWriter) bool {
			return true
		}
		assert.True(t, implementsInterface(&FakeFileWriter{}))
	})

	t.Run("Returns Error when provided", func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		filename := "out.txt"
		file, err := Os.Create(filename)
		assert.Nil(t, err)

		writer := FakeFileWriter{Error: ExpectedError}
		written, err := writer.Write(file, []byte("abc"))
		assert.Equal(t, 0, written)
		assert.Error(t, ExpectedError, err)
	})

	t.Run("Writes File to Written when no Error", func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		filename := "out.txt"
		file, err := Os.Create(filename)
		assert.Nil(t, err)
		assert.NotNil(t, file)

		written := bytes.NewBuffer(nil)
		writer := FakeFileWriter{Written: written}
		n, err := writer.Write(file, []byte("abc"))
		assert.Nil(t, err)
		assert.Equal(t, 3, n)

		out := make([]byte, 3)
		_, err = writer.Written.Read(out)
		assert.Nil(t, err)
		assert.Equal(t, []byte("abc"), out)
	})
}
