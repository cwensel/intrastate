# Verification — RDR 0011 `flow next` match-conditioned candidates

Phase 3a Chain-of-Verification. Each entry below designs an input that would
make a correct implementation VISIBLY violate a REQ, then runs it against the
built CLI (`make build` → `./bin/intrastate`, HEAD `0c2e6e5`). Probes were
authored independently of the Phase 1 test files, which were not read.

Model fixtures were authored fresh under a scratch directory (not committed);
each is described inline by the row shape it exercises, so a later reader can
reconstruct it from the description.

---

## CoVe pass 1 — independent probe of the implementation (no test files read)

**Result: no violations found.** 41 probes across 13 purpose-built scratch
models plus `models/rdr.toml`. Every probe designed to expose a breach of the
candidate predicate, the `--all` restoration, the `unknown` payload, the
demand-set widening, the `disposition` grid, and the exit/output contract
returned the behaviour the record specifies.

### Probes that would have exposed a breach, and what they returned

**The candidate predicate and its boundaries (REQ-1, 3–5, 7–14, 68, 77–80, 83).**

- `models/rdr.toml` at owned `stage=resolved` → exactly `prelock`,
  `resolve-abandon`, `resolve-route-back`; `outcomes` the full three-member
  alphabet; `readers` `["rdr-status"]`. `--all` → 21 candidates, same
  `outcomes`, same `readers`. A stripped-match build reports 21 by default,
  so the count discriminates (REQ-87/88/95).
- A two-row model over an OBSERVED match key: no `--tag` → both rows
  candidates each carrying `{mode, absent}`; `--tag mode=fast` → only the
  present-and-EQUAL row, `unknown` `[]`; the present-and-UNEQUAL row is
  gone (REQ-11/12/77/78/79).
- **Empty-string presence (REQ-4, ASSUMPTION A-1).** A reader reporting
  `mode = ""` against a row matching `mode eq "fast"` → the row is
  EXCLUDED (zero candidates), i.e. `""` is PRESENT and compared, not
  treated as absent. The same model with the key genuinely omitted →
  candidate carrying `{mode, absent}`. A non-empty presence test would
  have inverted the first case; it does not.
- **Dead row (REQ-9/80).** A row authoring `eq = "x"` and `in = ["y"]` on one
  key: candidate carrying `{dk, absent}` while `dk` is absent; excluded at
  `dk=x`, at `dk=y`, AND at `dk=z` — a value on neither side of the pairing.
  The exclusion is therefore the kernel's conjunction, not a CLI literal
  comparison.
- **Multi-key mixed-state row (REQ-103).** One row matching `k1 eq x`,
  `k2 eq x`, `k3 eq x` under `--tag k1=x --tag k2=y` (present-equal,
  present-unequal, absent simultaneously) → EXCLUDED by default; under
  `--all` a candidate whose `unknown` carries no match entry. This is the
  per-key-independence case A12's spike recorded as unmeasured.
- **Set-kinded match key (REQ-6/102/127).** A `set`-kinded observed key
  matched by `eq = ["green"]`, supplied as `--tag 'labels=["green","green"]'`
  → the duplicate collapses to the canonical `["green"]` and the row
  MATCHES; `--tag 'labels=["blue"]'` selects the sibling row instead. A
  filter-before-`KernelRow()` build reimplementing `seamValue` would
  mis-compare here; it does not.
- **Guard decided false with an absent match key** → excluded in BOTH modes,
  with a guard-true sibling surviving in both as the negative control
  (REQ-83/86).
- Source check: the only `Value ==` comparison in `flow_next.go` is the
  pre-existing `ClearSentinel` write-preview split; no literal comparison
  exists in `internal/cli` for match (REQ-3/5/10).

**`--all` restores the PREDICATE, not the reader set (REQ-36–42, 46–47, 50, 86, 97, 115; Q1/Q2).**

- `--all` registered on `next` only. `flow resolve|read-state|set-state
  --all` and `flow --all` and root `--all` each → `command-error`, exit 2,
  one NDJSON envelope. `-a` shorthand is unregistered (`unknown shorthand
  flag: 'a'`), confirming ASSUMPTION A-9.
