Model: claude-opus-5

# Critique Diff — RDR 0025, Pass A (opus-5, C-1..C-13) vs Pass B (sonnet-5, C-1..C-6)

Reconciled by passage anchor, not by id. Pass A carries 13 rows, Pass B 6.
Eleven distinct findings survive dedup. Every code-level claim either pass
makes about `main` was read against source; verdicts below carry a
`path::Symbol` anchor.

## 1. Reconciled ledger

| ID | RDR passage | Failure mode | In A | In B | Agreement | Code-verified |
|----|-------------|--------------|------|------|-----------|---------------|
| D-1 | `0025:A12`, `0025:C2` (argv0 resolves against "the model file's directory", recorded by the loader) | The loader has no path to record. `Load(src []byte, sourceID string)` does no file I/O; `sourceID` lives on `loader`, never on `Model`. The two openers are callers, and both pass the raw `--model` string, which is never absolutized. | C-1 | C-3 | **both** | **CONFIRMED** — `internal/table/load.go::Load` (sig + package doc "performs no file I/O and no path resolution (`0002:EIA`)"); `internal/table/load.go::loader` (`sourceID` field); openers `internal/cli/flow_input.go::selectModel`, `internal/cli/lint.go::runLint`; `internal/cli/flow_input.go::selectModelPath` returns `--model` verbatim |
| D-2 | `0025:A11`, `0025:C4` `detail:` line ("`refusalOf` — the constructor every *post-selection* refusal already passes through") | `refusalOf` is not the sole post-selection constructor. `refusalWithKeys` wraps it, and `Executor.Read` mints **every** read refusal through it — including `ClassExecutionFailure` from `invokeRead`. Hooking `refusalOf` alone reaches gate and write but drops the stderr tail on the whole read path. | C-2 | — | A-only | **CONFIRMED** — `internal/accessor/executor.go::refusalWithKeys` (calls `refusalOf`, adds `Keys`); `internal/accessor/executor.go::Executor.Read` uses `refusalWithKeys` at both refusal sites |
| D-3 | `0025:C4` "Rendering is a second hop" ("the applied-sense text wins the `CLIError.Detail` slot and the stderr tail renders beneath it") | `clierr.CLIError` has exactly one `Detail string`; `EmitText` renders exactly one `  detail:` line. There is no "beneath it". No Existing Infrastructure Audit row was opened for `clierr`. | C-3 | — | A-only | **PARTIAL** — one-slot claim **CONFIRMED** (`internal/cli/clierr/clierr.go::CLIError.Detail`, `::EmitText`); A's *corollary* that concatenation "breaks the one-line invariant" is **REFUTED** — the field doc says "May be multi-line", and `runLint` already ships a multi-line `err.Error()` through it (`internal/cli/lint.go::runLint`). The missing slot is real; the one-line invariant is not |
| D-4 | `0025:§approach`, `0025:C5`, `0025:C6`, `0025:F1`, `0025:F2`, `0025:MVV` step 2 — "`intrastate lint` validates the declaration" | `runLint` is `table.LoadWithAdvisories` + `graphlint.Run`. `accessor.Validate` — the function the RDR treats as the live validator — has zero production callers. And `Load` is fail-fast per *document*, so a five-defect model reports one finding per run. | C-4 | — | A-only | **CONFIRMED** — `internal/cli/lint.go::runLint` (no `accessor.Validate` call); `internal/accessor/validate.go::Validate` reached only from `internal/accessor/*_test.go` (14 test sites, 0 production); `internal/table/load.go::Load` doc: "load is fail-fast and returns one categorized error, never a list" |
| D-5 | `0025:C1`, `0025:A9`, `0025:D-selection-predicate` | A command entry yields `Path: ""` (plain `string`, no pointer). `flowbind.Registry` unconditionally builds `Reader{Path: acc.Path}` / `&Writer{...}` / `Gate{...}` with no carrier branch. `flowbind.load("")` gets `ENOENT`, which it deliberately treats as an **empty artifact, not an error** — so a discriminator bug fails open into silent success, not a crash. | C-5 | C-1 (partial: B names the missing branch, not the fail-open) | **both** (mechanism partly divergent) | **CONFIRMED** — `internal/table/load.go::Accessor.Path` is `string`; `internal/cli/flowbind/registry.go::Registry` — three unconditional `Path: acc.Path` loops, no `acc.Command` branch anywhere; `internal/cli/flowbind/flowbind.go::load` returns `store{}, nil` on `os.ErrNotExist`. Verified `os.ReadFile("")` does return an `ErrNotExist`-matching error, so the fail-open path is live |
| D-6 | `0025:S5b`, `0025:A4` ("a role with a path-backed and a command-backed reader, both admissible to today's first-match `readerFor`") | A: the S5b fixture will not load — `checkAccessorBindings` refuses an OWNED tag served by ≠1 reader and an OBSERVED tag served by >1, so two readers on one role are authorable only with disjoint key sets, under which the read-back reader is loader-determined and `readerFor`'s first-match is unobservable. B: `readerFor` is first-match with no ambiguity check and 0016 is Draft, so a migration silently verifies against the stale reader. | C-6 | C-5 | **both** (same passage, **incompatible** mechanism — see §4) | **A CONFIRMED, B PARTIAL** — `internal/table/load.go::checkAccessorBindings` refuses owned `readerCount != 1` and observed `readerCount > 1` per *tag key*; `internal/accessor/model.go::readerFor` is a bare first-match loop over `reg.Definitions`, matching on `Role` only. Both are true; the arity check is keyed on **tags**, `readerFor` selects on **role**, so B's two-reader migration is authorable exactly when the key sets are disjoint — which is the case A says makes S5b's assertion vacuous |
| D-7 | `0025:C2` substitution contract ("the substituted value must be an absolute path — refuses when the caller-bound artifact path is relative") | A user-facing behavioural asymmetry between two carriers of the same accessor, surfaced only at runtime, accepted in one sentence and never reaching Consequences, Failure Modes, or the MVV. Every fixture and doc example in the repo binds a relative `--artifact`. | C-7 | — | A-only | **CONFIRMED** — `internal/cli/flow_input.go::parseArtifacts` stores `out[role] = path` verbatim, no `filepath.Abs`, no absoluteness check. Confirmed absent from `0025:§consequences` and `0025:§failure-modes` |
| D-8 | `0025:C6` "gate site: the command binding's constructor… the flag therefore reaches the registry constructor, which is the one new parameter this clause adds"; `0025:S7` | The registry-constructor framing does not match the call graph. `flowbind.Registry` is a package-level free func with **one** production call site, inside `buildRequest`'s struct literal, shared by every verb. The flag has one place to be read and many places to be *registered*; S7's "all three `NewExecutor` callers refuse identically" tests one code path N times. | C-8 | C-1 | **both** | **CONFIRMED, and both passes undercount** — `internal/cli/flowbind/registry.go::Registry` is `func Registry(m *table.Model) accessor.Registry`; sole production call site `internal/cli/flow_exec.go::buildRequest`. But `buildRequest` has **four** callers across four verbs — `flow_next.go:207`, `flow_resolve.go:227`, `flow_state.go:135`, `flow_state.go:274` — not the three `NewExecutor` sites either pass counted. The registration-drift hazard is strictly larger than described |
| D-9 | `0025:C3`, `0025:S5` ("the inherited reserved-literal rule composes unchanged… `internal/accessor/executor.go:96`") | The line anchor is stale and the rule does not compose on the read-back path. The `<clear>`-is-unreadable rule is parameterized by `clearIsUnreadable`, which `Executor.Read` passes `true` and the read-back passes `false` — by deliberate design. | C-9 | — | A-only | **CONFIRMED** — `internal/accessor/executor.go:96` is a comment fragment ("establish the key's content: the binding cannot tell a removal from"); the rule is in `internal/accessor/executor.go::readOutcome.classify` gated on `clearIsUnreadable`; `Executor.Read` passes `true`, `Executor.Write`'s read-back passes `false` with the in-code rationale "The `<clear>`-literal rule is NOT applied here" |
| D-10 | `0025:C5` `precedence:` line ("ACROSS entries there is no order… no test may assert it") | Understated. Across *capability tables* the order is fixed: `loadAccessors` calls `accessorTable` for read, then write, then gate, returning on the first error. A model with a read defect and a gate defect deterministically reports the read one. | C-10 | — | A-only | **CONFIRMED** — `internal/table/load.go::loadAccessors` — three sequential `l.accessorTable(...)` calls (read, write, gate), each `if err != nil { return err }`. The within-table map-range nondeterminism C5 describes is real, but it is scoped to one table |
| D-11 | `0025:C6`, `0025:§consequences` (per-invocation flag as the complete v1 gate) | The flag-only shape is the one shape none of the RDR's own cited reversal-ledger peers (Consul, Hugo, beads) kept: all moved to a persistent gate. A CI pipeline with dozens of call sites papers it over with a Make target or alias on day one, evaporating the "explicit per-invocation friction is the point" safety property. | — | C-4 | **B-only** | n/a (judgement) — grounded in the RDR's own text: `0025:§consequences` concedes "a scripted caller repeats the flag" and `0025:C6` names the config-file successor; `0025:A10` (no config surface exists) is the real constraint |
| D-12 | `0025:A3`, `0025:§approach` ("that wrapper class is small and mechanical, not bespoke per integration") | A: the assumption was narrowed until the spike agreed with it; what the spike actually established (R3) is that **every** established-tool write needs a wrapper, so the write half — the half the Problem Statement exists for — is served by shell the model does not carry and lint cannot see; and the QOC row-2 score for O2 was set on the read/gate half without adjusting the total. B: the spike has n=1 (`git config`), a best-case tool; any non-positional, transactional, or multi-key write needs a wrapper that is not 8 lines. | C-12 | §3 (unnumbered) | **both** (converging from different directions) | **CONFIRMED as documented** — `0025:A3` Evidence states "narrowed from 'a useful class binds directly'" and "R3 **every** established-tool write needs a wrapper"; `0025:§decision-rationale` row 2 scores O2 a 4 and the prose openly says "The second row is scored on the *read and gate* half" with no total adjustment. The n=1 point is verifiable from the spike dir naming a single tool |
| D-13 | `0025:F6`, `0025:A3` hazard P-1 (`tee {artifact}` overwrites the artifact) | A spike-reproduced data-loss hazard mitigated only by prose. C5's six categories cannot catch it — a tool's stdin behaviour is not visible from its argv — there is no lint arm, no runtime guard, and `0025:F8` says bindings hold no state, so no undo. Data loss with a green lint. | C-13 | C-6 | **both** | **CONFIRMED as documented** — `0025:F6` classifies it "Silent risk", offers "The wrapper contract… is the answer" as the whole mitigation; `0025:C5`'s six categories contain no stdin-sink arm; `0025:F8` "no retry, no undo" |
| D-14 | `0025:§metadata` Profile `foundational`, `0025:G-proportionality` ("the sole author of at most one independent load-bearing contract") | Six normative contracts, six new load categories, a new `Refusal` field, a new `table.Model` field, a new CLI flag, a new package, and two signature changes land as one irreversible branch. The Profile field concedes the count and then declares no split. | C-11 | — | A-only | n/a (judgement) — the surface count is confirmable from `0025:§normative-contracts` C1–C6 and `0025:C5`'s six-category registration; the split question is editorial |
| D-15 | `0025:A11`, `0025:A12` vs `0025:C2`, `0025:C4` (Pending assumptions carrying settled-fact normative prose) | The Finalization Gate's own Assumption Verification rule forbids settled-fact prose depending on a `Pending` assumption. A11 and A12 are both `Pending`; C2 and C4 are written in fully normative, unhedged voice. This is a lock-blocking template violation visible on one read. | C-2/C-1 (as consequences) | C-2 + §1(b) | **B-only as a *gate* finding** | **CONFIRMED (documentary)** — `0025:A11` and `0025:A12` both `Status: Pending`; `0025:C2` and `0025:C4` carry the mechanisms in normative blocks with no conditional. Note this is the *meta*-finding whose object-level instances are D-1 and D-2, which A verified against source and B did not |

