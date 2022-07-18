package mwdump

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXmlDumpUnmarshall(t *testing.T) {
	raw := `
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
	var xmlDump XMLDump
	xml.Unmarshal([]byte(raw), &xmlDump)

	assert.Equal(t, xml.Name(xml.Name{Space: "http://www.mediawiki.org/xml/export-0.11/", Local: "mediawiki"}), xmlDump.XMLName)
	assert.Equal(t, 1, len(xmlDump.Pages))
}
