TODO
====

* Why use rst at all?
  Why don't we:
    pandoc to html,
    index the headers to a json object,
    write a simple javascript search page.
    score search results by header priority.
    we could even include a text-snippet if we wanted.

* `internal/renderers` should not own `Write()`.
   Generate the relpath, correct the links,
   but writing in the renderer would just duplicate work.

* rendered rst has directive `.. _first_title`
  I'm guessing this should be replaced with the actual page title.

* `internal/elements/page.invalidPathCh`
   we're rewriting the files, but surely we'll need to adjust the links within the rendered html

* `internal/elements.page.Page.pandocWikiToRst()`
   simplify error-handling using mv-on-write, this is bloated.

* `main.go`
   page rendering could easily be asynchronous
