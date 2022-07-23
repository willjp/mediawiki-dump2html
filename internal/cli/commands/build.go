package cli

import (
	"encoding/xml"
	"errors"
	"io/fs"

	"github.com/spf13/afero"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

var BuildFailedError = errors.New("Build Failed")

type BuildOpts struct {
	XMLDump string // path to xmldump
	OutDir  string // directory to write statichtml files
}

type Build struct {
	Opts         BuildOpts
	Renderer     interfaces.Renderer
	RenderWriter interfaces.RenderWriter
}

func (this *Build) Call() (err error) {
	raw := this.readXml()
	dump := this.parseXml(raw)
	this.createOutDir()

	errs := this.RenderWriter.DumpAll(this.Renderer, dump, this.Opts.OutDir)
	if errs != nil {
		return BuildFailedError
	}
	return nil
}

func (this *Build) readXml() []byte {
	Os := afero.Afero{Fs: appfs.AppFs}
	raw, err := Os.ReadFile(this.Opts.XMLDump)
	if err != nil {
		panic(err)
	}
	return raw
}

func (this *Build) parseXml(raw []byte) *mwdump.XMLDump {
	var dump mwdump.XMLDump
	err := xml.Unmarshal(raw, &dump)
	if err != nil {
		panic(err)
	}
	return &dump
}

func (this *Build) createOutDir() {
	_, err := appfs.AppFs.Stat(this.Opts.OutDir)
	if errors.Is(err, fs.ErrNotExist) {
		err = appfs.AppFs.MkdirAll(this.Opts.OutDir, 0755)
	}
	if err != nil {
		panic(err)
	}
}
