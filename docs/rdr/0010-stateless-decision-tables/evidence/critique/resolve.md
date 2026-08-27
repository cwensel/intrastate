Model: claude-opus-5[1m]

# critique — resolve, 0010 (stateless decision tables)

Origin ledger: [`diff.md`](diff.md)'s reconciled `D-1..D-18`, built from the two
parallel `--auto` passes ([`critique.md`](critique.md), `claude-opus-5`, 13 rows;
[`critique-modelB.md`](critique-modelB.md), `claude-sonnet-5`, 8 rows). The `C-N`
ids in the two input files are per-file; `D-N` is the stable key and every
disposition below traces to one. No finding was actioned that does not.

Grounding note: the diff pass re-derived pass A's source claims rather than
trusting them, and **corrected four counts** (five `len(Initial)` sites → four;
eleven/seven payload fields → thirteen; nine seed-literal fields → eight; three
"observed dimensions" sites → four). I re-verified each correction against `main`
before editing, and found a fifth: there are **seven** "observed dimensions"
occurrences, of which **five** are the defect and two are correct usage
(§problem-statement's "which row applies to the observed dimensions" and the
mermaid caption's "walks the observed dimensions" both describe resolve matching
on observed tag *values*, which is accurate). Fixing all seven would have been
the wrong unit — the mistake is "coverage is proven over observed dimensions",
not the word "observed".

## Dispositions

