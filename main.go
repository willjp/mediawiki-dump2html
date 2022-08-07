package main

import (
	"os"

	"github.com/willjp/go-logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/cli"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/log"
)

func main() {
	log.Log.SetLevel(logger.LvDebug)
	cli := cli.New()
	cmd := cli.Parse()
	err := cmd.Call()
	if err != nil {
		log.Log.Error("[ERROR] error during build")
		os.Exit(1)
	}
}
