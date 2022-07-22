package utils

import (
	"github.com/spf13/afero"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
)

// file expects a os.File
func RmFileOn(file afero.File, err error) {
	if err != nil {
		logger.Errorf("Error encountered, removing: %s", file.Name())
		appfs.AppFs.Remove(file.Name())
	}
}

func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func LogWarnOn(err error) {
	if err != nil {
		logger.Warnf("Ignored Error: %s", err)
	}
}
