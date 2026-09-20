Model: claude-opus-5[1m]

# 3amigo consolidation — cli/0030

Three isolated persona passes (PM, Implementer, QA), each run in a fresh
sub-agent context that saw only the RDR and its own persona block. All three
stamped `claude-fable-5-1`. Files: `persona-1-pm.md`,
`persona-2-implementer.md`, `persona-3-qa.md`.

## Hotspots (ids raised by >=2 isolated personas)

Overlap marks a hotspot PASSAGE, not a validated finding; a single-persona
finding is not thereby weaker.

| id | personas | what converged |
| --- | --- | --- |
| `0030:C2` | PM, Impl, QA (3) | the bound refusal: what its detail carries, the enum value slot, the remedy sentence, domain-vs-suffix order |
| `0030:C3` | PM, QA (3 refs) | the authoring guide's scope; row-identity publication |
| `0030:C1` | Impl, QA (2) | "exactly two atoms" vs retained bounds; subsumption untested; always-append suffix untested; duplicate/`#` refusal scope |
| `0030:D-identity` | Impl, QA (2) | the step point's sort tuple is unstated (`compareAtoms` takes an Atom; a step point is not one) |
| `0030:MVV` | PM, QA (2) | validates transparency, not the legibility/row-count outcome; vacuous-pass risk on S1's set equality |
| `0030:S3` | Impl, QA (2) | the enum sibling's value slot and its assumed `tier` domain |

## Merged findings

### High (QA)

- **`0030:C1` subsumption + `0030:§testing-strategy`** — the subsumption rule
  (one choice point per stepped tag; a local OR inherited match `in` rewritten
  per row to `match eq = <cell>`, never retained) has its only test in
  `0030:A10`'s prose. Absent from S1-S7 and the MVV. HOTSPOT `C1`.
  Blocks: the gating scenarios for local-`in`, inherited-`in`, and the
  one-suffix-element assertion.
- **`0030:C1` always-append + `0030:D-identity` + `0030:C3`** — C1 fences a
  code fork on the shared `suffixed` flag (a stepped row always carries its
  cell, even at exactly one admitted cell) and C3 requires every row-naming
  surface to publish `rule#cell`. No scenario asserts either; no fixture has a
  step rule admitting exactly one cell. HOTSPOTS `C1`, `D-identity`, `C3`.
  Blocks: single-cell suffix publication, multi-cell identity publication, and
  the single-member-`in` negative control.
- **`0030:C2` + `0030:S3` + `0030:F3`** — "the FIRST failing cell in the tag's
  own domain order, not the row's suffix order" and "in either direction" have
  no scenario. S3 has one failing cell; no scenario uses a negative step. The
  two orders genuinely diverge (suffixes sort by `slices.Compare` on strings).
  HOTSPOTS `C2`, `S3`. Blocks: multi-failure domain-order assertion (int and
  enum) and both negative-step refusals.
- **`0030:C1` duplicate/`#` member refusals + `0030:S5`** — the duplicate-member
  sentence reads unconditional while the `#` sentence is scoped to a step write;
  `0030:A3` records a `#` member loadable today. HOTSPOT `C1`. Blocks: the
  negative control that an UNSTEPPED `#`/duplicate domain loads as today.

### Medium

- **`0030:C2` remedy + contradiction (PM)** — `§risks-and-mitigations` says the
  detail "names the `guard.all` complement"; C2 owes only rule, tag, cell and
  stepped value; S3 pins "only the cell is new text". The author most likely to
  hit C2 wrote `unless tier = "large"` (not consulted per C1) and is told
  nothing about why. HOTSPOT `C2`. Blocks: what C2's detail is owed.
- **`0030:C2` enum value slot (Impl, QA)** — for an enum step leaving the
  domain there IS no stepped value; S3's enum sibling names only `tier` and
  `large`. S3's `large` also presumes a domain whose last member is `large`,
  which the MVV never declares. HOTSPOTS `C2`, `S3`. Blocks: the C2 detail
  format and the enum half of S3 / MVV row 4.
- **`0030:C1` "exactly two atoms" (Impl)** — C1 says the stepped tag carries
  "exactly two atoms per row", then retains "every OTHER authored atom — the
  `guard.all` bounds"; the `trace` witness (row 4) shows THREE. Reads as a
  countable invariant an implementer would assert and it is false whenever a
  bound is retained. HOTSPOT `C1`. Blocks: per-row atom emission in Phase 2 and
  any S-test asserting the stepped tag's atom set.
- **`0030:D-identity` sort tuple (Impl, QA)** — a step point is not an `Atom`,
  so `compareAtoms`' (key, block, operator, literal) does not apply as written,
  and its `guard.all eq` literal is per-cell (row-dependent). No scenario has
  two step points in one rule, or a step composed with an `in` on a different
  tag. HOTSPOT `D-identity`. Blocks: `retry#3#mid` vs `retry#mid#3`, hence
  `compareRows`, dump order, and every fixture identity.
- **`0030:§problem-statement` enum residual (PM)** — for an `enum` the only
  positive atom that can exclude the terminal member is `in = [<every member
  but the last>]` (ordered operators are int-only; `unless` not consulted;
  guard grammar out of scope). So the enum cap is a re-listing of the domain
  minus one, edited on every domain change — the "say it once" outcome holds
  for the stepped VALUE but not the cap. `§consequences` says only "its own
  positive atom". Blocks: whether the enum arm is claimed to deliver the
  tier-domain-change outcome, or the record discloses the residual.
