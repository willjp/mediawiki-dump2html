TODO
====


Cleanups
--------

* `writers.Dump()`
   simplify error-handling using mv-on-write, this is more complicated than it needs to be.

* `renderers/html/css.go`
   this is still an untestable hack. maybe we could use make this a renderer unto itself and use `writers/`?


Bugfixes
--------

* `internal/elements/html` should not exist. HTML cannot be parsed by XML.

* `internal/renderers/html` Some anchors are not being identified/replaced.
   For example, `programming.html` has links that include `:` characters.


Features
--------

* `general`
  now that we have a direction, write tests!

* `renderers.HTML`
  index headers on each page to JSON obj, and write simple javascript search page.

* `renderers.HTML`
  use templating to re-add sidebar?

* `renderers.RST`
  rendered rst has directive `.. _first_title`
  I'm guessing this should be replaced with the actual page title.

* `writers.DumpAll()`
   page rendering could easily be asynchronous
