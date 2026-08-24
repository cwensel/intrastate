# RDR 0005 — Stage 4 research citations (scoped re-entry: A3, A4, A5, A6)

Scope from the Status qualifier: `Draft [revised from Final 2026-06-24;
re-verify A3, A4, A5, A6]`. A1 and A2 carry forward as already Verified
(A1 re-checked opportunistically against project source; see below).

## Accepted citations

| Assumption | Method | Durable anchor | Verdict |
| --- | --- | --- | --- |
| A3 | Peer RDR | `JDR 0001 §D2` — "**Resolved: (b).** Gate, then count."; "refuses `guard_unevaluable` if any is undecidable, and only then applies exact-one matching and escape reachability" | selection procedure owned by §D2 + kernel, not RDR 0002 |
| A3 | Peer RDR | `JDR 0001 §D4` — "**Resolved: (b).** The kernel enforces the guard domain" | guard presence/unevaluability decided in kernel |
| A4 | Peer RDR | `JDR 0001 §D8` — "`flow resolve` and `flow next` MAY run declared read accessors over explicit `--artifact` bindings and MUST assemble `Input.Owned` from them"; "The purity fence moves to where it belongs: `internal/resolve`, not the verb" | verb owns load→decide; fence at the kernel package |
| A4 | Peer RDR | `JDR 0001 §D9` — "**Resolved: (c).** … Only the selected row's `gate` list runs"; "**All gates on the selected row run; deny-overrides; every result reported.**"; "**Deny is not an escape class.**" | gate site after exact-one selection, before the plan |
| A5 | Source Search | `JDR 0001 §D10` items 1/3 — "**Exit 3 = the environment could not be consulted; repair it and re-run the same request unchanged.**"; "**Exit 2 = the request or the model is wrong, or the model said no.**"; "No new exit group (0005's A-block)."; "`CLIError` gains exactly one `omitempty` structured field … `Findings []clierr.Finding` … `Finding{Code, Message, Param, Locator, Hint}`" | exit mapping + single new field |
| A6 | Design Decision | `JDR 0001 §D11` — "**`--write k=v` ↔ the write block; a repeatable `--clear <key>` ↔ the rule-level `clear` list.**"; "**A `set`-kind key's `--write` value is a JSON array literal**" | grammar pinned upstream |
| A6 | Spike | `JDR 0001 §D13` — "a set crosses as its canonical JSON array — members sorted, duplicate-free, compact encoding"; rendered bytes in `../spikes/d13-canonical-set.out` | canonical byte form fixtures |

## §D10 item 6 (model entry)

"an explicit file path (non-normative `--model <path>`; 0002's loader takes
bytes) beside the config-resolved `--flow <id>`; mutually exclusive;
`flow-model-not-found` either way."

## Rejected / narrowed branches

- **RDR 0002 `Normative Contracts` as A3's selection-procedure citation** —
  rejected. §D2 owns the procedure in its own voice and requires 0002 to
  *restate* it; 0002 is a follower. A3's re-anchor to §D2 + the kernel is the
  correct one.
- **Treating `Finding{Code, Message, Param, Locator, Hint}` as a JDR-normative
  spelling** — narrowed. The JDR header marks Go identifiers non-normative
  ("the RDRs own exact contracts"), so RDR 0005 is the normative home for the
  spelling; the JDR fixes the *shape* (exactly one omitempty field) only.
- **`--plan <file|->` on `set-state`** — deferred by §D11 as a seed, not part of
  this lock; RDR 0005 already records the deferral.

## RDR-owned extensions beyond JDR 0001's non-normative table

Not contradictions — the JDR defers exact code spellings to the RDRs — but
recorded so a later sweep does not mistake them for JDR citations:

1. `flow-write-duplicate` — §D11 says "the same key under `--write` and
   `--clear`, or twice under either, is a duplicate refusal" but mints no code.
   RDR 0005 mints the spelling.
2. Set-kind `--tag` values as JSON array literals with `flow-tag-invalid` on
   kind mismatch — §D11 pins the array literal for `--write`; §D13 fixes the
   kernel-seam form. Extending it to `--tag` is RDR 0005's call and is
   consistent with §D13's "no third encoding exists".
3. `flow-model-not-found` covering the *arity* error (neither/both of `--flow`,
   `--model`) as well as non-resolution — a defensible widening of §D10 item 6's
   "either way".

## Negative results

- `negative: A6 element-escaping — no JDR/RDR text fixes JSON string escaping
  inside a canonical set member` (searched JDR 0001 §D13, RDR 0002 set-literal
  clauses). See the spike finding; carried to the author's round.

