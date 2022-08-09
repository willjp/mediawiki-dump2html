
mediawiki-dump2html
===================

Automation around `pandoc` to to converts a mediawiki XML dump to a static html site.

Because you should never not have access to your notes.



Features
--------

* much faster than wget, and similar approaches
* incremental backups (only replace missing, outdated page revisions)
* syntaxhighlighting preserved



Requires
--------

* go_
* pandoc_

.. _go: https://go.dev/
.. _pandoc: https://github.com/jgm/pandoc



Install
-------

.. code-block:: bash

    go install github.com/willjp/mediawiki-dump2html@latest

    # append $GOBIN to your $PATH
    PATH=${PATH}:${GOBIN:=~/go/bin}



Usage
-----

.. code-block:: bash

   # dump your wiki
   php ${your_wiki}/maintenance/dumpBackup.php --current --quiet > dump.xml

   # generate statichtml
   mediawiki-dump2html -i dump.xml -o out/

