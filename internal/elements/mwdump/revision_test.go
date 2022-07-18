package mwdump

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRevisionUnmarshall(t *testing.T) {
	raw := `
	<revision>
	  <id>2222</id>
	  <parentid>2222</parentid>
	  <timestamp>2022-01-01T12:00:00Z</timestamp>
	  <contributor>
	    <username>Veronica</username>
	    <id>1</id>
	  </contributor>
	  <comment>/* Written Cleverly */</comment>
	  <origin>3333</origin>
	  <model>wikitext</model>
	  <format>text/x-wiki</format>
	  <text bytes="678" sha1="4fvil0i5inu7k04t7h43p3xn9rmijsl" xml:space="preserve">== Table of Contents ==</text>
	</revision>
	`
	var revision Revision
	xml.Unmarshal([]byte(raw), &revision)

	assert.Equal(t, "== Table of Contents ==", revision.Text)
	assert.Equal(t, time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC), revision.Timestamp)
}
