package mwdump

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPageUnmarshall(t *testing.T) {
	raw := `
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
	`
	var page Page
	xml.Unmarshal([]byte(raw), &page)

	assert.Equal(t, "Main Page", page.Title)
	assert.Equal(t, 2, len(page.Revision))
}

func TestLatestRevision(t *testing.T) {
	tcases := []struct {
		name         string
		revisions    []Revision
		expectedText string
	}{
		{
			name: "Identifies latest when order is ascending",
			revisions: []Revision{
				{
					Text:      "== My New Header ==",
					Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				{
					Text:      "== My Old Header ==",
					Timestamp: time.Date(2022, time.January, 1, 11, 0, 0, 0, time.UTC),
				},
			},
			expectedText: "== My New Header ==",
		},
		{
			name: "Identifies latest when order is descending",
			revisions: []Revision{
				{
					Text:      "== My Old Header ==",
					Timestamp: time.Date(2022, time.January, 1, 11, 0, 0, 0, time.UTC),
				},
				{
					Text:      "== My New Header ==",
					Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedText: "== My New Header ==",
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			page := Page{
				Title:    "Main Page",
				Revision: tcase.revisions,
			}
			latest := page.LatestRevision()
			assert.Equal(t, tcase.expectedText, latest.Text)
		})
	}
}