## 2. Convergence

Five findings were reached independently. Rank these up.

- **D-1 (A:C-1 / B:C-3) — the loader carries no path.** Both passes read
  `flowbind.Registry`'s signature and concluded C2's base dir is unreachable.
  They stop in different places: B stops at "A12 is Pending and the in-memory
  constructor set is unenumerated"; A goes further and shows the *verification
  plan itself* is malformed — A12 says "verify by reading the load site", and
  there is no load site inside the loader, only two callers, neither named,
  both passing an unabsolutized `--model`. Two independent readers landing on
  the same passage from different angles is the signal that C2's argv0 clause,
  not just A12's status, is the thing to fix.
- **D-5 (A:C-5 / B:C-1) — the missing carrier discriminator.** Both saw that
  `Registry` has no `if acc.Command != nil` anywhere. Only A followed the
  `Path: ""` through to `flowbind.load("")` and found the fail-open. B's row
  frames the consequence as "the implementer burns a day"; A's frames it as
  silent state loss. A's is the load-bearing half and the code confirms it.
- **D-6 (A:C-6 / B:C-5) — S5b / the two-reader scenario.** Both passes
  independently flagged the same scenario as unsound, with *opposite*
  mechanisms. See §4.
- **D-8 (A:C-8 / B:C-1) — C6's "registry constructor" does not exist.** Both
  independently counted `flowbind.Registry`'s one production call site and
  found the "one new parameter" framing wrong. Two passes converging on a
  single-call-site count is strong; and both still undercounted the
  registration surface (four `buildRequest` callers, not three).
