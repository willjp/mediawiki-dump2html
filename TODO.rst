TODO
====


Cleanups
--------


Bugfixes
--------

* `internal/writers/dump` should read a modified timestamp written in HTML meta tags,
   instead of using the file modified date. Pages may be written to mediawiki,
   between the time a dump is created, and program is run. (so we could miss pages)

* `renderers.HTML` pandoc's produced CSS is tied to the rendered page.
  if we don't want to use standalone for each page, we'll need to create a dummy page with the elements we want styled.


Features
--------

* `writers.DumpAll()`
   page rendering could easily be asynchronous

* `renderers.HTML`
  index headers on each page to JSON obj, and write simple javascript search page.

* `renderers.HTML`
  use templating to re-add sidebar?

* `renderers.HTML`
  we should have a table-of-contents on each page

* `renderers.RST`
  rendered rst has directive `.. _first_title`
  I'm guessing this should be replaced with the actual page title.

