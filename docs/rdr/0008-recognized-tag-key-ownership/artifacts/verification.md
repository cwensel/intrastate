# Verification — RDR 0008 Ownership of the recognized-outcome tag key name

Two independent passes, recorded in the order they were written. **Phase 3a**
(CoVe, below) worked from the RDR and `req-list.md`; **Phase 3b** (adversarial,
at the end of this file) worked from the record's Failure Modes section. Neither
read the other's findings or Phase 1's test suites. They converged on the same
three seams and **disagree on whether those seams violate this RDR** — see
"Three seams probed as candidate violations" for both readings side by side.

---

## Phase 3a — CoVe verification (REQ-derived probes)

Stage 8 Phase 3a, independent CoVe verification — 2026-08-26.

**Method.** For each REQ in `req-list.md` (106 REQs across `0008:C1`–`0008:C6`,
`0008:MVV`), a violating input was derived **from the record's spec text only** —
the Phase 1 test files and `coverage.md` were not read — and run against the
shipped implementation through throwaway probe harnesses in `internal/resolve`
and `internal/table` plus the real CLI binary. The probes were deleted before
this file was written; nothing they added survives in the tree.

**Verdict: PASS.** No REQ was observed violated. No `FAIL-N` entry.

---

## Probes run, and what each would have caught

### Kernel — `internal/resolve` (blocks 1, 4, 5, 6)

| Probe | REQs | Violating input sought | Observed |
| --- | --- | --- | --- |
| P1 | 1, 6, 74 | a view that binds the outcome under any other key | `assemble` binds `recognized` → `("round-clean", ProvenanceRecognized, true)`, `Len()==1`, read through the **real kernel**, not two literals |
| P2 | 3 | a `recognized` key present with no outcome in flight | absent; `Len()==0` |
| P3 | 36, 38, 42, 43, 46, 50, 71, 85–87, 89, 90 | any of 8 breach variants (owned / observed / duplicate-in-one-sequence / empty-`Recognized` / `RequiresOwned` / doubly-breaching) resolving instead of erroring | all 8 breach: non-nil error from **both** call sites, zero-valued `Result` (no `Plan`, no `Refusal`), single error, no new `RefusalKind` |
| P6 | 53, 55, 88, 99 | a conforming `Input` gaining an error | `Input{}`, outcome-only, and an input carrying `Recognized`/`RECOGNIZED`/`" recognized"`/`"recognized "` as owned+observed keys all return nil from both sites |
| P7 | 51, 96 | the documented kernel-side residual made unreachable | `missingOwned` still returns `["recognized"]` on the package-internal path; residual intact |
| P11 | 76–78 | guard-side and matcher-side views diverging in one resolve | guard atom over `recognized` received `"round-clean"` while a `Match` on the same key selected the row — same resolve, same binding |
| P12 | 4, 5 | this RDR refusing the empty-string outcome (over-implementation) | resolves to a plan, nil error, no `recognized` key — correctly left outside the obligation |
| P13 | 17 | a kernel refusal for a row matching under an innocent name | `no_match`, nil error — a data-level miss, never a kernel error |
| P14 | 53, 54, 94 | nondeterministic D3 precedence on the bypass path | `owned` > `observed` > `recognized` held identically over 100 repeats × 3 collision shapes; single outcome each |
| P15 | 66, 68 | `Resolve` running a second, independently-written check | **343** combinations of owned × observed × `RequiresOwned` keys: **0 disagreements** between `CheckInput` and `Resolve` |
| P16 | 41 | the predicate reading outside its stated three sequences | `Flow`, `Revision`, `Outcomes`, `RuleID`, `SourceLocator`, `Match`, `Guard`, `NextTags`, `Writes`, `Escape` all set to `recognized` → nil; only the three sequences move the verdict |
| P17 | 44, 45, 90 | the breach skipped because another condition fired first | detected alongside an unmodeled outcome, a nil guard seam, and a 0009-shaped escape+writes row; the check is first in `Resolve` and unconditional |
| MVV | 72 | either MVV kernel half failing | row matching on `recognized` fires against the outcome through the real kernel; owned / observed / `RequiresOwned` variants rejected at both call sites |

### Normalizer — `internal/table` (blocks 2, 3, and the advisory)

