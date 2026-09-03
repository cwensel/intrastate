# Charted — repeatability lens, RDR 0026

Findings real but out of scope for this record; dismissed from this loop citing where
they landed.

- **G-4 / D-7a — the binding methods' full signatures (`Reader.Read`, `Gate.Gate`,
  `(*Writer).Apply`).** All three runs could answer the generation prompt's question 1
  only partially, because A1 gives receiver + method names and line numbers but no
  parameter or return lists. Out of scope: this record does not change those signatures,
  and printing them would be documentation of an unchanged surface, not a contract this
  record locks. run-1's `Writer.Write` guess is simply wrong against A1's `(*Writer).Apply`
  (grounded on main) — a widening failure in that run, not an RDR silence.
  Successor: whichever record next changes a binding method's surface prints it there.

- **G-5 — the `Detail` string's exact separators and the `exited N` / `killed by signal N`
  spelling.** C1 `refusal:` fixes the three parts' ORDER (applied-sense, held-pipe reason,
  stderr tail) and the pipe-name ordering, and now names the fields the exit status is
  composed from. The literal punctuation between the parts is left to implementation on
  purpose: 0025:C4 owns the `Detail` composition idiom, and pinning separators here would
  fork it. Not a determinacy gap in this record's sense — no step order or field owner
  turns on it. Successor: none owed; if a golden ever pins the exact string, it belongs
  with 0025:C4.
