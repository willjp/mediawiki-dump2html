package utils

import (
	"github.com/spf13/afero"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
	"github.com/willjp/mediawiki-dump2html/internal/log"
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
