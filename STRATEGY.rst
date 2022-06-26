
Parsers
=======

Winning Solution
----------------

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

