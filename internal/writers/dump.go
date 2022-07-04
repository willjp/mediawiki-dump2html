package writers

import (
	"errors"
	"io/fs"
	"os"
	"time"

	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/renderers"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

func Dump(renderer renderers.Renderer, page *elements.Page, outPath string) error {
	rmFileOn := func(file *os.File, err error) {
		if err != nil {
			logger.Errorf("Error encountered, removing: %s", file.Name())
			os.Remove(file.Name())
		}
	}

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
		file, err := os.Create(outPath)
		utils.PanicOn(err)

		logger.Infof("Writing: %s\n", outPath)
		rendered, err := renderer.Render(page)
		if err != nil {
			rmFileOn(file, err)
			return err
		}
		_, err = file.WriteString(rendered)
		if err != nil {
			rmFileOn(file, err)
			return err
		}
	}
	return nil
}