- **Q1 discriminator (ASSUMPTION A-12).** A model whose key `k` is BOTH a
  match key and a guard key and is NOT in `RequiresOwned`: under `--all`,
  `{k, absent}` STILL appears — from the guard atom. A per-KEY filter would
  have deleted it. The filter is per-ATOM on `Block`, as C2 requires.
- **Filter/dedup ORDER (REQ-40/41/97).** Two independent fixtures: (i)
  `models/rdr.toml` over an empty artifact, where `stage` is both a
  `[rule.match]` and a `[rule.write]` key — `--all` still reports
  `{stage, absent}` from the owned-key walk; (ii) a purpose-built model
  whose owned key `mw` is matched AND written — `--all` still reports
  `{mw, absent}`. A build that dedups before filtering deletes both pairs;
  neither is deleted.
- **Q2 discriminator (ASSUMPTION A-13, REQ-18/115).** Over the match-only
  owned-key model, `next` and `next --all` report the IDENTICAL `readers`
  set `["md","st"]`. `--all` does not revert the demand set.

**The `unknown` payload: shape, vocabulary, dedup, ordering (REQ-22–31, 44, 71–73, 81–82, 84–85, 104, 113, 124–125).**

- Wire members are exactly `{"key","reason"}`, always emitted, on a list
  named `unknown` present as `[]` rather than omitted — confirmed by
  `keys` on the candidate object with and without `--evaluate-gates`
  (REQ-44, A-2). `gates` remains `omitempty` (0005-carried, unchanged).
- **`uncomparable`, the payload-only reason (REQ-84/112).** An `int`-kinded
  OWNED key delivered by a reader at the non-integer value `"abc"`, guarded
  `gt = 3` → `{num, uncomparable}`. (`--tag` cannot produce this: it
  conforms the value at parse and refuses `flow-tag-invalid`, so the
  reader route is the one reachable producer, exactly as A13 states.)
- **The `owned_state_unavailable` fidelity limit (REQ-31/85/112).** A row
  carrying BOTH that `uncomparable` guard atom AND an owned key its bound
  reader reports absent → candidate carrying `{extra, absent}` from the
  walk, and the `uncomparable` fact is NOT reported. The sibling row
  without the missing owned key reports `{num, uncomparable}` in the same
  run. This is the precedence limit pinned rather than a defect.
- **Dedup on the pair (REQ-24/29).** A guard atom over an ABSENT key
  produces `{num, absent}` exactly ONCE — the walk and the kernel payload
  agree and collapse.
- **Sort order (REQ-29/113).** A row minting facts deliberately out of
  order — match atoms authored `zz`, `mm`, `aa`; gate ids `zgate`,
  `agate`; plus an unestablished owned key `extra` — returns
  `aa/absent, agate/not-evaluated, extra/absent, mm/absent,
  zgate/not-evaluated, zz/absent`: all four sources interleaved and sorted
  by `(key, reason)`. An unsorted build fails here.
- **Gate ids (REQ-26/27/28/81).** Ids come from the ROW's own `gate` list,
  not a model-level set: two rows declaring different gates report their own
  id only. `not-evaluated` appears on no non-gate entry, and `uncomparable`
  on no match/owned/gate entry across every fixture (REQ-124). Gate ids are
  emitted on `owned_state_unavailable` too (REQ-28).
- **Set-kinded property.** The `unknown` list is a set of pairs: identical
  pairs from different sources collapse; distinct pairs on one key are not
  reachable (no key can carry two reasons), matching S7's stated limit.
- Text mode renders `candidates[0].unknown[0].key: dk` /
  `...reason: absent` through the shared flattener, both members present,
  no bespoke template (REQ-72/73).
- No schema or version field is emitted (REQ-125).

**Gates run only for REPORTED candidates (REQ-12/70).** The sharpest available
discriminator: a model whose match-EQUAL row declares gate `gfast` and whose
match-UNEQUAL row declares gate `gslow`, with the `gate` artifact role
deliberately UNBOUND. Under `--tag mode=fast --evaluate-gates` the run fails on
`gfast` only; under `--tag mode=slow --evaluate-gates` it fails on `gslow` only.
A build that ran an excluded row's gates would have failed on `gslow` in the
first case. It does not.

