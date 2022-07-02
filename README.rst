
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

First, create an xml dump of your wiki

.. code-block:: bash

   php ${your_wiki}/maintenance/dumpBackup.php --full --quiet > dump.xml


Now convert to statichtml

.. code-block:: bash

   go install willpittman.net/x/mediawiki-to-sphinxdoc@latest
   php ${your_wiki}/maintenance/dumpBackup.php --full --quiet > dump.xml
   mediawiki-to-sphinxdoc -i dump.xml -o out/
