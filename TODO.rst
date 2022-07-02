TODO
====

* `internal/elements/page.invalidPathCh`
   is pretty sparse, use a posix filepath char allowlist.

* `internal/elements/page.invalidPathCh`
   we're rewriting the files, but surely we'll need to adjust the links within the rendered html

* `internal/elements.page.Page.pandocWikiToRst()`
   simplify error-handling using mv-on-write, this is bloated.

* `main.go`
   page rendering could easily be asynchronous
