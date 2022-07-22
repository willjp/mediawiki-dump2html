package main

import (
	"encoding/xml"
	"errors"
	"io/fs"

	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	renderers "willpittman.net/x/mediawiki-to-sphinxdoc/internal/renderers/html"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/writers"

	"github.com/spf13/afero"
)

func main() {
	// configuration
	Os := afero.Afero{Fs: appfs.AppFs}
	logger.SetLevel(logger.LvDebug)
	outDir := "/home/will/out"
	raw, err := Os.ReadFile("/home/will/dump.xml")
	utils.PanicOn(err)

	_, err = appfs.AppFs.Stat(outDir)
	if errors.Is(err, fs.ErrNotExist) {
		err := appfs.AppFs.MkdirAll(outDir, 0755)
		utils.PanicOn(err)
	} else {
		utils.PanicOn(err)
	}

	var dump mwdump.XMLDump
	xml.Unmarshal(raw, &dump)

	renderer := renderers.New()
	writers.DumpAll(&renderer, &dump, outDir)
}
