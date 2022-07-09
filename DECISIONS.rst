Decisions
=========

* don't raise on fail.
  mediawiki's own parser uses a best-attempt-then-leave-it-to-browser
  strategy. In an outage, I'd rather have a partial backup than none.
  We can communicate failures with the process exitcode.

* html renderer uses `pandoc --standalone`.
  This embeds CSS into each HTML page, which is not efficient (duplication).
  Unfortunately, `--standalone` only renders as much CSS as is required,
  so to extract this we'd need to have a dummy page containing all possible HTML.

  technically, you can get the raw, templated stylesheet with
  `pandoc --print-default-data-file templates/styles.html`
  It may be worth rendering a dummy file with all possible syntax, and extracting the style.
  (there doesn't seem to be an option to render this template from cli)
