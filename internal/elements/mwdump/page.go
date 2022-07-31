package mwdump

// Represents a mediawiki <page> element
type Page struct {
	Title     string     `xml:"title"`
	Revisions []Revision `xml:"revision"`
}

func (page *Page) LatestRevision() Revision {
	latest := page.Revisions[0]
	for _, revision := range page.Revisions {
		if latest.Timestamp.Before(revision.Timestamp) {
			latest = revision
		}
	}
	return latest
}