- **`0030:MVV` outcome criterion (PM, QA)** — the MVV proves the mechanism is
  invisible; it has no criterion for the outcome `§problem-statement` names
  (source legibility / row count). `§background` asserts ~24 rows vs 8/10 and
  `§decision-rationale` row 6 decides on "source is one row"; no MVV item or
  S-scenario asserts `ladder-step.toml` is that shape. HOTSPOT `MVV`. Blocks:
  accepting the MVV as validation of the user's problem.
- **`0030:S1` + `0030:MVV` item 2 vacuous pass (QA)** — "identical finding sets
  … spanning both the blocking codes and the advisory `graph-idempotent-write`"
  has no MUST that the set is non-empty or carries a finding of each tier;
  nothing forces an idempotent write. Empty-equals-empty passes. HOTSPOT `MVV`.
  Blocks: the vacuous-pass guard.
- **`0030:§prerequisites` / 0012 sequencing (Impl)** — RDR 0012 is Draft;
  today's `guard.Evaluator` is `struct{}` with string `eq`, and `0012:C2` would
  make `eq`/`in` typed. Prerequisites name only 0002 as Implemented. Blocks:
  Phase 2's start and the extracted shim's signature (does the loader's admits
  call pass the declared kind).
- **`0030:S2` + `0030:A6` sweep criterion (QA)** — no criterion for cells no row
  claims (where `flow resolve` refuses rather than plans), and "state" for a
  two-tag ladder is not fixed. Blocks: the sweep harness's pass/fail on an
  unmatched cell, and the sweep's domain.
- **`0030:A11` + `0030:C1` both-carriers oracle (QA)** — A11 names the
  discriminating oracle (graph's per-row `next` / dump's next column carrying
  the stepped value) because a `Writes`-only implementation passes every
  `Fingerprint` assertion. Not among S1-S7. HOTSPOT `C1`. Blocks: the
  `Writes`-only regression test.
- **`0030:S5` zero-cell discriminator (QA)** — "the refusal is new rather than
  `0003:C6`'s existing check" is not observable from the envelope (both are
  `malformed_tag_declaration`). Blocks: the paired control (same `lt = 0`
  WITHOUT the step loads today).
- **`0030:C2` + `0030:S3` message-text assertion (QA)** — the cell's rendering
  is unspecified (`cell 9`, `#9`, `at 9`) and whether the structured
  `rule`/`key`/`literal` fields are populated is not stated. HOTSPOTS `C2`,
  `S3`. Blocks: a stable (non-substring) assertion.

### Low

- **`0030:C2` remedy sentence (PM)** — "a positive atom excluding the cell, or
  a literal row for it" reads as two alternatives; under C1 a literal row
  elsewhere does not change THIS rule's admitted cells. The literal row is the
  remedy for the coverage gap that FOLLOWS the exclusion. HOTSPOT `C2`.
- **`0030:C3` guide scope (PM)** — the guide obligation omits the two rules an
  author cannot infer from the form: authored `domain` order IS step order
  (reordering silently changes the model, no lint) and `unless` never excludes
  a cell. HOTSPOT `C3`.
- **`0030:§problem-statement` wording (PM)** — "mean it for every value the tag
  can hold" overstates the delivered "every cell the rule's positive atoms
  admit"; and when the cap equals `max`, the cap lives in two places.
- **`0030:C2` overflow arm (Impl)** — `conformDomain`'s int arm calls
  `strconv.Atoi` on the rendered literal, so a stepped value beyond int64
  reads "is not an int" rather than "is above max". HOTSPOT `C2`.
- **`0030:C1` copied `intWidth`/`IntDomain` (Impl)** — restated in `table`
  while the same cycle is solved for `Evaluator` by a MOVE into `resolve`; the
  decimal cell rendering must byte-match `guard/assignment.go`'s
  `strconv.Itoa`. HOTSPOT `C1`.
- **`0030:C1` step + `clear` on one key (Impl)** — silence. Today
  `renderWrites` lets the clear list overwrite `assignments[key]` silently.
  HOTSPOT `C1`.
- **`0030:§phase-2-expand-into-cells` scope (Impl)** — "two callers repointed"
  understates it: three non-test repoints plus ~20 `guard.Evaluator` test
  references including a field-count reflection test.
- **`0030:S5b` + `0030:A5` open expected set (QA)** — does not say whether other
  advisories may co-fire or whether exactly one redundant-row is required.
- **`0030:S4` atom position (QA)** — "a predicate literal `{ step = 1 }`" does
  not fix the arm (`match eq`, `guard.all eq`, or an `in` member).
- **`0030:S6` + `0030:C3` oracle (QA)** — "decodes with no new member" names no
  mechanism (strict decode? golden enumeration?). HOTSPOT `C3`.
- **`0030:S5` detail wording (QA)** — fixes no wording; the category is the only
  firm assertion.

## Recorded non-findings (so the next reader need not re-check)

- C1's three-valued "undecided excludes" arm is declared unreachable and
  untestable by design; the record says so (QA).
- `[initial]` and predicate literals already refuse a TOML inline table via
  `valueMembers`' default arm under their own categories; `sourceRule.Write` is
  `*map[string]any` so `{ step = n }` decodes without a struct change;
  `flow set-state` exists; `docs/model-authoring.md` has a write-block section;
  `checkIdempotentWrites` cannot fire on a non-zero step (Impl).
- No persona found the record solves a different problem than the one stated
  (PM).
