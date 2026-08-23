Model: claude-opus-5[1m]

# RDR 0002 — Stage 4 Resolve research log (scoped re-entry)

Scope: `Draft [revised from Final 2026-08-12; re-verify A2, A7]` plus A8
(`Pending`, MVV-critical). A1, A3, A4, A5, A6 carried forward as already
Verified; their anchors were re-checked for survival only (below), not
re-derived.

## Domain routing

Every claim in scope is about **this project's own code** (`internal/resolve`)
or about **peer design documents in this repo** (JDR 0001, RDR 0007/0008/0009).
No claim in scope is about external dependency, peer-tool, or standard
behavior, so no `arc` corpus search was run — the corpora in
`.rdr/resources.md` do not own any behavior A2/A7/A8 assert. The external
prior art the RDR cites (Sismic, `transitions`, Stateless, statechart
literature) backs A6 and the Decision Rationale, both carried forward.

negative: A2/A7/A8 — no corpus evidence sought in DevRef / StateMachineLit /
StateMachineRes / PapersFast; the owning source is project code and in-repo
peer RDRs, per the prompt's two-source-domain rule.

## Reuse audit (`.rdr/env.md` reuse-audit paths)

All nine reuse-audit paths exist. Findings:

- `internal/resolve/` — kernel present (`resolve.go`, six `_test.go`). No
  loader, no normalizer, no TOML. Confirms the Existing Infrastructure Audit
  row "Normalizer/table package … Not implemented → Introduce".
- `internal/cli/clierr/clierr.go::CLIError` — carries `Code`, `Message`,
  `Param`, `Detail`, `Hint`, `Group`, `Cause`. A5's anchors resolve.
- `internal/cli/respond/respond.go::Fail` — emits the envelope. A5 resolves.
- `internal/cli/config/config.go::Load` — emits `config-not-found` /
  `config-read-error`; `config-invalid` is a TODO. A5 resolves.
- **Reuse finding (new):** `config.go:98` carries `// TODO: parse data into
  cfg once a TOML library is chosen`. No TOML library is in `go.mod` today.
  RDR 0002's New Dependencies section picks `github.com/pelletier/go-toml/v2`;
  that choice, when it lands, also settles the open choice for
  `internal/cli/config`. Not a reason to route back — the RDR is introducing
  the capability, not duplicating one — but the two consumers should share one
  library rather than each picking independently.

No reuse finding refutes the approach. Nothing in the codebase already parses
or normalizes a transition table.

## A2 — row order is not part of successful edge selection

Three independent legs, all confirmed:

1. `docs/jdr/0001-resolve-kernel-seam.md` §D2 — "**Resolved: (b).** Gate, then
   count." Its "Lands in **0002**" clause obliges 0002 to restate the resolver
   flow as gate-then-count and to qualify Scenario 4 with "no sibling
   candidate is unevaluable". Both duties are discharged in the current draft
   (Normative Contracts gate-then-count block; Testing Strategy scenario 4).
2. `docs/rdr/0007-guard-predicate-totality.md` (Status: **Final**) states the
   ordering as the kernel's in a normative block, and verifies it itself by
   Source Search (its A20).
3. `internal/resolve/resolve.go::Resolve` — the shipped kernel already gates
   before counting:
   `selected, blocked := gate(candidates, in.Guards, view); if blocked != nil
   { return refuse(in, *blocked), nil }; switch len(selected)`.
   `::gate` runs prune-GuardFalse → `missingOwned` →
   undecidable-veto, then returns survivors. `::escapeOrRefuse` delegates to
   the same `gate` before `switch len(viable)`.

No first-match, priority, or early-return selection path exists.
`::missingOwned` sorts with `slices.Sort(missing)` and its doc comment states
"Row order is a normalization detail RDR 0002 owns and must not reach the
reported diagnosis." `adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
freezes the property.

Verdict: A2 holds, and its evidence is stronger than Design Decision — the
ordering is implemented and test-frozen, so the claim is Source-Search
verifiable rather than merely decided.

## A7 — deterministic expanded-table ordering is a format contract

