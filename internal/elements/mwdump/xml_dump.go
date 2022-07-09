package mwdump

import "encoding/xml"

// Represents a mediawiki <mediawiki> element (the root element in a mediawiki xml dump).
type XMLDump struct {
	XMLName  xml.Name `xml:"mediawiki"`
	SiteInfo string   `xml:"siteinfo"`
	Pages    []Page   `xml:"page"`
}
