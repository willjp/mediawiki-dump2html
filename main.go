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
)

func main() {
	logger.SetLevel(logger.LvDebug)
	// configuration
	sphinxRoot := "/home/will/out"
	raw, err := os.ReadFile("/home/will/dump.xml")
	utils.PanicOn(err)

	_, err = os.Stat(sphinxRoot)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(sphinxRoot, 0755)
		utils.PanicOn(err)
	} else {
		utils.PanicOn(err)
	}

	renderer := renderers.RST{}
	var dump elements.XMLDump
	xml.Unmarshal(raw, &dump)
	for _, page := range dump.Pages {
		renderer.Write(&page, sphinxRoot)
	}
}
