Model: claude-opus-5[1m]

# Charted to successor — critique, RDR 0024

Findings real but net-new scope for this record. Recorded here and dismissed
from this loop; none was folded into the draft. The first three RE-RAISE
classes already charted by the 3amigo pass (`../3amigo/Charted.md`) — they are
listed again because two independent critique passes reached them, which is
evidence the successor is worth scheduling, not evidence this record should
absorb them.

- **M-09 / M-13 (A C-9, C-13; B C-4, C-8) — merge-day coverage and the
  consumer seam test.** Both passes independently reached 3amigo's PM-3 and
  PM-7: this RDR ships capability, not coverage; the seed defect's headline
  failure (routing to a nonexistent command) is closed only by a cross-file
  seam test that lives in the sibling engine repo and is unscheduled.
  Unchanged disposition — intrastate has no command surface to check against,
  so the test cannot be written here. **Successor**: the consumer-repo record
  3amigo's Charted.md already names, now with two-pass corroboration.

- **M-10 remedy — multi-defect load reporting.** Accumulating findings at the
  load tier needs the total order `internal/table/load.go` deliberately
  withholds. Out of scope by C2's own reasoning. The finding's *residue* —
  that this is the first load category to fire N times over values the author
  did not edit — was NOT charted: it is folded into Consequences as one
  sentence, because it is a property of this record's own change.

- **M-11 (A C-11) — scalar-hollowing is unobservable from output.** A fully
  `scalar`-declared model is output-indistinguishable from a proven one, and
  Meyer's cited "flag the loopholes" condition is declined on the ground that
  `0006:C17` closes the advisory tier. The critique is correct that this is a
  procedural rather than substantive answer. It stays declined here — the
  finding class it needs does not exist and creating one reopens a locked Final
  contract. **Successor**: an RDR that reopens `0006:C17`'s advisory tier, or a
  non-finding surface (a `dump` column, a lint `--stats` mode) that reports
  declared-vs-scalar coverage without minting an advisory finding.

- **M-21 (B C-7) — the `dispositions` plural/shape mismatch.** The field is
  plural but maps each key to exactly one token. Real, cosmetic, and the name
  is already load-bearing in `JDR 0002 §D1`'s PLAN assignment and in 0023's
  draft. Renaming now costs a peer amendment for no behavioral gain.
  **Successor**: none scheduled; recorded so a future reader knows the shape
  was seen and kept deliberately.

- **M-14 / M-18 — accepted costs, not defects.** `int`'s non-canonical
  literals (`03`, `+5`, `-0`) are inherited from RDR 0003 by the A4 reuse and
  no shipped or motivating model emits an int; re-deriving them locally to get
  stricter parsing would fork the kind rules 0003 owns, which is the worse
  trade. The QOC spread (22 vs 18) rests on a correctness row that is
  conditional on the seam test — already priced by M-13's charted successor —
  and option B remains rejected on the two proven seed defects it cannot catch
  regardless of scoring.
