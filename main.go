package main

import (
	"encoding/xml"
	"errors"
	"io/fs"
	"os"

	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/renderers"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/writers"
)

func main() {
	// configuration
	logger.SetLevel(logger.LvDebug)
	outDir := "/home/will/out"
	raw, err := os.ReadFile("/home/will/dump.xml")
	utils.PanicOn(err)

	_, err = os.Stat(outDir)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(outDir, 0755)
		utils.PanicOn(err)
	} else {
		utils.PanicOn(err)
	}

	var dump elements.XMLDump
	xml.Unmarshal(raw, &dump)

	renderer := renderers.HTML{}
	writers.DumpAll(&renderer, &dump, outDir)
}