**The demand-set widening and its effect on `flow resolve`, in BOTH directions
(REQ-16–21, 92–93, 114–118, 120, 126).** A model declaring an owned key `mode`
that one row MATCHES on but no row writes, clears, or guards, served by exactly
one reader:

- `next` invokes `md`, carries `mode` in the view, and match-DECIDES the row
  (no `{mode, absent}`) — the discriminating oracle for the term (REQ-114).
  Same `readers` under `--all` (REQ-115).
- `flow resolve --outcome go`, all three named arms: reader bound and
  answering → PLAN, exit 0; role unbound → `flow-artifact-missing`, exit 2;
  reader refusing (malformed artifact) → `flow-accessor-failed`, exit 3.
  `flow-no-match` in NONE of the three (REQ-19/93).
- **The BREAKING arm reproduced (REQ-20/116).** With a second ordinary row
  for the same outcome that does NOT match on `mode`, and `md` unbound: the
  run that would have returned `plain-row`'s plan instead exits 2
  `flow-artifact-missing`. The union-over-rows demand pulls the reader in
  above the kernel. This is the class C1/Consequences NAME and ACCEPT; it
  fires exactly here and nowhere observed beyond it.
- **Neither over- nor under-fires.** Over a model whose match key is
  OBSERVED rather than owned, the term is empty and `readers` is unchanged
  between `next` and `resolve` (REQ-21). Over `models/rdr.toml`, `flow
  resolve` at `advance`/`revise`/`abandon` still plans and an unmodeled
  outcome still refuses `flow-unmodeled-outcome` (REQ-118).
- **Both verbs agree (REQ-117).** `next` and `resolve --outcome go` over one
  model return byte-equal `readers`.
- **Escape-row scoping.** `resolve --outcome go` does NOT demand the reader
  serving an escape row bound to `stop`; `resolve --outcome stop` DOES. `next`
  takes the union over all rows and demands both. Consistent with the demand
  set being a union over the OUTCOME's rows and `next`'s DEV-1 reading (a);
  not a violation.

**The three pins licensing the presence test (REQ-63/106).** A model declaring
one owned key served by two readers is refused at load
(`malformed_accessor_binding: owned tag mode is served by 2 readers; want
exactly one`); `--tag mode=fast --tag mode=slow` → `flow-tag-duplicate`, exit 2;
`--tag phase=b` on an owned key → `flow-tag-owned`, exit 2. All three precede
any probe.

**Structural REQs checked at source rather than at runtime.** `internal/resolve`
carries an EMPTY diff across the branch (`git diff --name-only main...HEAD`),
and the only production files changed are `internal/cli/flow_next.go` and
`internal/cli/flow_exec.go` — REQ-34/107/120 hold in the one mechanically
checkable form MVV 9 names. No occurrence of `unresolved`/`Unresolved` remains
in production or in the three 0005 test files' prose (REQ-55/60); the
`flow_next_0005_test.go` header names RDR 0011 as the source of the default
(REQ-53). `docs/cli-output-contract.md` documents the payload, the `unknown`
reason table, and the `--all` invocation, and `README.md:43` no longer states
the overridden default (REQ-56/122). `go test ./...` green.

### Coverage and what was sampled rather than probed

Probed directly, with a designed-to-break input: REQ-1..21, 22..31, 32..44,
46..50, 63..64, 66..73, 74..86, 87..93, 95..106, 112..118, 120, 122..127. That
is the candidate predicate, the demand set, the `unknown` payload, `--all`, the
`disposition` grid in full, the MVV, and the exit/output contract.

Sampled or checked structurally rather than by probe (they are review-time or
diff-time obligations no runtime input can express, and the record says so):
REQ-45 (request identity), REQ-51..62 (help wording and the test re-homing
census — the header, the production rename, and the prose sweep were checked;
the five per-line test reads were NOT, since the files are Phase 1 artifacts
this pass is forbidden to read), REQ-65, REQ-94, REQ-107..111, REQ-119,
REQ-121, REQ-128..130.