## Project-source verification (own-code domain)

Reuse audit over the `$RDR_ENV` paths — **no reuse finding**; the surface is
greenfield and the RDR builds anew correctly:

| Behavior the approach introduces | Already provided? |
| --- | --- |
| `flow` group / `next`/`resolve`/`read-state`/`set-state` | ABSENT — `internal/cli/root.go` registers only `newVersionCmd` |
| transition-model TOML loader | ABSENT — no TOML dependency in `go.mod`; `internal/cli/config::Load` is a generic config loader |
| accessor executor (read/gate/write) | ABSENT from product code (only an RDR 0004 evidence spike) |
| `Finding` type / `findings` field in `clierr` | ABSENT — `internal/cli/clierr::CLIError` carries `Code`, `Message`, `Param`, `Detail`, `Hint`, `Group`, `Cause` only |
| text-payload rendering through `respond.OK` | ABSENT (partial) — `internal/cli/respond::OK`'s `ModeText` branch emits only notes/warnings to stderr and drops `Success.Data` |

A5 source search, observed mapping in `internal/cli/clierr::ExitCodeFor`:
`GroupSuccess`/`GroupWarning` → 0, `GroupUserEnv`/`GroupInternal` → 2 (one
case arm), `GroupEnvUnavailable` → 3, `GroupSignalCancel` → 130, non-`CLIError`
→ 1. Matches the RDR's claim. `internal/cli/clierr::ErrorCode` exposes the code
for branching. `Findings []Finding` does not exist yet, as the RDR states.

A1 anchors re-checked and still resolving: `internal/cli/respond::Success`,
`::OK`, `::Fail`, `::ValidateMode`, `internal/cli/clierr::CLIError`,
`internal/cli::ExecuteAndEmit`. `internal/cli::newVersionCmd` still writes text
success with `cmd.Println` — the drift the RDR records is real.

Kernel: `internal/resolve::Resolve` is pure (package imports only `slices` and
`strings`; enforced by
`internal/resolve/resolve_test.go::TestReq11_KernelImportsNoCLIOutputOrPersistenceFacility`),
takes `Input{Flow, Table, Owned []Tag, Observed []Tag, Recognized string,
Guards GuardEvaluator}`, and `internal/resolve::RefusalKinds` is the closed five
(`no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`,
`unmodeled_outcome`), closed-set-enforced by
`::TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`. It never executes
accessors to fill missing owned state
(`::TestReq13_KernelDoesNotExecuteAccessorsToFillMissingOwnedState`).

**Naming hazard for implementation (not a defect in the RDR).** The kernel has
an unexported `gate` that is a *pre-selection guard filter* (prune-on-false,
then owned-state check, then undecidable check) delegating atom truth to the
`Guards` seam. RDR 0005's "gates" are the *post-selection* accessor gates of
JDR §D9 and have no counterpart in `internal/resolve`. Two mechanisms sharing
one word; implementation must not wire 0005's gate results into the kernel's
`gate`.

`Input.Recognized` is a `string`, not a `Tag` — the kernel injects it under the
reserved `recognized` key itself. Consistent with the RDR's `flow-tag-reserved`
refusal at the CLI boundary.

## Peer-RDR grounding (siblings as they stand today)

Nine of ten claims VERIFIED against sibling text:

