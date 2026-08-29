model: claude-opus-5[1m]

# Repeatability diff — RDR 0023 (Resolve-envelope projection opt-out)

Runs compared: `run-1.md` (claude-opus-5[1m]), `run-2.md` (claude-sonnet-5),
`run-3.md` (claude-fable-5). All three declare `variant: full (profile:
foundational)`; run-1's header reads `full`, so three runs is the intended
coverage, not a mismatch.

RDR read via `rdr inspect` only: `0023:C1`, `0023:C2`, `0023:S1`–`S5`,
`0023:MVV`, `0023:A2`, `0023:A9`, `0023:D-identity`, `0023:D-naming`,
`0023:D-wire-byte-format`, `§approach`, `§source-authority-census`,
`§disposition`, `§oracle-discriminability`, `§illustrative-code`,
`§prerequisites`, `§phase-1-projection-on-the-verb`, `§phase-2-oracles`,
`§phase-3-docs-and-help`.

---

## 1. Disagreements

### D1 — Where the flag is read, and whether one function assembles-and-projects
**3-way split.** Lands on `0023:C1` (report-only clause) and
`§phase-1-projection-on-the-verb`.

| Run | Rendering |
| --- | --- |
| run-1 (opus) | Flag read at the TOP of `runFlowResolve` (pseudo-code line 2, `planOnly = cmd.Flags().GetBool("plan-only")`), stored, then consumed at the projection site far below. Prose simultaneously asserts "read once, at the projection site". |
| run-2 (sonnet) | Flag read at the top, then **passed as a parameter into the assembly function**: `buildResolvePayload(result, req, planOnly)` — assembly and projection are fused into one helper, and the projection is inlined field-by-field (`payload.Model = nil`, …) rather than a distinct call. |
| run-3 (fable) | Flag read **at the projection site itself**, inline in the `if` statement: `if planOnly := cmd.Flags().GetBool("plan-only"); planOnly:`. |

Divergence source: Phase 1 says "The flag is read once, at the projection site,
and nowhere in gate evaluation, rule selection, or refusal construction." The
clause conflates two different properties — *lexical read location* and *which
code paths may consult it* — and names only the negative set. Reading it at the
top of `runFlowResolve` and threading a bool satisfies the negative list while
violating the literal "at the projection site"; run-2's fusion into
`buildResolvePayload` makes the flag an assembly-function *parameter*, which is
the reading most at risk of the branch migrating into assembly.

run-2 is the outlier that matters: fusing assembly and projection means there is
no standalone projection function to implement "FROM the declaration table"
(C2), and the record's own §Phase 1 site rule ("applied after payload assembly,
before the gateway") is no longer structurally checkable.

**Candidate rewrite (C1 or §phase-1):** state the read as one clause with both
halves — the flag is read at exactly one lexical site, that site is the
projection branch on the success path, and it is never passed as a parameter
into payload assembly.

---

### D2 — The Go type of `observed` and `owned`
**3-way split.** Lands on `0023:§source-authority-census` and `0023:C2`.

| Run | `observed` | `owned` |
| --- | --- | --- |
| run-1 | `*map[string]string` (marked GUESS) | `*map[string]string` (GUESS) |
| run-2 | `*observedView` — a **named struct type**, invented | `*map[string]any` |
| run-3 | `*map[string]string` (marked GUESS) | `*map[string]string` (GUESS) |

Divergence source: the census gives writers (`observedTagMap(req.observed)`,
`tagMap(owned)`) and JSON keys, and never a Go type. A2 fixes only
"pointer-valued echo fields with `omitempty`". run-2 escalated the same silence
into a *new named type* and a `map[string]any` element type, which is a
materially different wire shape from `map[string]string` — `any` values
serialise nested structures the two `string`-typed runs cannot produce.

This matters beyond naming: the census's `tagMap` / `observedTagMap` writers are
the same two functions A9 names as the nil-container hazards, and A9's Phase 1
oracle ("`{}`/`[]`, never `null`") is written against whatever their value type
is.

**Candidate rewrite (`§source-authority-census`):** add a Go-type column to the
14-field table, or state in A2/C2 that the conversion is `T` → `*T` with `T`
unchanged from `main`.

---

