package cli

import (
	"errors"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	stubs "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

var ExpectedError = errors.New("Expected")

func sampleXml() string {
	return `
	<mediawiki xmlns="http://www.mediawiki.org/xml/export-0.11/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.mediawiki.org/xml/export-0.11/ http://www.mediawiki.org/xml/export-0.11.xsd" version="0.11" xml:lang="en">
	  <siteinfo>
	    <sitename>My Wiki</sitename>
	    <dbname>foodb</dbname>
	    <base>https://example.com/index.php/Main_Page</base>
	    <generator>Mediawiki 1.37.2</generator>
	    <case>first-letter</case>
	  </siteinfo>
	  <page>
	    <title>Main Page</title>
	    <ns>0</ns>
	    <id>1</id>
	    <revision>
	      <id>1111</id>
	      <parentid>2222</parentid>
	      <timestamp>2022-01-01T12:00:00Z</timestamp>
	      <contributor>
	        <username>Eliot</username>
	        <id>1</id>
	      </contributor>
	      <comment>/* Written Cleverly */</comment>
	      <origin>3333</origin>
	      <model>wikitext</model>
	      <format>text/x-wiki</format>
	      <text bytes="678" sha1="4fvil0i5inu7k04t7h43p3xn9rmijsl" xml:space="preserve">== Table of Contents ==</text>
	    </revision>
	    <revision>
	      <id>2222</id>
	      <parentid>2222</parentid>
	      <timestamp>2022-01-01T12:00:00Z</timestamp>
	      <contributor>
	        <username>Margo</username>
	        <id>1</id>
	      </contributor>
	      <comment>/* Written Cleverly */</comment>
	      <origin>3333</origin>
	      <model>wikitext</model>
	      <format>text/x-wiki</format>
	      <text bytes="678" sha1="4fvil0i5inu7k04t7h43p3xn9rmijsl" xml:space="preserve">== Table of Contents ==</text>
	    </revision>
	  </page>
	</mediawiki>
	`
}

func TestBuildImplementsInterface(t *testing.T) {
	implementsInterface := func(interfaces.CliCommand) bool {
		return true
	}
	assert.True(t, implementsInterface(&Build{}))
}

func TestBuildCall(t *testing.T) {
	setup := func(t *testing.T, xmlpath string, outdir string) *afero.Afero {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		Os.MkdirAll(path.Dir(xmlpath), 0755)
		Os.WriteFile(xmlpath, []byte(sampleXml()), 0644)
		return &Os
	}

	t.Run("Builds OutDir", func(t *testing.T) {
		xmlpath := "/home/you/dump.xml"
		outdir := "/var/tmp/out"
		Os := setup(t, xmlpath, outdir)
		renderer := stubs.FakeRenderer{}
		writer := stubs.FakeRenderWriter{}
		opts := BuildOpts{XMLDump: xmlpath, OutDir: outdir}
		cmd := Build{Opts: opts, Renderer: &renderer, RenderWriter: &writer}

		err := cmd.Call()
		assert.Nil(t, err)

		exists, err := Os.DirExists(outdir)
		assert.Nil(t, err)
		assert.True(t, exists)
	})

	t.Run("Calls DumpAll", func(t *testing.T) {
		xmlpath := "/home/you/dump.xml"
		outdir := "/var/tmp/out"
		setup(t, xmlpath, outdir)
		renderer := stubs.FakeRenderer{}
		writer := stubs.FakeRenderWriter{}
		opts := BuildOpts{XMLDump: xmlpath, OutDir: outdir}
		cmd := Build{Opts: opts, Renderer: &renderer, RenderWriter: &writer}

		err := cmd.Call()
		assert.Nil(t, err)

		assert.True(t, writer.DumpAllCalled)
	})

	t.Run("If DumpAll returns errors, returns BuildFailedError", func(t *testing.T) {
		xmlpath := "/home/you/dump.xml"
		outdir := "/var/tmp/out"
		setup(t, xmlpath, outdir)
		renderer := stubs.FakeRenderer{}
		writer := stubs.FakeRenderWriter{DumpAllErrors: []error{ExpectedError}}
		opts := BuildOpts{XMLDump: xmlpath, OutDir: outdir}
		cmd := Build{Opts: opts, Renderer: &renderer, RenderWriter: &writer}

		err := cmd.Call()
		assert.Error(t, BuildFailedError, err)
		assert.True(t, writer.DumpAllCalled)
	})
}
