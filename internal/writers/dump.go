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

func DumpAll(renderer renderers.Renderer, dump *mwdump.XMLDump, outDir string) error {
	renderer.Setup(dump, outDir)

	for _, page := range dump.Pages {
		outPath := path.Join(outDir, renderer.Filename(&page))
		Dump(renderer, &page, outPath)
	}
	return nil
}

func Dump(renderer renderers.Renderer, page *mwdump.Page, outPath string) error {
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
		rendered, err := renderer.Render(page)
		if err != nil {
			return err
		}

		file, err := os.Create(outPath)
		defer file.Close()
		utils.PanicOn(err)

		_, err = file.WriteString(rendered)
		if err != nil {
			utils.RmFileOn(file, err)
			return err
		}
	}
	return nil
}