(a) **Fixed.** JDR 0001's withdrawn JD-11 defect (dump clause omitted row kind
and escape failure classes while the round-trip invariant required them) is
repaired in the current text. Both lists now carry row identity, source
locator, row kind, outcome, predicate atoms with block, writes, required-owned
keys, and escape failure classes. The dump list additionally leads with
"model id", which the round-trip list omits; model id is inside row identity
per the Identity decision, so this is redundant wording, not a semantic gap.

(b) **Gap in the contract text.** "sorting rows by row identity" does not state
the field order inside the identity tuple, that the compare is lexicographic,
whether the source locator participates, or how a present-vs-absent expansion
suffix ranks. Across a multi-model dump the clause as written does not
determine a total order. RDR 0007 sets the drafting precedent by spelling its
own payload tuple field by field — `(RuleID, SourceLocator, key, block,
operator token, literal)` — and arguing its totality by construction.

(c) **The spike does not witness the contract.** `evidence/spikes/main.go::Row`
carries six fields (`Model, Rule, Kind, Source, Match, Write`). Against A7's
nine-field list:

| A7 field | Spike |
| --- | --- |
| model id | emitted |
| row identity | partial — `(model, rule)`; no expansion suffix path exists |
| source locator | emitted (`source=rdr:prelock`) |
| row kind | emitted (`kind=transition` / `kind=escape`) |
| outcome | **absent as a row field** — still an atom in `match=[…]` |
| predicate atoms with block | emitted, but as flattened `all:`/`unless:` string prefixes, not `(key, operator, literal, block)` tuples |
| writes incl. `<clear>` | emitted (`prelock_lens=<clear>`) |
| required-owned keys | **absent entirely** |
| escape failure classes | emitted in the **wrong container** — `write=[escape=no_match]` |

(d) **Sorts deterministically, on the wrong keys.** `sort.Slice(rows, …)`
compares `Model` then `Rule` — no locator, no suffix. `sort.Strings(match)`
sorts whole formatted strings whose block prefix leads, so the effective order
is (block, key, operator, literal), inverting the contract's stated
key-then-block order. Writes sort as strings including the injected `escape=`
pseudo-entry.

The escape-row rendering additionally breaches RDR 0009 (Final): "a Row with a
non-empty Escape list MUST have an empty Writes slice", which 0002's own
escape-rule clause restates.

Verdict: A7's decision stands as stated in the contract text, but the spike is
a Phase-1 parse spike, not a normalizer conforming to the current dump
contract. It cannot serve as the ordering witness. Item (b) is a genuine
under-specification in the RDR's own text.

## A8 — every normalized row binds exactly one outcome

Source Search against `internal/resolve/resolve.go`, all four legs confirmed:

1. `::Resolve` filters candidates on `if row.Outcome != in.Recognized {
   continue }`, after skipping escape rows and before `view.matches(row.Match)`.
2. `::escapeOrRefuse` applies the same filter to escape rows:
   `if row.Outcome != in.Recognized || !view.matches(row.Match) { continue }`.
3. `::assemble` binds the outcome under `const recognizedTagKey = "recognized"`
   only when `in.Recognized != ""`, and writes it first so `Observed` then
   `Owned` take precedence.
4. `::Resolve` refuses `unmodeled_outcome` before any row is consulted:
   `if !in.Table.models(in.Recognized) { return refuse(in, Refusal{Kind:
   KindUnmodeledOutcome}), nil }`.

Outcome-less rows: with `Input.Recognized` non-empty, a row whose `Outcome` is
`""` never matches on either path. With `Input.Recognized == ""` and the
alphabet not containing `""`, `unmodeled_outcome` fires before rows. **One
residual path:** if `Input.Recognized == ""` *and* the alphabet contains `""`,
`Table.models("")` passes, `assemble` skips the binding (view has zero tags),
and a row with `Outcome: ""` matches and emits a Plan with no bound outcome.
RDR 0002 already forbids this upstream — its Testing Strategy lists "an
alphabet containing the empty string" as a rejected malformed variant, and its
recognized-totality clause requires the alphabet be non-empty, duplicate-free,
and free of the empty string. The kernel does not enforce the totality itself;
it relies wholly on the normalizer. That dependency is now stated explicitly
on A8 rather than left implied.

Field names confirmed verbatim on main, no renames: `Row.Outcome string`,
`Input.Recognized string`, `Row.RequiresOwned []string`, `Table.Outcomes
[]string`, `Table.models`, `Resolve`, `assemble`, `escapeOrRefuse`.