### D3 — Is the reader pass upstream of, or interleaved with, the last refusal?
**3-way split.** Lands on `0023:C1` (report-only / "skips work whose only
consumer is a projected-away field") and `0023:S1`.

| Run | Rendering |
| --- | --- |
| run-1 | Reader pass at step 7, **before** `resolve.Resolve` and before gate evaluation; its own `respond.Fail` arm at step 8. Three further `respond.Fail` returns follow it. |
| run-2 | Reader pass (`invokeReaders`) placed **after** `kernel.Resolve` and after the only refusal return — i.e. entirely on the success path. |
| run-3 | Reader pass **before** `resolve.Resolve`, but with no refusal arm of its own; `owned` assembled after it. |

Divergence source: the record fixes only (a) that the projection sits after the
last `respond.Fail`, and (b) that the reader pass must actually execute. It
never orders the reader pass against the refusal returns. This is consequential
for S1: if the reader pass is on the success path only (run-2's placement), the
"observed reader execution ± flag" assertion cannot discriminate a refusing run,
and C1's "skips work whose only consumer is a projected-away field" breach has a
narrower reachable surface than run-1's and run-3's layouts imply. `owned` is
also written by the reader pass per the census ("reader pass over the caller's
own request"), which run-2's ordering contradicts — it builds `ownedView`
independently of `invokeReaders`.

All three flagged this ordering as a GUESS in some form (run-1 G10, run-3's
closing note, run-2 implicitly), so the silence is acknowledged by every run.

**Candidate rewrite (`§source-authority-census` or C1):** name the reader pass's
position relative to the refusal returns, and state whether `owned` and
`readers` come from the same pass.

---

### D4 — The declared assignment table: name, type, and file
**3-way split.** Lands on `0023:C2` (carrier clause) and `0023:S3`.

| Run | Table |
| --- | --- |
| run-1 | `resolvePayloadSides`, `map[string]payloadSide` with `sideEcho`/`sidePlan` consts, in a new file `internal/cli/flow_partition.go` |
| run-2 | Not named as a symbol at all — described only as "a standalone table keyed by field name, in its own declaration site" in the data-model prose. No Go declaration is rendered. |
| run-3 | `resolveFieldGroups`, `map[string]fieldGroup` "(or ordered table)", home unspecified |

Divergence source: C2 constrains the carrier's *shape* ("standalone table keyed
by field name") and *site* ("its own declaration site", NOT a struct marker) but
names no symbol and no file. Note C2 is internally in tension here: an earlier
paragraph offers the carrier as "a per-field marker on `resolvePayload` **or** a
table keyed by field name", and a later paragraph forbids the marker outright
("It MUST be a standalone table … NOT a per-field marker"). run-2's failure to
produce any declaration is plausibly downstream of that unresolved either/or.

run-3 additionally hedges the *shape* — "(or ordered table)" — which C2 does not
license; ordering is not a property the table needs, and an ordered table is a
different structure from a map keyed by field name.

**Candidate rewrite (C2):** delete the superseded "a per-field marker … or a
table" phrasing so the carrier is stated once, in its constrained form.

---

### D5 — The total command-tree walker's signature and return
**3-way split.** Lands on `0023:C1` (WHOLE-means-total paragraph) and
`0023:S4`.

| Run | Walker |
| --- | --- |
| run-1 | `walkCommandTreeTotal(root *cobra.Command) []*cobra.Command` — returns every command; the caller filters for registrants |
| run-2 | `walkCommandTree(root *cobra.Command) map[string]bool` — **returns the set of commands that register `plan-only`**, i.e. the filter lives inside the walker. Reuses the *name* of the shipped partial walker it must not reuse. |
| run-3 | `walkAllCommands(root *cobra.Command) []*cobra.Command`, explicitly homed in the S4 test file, not shipped code |

Divergence source: C1 and S4 fix the walker's *semantics* exhaustively (no
name-skip, no `Hidden` gate, `help` and `completion` materialized via cobra's
own initializers and asserted present) and name neither symbol nor signature.
run-2's shape is the risky one twice over: naming it `walkCommandTree` collides
with `internal/cli/help_all.go::walkCommandTree`, which C1 explicitly forbids
reusing, and folding the `plan-only` filter into the walker means the "walked
set CONTAINS `completion`" precondition S4 mandates cannot be asserted against
the walker's own return value — a `map[string]bool` of registrants will never
contain `completion` in a passing build.

run-1 and run-3 both return the full command list, which does support the
precondition assertion.

**Candidate rewrite (S4):** state that the walker returns the walked command
set, and that the `help`/`completion` presence precondition is asserted against
that return.

---

### D6 — Where the projection is applied, relative to `respond.Fail`
**2-way split, and it splits along a model boundary: run-2 (sonnet) vs.
run-1/run-3.** Lands on `0023:C1` and `§phase-1-projection-on-the-verb`.

- run-1: pseudo-code carries an explicit `# ---- LAST respond.Fail is above this
  line. Success path only below. ----` marker, plus the record's prohibition on
  a defensive `if refusing, skip projection` guard.
- run-3: same — "last Fail return — projection is unreachable above", and
  "success path only, before gateway".
- run-2: the projection is a `if planOnly:` block inside `buildResolvePayload`'s
  described responsibility, reached via a single refusal return placed early.
  run-2 never states the *structural* placement rule (after the last
  `respond.Fail`) and never mentions the forbidden defensive guard, rendering
  A6's flag-blindness as a property of the refusal payload's shape ("carries no
  echo-group member") rather than of the site.

Divergence source: the structural placement rule and the guard prohibition live
only in `§phase-1-projection-on-the-verb`, not in C1's normative text — C1 says
only "applied to the verb-specific result BEFORE `respond.OK`, never in the
respond gateway and never per output mode". A reader who widens into `§approach`
(which repeats the weaker "before `respond.OK`" form) but not into Phase 1 gets
run-2's reading, which is exactly what happened: run-2's widening list names
`§approach` and omits `§phase-1`.

**Candidate rewrite (C1):** lift the "after the last `respond.Fail` return, and
no defensive refusal guard" sentence out of Phase 1 into C1's report-only
clause. It is a normative structural requirement, not an implementation-plan
detail.

---

### D7 — `escape_class` Go type and presence mechanism
**2-way split: run-2 (sonnet) vs. run-1/run-3** — but note run-1 and run-3
agree only on the mechanism, not on the classification. Lands on `0023:C1`
(presence-rule paragraph) and `§source-authority-census`.

- run-1: `string` with `,omitempty` (explicitly flagged G7 — "the encoding
  mechanism is not stated").
- run-3: `omitempty` per its producer; "the 14th field / conditional 14th key".
- run-2: `*string` with `json:"escape_class,omitempty"` — a **pointer**, putting
  `escape_class` mechanically in the same class as the five echo fields.

Divergence source: C1 states the presence *rule* (`0005:A-3`, omitted when
unescaped and when unprobeable) and never the encoding. This matters because a
`*string` `escape_class` is indistinguishable at the type level from a projected
echo field, and C2 puts `escape_class` in the always-keep core. A reflective
oracle (S3) or a future second projection mode that keys off "pointer +
omitempty ⇒ projectable" would silently sweep an always-keep field.

**Candidate rewrite (`§source-authority-census`):** state that only the five
echo fields become pointer-valued, and that `escape_class`'s presence rule is
realised by a mechanism distinct from the projection mechanism.

---

### D8 — Which envelope/measurand facts are recovered at all
**2-way split: run-2 (sonnet) vs. run-1/run-3.** Lands on `0023:C1`
(STRICTLY SHORTER measurand) and `0023:S5`.

- run-1 and run-3 both reconstruct the width measurand precisely: the FULL
  EMITTED LINE including `{"type":"ok","data":{…}}`, trailing newline excluded,
  asserted in BOTH modes, with the text half resting on Pending A8 and narrowing
  to JSON before lock if refuted. Both also recover the text-subset rule as *set
  membership, not subsequence*, because wire order is declaration order while
  text is `sort.Strings`-flattened.
- run-2 recovers neither. It states "Projected text lines are a byte-identical
  subset (per-line) of the default-mode text lines" and never mentions the
  strict-width clause, the envelope-inclusive unit, A8's Pending status, or the
  membership-vs-subsequence distinction.

Divergence source: not RDR silence — C1 is emphatic on all of it. run-2 read
C1's normative block and dropped its longest, most load-bearing paragraph. This
is a *reconstruction* gap rather than a *record* gap, and it is the strongest
model-boundary signal in the set: the alt-model run lost the one clause the
record spends the most words fixing.

**No rewrite indicated** — recorded because a diff should not read run-2's
silence as the record's silence. If anything it argues C1 is too long to survive
a single read: the width measurand could move to its own contract id.

---

### D9 — `gates` element type
**2-way split: run-1 vs. run-2/run-3** (weak). Lands on
`§source-authority-census`.

- run-1: `[]…` with an explicit GUESS (G8) that the element type is never given.
- run-2: `[]gateResult` — a concrete invented type name.
- run-3: `Gates: gateResults` in the assembly literal, type unstated.

Divergence source: the census names `req.runGates` as the writer and no type.
Same root cause as D2; folded here so it is not double-counted.

---

## 2. GUESS clusters

Contracts two or more runs independently marked GUESS. These are the candidate
rewrites, strongest first.

### GC1 — The projection function's symbol (3/3 GUESS)
run-1 G1 (`projectPlanOnly`, value receiver, non-pointer return), run-3
("GUESS at name/signature: `projectPlanOnly(p resolvePayload) resolvePayload`"),
run-2 ("name is a GUESS — the record does not name the assembly function").
Lands on **`0023:C2`**. C2 requires the projection be "implemented FROM that
declaration" but names no function, so run-2 was free to dissolve it into an
assembly helper (D1). Naming the function — or stating that the projection is a
distinct function taking the assembled payload — closes D1 and GC1 together.

### GC2 — Concrete Go types for the 14 payload fields (3/3 GUESS)
run-1 G6 and G8, run-2 ("GUESS on exact field naming/tags beyond what C1/C2/A2
state"), run-3 ("GUESS at exact Go types"). Lands on
**`§source-authority-census`**. Root cause of D2, D7, and D9. The census already
has the table; a Go-type column is the minimal fix.

### GC3 — The declaration table's symbol and home (3/3 GUESS)
run-1 G2/G3, run-3 ("GUESS at name/shape … `resolveFieldGroups`"), run-2
(renders no symbol, describing the constraint only). Lands on **`0023:C2`**.
See D4 — compounded by C2's unresolved "marker or table" phrasing.

### GC4 — The total walker's symbol and signature (3/3 GUESS)
run-1 G9, run-2 (renders it as new but reuses the forbidden name), run-3 ("GUESS
at name/home"). Lands on **`0023:S4`**. See D5. run-3 additionally guesses the
*home* (test file, not shipped code) — a placement no run disputes but no
element states.

### GC5 — Runtime step order upstream of the projection (3/3 GUESS)
run-1 G10 ("the record fixes only *where the projection sits*"), run-3's closing
parenthetical ("GUESS: exact assembly order of readers/owned/kernel-call"),
run-2's ordering differs without being flagged but its §1 header disclaims
"the record fixes behavior, not Go syntax". Lands on
**`§source-authority-census`**. See D3.

### GC6 — The reader-execution recording seam (2/3 GUESS)
run-2 names it as helper #3 with an explicit GUESS
(`recordReaderInvocation` "or a counting seam around `flow_exec.go`'s reader
pass"); run-1 and run-3 both describe the *requirement* (observed execution at
the invocation site, never recomputed, never read from the payload) without
naming a seam, and neither lists it among their three most important helpers.
Lands on **`0023:S1`**. The record specifies the measurand exhaustively and the
*mechanism* not at all ("a counting or recording seam around the reader pass" is
C1's only gesture). Notable inversion: run-2, the run that lost the most C1
detail elsewhere (D8), is the only one that surfaced this as a first-class
helper.

### GC7 — The flag's landing field (1/3, recorded for completeness)
Only run-1 flags it (G5, `req.planOnly`); run-2 threads a bare local through a
parameter; run-3 never stores it. Single-run, so not a cluster — but it is the
same silence D1 turns on.

---

## 3. Agreement

All three runs rendered identically, with no divergence worth a finding: the
CLI surface (`--plan-only`, boolean, default `false`, long form only, no
shorthand, registered on `newFlowResolveCmd` only and never on
`registerSelectionFlags` / the flow group's persistent set / the root's
persistent set); the disposition table (exit 2 + `command-error` for the flag on
a sibling verb and on a pre-RDR binary; refusals byte-identical and flag-blind;
no new refusal code or exit group minted); the ECHO membership
(`model`, `observed`, `owned`, `readers`, `outcome`) and PLAN membership
(`revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped`,
`escape_class`); A2's mechanism (pointer echo fields with `omitempty`, populated
in default mode, nilled at projection, keys ABSENT never `null`/`{}`/`""`); the
14-field `NumField()` pin and that the conversion changes types only; C2's
carrier-independence rationale (a struct marker is edited by the same commit
that retypes the fields, so S3 would assert the implementation agrees with
itself); `emit` present as `{}` per `0010:C4`; `escape_class` compared against
the same run's default output rather than an unconditional literal; `revision`
as a plan-side, contracted-but-vacant identity slot; that nothing new is
persisted and the boundary is the NDJSON stdout record; and that both output
modes render one projected payload through the shipped `respond/text.go::flatten`
so `0005:C1` holds by construction.
