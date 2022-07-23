package writers

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/log"
	stubs "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

var ExpectedError = errors.New("Expected")

func samplePage(revDate time.Time, title string) mwdump.Page {
	return mwdump.Page{
		Title: title,
		Revision: []mwdump.Revision{
			{Text: "== My New Header ==", Timestamp: revDate},
		},
	}
}

func TestDumpAll(t *testing.T) {
	setup := func(t *testing.T) {
		stubLog := logger.NewStubLogger()
		log.Log = &stubLog
		appfs.AppFs = afero.NewMemMapFs()
		outDir := "/var/tmp"
		err := appfs.AppFs.MkdirAll(outDir, 0755)
		assert.Nil(t, err)
	}

	t.Run("Performs Renderer.Setup()", func(t *testing.T) {
		setup(t)
		revDate := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
		pages := []mwdump.Page{samplePage(revDate, "one"), samplePage(revDate, "two")}
		dump := mwdump.XMLDump{Pages: pages}
		renderer := stubs.FakeRenderer{}
		outDir := "/var/tmp"
		errs := DumpAll(&renderer, &dump, outDir)
		assert.Nil(t, errs)
		assert.True(t, renderer.SetupCalled)
	})

	t.Run("Dumps pages", func(t *testing.T) {
		var err error
		var exists bool

		setup(t)
		Os := afero.Afero{Fs: appfs.AppFs}
		revDate := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
		pages := []mwdump.Page{samplePage(revDate, "one"), samplePage(revDate, "two")}
		dump := mwdump.XMLDump{Pages: pages}
		renderer := stubs.FakeRenderer{}
		outDir := "/var/tmp"
		errs := DumpAll(&renderer, &dump, outDir)
		assert.Nil(t, errs)

		exists, err = Os.Exists("/var/tmp/one")
		assert.Nil(t, err)
		assert.True(t, exists)

		exists, err = Os.Exists("/var/tmp/two")
		assert.Nil(t, err)
		assert.True(t, exists)
	})

	t.Run("Logs, Continues, and returns errs if error encountered during dump", func(t *testing.T) {
		setup(t)
		revDate := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
		pages := []mwdump.Page{samplePage(revDate, "one"), samplePage(revDate, "two")}
		dump := mwdump.XMLDump{Pages: pages}
		renderer := stubs.FakeRenderer{RenderErrors: []error{ExpectedError}}
		outDir := "/var/tmp"
		errs := DumpAll(&renderer, &dump, outDir)
		assert.Nil(t, errs)

		stubLog := log.Log.(*logger.StubLogger)
		assert.Equal(t, 2, len(stubLog.ErrorMsgs))
		assert.Equal(t, "Error dumping '/var/tmp/one' -- Expected", stubLog.ErrorMsgs[0])
		assert.Equal(t, "Error dumping '/var/tmp/two' -- Expected", stubLog.ErrorMsgs[1])
	})
}

func TestDump(t *testing.T) {
	setup := func(t *testing.T) {
		stubLog := logger.NewStubLogger()
		log.Log = &stubLog
		appfs.AppFs = afero.NewMemMapFs()
		outDir := "/var/tmp"
		err := appfs.AppFs.MkdirAll(outDir, 0755)
		assert.Nil(t, err)
	}

	t.Run("Writes Render to correct file when file does not exist", func(t *testing.T) {
		setup(t)
		Os := afero.Afero{Fs: appfs.AppFs}
		renderer := stubs.FakeRenderer{}
		page := samplePage(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), "file")
		outPath := "/var/tmp/file.txt"

		Dump(&renderer, &page, outPath)
		exists, err := Os.Exists(outPath)
		assert.Nil(t, err)
		assert.True(t, exists)

		stubLog := log.Log.(*logger.StubLogger)
		assert.Equal(t, 1, len(stubLog.InfoMsgs))
		assert.Equal(t, "Writing: /var/tmp/file.txt", stubLog.InfoMsgs[0])
	})

	t.Run("Writes Render to file when it exists, but is older than latest revision", func(t *testing.T) {
		setup(t)
		var err error
		outPath := "/var/tmp/file.txt"
		fileDate := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
		revDate := time.Date(2022, time.February, 1, 12, 0, 0, 0, time.UTC)

		Fs := appfs.AppFs
		Os := afero.Afero{Fs: Fs}
		Os.WriteFile(outPath, []byte("abc"), 0644)
		Fs.Chtimes(outPath, fileDate, fileDate)

		renderer := stubs.FakeRenderer{}
		page := samplePage(revDate, "file")

		Dump(&renderer, &page, outPath)
		finfo, err := Os.Stat(outPath)
		assert.Nil(t, err)
		assert.NotEqual(t, revDate, finfo.ModTime())

		stubLog := log.Log.(*logger.StubLogger)
		assert.Equal(t, 1, len(stubLog.InfoMsgs))
		assert.Equal(t, "Writing: /var/tmp/file.txt", stubLog.InfoMsgs[0])
	})

	t.Run("Does not write Render to file when it exists, but is newer than the latest revision", func(t *testing.T) {
		setup(t)
		outPath := "/var/tmp/file.txt"
		Fs := appfs.AppFs
		Os := afero.Afero{Fs: Fs}
		Os.WriteFile(outPath, []byte("abc"), 0644)
		oldFinfo, err := Os.Stat(outPath)

		renderer := stubs.FakeRenderer{}
		page := samplePage(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), "file")

		Dump(&renderer, &page, outPath)
		finfo, err := Os.Stat(outPath)
		assert.Nil(t, err)
		assert.Equal(t, oldFinfo.ModTime(), finfo.ModTime())

		stubLog := log.Log.(*logger.StubLogger)
		assert.Equal(t, 0, len(stubLog.InfoMsgs))
	})

	t.Run("Removes file if failure during write", func(t *testing.T) {
		setup(t)
		writeFileString = func(file afero.File, s string) (ret int, err error) {
			return 0, ExpectedError
		}
		Fs := appfs.AppFs
		Os := afero.Afero{Fs: Fs}
		renderer := stubs.FakeRenderer{}
		page := samplePage(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), "file")
		outPath := "/var/tmp/file.txt"

		errs := Dump(&renderer, &page, outPath)
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])

		exists, err := Os.Exists(outPath)
		assert.Nil(t, err)
		assert.False(t, exists)

		stubLog := log.Log.(*logger.StubLogger)
		assert.Equal(t, 1, len(stubLog.InfoMsgs))
		assert.Equal(t, "Writing: /var/tmp/file.txt", stubLog.InfoMsgs[0])
	})
}