| ID | Disposition | Ledger entry | Section touched |
|---|---|---|---|
| D-1 | **fixed** | C5 fence carries no `reason`; 0006's set is closed | `C5`, `A14` (new), `§technical-design`, `§metadata` Overrides, Prerequisites |
| D-2 | **fixed** | A4's census was TOML-only; Go literal missed | `A4`, `§testing-strategy` Done clause |
| D-3 | **fixed** | A10's silence mechanism misattributed for two codes | `A10` |
| D-4 | **fixed** | "exactly one place" false against four shipped sites | `§technical-design` item 1, `§authority` cue |
| D-5 | **fixed** | `emit` payload position unspecified | `C4` |
| D-6 | **fixed** | fence lands in an unreachable branch; no escape precedence | `C5`, `A13` |
| D-7 | **fixed** | `emit` string-only vs A6's own evidence | `A6` (Verified → **Pending**) |
| D-8 | **fixed** | Pending assumptions under an unconditional `MUST` | `C5`, `§risks` |
| D-9 | **fixed** | one-directional contract vs three "both ways" passages | `§approach` ×2, `§decision-rationale` scorecard, `ALT2` |
| D-10 | **fixed** | no test for emit on an expanding (`in`-atom) rule | `S3`, `§oracle`, `§trace` (new row 4′) |
| D-11 | **fixed** | licensed-diff rule forbids its own 103 required edits | `§testing-strategy` Done clause |
| D-12 | **dismissed-with-cite** | mitigation parity vs C5's fence | `A9` (reason recorded) |
| D-13 | **fixed** | coverage "over observed dimensions" ×5 incl. inside C5 | `§problem-statement`, `§approach`, `§technical-design` data-flow, `C5`, C5's ⇒ line |
| D-14 | **dismissed-with-cite** | fence message not required to be actionable | — |
| D-15 | **fixed** | four-reader agreement unenforced | `§authority` cue |
| D-16 | **dismissed-with-cite** | "otherwise" idiom doc-only | — (A9's "If wrong" already names and accepts it) |
| D-17 | **dismissed-with-cite** | `--outcome` friction parked as a kata | — (`BR5`, a decided call) |
| D-18 | **fixed** | single-consumer hit policy, no reversibility statement | `A5` |

### The two hardest — how the forks collapsed

**D-1 + D-6 (C5's fence).** These are one defect seen from two sides, and the
tiebreaker-reduction gate collapsed them without escalation because the evidence
is decisive rather than balanced. `internal/graphlint/taxonomy.go` declares the
`reason` set **closed and append-only**, and
`findings_0006_test.go::TestReq80_…` asserts it by sorted literal *and* asserts
every `graph-unprovable-coverage` emission carries a member of it. All three
members are dimension-scoped; a group-level emission has no dimension. So C5's
"not a new code" was true and beside the point — the code is reusable, the
**discriminator is not**. There is no reading on which the fence lands additively.
Rather than soften the fence (which would restore the silent green the RDR exists
to close) or invent a fourth value silently, the append is stated explicitly as
`no-participating-dimension`, booked as **A14**, recorded in `Overrides` as *the
one non-additive clause in this RDR*, and added to Prerequisites as a dependency
on 0006's assent. Site and precedence are pinned in the same edit: the arm goes
inside `checkCoverage`'s `len(dims) == 0` branch **ahead of** `emitCoverageArms`
(everything A13 previously pointed at sits below that branch's `return`, and
`emitStructurallyUnprovable` loops `guard.Dimensions`, empty in exactly this
case), and this finding beats `graph-coverage-closed-by-escape` where both apply.

**D-9 (one-directional vs both ways).** C1 is right on the merits and its own
parenthetical says why: a bidirectional check would make every rootless state
machine a load failure, breaking `0006:C18` and making this RDR's own MVV step 6
control unwritable. So the three narrative passages are the defect, not the
contract — including the scorecard cell the Decision Rationale says approach C
"wins on". Fixing the cell rather than the contract is the only option that
leaves the recorded decision standing on a true statement.

## Dismissed, with reasons

- **D-12** — mitigation parity between the N-escape-rows consequence (doc) and
  the match-atom mistake (a `MUST` fence). The two are not the same shape at
  lint, which is the whole basis for treating them differently:
  `coverage.go::bareEscapeFor` is scoped per group and per rescued class, so an
  unrescued outcome finds no rescue and its group **already reports**
  `graph-coverage-gap`. The match-discriminated table reports *nothing* — a
  silent green. Lint names the symptom here and the doc names the remedy; a
  dedicated "outcome has no rescue row" code would be a 0006 taxonomy addition
  for a case 0006 already reports. Reason written into `A9` so the asymmetry is
  no longer unexplained.
- **D-14** — that `graph-unprovable-coverage`'s message must name the remedy
  ("author discriminators as guard atoms"), not just the symptom. Real as UX,
  but message text is 0006's surface, and no contract in this repo requires
  actionable message text. Charting it as a message-quality change to 0006 would
  be net-new scope; C5 already forces the finding to *fire*, which is the part
  this RDR owns. The authoring doc (Phase 4) carries the remedy.
- **D-16** — the "otherwise" idiom is Phase-4-doc-only. `A9`'s own "If wrong"
  arm already names this exact outcome and accepts it; re-raising it is a
  re-litigation of a decided call, not a new defect.
- **D-17** — `--outcome` stays required for a class with nothing to select from.
  `BR5` decided this at triage ("applies to both classes; parked as a plain
  kata"), so it is a decided scope call. Noted, not reopened.

## Compute, don't argue

- The `len(Initial) == 0` site count was **executed**, not reasoned:
  `grep -rn "Initial) == 0" internal/ | grep -v _test` returns exactly four —
  `reach.go:91`, `analysis.go:113`, `:233`, `:515`. Pass A's "five" counted a
  `range` as a test; the RDR now states four and names which two become
  class-aware and which two must not.
- `resolvePayload`'s field count was read off the struct: **13** fields
  (`flow_resolve.go:30-45`), not the eleven pass A asserted.
- The "observed dimensions" sweep was run over the file (seven occurrences), not
  estimated, and each was judged individually — five fixed, two correct as-is.
- Cross-lens check against the shared draft: cove pinned C5's terminal-arm
  scoping, 3amigo added C5's zero-dimension fence, and this pass amends C5 a
  third time (reason, site, precedence) plus its ¶1 dimension wording. No field
  is pinned two ways — the terminal-arm scoping and the root predicate are
  untouched here, and the fence I amended is the one 3amigo added, made
  implementable rather than redefined.

## Amendment sweep (§amendment-sweep)

- The `no-participating-dimension` append swept every "no taxonomy change" /
  "every code already exists" claim: one site (`§technical-design` item 2),
  corrected, plus `Overrides` and Prerequisites gaining the dependency.
- The one-directional correction swept all four "both ways" sites (§approach ×2,
  scorecard, ALT2) — grep-confirmed zero survivors of "both ways"/"either
  direction"/"both directions" in that sense.
- The guard-dimension correction swept seven "observed dimensions" sites; five
  amended, two left as correct usage with the distinction stated.
- The licensed-diff rule gained two shapes, and `A4` gained the Go-side census
  it had missed — the two are one claim stated in two places, so both moved.

## Needs (re)verification — carried to Stage 6

- **A14 (new, Pending)** — appending `no-participating-dimension` to 0006's
  closed `reason` set without breaking `TestReq80`, **and 0006's assent to the
  append**. Method: MVV Test. This is the one clause of this RDR that is not
  additive on its owner's grammar; if refused it is a route-back 0006 owns.
- **A13 (Pending, restated)** — unchanged in status, but its verification target
  was *wrong* and is now corrected: the previous target
  (`emitStructurallyUnprovable`'s per-dimension shape) is unreachable in the
  zero-dimension case, so the question as posed was unanswerable. It now asks
  whether emission from inside the `len(dims) == 0` branch, ahead of
  `emitCoverageArms`, yields the finding without double-reporting.
- **A6 (Verified → Pending)** — the supporting reading was refuted: a flat string
  does **not** carry `Next: /rdr-prelock 0046 critique` whole, because `0046` is
  caller-supplied and unknowable at authoring time. The claim is narrowed to the
  checkable one (every value the consumer needs is a fixed string known at
  authoring time; per-invocation data comes from the caller) and re-verified at
  Stage 6 against rdr#tmxk's actual emit keys.
- **A12 (Pending, unchanged)** — carried from cove.
- **A10 (Verified, unchanged in status)** — its evidence was corrected, not
  weakened: two codes are silent by the `len(Initial)` early return, not by the
  reasons A10 gave. The fourteen-code partition it asserts is untouched, and the
  correction ships with the source read, so it stays Verified.
- No other previously Verified assumption was invalidated.

## Charted to successor

None. Every actioned finding landed inside this RDR's existing contract surface.
The two candidates for net-new scope were both dismissed rather than absorbed:
D-14's message-actionability (0006's message surface) and D-12's proposed
"unrescued outcome" code (a 0006 taxonomy addition for a case 0006 already
reports). The `reason` append (A14) is *not* charted — it is a declared
dependency of this RDR's own contract, recorded in `Overrides` and Prerequisites
where implementation will meet it.

## Mini-checks

Not re-read: the cue read was discharged by the first lens pass of this RDR
(cove), whose fired tables — `Authority`, `Oracle`, `Trace` — persist in the
draft and were inherited here. This pass's fixes added no new cue; they extended
the existing tables (Oracle gained the expanded-row control, Trace gained row 4′
and an amended row 2″, Authority's cue gained the enforcement note).
**Mini-checks: none fired (inherited from cove).**

## Convergence

One iteration. All 18 reconciled entries dispositioned: 13 fixed, 4
dismissed-with-cite, 1 (D-15) folded into a fix. No entry re-opened by another's
fix. The edits are pins, corrections, and one booked dependency against a stable
frame — not a rewrite — so no re-run is owed under the substantial-fix rule. The
one structural regression the edits did introduce (a numbered list inside Testing
Strategy minting spurious `S` elements, lint findings 1 → 17) was caught by the
mechanical gate and fixed to prose; `rdr lint 0010` is back to its single
pre-existing `gate:inline` finding, which is Stage 7's to clear.

Pass health, per the lens's Expected signal: **both healthy** — concrete named
passages, real projector anchors, origins spanning §1/§2/§3/premortem/AT. Pass A
carried a counting-discipline problem (four wrong counts, all corrected at the
barrier before any reached the record); pass B made almost no independent source
claims but contributed the authoring/enforcement axis A missed. Independent
convergence on D-7, D-8, D-12, D-18 without shared source access is the
strongest signal in the set.
