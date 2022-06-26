package elements

import "time"

// Represents a mediawiki <revision> element (page revision).
//   Each revision is a standalone page version (not a diff).
//   The last revision is the latest.
type Revision struct {
	Text      string    `xml:"text"`
	Timestamp time.Time `xml:"timestamp"`
}