## Cross-cutting findings surfaced while verifying (not in scope to fix here)

These are recorded because they bear on lock readiness, not resolved by this
stage:

- **F1 — `resolve.OpExists` / `LiteralTrue` / `LiteralFalse` do not exist.**
  Full-repo grep finds them only in RDR prose (0007 specifies them normatively)
  and in RDR 0003's spike harness. Nothing under `internal/`. RDR 0002's
  existence-constant contract and its Testing-Strategy scenario 6
  ("byte-for-byte") therefore point at unwritten code. RDR 0007 is Final but
  unimplemented; this is a forward reference, not a defect in 0002's design,
  but no Source Search can attest it today.
- **F2 — `Row.Guard` is still `string`.** JDR 0001 §D1 resolved that `Row`
  carries parsed atoms, and RDR 0007 absorbs the `Row` change. That change has
  not landed. So 0002's per-atom `block` retention has no kernel type to
  target yet. Same class as F1: sequencing, not contradiction.
- **F3 — `Row` splits `NextTags` and `Writes`,** and has no `Kind` field
  (escape is inferred from a non-empty `Escape` list). 0002's normalized-row
  field list names "writes" and "row kind" without mapping either onto the
  kernel's actual field split. Implementation will have to state the mapping.
- **F4 — `recognizedTagKey` is unexported** (`resolve.go:110`). The normalizer
  must agree byte-for-byte with a constant it cannot import. RDR 0008 owns the
  reserved key; whether it is exported is that RDR's call.
- **F5 — `assemble` writes `recognized` first,** so an owned or observed tag
  keyed `recognized` shadows the outcome in the view. The `Row.Outcome` filter
  reads `in.Recognized` directly and is unaffected, but a guard reading the
  `recognized` key would see the shadow. 0002's normalizer rejection of an
  owned/observed declaration named `recognized` is currently the only barrier.
- **F6 — spike fixtures write undeclared tags.** `rdr-fixture.toml`'s
  `reconcile-rewind` writes `rewind_scope` and clears `prelock_lens`; neither
  has a `[tags.<tag>]` declaration. 0002 requires the model to declare every
  tag it matches or writes, and lists "write to non-owned tag" as a validation
  category. The fixtures are declared canonical examples implementation tests
  must promote, so they must be conformant before they are promoted.
- **F7 — `[dump].order` in both fixtures is `["model","rule","match","write"]`,**
  which predates the current nine-field dump contract.

---

# Round 2 — deep research for the author's round (2026-08-22)

## Go language/stdlib grounding for the dump-ordering clause

Verified against `$GOROOT` (`/opt/local/lib/go-1.26`) and the `langref`
peer-CLI checkout set. Every claim below carries a source citation; none is
from memory.

### Byte-lexicographic ordering is spec-guaranteed

Go spec, Comparison operators (`go_spec.html:5145-5146`): "String types are
comparable and ordered. **Two string values are compared lexically
byte-wise.**" Restated for `min`/`max` at `go_spec.html:7799-7801`.

The whole comparison chain inherits it:
- `<` on strings — the spec statement itself.
- `cmp.Compare[string]` — `$GOROOT/src/cmp/cmp.go:41-58`, implemented purely
  via `x < y` / `x > y`; the NaN branches are dead for strings.
- `strings.Compare` — `$GOROOT/src/strings/compare.go:15` delegates to
  `internal/bytealg.CompareString`; byte-by-byte then length
  (`compare_generic.go:11-36`).

Locale-independent and stable across platforms, architectures, and Go
versions — no `unicode/collate` or locale input is anywhere in this path, and
changing it would be a spec break. **"byte-lexicographic" is safe to write
into a normative contract**; it is the spec's own wording.

Two precision notes for the contract text:
- For valid UTF-8, byte order and code-point order coincide. They diverge only
  for invalid UTF-8, so "byte-lexicographic" is correct either way.
- Byte order is not human-alphabetical: all uppercase ASCII precedes all
  lowercase (`Z`=0x5A < `a`=0x61), with `_`=0x5F between. Stable, but a
  mixed-case model/rule id set will read oddly in a review dump.