- **D-12 / D-13 (A:C-12,C-13 / B:§3,C-6) — the wrapper and the stdin sink.**
  Both passes named A3's wrapper claim as the assumption that fails first, and
  both named F6's `tee {artifact}` hazard as spike-proven and prose-mitigated.
  A attacks the assumption's *provenance* (narrowed until the spike agreed,
  and the QOC total never re-scored); B attacks its *generality* (n=1, and
  `git config` is the friendliest possible tool). Neither reading excludes the
  other; together they are the stronger case.

## 3. Divergence

### A-only, correctly found (B missed these)

- **D-2 (`refusalWithKeys` bypasses `refusalOf`).** The single highest-value
  finding in either pass, and B did not reach it. B's C-2 gestures at A11
  ("if `refusalOf`'s call sites don't all cleanly supply-or-nil the error")
  but treats it as an open risk; A read the code and found the concrete
  bypass. Reads are the capability the RDR says binds *directly* to
  established tools, so this is the diagnosis gap on the most common failure.
- **D-3 (no second rendering slot).** B never opened `clierr`.
- **D-4 (`lint` does not call `accessor.Validate`).** B never checked what
  `runLint` actually does. This is the RDR's most-repeated claim — six
  passages — and it is unanchored. A big miss for B.
- **D-7 (relative `--artifact` regression).** B missed it entirely; it is a
  user-facing carrier asymmetry with no Consequence and no Failure Mode.
