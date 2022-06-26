package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

// A mediawiki xml dump.
type XMLDump struct {
	XMLName  xml.Name `xml:"mediawiki"`
	SiteInfo string   `xml:"siteinfo"`
	Page     []Page   `xml:"page"`
}

// A page.
type Page struct {
	Title    string     `xml:"title"`
	Revision []Revision `xml:"revision"`
}

// A page revision.
//   Each revision is a standalone page version (not a diff).
//   The last revision is the latest.
type Revision struct {
	Text      string `xml:"text"`
	Timestamp string `xml:"timestamp"`
}

func main() {
	raw, err := os.ReadFile("/home/will/dump.xml")
	panicOn(err)

	var dump XMLDump
	xml.Unmarshal(raw, &dump)
	for _, page := range dump.Page {
		revision := page.Revision[len(page.Revision)-1]
		fmt.Println(revision)
	}
}