| Claim | Anchor | Verdict |
| --- | --- | --- |
| A3 — 0002 defines the recognized-outcome alphabet + normalized rows with rule ids/locators | `0002::Technical Design` — "Recognized outcome alphabet: the closed set of outcome tags that the recognizer may emit for the flow"; "carries the source identity (`(model id, rule id)` plus any expansion suffix, and the source locator)" | VERIFIED (locator is derived, not authored) |
| A3 — 0002 relies on rather than owns selection; citation move to §D2 is correct | `0002::Normative Contracts` — "The selection procedure itself — gate-then-count, and which refusals are escapable — is **JDR 0001 §D2's and the kernel's, and MUST NOT be restated here**" | VERIFIED |
| A3 — 0003 makes guard facts symbolic tag-set predicates owned outside the CLI | `0003::Approach` — "make every guard a symbolic tag-set predicate that lint can reason about"; "This RDR's evaluator owns **value semantics over a present value** — nothing more" | VERIFIED |
| A4 — 0004 keeps read/gate/write distinct over caller-supplied roles | `0004::Normative Contracts` — "Every accessor definition MUST declare exactly one capability: read, gate, or write"; "Accessors MUST operate on caller-supplied artifact roles" | VERIFIED |
| A4 — 0004 owns read-back incl. non-owned tags unchanged | `0004::Normative Contracts` — "verify each planned owned tag for equality … and that observed and recognized tag values present before the write are unchanged" | VERIFIED |
| 0002 declares `[read.<id>]`/`[write.<id>]`/`[gate.<id>]` with `keys` | `0002::Technical Design` — "`keys` is the only binding between tags and accessors" | VERIFIED |
| 0002's loader takes bytes + source id, no file I/O | `0002::Existing Infrastructure Audit` — "takes **already-read bytes plus a source id** for the locator, not a filesystem path" | VERIFIED |
| 0002 defines a load-category taxonomy | `0002::Conditional Mini-Checks` — 25 categories | VERIFIED |
| 0008 defines `reserved_tag_key`, `recognized` is the reserved key | `0008::Normative Contracts` — "a data-level validation failure in the `reserved_tag_key` category, reported at table load/lint before any resolution" | VERIFIED |
| 0004 classifies gate timeout/exec-failure as an *environment* failure | `0004::Disposition Table` — "`exit code` is out of scope here — RDR 0005 owns the CLI mapping" | NOT SUPPORTED BY 0004 |

**Exit-3 provenance (A5).** RDR 0004 mints the refusal *classes* (`timeout`,
`execution_failure`, `incomplete_read`, `read_back_incomplete`) but explicitly
declines to assign exit semantics. RDR 0005 does not miscite it — A5 sources
exit 3 from `JDR 0001 §D10` — but 0004 supplies no *independent* support, so
§D10 is the sole authority for the exit-3 rule. Recorded so a later sweep does
not look for a second leg that was never there.

**`omitempty` — reconciled upstream, not a live contradiction.** `0006::Normative
Contracts` says "the `findings` field MUST NOT be `omitempty`"; RDR 0005 pins
"exactly one `omitempty` structured field". `JDR 0001 §D10` item 3 resolves it:
"The `omitempty` tension with 0006 is apparent only: a lint *failure* always
carries at least one blocking finding, so the field is never empty on
`CLIError`; 0006's 'empty list is the receipt' binds the *success* payload's own
non-omitempty `data.findings`. Recorded as a 0006 citation repair, not a
reopening." RDR 0005's text is consistent with that reading. Stale anchor on the
0006 side, not RDR 0005's defect.

## Open — carried to the author's round

**OPEN-1 `clierr.Finding` field set (blocks A5).** RDR 0005:419 pins
`Finding{Code, Message, Param, Locator, Hint}` and calls it "shared with RDR
0006's lint findings". RDR 0006 is **Final** and normatively requires strictly
more of that same type: `0006::Normative Contracts` — "Every blocking finding
MUST carry a stable code, model identity, severity, human-readable message, and
the source rule/context id or source span … A finding attributed to one guard
atom MUST carry that atom's `Key`, `Operator`, `Literal`, and `Block`; a finding
scoped to an escape population MUST carry the failure class." `0006::Technical
Design` names the concrete Go fields: "`Key`, `Operator`, `Literal`, `Block`,
and `Class` are `string`". `severity`, model identity, and `Class` have no home
in the five. JDR §D10 item 3 called the record "non-normative" and the change "a
0006 citation repair", but 0006's text is unamended and Final — so the five-field
spelling in RDR 0005 cannot be the whole type.

**OPEN-2 set-member JSON escaping (exactness edge on A6).** `JDR 0001 §D13`
fixes "members sorted, duplicate-free, compact encoding" but is silent on string
escaping inside a member. No RDR or JDR text fixes it (grepped `docs/rdr/*.md`,
`docs/jdr/*.md` for `SetEscapeHTML`/`EscapeHTML`/`u003c`). Both existing product
output paths use `json.Marshal`
(`internal/cli/clierr/clierr.go:135`, `internal/cli/respond/respond.go:183`),
which HTML-escapes `<`, `>`, `&`. Spike
`../spikes/d13-canonical-set.out` shows the divergence: `json.Marshal` renders
`["a<b","p>q","plain","x&y"]` where an
`Encoder`+`SetEscapeHTML(false)` renders `["a<b","p>q","plain","x&y"]` — not
byte-identical. RDR 0005's "copy-through and read-back equality are byte
equality" (`:501-502`, `:542`) holds only if one encoder is fixed. The spike
confirms copy-through IS byte-stable through `json.Marshal` end-to-end
(plan → `Tag.Value` string → payload → re-render all agree), so pinning
`json.Marshal` closes it.