| Probe | REQs | Violating input sought | Observed |
| --- | --- | --- | --- |
| T1, T2 | 10, 20, 57 | the quoted form treated as a distinct name | `[tags.recognized]` and `[tags."recognized"]` both load clean — identical post-parse key |
| T3–T6, G8 | 11, 18, 25, 26, 29, 30, 91, 92 | a missing or wrong-direction payload field | both directions refuse in `reserved_tag_key` carrying all three fields: `kernel-owned` → offending as authored + remedy `recognized`; `author-must-rename` → offending `recognized` + remedy **the empty string** (never the reserved key) |
| T7, E1, E5, E7 | 19, 56, 57, 59, 60, 79, 80 | folding or trimming leaking into the identity rule | `Recognized`, `RECOGNIZED`, `" recognized"`, `"recognized "`, `ReCoGnIzEd`, `"\trecognized"` all stay **ordinary unreserved names** and each raises the advisory; a recognized-provenance declaration under a near-miss name still fails `kernel-owned`, proving the near-miss is not a name variant |
| T8, E6 | 80 | a blanket advisory | `[tags.result]` and `[tags.recognize]` raise none — targeted, not blanket |
| T11, T13, E8 | 61, 63, 81 | the advisory altering the load/lint verdict | a clean model + near-miss loads clean **and** advises; a refusing document still advises; `Load` and `LoadWithAdvisories` return identical verdicts on every fixture |
| E4 | 63 | a category discriminator on the advisory | `Advisory{Authored, Reserved, Rule}` — no category field; cannot participate in dispatch |
| T10 | 13, 84 | two `[tags.recognized]` reaching the naming rule | `malformed_toml` duplicate-key, before the rule — as specified |
| E2 | 12, 82, 83 | more than one failure, or an unstable category | 200 loads of two wrongly-named recognized declarations → **exactly one** failure every time, category `reserved_tag_key` stable, offending name varying (which is what REQ-82 leaves unspecified); no separate cardinality check exists in `load.go` |
| G1, G2 | 22 | a `reserved_tag_key` check bolted onto predicate positions | `[rule.guard.all.recognized]` / `[rule.guard.unless.recognized]` → `malformed_outcome_binding` (0002's), not this category |
| G3 | 23 | a write-position check added here | `[rule.write] recognized` → `write_to_non_owned_tag` (0002's) |
| G4 | 24 | the stray-predicate case closed here | `[rule.match.Recognized]` → `unknown_tag`; the intent channel stays open, as the record says it must |
| G5, G6, T12 | 106 (deviation D2) | a capability `keys` entry admitting the reserved key | all three of `[read.*]`, `[write.*]`, `[gate.*]` refuse under `malformed_accessor_binding` — 0002's category, which is what D2 asks for |
| G7 | 15, 16 | the lower bound re-homed into `reserved_tag_key` | a model with no recognized declaration → `malformed_model_declaration`; 0002's clause left where it is |

### CLI — `internal/cli` (§JD-8 / §D8 citations)

| Probe | REQs | Observed |
| --- | --- | --- |
| `--tag recognized=x` | 104 | refused with `flow-tag-reserved`, `GroupUserEnv`, exit 2 — and **before** any accessor runs (it precedes the `flow-artifact-missing` refusal on the same invocation) |
| `--tag` routing | 105 | `parseTags` output reaches only `Input.Observed` (`flow_resolve.go:109,221`, `flow_next.go:287`); no path routes it into `Recognized` or `Owned` |
| load failure envelope | 102, 103 | `reserved_tag_key` surfaces as `flow-model-invalid` carrying `findings[{code: "reserved_tag_key", …}]` — §D10's shape; `clierr.Finding` carries `Code`/`Message`/`Param`/`Locator`/`Hint` as §D10 fixes |

### Structural / negative REQs

- **REQ-69** (scope ceiling): the whole-package export sweep shows **exactly one**
  new exported symbol on `internal/resolve` — `CheckInput`. The `internal/table`
  additions (`Advisory`, `LoadWithAdvisories`, `Failure.{Offending,Remedy,Rule}`,
  the three rule-id constants) are the latitude REQ-34/64 grant.
- **REQ-40 / REQ-70**: RDR 0009 fixes its predicate as a **no-argument method on
  `Table`** named `CheckValid`/`Validate` (`0009:1000-1013`, `:1611`).
  `resolve.CheckInput(Input) error` is a package-level function over a different
  receiver-less subject; no 0009 symbol name is claimed.
- **REQ-67**: `Resolve(in Input) (Result, error)` — signature byte-identical to
  the pre-0008 tree.
- **REQ-48 / REQ-52**: `missingOwned` is untouched in the 0008 diff, and no
  source-lint rule over `requires_owned` was added.
- **REQ-49**: no `reserved_tag_key` failure is emitted anywhere for a
  `RequiresOwned` entry — the only `RequiresOwned` walk in `load.go` is 0002's
  writer-arity check.
- **REQ-65**: no sigil-guarded spelling (`_recognized`, `<recognized>`) exists in
  the tree; the canonical name is the shipped `recognizedTagKey`.
- **REQ-9**: the pointer comment on `internal/resolve/resolve.go::recognizedTagKey`
  cites RDR 0008 `0008:C1` and names `CheckInput` as the enforcement point.
- **REQ-53**: `RefusalKinds()` still returns exactly the five RDR 0001 kinds.
- **REQ-95 / REQ-98**: neither the A5 unreachability half nor a source-lint half
  was written as a test — correctly absent.

---


## Three seams probed as candidate violations — and why 3a did not file them

Phase 3a and Phase 3b ran independently and converged on the **same three
seams**. Phase 3b filed all three as ADV failures against the record's **Failure
Modes** section; Phase 3a judged each to fall outside what the **REQ set** binds
and filed none. The evidence is not in dispute — only which document governs.
Both readings are recorded here so the implementer resolves the disagreement
deliberately rather than by whichever artifact is read second.

**S1 — mixed-direction model reports a nondeterministic rule identifier.**
(3b's ADV-3.) Confirmed independently: 200 loads of identical bytes carrying
both a wrongly-named recognized declaration and an owned declaration named
`recognized` yield either `(outcome, "recognized", …/kernel-owned)` or
`(recognized, "", …/author-must-rename)` — opposite remedies. Cause:
`internal/table/load.go:163` ranges a Go map.
*3a's reading*: RDR 0002's fail-fast clause is explicit that "the order in which
independent defects are checked is deliberately **unspecified** … a document
tripping two categories MAY be refused with either" (`0002:774-782`), and RDR
0008 states no determinism obligation over the load payload — its only "same
input yields the same disposition on repeat runs" clause (TS-8 / REQ-94) governs
the **kernel's** D3 backstop, which P14 verified deterministic. REQ-33 fixes the
**category** as "the stable one" and the rule id as "the finer key"; the category
held across all 400 mixed- and same-direction loads.
*3b's reading*: REQ-27/33 presuppose that identical bytes yield an identical
payload, and REQ-82 licenses only a choice between two *same-direction*
declarations.
*Note*: 3b's remedy — sorting the declaration keys before the scan — is a
two-line change that satisfies both readings and forecloses the question. That
asymmetry is the practical argument for doing it regardless of who is right.

**S2 — the near-miss advisory has no production caller.** (3b's ADV-1.)
Confirmed independently: `table.LoadWithAdvisories` is called by nothing outside
`internal/table`'s own tests; `internal/cli/lint.go:94` calls `table.Load`. A
model whose only irregularity is `[tags.Recognized]` lints clean with no
advisory at any user-observable surface.
*3a's reading*: this is precisely the disposition Phase 0 recorded as **Q1** and
deviation **D5** ratified — the record disclaims the surface ("which user-facing
command surfaces it, and under which exit code, is RDR 0005's mapping decision
and is not settled here"), and REQ-63 forbids the advisory altering the verdict,
which rules out the failure envelope as its carrier. The table-side obligation
REQ-58..64 states is met.
*3b's reading*: REQ-58's "normative, not deferred" is discharged by assertion if
the channel reaches no one.
*Note*: both agree the table-side contract is satisfied and that the gap is a
**delivery** gap. It belongs to RDR 0005/0006, which own the clean-load output
contract — but until one lands, the mitigation reaches no author.

**S3 — the CLI drops the remedy name and rule identifier.** (3b's ADV-2.)
Confirmed independently: `loadFailure` (`internal/cli/flow_input.go:144-152`)
populates only `Code`/`Message`/`Locator`, leaving `Param` and `Hint` empty; the
separate `lint` verb refuses with `model-invalid` + a `detail` string rather
than §D10's `flow-model-invalid` + `findings[]`.
*3a's reading*: REQ-26/28 bind the three fields "**at the data level**", which
`table.Failure` satisfies, and REQ-93 is explicit that "the consumer is a test
stub, **not RDR 0005's exit-code map**: this RDR asserts the payload contract,
and 0005 owns whatever mapping it later adds."
*3b's reading*: with the direction absent from the wire, REQ-30's "a renderer
MUST NOT present the reserved key as the required name in this direction" is
unenforceable.
*Note*: S3 and S2 share a carrier. Whatever verb is chosen to deliver the
advisory is also the natural place to widen the finding.

---

## Coverage note

Every REQ carrying a testable obligation was exercised. The REQs that are
latitude grants (7, 34, 64, 70), negative scope fences (5, 12, 24, 52, 95, 98),
or citation repairs (102–105) were verified by confirming the implementation
does **not** do the forbidden thing, or that the cited site exists — not by a
positive assertion, which for those clauses would over-implement. Deviation D3's
record-side citation half (a documentation edit on
`docs/rdr/0008-recognized-tag-key-ownership.md`) remains owed and is outside
this phase, as Phase 1 and Phase 2 both flagged.

---

## Phase 3b — adversarial review (ADV entries)

Independent pass. Brief: the record's **Failure Modes** section, `req-list.md`,
and the implementation tree. Phase 1's `*_0008_test.go` suites and Phase 3a's
findings were not read. Three failure modes named, three tests added, **all
three fail against the implementation at `f120171`**.

### ADV-1 — the near-miss advisory is normative but unreachable

**Failure mode.** The Load-Bearing Decisions fix the advisory as
"**Advisory warning — normative, not deferred**" (REQ-58/59/61). The
implementation puts it on `table.LoadWithAdvisories`
(`internal/table/advisory.go:41`) — and **nothing outside `internal/table`'s
own tests calls that function**. Every real caller takes `table.Load`:
`internal/cli/lint.go:94` and `internal/cli/flow_input.go:121`. A model
declaring `[tags.Recognized]` therefore lints entirely clean with no advisory
at any surface a user can observe, which turns REQ-58's "not deferred" into a
statement about an unreachable function. REQ-61 makes the channel advisory
*precisely so it can be shown without blocking*; showing it nowhere discharges
the requirement by assertion. JDR 0001 §D10 (REQ-102) names the intended
carrier — the advisory rides `Finding.Hint`.

**Test.** `internal/cli/reserved_key_adv_0008_test.go`
`TestAdv0008_LintSurfacesTheNearMissAdvisory`. Loads a model whose only
irregularity is a `[tags.Recognized]` declaration, asserts via
`LoadWithAdvisories` that the fixture really is a near-miss and really loads
clean (REQ-63: the advisory must not alter the verdict), then drives
`intrastate lint --as=json` and looks for a finding carrying rule
`reserved-tag-key/near-miss`.

**Currently fails.** Yes — `findings: []`. The advisory reaches no user.

### ADV-2 — `0008:C3`'s payload is dropped at the CLI carrier

**Failure mode.** The Failure Modes section calls the breach "Visible" and
states "the failure data carries the direction-specific rule identifier, the
offending name, and the remedy name". REQ-28 sharpens it: "the guidance
travels in the failure data, not the renderer." `internal/table/load.go:165-184`
populates `Failure.Offending` / `.Remedy` / `.Rule` correctly — and then
`internal/cli/lint.go:94-103` maps every load refusal onto a bare
`model-invalid` `CLIError` with `Detail: err.Error()` and **no `Findings` at
all**; `internal/cli/flow_input.go:135-152` (`loadFailure`) does emit one
finding but populates only `Code`/`Message`/`Locator`, dropping the remedy and
the rule identifier into prose. A consumer cannot branch on
`reserved-tag-key/kernel-owned` vs `reserved-tag-key/author-must-rename`, so
REQ-30's "a renderer MUST NOT present the reserved key as the required name in
this direction" is unenforceable — the direction is not on the wire. REQ-103
(JDR 0001 §D10) pins `reserved_tag_key` as a `findings[]` code under one CLI
code, which `lint` does not emit.

**Test.** `internal/cli/reserved_key_adv_0008_test.go`
`TestAdv0008_LintCarriesReservedKeyPayloadOnTheWire`. Builds the kernel-owned
direction (a recognized-provenance declaration named `outcome`), confirms
`table.Load` refuses in `reserved_tag_key` so a green result could only mean
the CLI dropped the payload, then asserts the JSON failure line carries a
finding coded `reserved_tag_key` whose `rule` is byte-for-byte
`reserved-tag-key/kernel-owned` (REQ-27/29/92) and whose `hint` names the
remedy.

**Currently fails.** Yes — the wire carries
`{"code":"model-invalid", …, "detail":"reserved_tag_key: tag outcome declares
provenance recognized under another name"}` and zero findings.

### ADV-3 — the three-field payload is nondeterministic across the two directions

**Failure mode.** `internal/table/load.go:165` ranges `decls`, a
`map[string]TagDecl`, and returns on the first branch that fires. Go
randomizes map iteration per range. A model that breaches **both** directions
at once — `[tags.recognized]` declared `owned` *and* `[tags.outcome]` declared
`recognized`, both legal TOML, distinct keys, so neither REQ-13's duplicate-key
rule nor REQ-15's lower bound preempts — reports a payload chosen by the
runtime hash seed:

```
offending=recognized  remedy=""            rule=reserved-tag-key/author-must-rename
offending=outcome     remedy="recognized"  rule=reserved-tag-key/kernel-owned
```

The category is stable (REQ-31), but the three fields `0008:C3` adds are not.
REQ-27 makes the rule identifier "a comparable token… a golden test asserts it
byte-for-byte" and REQ-33 makes it the value for "remediation lookup" — both
presuppose identical bytes yield an identical payload. The two directions
prescribe **opposite edits**, so a flipping remedy is worse than none. REQ-82
(TS-5) licenses exactly one nondeterminism — which of *two same-direction*
declarations is named — and says nothing about a coin-flip between
contradictory remedies. Sorting the declaration keys before the scan fixes it.

**Test.** `internal/table/reserved_key_adv_0008_test.go`
`TestAdv0008_DoublyBreachingModelReportsOneStableDirection`. Loads the same
bytes 200 times and asserts the `(Offending, Remedy, Rule)` triple is
single-valued.

**Currently fails.** Yes — reproducibly two distinct payloads (e.g. 161/39,
150/50 across runs).

### Modes examined and NOT filed

- **REQ-35's `unknown_tag` extension is vacuous, not broken.** REQ-26 extends
  the payload to "any undeclared-tag failure whose offending key is the
  reserved name", and none of the four `CatUnknownTag` sites carries one. It is
  unreachable through `Load`: `loadTags` enforces the `[tags.recognized]` lower
  bound (`load.go:186-188`) before `accessorTable`, `atomsFromBlock`, and
  `loadInitial` ever run, so the reserved key is never undeclared at those
  sites. Recorded as latent; no test added rather than a vacuous one.
- **A10 (accessor-produced owned keys) holds.** `accessor/executor.go:133-147`
  `classify` builds every `KeyValue` from the *requested* key list, never from
  artifact-supplied text, so `OwnedSnapshot` cannot introduce a spelling of its
  own and no accessor-borne user data is reclassified as a programmer mistake.
- **`CheckInput` at `Resolve`'s entry is not preemptible.** It is the first
  statement of `Resolve` (`resolve.go:413-416`), ahead of `assemble` and
  `Table.models`, so REQ-44's "detected whenever present — never skipped
  because another precondition also fired" holds for every precondition that
  exists at HEAD. No order is encoded, per REQ-45 / JDR 0001 §JD-5.
- **Near-miss folding is correct.** `strings.EqualFold(strings.TrimSpace(k),
  "recognized")` fires on `Recognized`, `RECOGNIZED`, `" recognized"`,
  `"recognized "`, `"\trecognized"` and not on `result` — the disjunctive
  trigger REQ-59/60 require.
