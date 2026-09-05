# Charted — repeatability lens, RDR 0019

Findings real but out of scope for this RDR, recorded here and dismissed
from the lens loop.

- **D8 — `--flow` operand type and `--as` admitted value set.** The three
  runs produced three synopses (`--flow <id>` vs `--flow <path>`; `--as
  json` vs `--as json|text`). 0019 does not own either: `--flow`'s operand
  and the `--as` value set are 0005:C1's shared selection/envelope surface,
  which 0019 extends by one verb without redefining. Out of scope because
  fixing it here would state a peer's contract in a second place — the
  drift 0019 is elsewhere single-sourcing away. Suggested successor: a
  0005 amendment or a doc pass, if the shared surface is genuinely
  under-stated there.

- **G3 — internal helper names and signatures.** Three runs, nine helper
  names, zero overlap (`allBoundStoresEmpty`/`isStoreEmpty`/`storeIsEmpty`,
  etc.). Not a defect: private helper names are the implementer's, and an
  RDR that fixed them would over-specify. No successor.

- **G5 — Go/CLI framework identity (cobra).** Two runs guessed cobra and
  both flagged the guess. Environmental, not contractual; the RDR names
  the shipped verbs it sits beside and that is enough. No successor.