One class could not be probed from the CLI and is not a gap: the
present-but-CONFLICTED match key (REQ-64, A3), which the three pins above
refuse before a probe exists — verified negatively by confirming all three
refusals fire.

### Verdict

No `FAIL-N` entry is owed from this pass. Every probe designed to distinguish
the specified behaviour from its nearest wrong neighbour — a filter applied
per-key instead of per-atom, a dedup that runs before the `--all` filter, a
demand set that reverts under `--all` or unions the wrong rows, a presence test
that treats `""` as absent, a payload read off the wrong refusal kind, gates run
for an excluded row — returned the specified behaviour.

## Phase 3b — adversarial failure-mode review

Three failure modes, each anchored in `## Trade-offs > ### Failure Modes`,
with the test that catches it. All three currently FAIL against the shipped
build. Tests live in `internal/cli/flow_adversarial_0011_adv_test.go` over
one new fixture, `flowEscapeMatchOwnedModel`.

The common root: C1's demand-set term adds each row's MATCH-block owned keys
to `invokedReaders`, applied to EVERY row of the model — escape rows
included. `0002:C4` makes a match block MANDATORY on every rule, escape
rules included (`normalizeRule` refuses one without), so an escape row
matching an owned key is not an exotic authoring choice a model author can
avoid. DEV-8's guard term had the same shape but a guard block is optional
on an escape rule; a match block is not.

### ADV-1 — `flow next` refuses for an escape row it never reports

- **Failure mode.** `flow next` exits 2 `flow-artifact-missing` over a model
  whose only unbound-reader demand comes from an ESCAPE row. `flow_next.go`'s
  row loop skips escape rows structurally (`if len(row.Escape) != 0 {
  continue }`) before `probeRow`, so no probe is ever built from one and no
  candidate is ever reported for one. The verb refuses to run at all for a
  fact no reported row could consume.