- **D-9 (stale `executor.go:96` anchor, `clearIsUnreadable` asymmetry).**
  Purely a code-reading finding; B did no line-level anchor checking.
- **D-10 (cross-table defect ordering is deterministic).** Narrow but real;
  it is the reason C5's precedence clause will be edited.
- **D-14 (proportionality).** Editorial. B did not raise scope-splitting.

### B-only

- **D-11 (the per-invocation gate gets wrapped away).** B's best independent
  contribution and A missed it. It is not a code claim — it is an argument
  from the RDR's *own* reversal ledger, which B correctly notices is a list of
  projects that all ended at a persistent gate, cited approvingly in support
  of the one shape none of them kept. A's C-8 covers the *plumbing* of the
  flag; only B covers whether the flag is a control at all. Ranked as a real
  finding, not a decline.
- **D-15 (Pending assumptions carry normative prose — a gate violation).** A
  raises the same facts but as *object-level* code errors (C-1, C-2); B raises
  it as a *process* violation of the Finalization Gate's own written rule,
  citing the template text. This is worth keeping separate: even if D-1 and
  D-2 were resolved tomorrow, the record as written would still have locked
  contracts on unverified assumptions, and that is the reusable lesson.

### Correctly declined

- B's §1(b) worry that A11 might fail on `Executor.Gate`/`Executor.Write`'s
  inline error handling is, on the code, the *wrong* place to worry: gate and
  write both route through `refusalOf` cleanly. The bypass is on `Read`.
  A found the right one; B guessed the wrong one from the same starting
  observation. Not a miss by A — a mislocation by B.
- A's C-11 (proportionality) is the kind of finding B reasonably declined; it
  does not change any contract.

## 4. Contradictions

### 4.1 D-6 — S5b: "the fixture will not load" vs "the fixture loads and picks the wrong reader"

