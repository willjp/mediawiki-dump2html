package writers

import (
	"errors"
	"io/fs"
	"path"
	"time"

	"github.com/spf13/afero"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/log"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

func DumpAll(renderer interfaces.Renderer, dump *mwdump.XMLDump, outDir string) (errs []error) {
	errs = renderer.Setup(dump, outDir)

	for _, page := range dump.Pages {
		outPath := path.Join(outDir, renderer.Filename(page.Title))
		new_errs := Dump(renderer, &page, outPath)
		if new_errs != nil {
			errs = append(errs, new_errs...)
			for _, err := range new_errs {
				log.Log.Errorf("Error dumping '%s' -- %s", outPath, err)
			}
		}
	}
	return nil
}

func Dump(renderer interfaces.Renderer, page *mwdump.Page, outPath string) []error {
	var fileModified time.Time
	stat, err := appfs.AppFs.Stat(outPath)
	switch {
	case err == nil:
		fileModified = stat.ModTime()
	case errors.Is(err, fs.ErrNotExist):
		fileModified = time.Unix(0, 0)
	default:
		panic(err)
	}

	revision := page.LatestRevision()
	if revision.Timestamp.After(fileModified) {
		log.Log.Infof("Writing: %s", outPath)
		rendered, errs := renderer.Render(page)
		if errs != nil {
			return errs
		}

		file, err := appfs.AppFs.Create(outPath)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		_, err = writeFileString(file, rendered)
		if err != nil {
			utils.RmFileOn(file, err)
			errs = append(errs, err)
			return errs
		}
	}
	return nil
}

// test seam, writes a string to a file
var writeFileString = func(file afero.File, s string) (ret int, err error) {
	return file.WriteString(s)
}
