
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
