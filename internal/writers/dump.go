package writers

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"time"

	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/renderers"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

func DumpAll(renderer renderers.Renderer, dump *mwdump.XMLDump, outDir string) (errs []error) {
	renderer.Setup(dump, outDir)

	for _, page := range dump.Pages {
		outPath := path.Join(outDir, renderer.Filename(page.Title))
		new_errs := Dump(renderer, &page, outPath)
		if new_errs != nil {
			errs = append(errs, new_errs...)
			for _, err := range new_errs {
				logger.Errorf("Error dumping '%s' -- %s", outPath, err)
			}
		}
	}
	return nil
}

func Dump(renderer renderers.Renderer, page *mwdump.Page, outPath string) []error {
	var fileModified time.Time
	stat, err := os.Stat(outPath)
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
		logger.Infof("Writing: %s\n", outPath)
		cmd := renderer.RenderCmd()
		rendered, errs := renderer.RenderExec(cmd, page)
		if errs != nil {
			return errs
		}

		file, err := os.Create(outPath)
		defer file.Close()
		utils.PanicOn(err)

		_, err = file.WriteString(rendered)
		if err != nil {
			utils.RmFileOn(file, err)
			errs = append(errs, err)
			return errs
		}
	}
	return nil
}
