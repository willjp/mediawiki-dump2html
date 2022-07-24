package utils

import (
	"github.com/spf13/afero"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/log"
)

// file expects a os.File
func RmFileOn(file afero.File, err error) {
	if err != nil {
		log.Log.Errorf("Error encountered, removing: %s", file.Name())
		appfs.AppFs.Remove(file.Name())
	}
}

func LogWarnOn(err error) {
	if err != nil {
		log.Log.Warnf("Ignored Error: %s", err)
	}
}
