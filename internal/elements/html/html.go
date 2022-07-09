package html

import "encoding/xml"

type Html struct {
	XMLName xml.Name `xml:"http://www.w3.org/1999/xhtml html"`
	Head    Head     `xml:"head"`
}
