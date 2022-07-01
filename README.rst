
mediawiki-to-sphinxdoc
======================

Converts a mediawiki xml dump to a sphinx-doc static-html website.


Features
--------

* searchable
* support for syntaxhighlighting
* fast, offline
* incremental backups


Requires
--------

* pandoc


Usage
-----

.. code-block:: bash

   go install willpittman.net/x/mediawiki-to-sphinxdoc@latest
   php ${your_wiki}/maintenance/dumpBackup.php --full --quiet > dump.xml
   mediawiki-to-sphinxdoc -i dump.xml -o out/
