package mwdump

// Represents a mediawiki <page> element
type Page struct {
	Title    string     `xml:"title"`
	Revision []Revision `xml:"revision"`
}

func (page *Page) LatestRevision() Revision {
	latest := page.Revision[0]
	for _, revision := range page.Revision {
		if latest.Timestamp.Before(revision.Timestamp) {
			latest = revision
		}
	}
	return latest
}
