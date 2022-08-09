TODO
====


Cleanups
--------


Bugfixes
--------

* delete pages not present in dump


Features
--------

* option to read from stdin, so can `php maintenance/dumpBackup.php --current --quiet | mediawiki-dump2html -o out/`
  without wasting disk writes on a tempfile

* `renderers.HTML`
  index headers on each page to JSON obj, and write simple javascript search page.

* `renderers.HTML`
  use templating to re-add sidebar?

* `renderers.HTML`
  we should have a table-of-contents on each page

* `renderers.RST`
  rendered rst has directive `.. _first_title`
  I'm guessing this should be replaced with the actual page title.

