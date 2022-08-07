package main

import (
	"os"

	"github.com/willjp/go-logger"
	"github.com/willjp/mediawiki-dump2html/internal/cli"
	"github.com/willjp/mediawiki-dump2html/internal/log"
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