### Absent-vs-present expansion suffix

`""` sorts before every non-empty string, and this is a consequence of the
spec rule above rather than an implementation accident
(`internal/bytealg/compare_generic.go:28-35`: after a common prefix is
exhausted, shorter returns `-1`; `""` shares a zero-length prefix with
everything).

The Go property is airtight. **The risk is modeling, not Go**: `""` means both
"no suffix" and "present but empty suffix," and since the tuple is a row
*identity*, that collision silently merges two distinct rows. In this RDR the
suffix is the outcome literal, and the `outcomes` alphabet is already required
to be non-empty, duplicate-free, and free of the empty string — so the empty
suffix is unrepresentable by construction. That makes `""` safe here, but the
invariant should be stated rather than left implicit in the sentinel choice.

Alternative encodings, for the record: `*string` (needs explicit nil-ordering,
`cmp.Compare` won't take it) or an explicit presence `bool` compared first
(Go orders `false < true`, so absent-first falls out). The langref precedent
for explicit presence comparison is
`etcd/tests/robustness/model/non_deterministic.go:99-104`.

### Sorting idiom

`slices.SortFunc` + `cmp.Compare` is current idiom; `sort.Slice` is
*discouraged but not deprecated* — no `Deprecated:` marker exists anywhere in
`$GOROOT/src/sort`, and its doc (`sort/slice.go:22-24`) says only "in many
situations, the newer [slices.SortFunc] function is more ergonomic and runs
faster." Do not write "deprecated" into the design.

- `slices.SortFunc` — `$GOROOT/src/slices/sort.go:30`; cmp must be a strict
  weak ordering.
- `slices.SortStableFunc` — `$GOROOT/src/slices/sort.go:37`.
- The `modernize` analyzer already rewrites `sort.Slice` → `slices.Sort`
  (`$GOROOT/src/cmd/vendor/.../modernize/sortslice.go:33-49`).

Field-by-field tuple comparison with early return is the idiom, e.g.
`langref/etcd/tests/robustness/model/non_deterministic.go:89-117`,
`langref/goreleaser/internal/redact/redact.go:53-58`.

**`SortFunc` vs `SortStableFunc` is a determinism tell, not a perf choice.**
If the comparator is a true total order over row identity, ties mean identical
rows and `SortFunc` is correct and faster. If two *distinct* rows can compare
equal, `SortFunc` yields non-deterministic output and `SortStableFunc` only
masks it by making input order load-bearing. Either way the fix is a total
comparator. A cheap check: sort, then assert no adjacent pair compares 0 —
turning the totality assumption into a tested invariant.

**In-repo precedent:** `internal/resolve/resolve.go:546` already uses
`slices.SortFunc(out, compareRefs)`, and `::compareRefs` (`:528-533`) is the
field-by-field `strings.Compare` chain. RDR 0002 should mirror the kernel's
existing idiom rather than invent one.

### Map iteration and encoder key order

Randomized iteration is spec-mandated (`go_spec.html:6750-6751`) and
deliberately implemented (`internal/runtime/maps/table.go:647-650`). The Go
1.23+ idiom is `slices.Sorted(maps.Keys(m))` — `maps.Keys` returns
`iter.Seq[K]` (`$GOROOT/src/maps/iter.go:25`), `slices.Sorted`
(`$GOROOT/src/slices/iter.go:66`). **`maps.Sorted` does not exist** (proposal
#68598 not accepted); do not write that spelling into the design. 32 call
sites of `slices.Sorted(maps.Keys(` across 10 langref repos, e.g.
`golangci-lint/pkg/goanalysis/runner.go:163` — `// for determinism`.

**Encoder ordering is the sharp edge:**
- `encoding/json` v1 sorts map keys (`encode.go:794-796`, `strings.Compare`).
- `encoding/json/v2` does **not** sort by default (`v2/arshal.go:112-113`);
  it needs the `Deterministic` option, whose own doc warns that "different
  versions of the same program are not guaranteed to produce the exact same
  sequence of bytes" (`v2/options.go:125-129`). The v1 API stays safe under
  `GOEXPERIMENT=jsonv2` because `DefaultV1Flags` includes `Deterministic`.
- `BurntSushi/toml` v1.6.0 sorts, but partitions direct values before
  sub-tables first (`encode.go:412-420`), so its output is sorted
  *within group*, not globally key-sorted.

**Conclusion: never rest a byte-stability contract on an encoder's map-key
ordering.** Emit from a pre-sorted slice of structs, never from a `map`. That
keeps the guarantee owned by this RDR and immune to encoder-version drift.
This directly supports the existing Performance Expectations sentence
("sorting stable model/rule/atom/write keys rather than relying on map
iteration") and is worth strengthening from "should" to the normative clause.

### Toolchain note

The RDR's Technical Environment and `go.mod` say Go 1.26.3; this machine's
toolchain is 1.26.6. Nothing above turns on the patch version (all findings
are spec-level or long-stable stdlib), but the stated version and what
actually builds should be reconciled.

## In-repo identity/ordering precedent (decisive for the sort question)

The kernel already has a row-ordering comparator, and RDR 0009 already
reasoned about its totality:

- `internal/resolve/resolve.go::RowRef` is `{RuleID, SourceLocator}` — **no
  model id field at all** (`resolve.go:275-279`).
- `::compareRefs` (`:528-533`) compares `RuleID` then `SourceLocator`.
- **RDR 0009 A8 (Verified, Source Search):** "Collapsing equal `RowRef`
  identities is the only tiebreak available … because `compareRefs` does not
  totally order distinct rows and RDR 0001's REQ-2/REQ-10 forbid breaking the
  tie by table position." 0009 records that rows with equal `(RuleID,
  SourceLocator)` are constructible and already present in the frozen suite
  (`fixup_test.go:79,101,125`), and that a positional tiebreak
  (`SortStableFunc`) is excluded by
  `adversarial_test.go::TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder`
  and `fixup_test.go::TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder`.

So the corpus has **already decided** that row ordering must be a function of
the input tuple and never of table position, and has already discovered that
`(RuleID, SourceLocator)` alone is not total. RDR 0002 does not get to invent
a different answer; it must supply an identity that *is* total.

**The locator subsumes model id.** RDR 0002 line 305: "The locator must
identify at least the model id and rule id; byte line/column coordinates are
optional diagnostic detail." So `(RuleID, SourceLocator)` does order across
models — the separate "model id" field in the dump list is redundant with the
locator, which is why the Round-Trip list can omit it without losing
information.

**But the locator is a hazardous sort key.** Because line/column coordinates
are permitted inside it, sorting on the locator makes dump order depend on
byte positions in the authored TOML: adding a comment renumbers lines,
changes locators, and reorders the dump — defeating the stability A7 claims.
This is a defect the corpus has not previously recorded.

## Corpus check on the empty-alphabet residual (A8)

- RDR 0001 never mentions `Outcomes` or "alphabet" — grep returns zero hits.
  Its `unmodeled_outcome` is defined as a *membership* test ("the recognized
  outcome is not modeled by the table"), never a well-formedness test on the
  alphabet itself.
- JDR 0001 §JD-10 left recognized-tag totality open ("whether a declared
  `recognized` tag is total … or partial"); RDR 0002's refine closed it in its
  own voice via the recognized-totality clause.
- No document obliges the **kernel** to validate its own `Table.Outcomes`.

So the A8 residual is genuinely unassigned: RDR 0002's normalizer is the sole
barrier against an empty-string alphabet, by default rather than by decision.

## Corpus check on fixture conformance (decisive for the fixture question)

RDR 0008 (Final) audited RDR 0002's spike fixtures **directly and by path**,
and treated their canonical status as the reason enforcement reaches them:

> *First clause* … **three violations at HEAD** … `0002-…/evidence/spikes/
> rdr-fixture.toml:25` and `kata-fixture.toml:22` (both `[tags.outcome]`),
> and `0003-…/evidence/spikes/guard-fixture.toml:36` (`[tags.rewind_target]`).
> RDR 0002 Load-Bearing Decisions names these "the canonical examples
> implementation tests must promote" … so the rule's first enforcement lands
> on the peer's own worked examples.

0002's two fixtures were repaired at 37a5bdf (the `[tags.recognized]` rename).
**RDR 0003's `guard-fixture.toml:36-37` still declares `[tags.rewind_target]`
with `provenance = "recognized"` and is still unrepaired** — a live breach of
a Final rule, belonging to RDR 0003 (Draft), not to this stage.

Precedent established: this corpus **does** hold RDR 0002's spike fixtures to
the RDR's own normative rules. The undeclared-tag defect found this stage
(F6) is the same class as the one RDR 0008 caught and required fixing.

Counter-precedent on urgency: the **shipped kernel test fixture** has the same
class of defect. `internal/resolve/fixtures_test.go::escapeRow` (`:249-259`)
sets `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}` and
`RequiresOwned: []string{"status"}` on a row with a non-empty `Escape` list —
breaching RDR 0009's Final precondition and RDR 0007 A21. RDR 0007 records
this as A25, a known collision to be repaired when implementation sequences
0007 and 0009, and explicitly declines to settle the ordering. So
non-conformant fixtures pending a Final rule's implementation are an accepted,
recorded state in this project — they are repaired at implementation, not
treated as lock blockers.

## Corpus check on NextTags (new unhomed-producer finding)

`NextTags` appears in RDR 0008 and RDR 0009 only — **zero times in RDR 0002**,
which owns normalization into the kernel row. The kernel `Row` carries both
`NextTags []Tag` ("the next state the row transitions to") and `Writes []Tag`
("the owned-tag writes the accessor layer applies") as *separate* fields
(`resolve.go:186-190`), and RDR 0009 A4 settled that the two are distinct and
that only `Writes` reaches the accessor.

RDR 0002's normalized-row field list says "writes including rendered
`<clear>` entries" and nothing about `NextTags`. So no document tells the
normalizer how to populate `NextTags`. This is structurally the same
unhomed-producer defect JDR 0001 §JD-3 recorded for `RequiresOwned` — which
0002's refine has since closed with an explicit producer clause. `NextTags`
has no such clause.

## Corpus check on RDR 0007's filed duties (both now discharged)

RDR 0007's open-item list carries: "Two duties are added to RDR 0002's
Refinement Context Direction list, which currently carries neither … (a) the
normalizer emits the kernel's exported existence token and its two boolean
literal forms (A16); (b) canonicalization of authored tag-key spellings is
bound as a normative clause (A22)."

Both are now discharged by the refine:
- (a) — 0002's existence-atom clause requires emitting `OpExists` /
  `LiteralTrue` / `LiteralFalse` verbatim and rejecting any other existence
  literal at load.
- (b) — 0002's tag-key identity clause binds exact byte equality with no
  folding.

Note (b) discharges A22 by a *different mechanism* than A22's named plan
requested. A22 asked 0002 to canonicalize (fold case/namespace/whitespace);
0002 instead requires exact match everywhere and makes the declaration key the
canonical spelling. This is **stronger, not weaker**: A22's stated failure
mode was "if 0002's normalizer ever matched declarations case-insensitively,
two authored spellings would resolve to one declaration yet reach the kernel
as distinct `Tag.Key` strings." Exact-match lookup makes that unrepresentable.
RDR 0007 should be told its A22 can be re-verified against 0002's clause, but
that is 0007's re-lock, not this stage.

## Corpus sweep — has any of this already been decided?

Four questions swept across `docs/rdr/0001`–`0009`, the postmortems,
`docs/jdr/0001`, and `docs/rdr/cluster-reconcile/`.

### Sort-tuple totality — DECIDED in RDR 0007, except the suffix rank

RDR 0007 §Normative Contracts (`:1567-1584`) already reasoned this exact
problem for its refusal payload and states the answer as a *total* tuple:

> The sort key MUST be TOTAL over payload entries. Row identity then key is
> NOT total… The ordering is therefore the tuple `(RuleID, SourceLocator,
> key, block, operator token, literal)`, compared field by field in that
> order.

So (i) field order, (ii) byte comparison, and (iii) locator participation are
all settled precedent RDR 0002 should **mirror rather than invent**. Only
(iv) — how an absent expansion suffix ranks against a present one — is
unaddressed anywhere in the corpus (`expansion suffix` appears twice, both in
0002).

The two Final peers disagree in a way that matters here, each locally correct:
0007 calls its six-tuple "total by construction" (atoms within one row);
RDR 0009 A8 proves `(RuleID, SourceLocator)` is **not** total for distinct
rows and shows collisions already in the frozen suite. 0009 escapes by
*collapsing* equal identities — a remedy unavailable to a dump that must emit
every row. So RDR 0002 inherits 0009's problem without 0009's remedy, and must
supply an identity that is genuinely total.

### Empty-string alphabet — DECIDED by 0002, but 0008 (Final) says otherwise

RDR 0002's recognized-totality clause forbids the empty string in the
alphabet. **RDR 0008 (Final) §Normative Contracts (`:1228-1237`) asserts the
opposite about 0002's own text:**

> A table whose declared `outcomes` alphabet contains the empty string would
> pass the alphabet gate … and resolve with no reserved key in the view;
> **nothing in this RDR or RDR 0002 forbids that alphabet entry. It is left
> unforbidden rather than silently assumed away.**

Two things follow. First, RDR 0008 independently discovered the **same kernel
path** this stage's A8 verification found, traced it to the same symbols
(`Table.models` as a plain `slices.Contains`; `assemble` injecting only for
non-empty `Input.Recognized`), and made a deliberate choice to leave it
unforbidden. Second, 0008's factual claim about 0002 was true when written and
is now stale — 0002's refine added the forbidding clause. 0002 is Draft and
repairable; 0008 is Final and this project does not amend RDRs.

Note also that 0002's clause does **not** actually close JD-10. Forbidding
`""` from the alphabet leaves untouched what JD-10 actually asks — whether a
declared `recognized` tag is total when no outcome is in flight. RDR 0008
§"Boundary against JDR 0001 §JD-10 (open)" (`:1318-1336`) is explicit that its
own blocks "do **not** settle it and must not be read as doing so." JD-10
remains open corpus-wide.

No document obliges the kernel to validate its own alphabet. Validation is
assigned solely to 0002's normalizer, by default rather than by decision.

### Spike fixtures as normative artifacts — no convention; 0002 is the outlier

RDR 0002 is the **only** RDR in the corpus that promotes on-disk
`evidence/spikes/` artifacts to canonical status. The corpus norm runs the
other way:

- RDR 0003 calls its own fixture defective (`:1567-1580`) and bars
  self-citation (`:2521-2523`).
- RDR 0007 stamps its harness encoding "vector-local, non-normative"
  (`:1756-1759`).
- RDR 0004 disclaims its spike as "fixture convenience, not the contract."
- RDR 0009 makes its *in-document MVV scenarios* normative (`:1561-1569`) —
  not spike files.

Yet RDR 0008 (Final) audited 0002's fixtures **by path** and required the
`[tags.recognized]` rename precisely *because* 0002 designates them canonical:
"the rule's first enforcement lands on the peer's own worked examples." That
rename landed at 37a5bdf. Three conformance defects survive it: no lifted
outcome field, a write-bearing escape row, and two undeclared written tags
(`rewind_scope`, `prelock_lens`). RDR 0008's own re-run
(`0008-…/evidence/spikes/a9-normalization.md:34-37`) independently confirms
"none produces an outcome field."

So the promotion is self-undermining as written: 0002 promotes the fixtures to
canonical *and* forbids what they contain. If they are normative they are a
second spec surface needing validation against the first — an obligation no
document states.

### 0007-before-0002 implementation order — OPEN, and the precedent runs the other way

Nothing in the corpus fixes a 0007-before-0002 order. Decisively, the
cluster's only explicit cross-RDR implementation-ordering prerequisites run
**0002-first**:

- RDR 0008 §Prerequisites (`:2090-2092`): "RDR 0002 implementation underway
  (the enforcement point is its normalizer's load/lint path; **this RDR's
  checks land inside that work, not before it**)."
- RDR 0009 (`:1345-1353`): the authored path "is closed by RDR 0002's
  normalizer… it becomes an *authoring-time* guarantee only when RDR 0002 is
  built."

RDR 0002's own Technical Environment reinforces it: "The kernel this RDR
produces rows for **is implemented**" — 0002 targets RDR 0001's shipped
kernel, not 0007's reshaped one.

RDR 0003 supplies the exact drafting precedent for booking the residual
dependency without blocking lock (`:2117-2123`): "RDR 0007 is `Final`, so the
*specification* is settled — but its kernel reshape is unimplemented… Phase 1
cannot begin against the atom-slice shape until that reshape lands, so **this
item gates implementation sequencing, not lock**."

RDR 0002 carries no equivalent item. Its Prerequisites gate on lock status
only ("RDR 0007, 0008, 0009 are Final") and never draw the consequence,
despite the same table distinguishing `Final` from `Implemented` for RDR 0001.

### Method note on A7

RDR 0001 (`:161-165`) — the same rubric 0002 carries — requires that "any
exactness claim such as … canonical, deterministic, or stable order must be
covered by a Critical Assumption Evidence Record or by the Minimum Viable
Validation." A7 is stamped `Design Decision`, which is a scoping choice being
*made*, not a property being *verified*. Under the rubric that is admissible
only if the MVV covers the ordering — which Testing Strategy scenario 5 does.
The MVV coverage is therefore load-bearing for A7 and must be named on its
Evidence line, not left implicit.

---

# Author's round — outcome (2026-08-22)

One consolidated round, four items, all four settled by the author.

1. **A7's witness** — rebuild the spike to the current contract rather than
   citing a stale one or downgrading A7. Done: `main.go` now lifts the
   outcome, derives `RequiresOwned`, renders escape rows write-free, retains
   per-atom block as a field, emits the existence constants, and sorts on the
   contract's stated keys. A7's Method moves `Design Decision` → `Spike`.
2. **Row sort order** — mirror RDR 0007's field-by-field drafting but
   **exclude the source locator**, because 0002 permits line/column detail
   inside it and ordering on it would churn the dump when unrelated source
   text is edited. Identity tuple is `(model id, rule id, expansion suffix)`.
   The divergence from 0007's tuple is stated in the clause with its reason.
3. **A8 residual** — record on A8 that this RDR's load-time rejection is the
   sole barrier against an empty-string alphabet, note RDR 0008's stale
   "nothing in … RDR 0002 forbids" sentence as superseded, and note that JDR
   §JD-10 stays open. No peer edits.
4. **Sequencing + `NextTags`** — add a Prerequisites item in RDR 0003's idiom
   ("gates implementation sequencing, not lock"), stating that it does *not*
   invert the cluster ordering, plus a normative clause giving `NextTags` and
   `Writes` a producer.

## Normative fixtures approved this round

- `evidence/spikes/output.txt` (SHA-256
  `c4be7447a241fb724632c53a9e5f39c7a9b5c7cc78e0a1e74ae4132271432a1b`) — the
  expanded-table value for the RDR and kata models. Named on A7 and in
  Testing Strategy scenario 2.
- `evidence/spikes/negative-cases.txt` — nine load-time refusals, the
  key-order-independence pair, and the three-run determinism transcript.
  Named in Testing Strategy scenarios 3 and 5.

## Defect found while rebuilding, and closed in the same pass

The totality argument for the new row ordering rests on rule ids being unique
within a model — which the RDR asserted nowhere. "Stable" and "must not be
reused for a different edge" are review guidance, not a load check. Added as
a normative clause, added to the validation category list, and witnessed by a
tenth spike negative case (`neg-duplicate-rule-id.toml` → `refused: duplicate
rule id "continue-prelock"`).

## Profile recount

`Seam Lineage: no prior accretion`, so no accretion floor applies and the
contract count governs. Independent load-bearing contracts: the sparse TOML
wire format; the normalization semantics that mint kernel rows (outcome
binding, `RequiresOwned`, next-state tags and writes, existence constants,
per-atom block); the expanded-table dump format and its total ordering; and
the validation category taxonomy. That is ≥2 independent contracts and a
format-plus-grammar-plus-taxonomy lock, which alone reads `large` — but the
consumers decide it: RDRs 0003, 0006, 0007, 0008, and 0009 all read fields
this RDR produces, and 0008/0009 name its normalizer as the enforcement point
their own checks land inside. That is a cross-RDR producer spanning modules.

**Profile: `large` → `foundational`.**

The ≥2-contract split signal is noted and deliberately not taken: the
contracts are not separable seams but one data format and the normalizer that
realizes it — splitting them would put the wire format and its only producer
in different documents, which is the coupling JDR 0001 §JD-3 and the
declaration-model rehoming were both written to repair.