- **RDR anchor.** `0011:F6` — the record's ONE named behaviour-change class
  for the demand-set term, scoped to `flow resolve` ("the requested outcome
  has a row on that key"). C1 says the same in its own voice: "The
  `flow resolve` VERB changes in exactly the one class the demand-set
  paragraph names, and nowhere else", justified by "a row cannot be
  match-decided without the key". An escape row is never match-decided by
  `next`. No Failure Mode contemplates `flow next` REFUSING; `0011:F7` says
  "Recovery: none needed" — and there is no recovery here, because C1 itself
  records that an owned key cannot be supplied (`--tag` → `flow-tag-owned`,
  "no other flag or ambient channel").
- **Test.** `TestAdv1_NextDoesNotRefuseForAnEscapeRowsMatchOwnedKey`.
- **Currently.** FAILS (`flow-artifact-missing`, role `side`).

### ADV-2 — `flow resolve` loses a plan to an escape row's rescue-phase reader

- **Failure mode.** `flow resolve --outcome go` that returns a plan through
  an ordinary row now exits 2, because an ESCAPE row binding the same
  outcome matches on an owned key whose reader is unbound. The rescue phase
  that would consult that row is provably unreachable on this run:
  `internal/resolve.escapeOrRefuse` is called only from the
  `len(selected) == 0` and `default:` arms of `Resolve`'s selection switch,
  and this run selects exactly one ordinary row.
- **RDR anchor.** `0011:F6`, at its boundary. F6's accepted class carries a
  stated warrant — the old refusal HID a fact ("`resolve` … refuses
  `flow-no-match` for every row matching on that key, over an artifact that
  HOLDS the fact"). Nothing is hidden here: the plan is correct and complete
  without the key. `invokedReaders`'s own doc comment already names this
  masking hazard for escape rows ("masking a valid plan") and fences it by
  OUTCOME alone — a fence an escape row binding the requested outcome
  passes, and one that is entirely empty for `next` (`outcome == ""`).
- **Test.** `TestAdv2_ResolveKeepsItsPlanWhenOnlyAnEscapeRowDemandsTheReader`.
- **Currently.** FAILS (`flow-artifact-missing`, role `side`).

### ADV-3 — `--all` is not the escape hatch F1 prescribes

- **Failure mode.** F1's remedy is a two-run protocol (run default, run
  `--all`, diff) and presumes both runs emit a payload. Over this model
  class neither does: C1 requires the demand set to be "a property of the
  MODEL … not of the mode … identical under --all", so `--all` inherits the
  same pre-row-loop refusal. The mode-independence C1 requires is exactly
  what denies the remedy F1 promises.
- **RDR anchor.** `0011:F1` — "Diagnose with `--all` … That eye-comparison
  is the whole remedy today: `intrastate` ships no model-inspection verb …
  and this RDR adds none." The caller has no verb left; `dump` is asserted
  FOREIGN to the flow group at `flow_surface_0005_test.go:105`.
- **Test.** `TestAdv3_AllRemainsF1sDiagnosticOverAnEscapeRowDemand`, which
  asserts only that each mode returns a payload with a `candidates` key —
  never which rows it holds — so it cannot be satisfied by weakening the
  predicate.
- **Currently.** FAILS in BOTH subtests (`default` and `all`).

### Confirmed handled — modes probed that the implementation gets right

These were investigated as candidate failure modes and the build proved
correct; no test was added, because a test that passes catches nothing.

- **A present-but-CONFLICTED key reaching the kernel.** C1's exactness rests
  on this being unreachable (A3). `TagSet.matches` returns false for a
  conflicted key WITHOUT comparing it, which would drop a row with nothing
  in `unknown` explaining it — A1's "If wrong" exactly. Traced every channel
  into `Owned`/`Observed`: `checkAccessorBindings` gives each owned key one
  reader; `Executor.Read` classifies over `def.RequestedKeys()`, so one
  reader cannot emit a key twice; `parseTags` refuses a repeated `--tag`
  (`flow-tag-duplicate`), an owned one (`flow-tag-owned`), and the reserved
  one (`flow-tag-reserved`). Unreachable as claimed.
- **`probeRow`'s `err != nil` swallow** returning a zero `resolve.Result`,
  which `excluded` reads as non-excluding and `summarize` reads as carrying
  no kernel facts. Both error returns are vacuous on the probe as A1 claims:
  `Table.CheckValid` keys on `len(Escape) != 0 && len(Writes) != 0` and the
  probe strips `Escape`; `CheckInput`'s three channels need a `recognized`
  key in `Owned`/`Observed`/`RequiresOwned`, and `renderWrites` refuses a
  write to a non-owned tag while `outcomeBinding` lifts every `recognized`
  match atom out of `Row.Atoms`.
- **View precedence divergence.** `internal/cli::assembledView` merges owned
  then observed (observed wins); `internal/resolve::assemble` merges
  observed then owned (owned wins). Harmless: `flow-tag-owned` keeps the two
  key sets disjoint, so the union is identical and no key takes a different
  value.
- **`--all` filter vs. a key that is BOTH a match atom and a guard atom on
  one row.** The filter is per-ATOM on `atom.Block`, so the guard atom's
  `{key, absent}` survives under `--all` — verified by direct run.
- **`RequiresOwned` reached through a CLEAR list rather than a write block.**
  The owned-key walk's `{key, absent}` survives `--all` there too, the same
  way REQ-40's oracle proves it for a written key.
- **A gate id colliding with a match key name.** Emits both
  `{k, absent}` and `{k, not-evaluated}` — correct, and precisely why C1
  dedups on the `{key, reason}` PAIR rather than on the key.

## Phase 3c — fixup

**Outcome: all three Phase 3b adversarial oracles resolved; `make check`
green; 130/130 REQ oracles still green; `internal/resolve` untouched
(`git diff --stat internal/resolve` empty).**

### The defect

C1's demand-set term added each row's match-block owned keys to
`internal/cli/flow_exec.go::invokedReaders` over EVERY row of the model,
escape rows included. `0002:C4` makes a match block mandatory on every
rule, escape rules included, so the class is reachable by ordinary
authoring rather than exotic.

### The reading taken

The match-owned demand term runs over the rows the INVOKING VERB can
actually consult: for `next`, its non-escape rows; for `resolve`, every
row of the requested outcome, escape rows included. Implemented as one
predicate in one function — `if outcome != "" || len(row.Escape) == 0`.
Recorded as DEV-8 (`SPEC-UNDER`).

**Evidence quoted from the record.** For the `next` half:

- C1 fixes the predicate over "each **non-escape** row of the requested
  model".
- C1 justifies the term by "the assembled view MUST actually carry the
  keys **the predicate reads**" and "a row cannot be match-decided without
  the key".
- The Joint-check line: "`runFlowNext` never lists escape rows."
- F7: "**Recovery**: none needed — `next` is effect-free in both modes" —
  and C1 itself records the key is owned, so `--tag` is refused
  `flow-tag-owned` and "no other flag or ambient channel can supply it".
- F1's diagnostic is a two-run default-vs-`--all` diff, which a
  pre-row-loop refusal denies in both modes.

For the `resolve` half — escape rows STAY in the demand set:

- C1: "the reader is invoked **even when the row `resolve` would select
  does not itself match on the key**: over an UNBOUND or REFUSING reader
  such a run turns from a plan into exit 2/3, before the kernel and
  **above the escape phase** — a class this contract NAMES AND ACCEPTS
  here … it is not denied."
- F6 scopes that class as "the requested outcome has a row on that key";
  `rescue-row` binds `go` and matches on `mode`, so the ADV-2 fixture is
  inside F6's scope.
- Independently: `internal/resolve.escapeOrRefuse` evaluates
  `view.matches(row.Match)` on every escape row binding the requested
  outcome, that phase's reachability is unknowable before the kernel runs,
  and `TagSet.matches` is two-valued — so narrowing `resolve` would let an
  absent key silently fail a rescue into `no_match` over an artifact that
  holds the fact, reinstating the hidden-fact defect the term removes.

### Per-oracle disposition

- **ADV-1** — FIXED in production. `flow next` no longer refuses
  `flow-artifact-missing` for an escape row's match-owned key; `readers`
  excludes `side`; `candidates` is exactly `[plain-row]`.
- **ADV-3** — FIXED by the same term (both subtests). Default and `--all`
  each emit a payload carrying `candidates`, restoring F1's two-run
  diagnostic.
- **ADV-2** — TEST-FIXTURE correction (DEV-9). The assertion contradicted
  the clause it cited; C1 names and accepts the class verbatim. The oracle
  is re-anchored to pin the BOUNDARY: `resolve` still demands the escape
  row's reader, and the refusal is `flow-artifact-missing` naming the
  `side` role — the remedy F6 prescribes. Fixture, model, and run
  unchanged; the oracle still fails a narrowed build and a build whose
  refusal does not name the role.

### Contract-collision check

None. No REQ oracle changed disposition. REQ-17's "MUST NOT be scoped to
`next`" is preserved verbatim — the term still binds both callers; it is
the escape-row CLASS that is scoped, not the verb — and
`TestReq17And117`'s fixture carries no escape rows, so its reader-set
agreement is unaffected. C1's mode-independence holds: `--all` and the
default consult the same row set and only the predicate over it differs.

### Pre-existing failure found and repaired

`TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec` was **already
red at the Phase 3b commit** (`2271b75`), verified by stashing the Phase 3c
edit and re-running. Its cause was Phase 3b's file NAME:
`flow_adversarial_0011_adv_test.go` ends `_adv_test.go` and matched neither
REQ-120's allow-list nor its `_0011_test.go` exemption. Renamed to
`flow_adversarial_adv_0011_test.go`; the oracle and its allow-list are
untouched. Recorded under DEV-9.

### Green status

`make check` passes end to end (fmt-check, vet, `golangci-lint` 0 issues,
build, `lint --model models/rdr.toml` 0 findings, `go test -race` all
packages). `internal/cli` coverage 89.2%. No new public surface.
