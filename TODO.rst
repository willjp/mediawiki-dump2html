TODO
====


Cleanups
--------


Bugfixes
--------

* `internal/elements/html` should not exist. HTML cannot be parsed by XML.

* `renderers.HTML` pandoc's produced CSS is tied to the rendered page.
  if we don't want to use standalone for each page, we'll need to create a dummy page to render with most syntax.


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
