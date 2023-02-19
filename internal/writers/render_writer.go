package writers

import (
	"path"
	"time"

	"github.com/spf13/afero"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
	"github.com/willjp/mediawiki-dump2html/internal/log"
	"github.com/willjp/mediawiki-dump2html/internal/utils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Writes output from a Renderer to disk
// Implements interfaces.RenderWriter
type RenderWriter struct{}

// Write all pages from provided xmldump
func (this *RenderWriter) DumpAll(renderer interfaces.Renderer, dump *mwdump.XMLDump, outDir string) (errs []error) {
	errs = renderer.Setup(dump, outDir)

	for _, page := range dump.Pages {
		outPath := path.Join(outDir, renderer.Filename(page.Title))
		new_errs := this.Dump(renderer, &page, outPath)
		if new_errs != nil {
			errs = append(errs, new_errs...)
			for _, err := range new_errs {
				log.Log.Errorf("Error dumping '%s' -- %s", outPath, err)
			}
		}
	}
	return nil
}

// Write a single page from an xmldump if it's revision timestamp is newer than the one on disk.
// If page does not exist, renders a new page.
func (this *RenderWriter) Dump(renderer interfaces.Renderer, page *mwdump.Page, outPath string) (errs []error) {
	renderedRevisionDate := this.renderDate(outPath)
	revision := page.LatestRevision()
	if revision.Timestamp.After(renderedRevisionDate) {
		log.Log.Infof("Writing: %s", outPath)
		rendered, errs := renderer.Render(page)
		if errs != nil {
			return errs
		}

		file, err := appfs.AppFs.Create(outPath)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		_, err = writeFileString(file, rendered)
		if err != nil {
			utils.RmFileOn(file, err)
			errs = append(errs, err)
			return errs
		}
	} else {
		log.Log.Debugf("Skipping Up To Date: %s", outPath)
	}
	return nil
}

// Get the mediawiki published date for a rendered html file.
// Returns unix epoch if file doesn't exist, or any error encountered.
func (this *RenderWriter) renderDate(filepath string) (revisionDate time.Time) {
	unixEpoch := time.Unix(0, 0)
	Os := afero.Afero{Fs: appfs.AppFs}
	exists, _ := Os.Exists(filepath)
	if !exists {
		return unixEpoch
	}

	file, _ := Os.Open(filepath)
	defer file.Close()
	node, _ := html.Parse(file)
	metaNode := utils.FindFirstChild(node, func(node *html.Node) *html.Node {
		return utils.HasParentHeirarchy(node, []atom.Atom{atom.Head, atom.Meta})
	})
	if metaNode != nil {
		for _, attr := range metaNode.Attr {
			if attr.Key != "content" {
				continue
			}
			revisionDate, err := time.Parse(time.RFC3339, attr.Val)
			if err == nil {
				return revisionDate
			}
		}
	}
	return unixEpoch
}

// test seam, writes a string to a file
var writeFileString = func(file afero.File, s string) (ret int, err error) {
	return file.WriteString(s)
}
