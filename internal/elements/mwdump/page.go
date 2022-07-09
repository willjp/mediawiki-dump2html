package mwdump

// Represents a mediawiki <page> element
type Page struct {
	Title    string     `xml:"title"`
	Revision []Revision `xml:"revision"`
}

func (page *Page) LatestRevision() Revision {
	return page.Revision[len(page.Revision)-1]
}
