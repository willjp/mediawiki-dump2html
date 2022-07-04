
Parsers
=======


Winning Solution-2
------------------

Why bother with an intermediate ReStructuredText, and another system if we just can do it here.

1. Iterate over latest pages from XML dump

2. Create a lookup table of all to-be-renamed pages, from their title, to their renamed posix-filepath

3. Use pandoc to convert to HTML

4. Increment headers, add page-title as h1

5. Replace targets in links that have been renamed for posix

6. Index headers in json blob (so we can search later)

7. Write static-html files


.. warning::
   pandoc actually supports syntax-highlighting with it's `--standalone` flag.
   `pandoc --print-highlight-style pygments` will dump a json fed to skylight
   `pandoc --list-highlight-styles` will show available colourschemes.

   I don't think you can dump a CSS file, so maybe render a page (with `--standalone`), and extract/point to the defined stylesheets?

.. note::
    Real HTML parsing would be preferable to regex, but it should do for now.
    We might be able to avoid the posix-rename if we esacape URLs, and filepaths -- which would make this faster, but dirtier to work with.


Winning Solution-1
------------------

Steps
.....

1. create a single xml file per page (remove all surrounding mediawiki xml)

2. use pandoc to convert mediawiki to rst (in-memory)


3. patch <br>s
  ::

    prepend to each page
        .. role:: raw-html(raw)
          :format: html

    and replace every `<br>` with
        :raw-html:`<br/>`

4. prepend title into the page.
   ::

     =========
     Main Page
     =========

5. Write to disk

6. Once whole wiki created in rst, build the sphinx index.
   (after being built once, the site is entirely static)


Considerations
...............

* incremental updates.
  mediawiki knows when an article was updated,
  we could compare this to the timestamp of the file and determine
  if it needs to be written to disk

* if this is written to the samba server, it could be used as-is
  no distribution required, no wearing out writes on all machines.
  (and still easy to retrieve to a USB if we need it in a pinch)

* redirects aren't a thing in rst,
  but we could create the redirects in raw html


Features
........

* searchable
* support for syntaxhighlighting
* fast, offline
* incremental backups