**A (C-6)** claims S5b's model does not load, because `checkAccessorBindings`
refuses a role served by two readers. **B (C-5)** claims it loads fine and
`readerFor` silently picks the first, verifying against a stale reader for two
release cycles.

**Both cited facts are real; A's framing of the arity rule is imprecise, and
that imprecision is load-bearing.** Read
`internal/table/load.go::checkAccessorBindings`: `readerCount` is keyed on
**tag key**, not on role. It refuses an OWNED tag served by ≠1 reader and an
OBSERVED tag served by >1. `internal/accessor/model.go::readerFor` selects on
**`Accessor.Role`**, matching the first read definition with that role.

So: two readers on one *role* load fine as long as their *key sets are
disjoint*. B's migration scenario — an old path-backed reader kept as rollback
alongside a new command-backed one — is authorable exactly in that shape, and
`readerFor` will pick by name-sort order with no ambiguity check. B is right
that the hazard is live.

But A is right about what this does to S5b. The scenario's stated purpose is
to pin `readerFor`'s first-match "so 0016's landing is a visible change." For
the *planned owned key* of a command-backed write, `checkAccessorBindings`
guarantees exactly one reader serves it — so on the case S5b cares about, the
read-back reader is loader-determined and first-match is unobservable. The
fixture that *does* load (disjoint keys) exercises a role-level ambiguity that
never touches the write's read-back path.

**Winner: both, on different halves — and the composite is worse than either
row.** S5b as written is unsound (A), *and* the underlying two-reader hazard
it was meant to guard is real and unguarded (B). Resolving A's row by
respecifying the fixture would not close B's hazard; resolving B's row by
waiting on 0016 would not make S5b's assertion meaningful. Evidence:
`internal/table/load.go::checkAccessorBindings` (tag-keyed arity),
`internal/accessor/model.go::readerFor` (role-keyed first match).

### 4.2 D-8 — how many sites must register the flag

**A (C-8)** frames the risk as "the flag must be registered on all three verbs
or `GetBool` returns false." **B (C-1)** says "however many verb definitions
actually call into `flowRequest`."

**B's framing is right and A's count is wrong — but B never produced the
count either.** `internal/cli/flow_exec.go::buildRequest` has **four**
production callers: `flow_next.go:207`, `flow_resolve.go:227`,
`flow_state.go:135`, `flow_state.go:274`. The RDR's own C6 text names "today:
set-state and the two flow-exec paths" — also three. The registration surface
is larger than the RDR, Pass A, or Pass B states. Evidence:
`internal/cli/flow_exec.go::buildRequest` call sites.

### 4.3 D-3 — the `Detail` field's multi-line invariant

**A (C-3)** warns that an implementer "concatenates a 4 KiB tail into a field
`sanitizeLine`'s siblings assume is one line."

**REFUTED on that sub-claim.** `internal/cli/clierr/clierr.go::CLIError`
documents `Detail` as "May be multi-line", and
`internal/cli/lint.go::runLint` already passes a multi-line `err.Error()`
through it. A's *primary* claim — one slot, one render line, no "beneath" —
stands. The one-line-invariant corollary does not, and the fix half must not
edit the draft to satisfy it.

## 5. Refuted rows

Only one sub-claim in either pass fails against `main`.

- **A:C-3, the one-line-invariant corollary.** `CLIError.Detail` is
  explicitly multi-line-capable and is already used that way. **Do not** amend
  the RDR to add a one-line constraint on `Detail`. The refuted piece is the
  parenthetical hazard, not the row: the missing second slot (D-3, PARTIAL) is
  confirmed and still needs an answer — either an EIA row for `clierr` plus a
  second-slot contract, or the "renders beneath" sentence deleted in favour of
  "the applied sense wins outright and the tail is dropped."

Two rows need their *count* corrected rather than refuted:

- **A:C-8 / B:C-1** — "three `NewExecutor` callers" / "three verbs" is four
  `buildRequest` callers. The RDR's C6 parenthetical ("today: set-state and
  the two flow-exec paths") carries the same undercount and should be
  corrected with them.
- **A:C-6** — `checkAccessorBindings` refuses on **tag arity**, not role
  arity. A's row reads as if two readers on one role are categorically
  unauthorable; they are authorable with disjoint key sets. The row's
  conclusion about S5b survives; its stated mechanism needs the qualifier.

No row in either pass cites code that does not exist.
