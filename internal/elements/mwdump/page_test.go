package mwdump

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
