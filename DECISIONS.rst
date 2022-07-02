Decisions
=========

* don't raise on fail.
  mediawiki's own parser uses a best-attempt-then-leave-it-to-browser
  strategy. In an outage, I'd rather have a partial backup than none.
  We can communicate failures with the process exitcode.
