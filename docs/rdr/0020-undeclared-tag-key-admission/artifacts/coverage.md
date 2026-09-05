# Coverage — RDR 0020 Undeclared `--tag` key admission

Phase 1 artifact. Column 1 is the REQ id from
`artifacts/req-list.md`; column 2 names the test that would FAIL if that
clause were broken. An EMPTY column-2 cell is an orphan mark.

Tests live in the source tree beside the surface they cover:

- `internal/cli/flow_carrier_fixtures_0020_test.go` — fixture corpus and oracles
- `internal/cli/flow_carrier_0020_test.go` — `0020:C1` carrier rule + closed refusal list
- `internal/cli/flow_carrier_empty_0020_test.go` — the hoisted empty-value arm, its guard, precedence
- `internal/cli/flow_carrier_mvv_0020_test.go` — `0020:MVV` end to end, plus fixtures F1–F3 and the disposition rows
- `internal/cli/flow_carrier_scope_0020_test.go` — seam scope, comments, docs

Red/green at authoring time: every REQ the record marks *(new)* is RED
against HEAD. Every REQ marked *(pin)* or *(move — behaviour-preserving)*
is GREEN, which is the correct state for a clause whose obligation is
"this must not change" — its failure mode is regression, not absence
(req-list standing note S2). The distinction is recorded in the
`state` column so no green row is mistaken for a tautology.

| REQ | Test | state |
| --- | --- | --- |
| REQ-1 | `TestReq1And31And49_0020_UndeclaredScalarAndArrayAreAdmittedVerbatim` | red (new) |
| REQ-2 | `TestReq2And6And26_0020_NoCanonicalisationFoldingNormalisationOrReencoding` | red (new) |
| REQ-3 | `TestReq3_0020_DeclarednessIsByteExactOnTheKeyWithNoCaseFolding` | green (identity rule already holds for near-miss keys at HEAD) |
| REQ-4 | `TestReq4And5_0020_DeclaredSetIsCanonicalisedAndTheCarrierIsNot` | red (new) |
| REQ-5 | `TestReq4And5_0020_DeclaredSetIsCanonicalisedAndTheCarrierIsNot` | red (paired with REQ-4 by construction) |
| REQ-6 | `TestReq2And6And26_0020_NoCanonicalisationFoldingNormalisationOrReencoding` | red (new) |
| REQ-7 | `TestReq7_0020_TheOnlyRefusalsForAnUndeclaredKeyAreTheClosedList` | red (the shape arms are reachable at HEAD) |
| REQ-8 | `TestReq8_0020_TheGrammarArmStillRefusesAndTheFirstEqualsSplitHolds` | green (pin) |
| REQ-9 | `TestReq9And29And36_0020_ProvenanceGuardsStillFireFirstForAnyValue` | red (the array-valued duplicate leg refuses on conformance at HEAD) |
| REQ-10 | `TestReq10And32And54_0020_AnUndeclaredKeysEmptyValueStillRefuses` | green (move — behaviour-preserving) |
| REQ-11 | `TestReq11And14And34_0020_TheKindArmIsReachableOnlyUnderARealDeclaration` | red (new) |
| REQ-12 | `TestReq12_0020_NoFlowTagUndeclaredCodeIsMinted` | green (pin — negative) |
| REQ-13 | `TestReq13_0020_TheModelSideClosedWorldStillRefusesAnUndeclaredTag` | green (pin) |
| REQ-14 | `TestReq11And14And34_0020_TheKindArmIsReachableOnlyUnderARealDeclaration` | red (paired with REQ-11) |
| REQ-15 | `TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind` | green (move) |
| REQ-16 | `TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind` | green (move — behaviour-preserving) |
| REQ-17 | `TestReq17And20_0020_ACarrierNeverReceivesTheSetConformanceMessage` | red (the malformed-array legs refuse at HEAD) |
| REQ-18 | `TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind` | green (the guard's witness; HEAD satisfies it accidentally, the hoist must preserve it deliberately) |
| REQ-19 | `TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind` | green (pin) |
| REQ-20 | `TestReq17And20_0020_ACarrierNeverReceivesTheSetConformanceMessage` | red |
| REQ-21 | `TestReq21And24And37_0020_AdmissionRefusalPrecedenceIsUnchanged` | green (pin + move) |
| REQ-22 | `TestReq22_0020_TheCarrierBranchSitsWhereTheConformanceArmsWouldHaveRun` | red (new) |
| REQ-23 | `TestReq23_0020_TheHoistedArmLandsAfterTheDuplicateMarkNotBefore` | green (pin — the anti-inversion rule) |
| REQ-24 | `TestReq21And24And37_0020_AdmissionRefusalPrecedenceIsUnchanged` | green (latitude bound — no splice position asserted) |
| REQ-25 | `TestReq25_0020_EitherImplementationShapeMustHoldBothBounds` | red (new) |
| REQ-26 | `TestReq2And6And26_0020_NoCanonicalisationFoldingNormalisationOrReencoding` | red |
| REQ-27 | `TestReq27And28_0020_TheWidenIngDoesNotReachWriteOrClear` | green (pin — out of scope) |
| REQ-28 | `TestReq27And28_0020_TheWidenIngDoesNotReachWriteOrClear` | green (pin) |
| REQ-29 | `TestReq9And29And36_0020_ProvenanceGuardsStillFireFirstForAnyValue` | red (paired with REQ-9) |
| REQ-30 | `TestReq30And38And57_0020_UndeclaredScalarsStillPassAndSitBesideDeclaredKeys` | green (pin) |
| REQ-31 | `TestReq1And31And49_0020_UndeclaredScalarAndArrayAreAdmittedVerbatim` | red (the change) |
| REQ-32 | `TestReq10And32And54_0020_AnUndeclaredKeysEmptyValueStillRefuses` | green (move) |
| REQ-33 | `TestReq33_0020_AMisspelledDeclaredKeyPassesSilentlyAsACarrier` | red (optional coverage per D2/Q2 — NOT an MVV gate) |
| REQ-34 | `TestReq11And14And34_0020_TheKindArmIsReachableOnlyUnderARealDeclaration` | red (paired leg) |
| REQ-35 | `TestReq35And55_0020_FixtureG_DeclaredSetEmptyKeepsItsConformanceMessage` | green (pin — the guard's witness) |
| REQ-36 | `TestReq9And29And36_0020_ProvenanceGuardsStillFireFirstForAnyValue` | red (paired with REQ-9) |
| REQ-37 | `TestReq21And24And37_0020_AdmissionRefusalPrecedenceIsUnchanged` | green (pin) |
| REQ-38 | `TestReq30And38And57_0020_UndeclaredScalarsStillPassAndSitBesideDeclaredKeys` | green (pin — the diagnosis leg) |
| REQ-39 | `TestReq39_0020_TheGuardedLookupCommentCitesTheCarrierContract` | red (new) |
| REQ-40 | `TestReq40And41_0020_TheAuthoringDocNoLongerRefusesUndeclaredArrays` | red (new — lands in `docs/model-authoring.md` per D1/Q1 reading (b)) |
| REQ-41 | `TestReq40And41_0020_TheAuthoringDocNoLongerRefusesUndeclaredArrays` | red (same obligation, Background's phrasing) |
| REQ-MVV | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (new) |
| REQ-42 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (MVV step 1, asserted not assumed) |
| REQ-43 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (MVV step 2 — the round-trip fidelity leg) |
| REQ-44 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (MVV step 3) |
| REQ-45 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (MVV step 4, adjacent to step 2 by construction) |
| REQ-46 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (MVV step 5) |
| REQ-47 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (steps 2–5 in one runnable test) |
| REQ-48 | `TestMVV_0020_UndeclaredTagKeyAdmission` | red (the undeclared/declared pair adjacent in ONE test) |
| REQ-49 | `TestReq1And31And49_0020_UndeclaredScalarAndArrayAreAdmittedVerbatim` | red (`0020:S1`) |
| REQ-50 | `TestReq50And51_0020_FixtureF1AndTheArrayLegsWireForm` | red (array leg; F1 scalar leg green) |
| REQ-51 | `TestReq50And51_0020_FixtureF1AndTheArrayLegsWireForm` | red (`map[string]string` wire form) |
| REQ-52 | `TestReq52_0020_FixtureF2PinsTheAbsoluteSelectionInBothRuns` | red (with-carriers leg) |
| REQ-53 | `TestReq53_0020_FixtureF3_DeclaredScalarGivenAnArrayLiteral` | green (pin — F3 is HEAD behaviour, now truthful) |
| REQ-54 | `TestReq10And32And54_0020_AnUndeclaredKeysEmptyValueStillRefuses` | green (pin — F4 byte-identical to HEAD) |
| REQ-55 | `TestReq35And55_0020_FixtureG_DeclaredSetEmptyKeepsItsConformanceMessage` | green (pin — fixture G byte-identical to HEAD) |
| REQ-56 | `TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind` | green (the negative test obligation, satisfied by the third row) |
| REQ-57 | `TestReq30And38And57_0020_UndeclaredScalarsStillPassAndSitBesideDeclaredKeys` | green (pin — no migration surface; whole-suite check confirms no non-0020 test turned red) |
| REQ-58 | `TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind` | green (pin — F4/E/G all three rows) |
| REQ-59 | `TestReq59_0020_ByteIdenticalMeansMessageStringEqualityNotAHash` | green (scoping — asserted by construction; no digest anywhere in the suite) |

## Orphans

None. Every REQ-1..59 and REQ-MVV carries a test.

## Notes on the green rows

The tautology question the phase gate asks is "is this test green because
the behaviour is already correct, or because it asserts nothing?". Each
green row above was checked against HEAD explicitly:

- **REQ-8, 10, 13, 19, 21, 23, 24, 27, 28, 30, 32, 35, 37, 38, 53, 54, 55,
  56, 58** — pin/move clauses. Their obligation is that the change does
  NOT disturb behaviour that already ships, and the normative fixtures
  (F3, F4, E, G) were read AT HEAD by the record's own spike. A green
  result is the required state; the test earns its place by turning red
  if the implementation disturbs the arm.
- **REQ-3** — HEAD's zero-decl lookup already carries a near-miss key, so
  the identity rule's observable consequence holds today. The test is a
  pin against a future implementation that added case folding or a
  fuzzy registry, both of which `0020:D-identity` forbids by name.
- **REQ-12** — a negative over a code that does not exist. Green is the
  only correct state; it turns red the moment `flow-tag-undeclared` is
  minted, which is the whole obligation.
- **REQ-59** — a scoping clause. Its content is that message-string
  equality, not a digest, is what "byte-identical" means here; the test
  asserts the two refusal sites share one message template and computes
  no hash, which is the clause satisfied by construction.
- **REQ-15/16/18** — the hoist. HEAD reaches all three rows of the
  three-way `value == ""` partition through `canonicalValue`'s kind
  branch, so the messages already agree. The obligation is that the
  hoist PRESERVES that agreement; the third row (declared SET) is the
  one an unconditional `value == ""` hoist breaks, and it is asserted on
  the full message string including its trailing space.

## REQ-MVV runner

The fixture model is `internal/cli/flow_carrier_fixtures_0020_test.go`'s
`carrierTableModel`, written to a file verbatim. `$IS` is a binary built
from the worktree at the implementing commit.

```sh
$IS flow resolve --model carrier.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' \
  --tag extra=plain --tag 'extras=["a","b"]' --as=json          # MVV 2
$IS flow resolve --model carrier.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' --as=json          # MVV 3
$IS flow resolve --model carrier.toml --outcome decide \
  --tag 'tier=["a","b"]' --tag 'labels=["security"]' --as=json   # MVV 4
$IS flow resolve --model carrier.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' --tag extra= --as=json  # MVV 5
```

## REQ-MVV output

MVV 1 — the fixture model declares scalar `tier` and set `labels`, and
declares nothing named `extra` or `extras` (`grep '^\[tags\.'`):

```
[tags.recognized]
[tags.tier]
[tags.labels]
```

MVV 2 — both carriers admitted, exit 0, `observed` byte-for-byte as given
(the array literal rides as the raw argument text, a JSON STRING — fixture
F1's shape extended by the array leg, REQ-50/REQ-51):

```
{"type":"ok","data":{"model":"carrier.toml","revision":"","observed":{"extra":"plain","extras":"[\"a\",\"b\"]","labels":"[\"security\"]","tier":"free"},"owned":{},"readers":[],"outcome":"decide","rule":"free","gates":[],"emit":{"plan":"basic"},"dispositions":{},"next":{},"writes":{},"clear":[],"escaped":false}}
exit=0
```

MVV 3 — the same invocation without the two carrier flags. Fixture F2:
`"rule":"free"`, `"emit":{"plan":"basic"}`, `gates`/`next`/`writes`/`clear`
empty and `escaped:false` in BOTH runs; the sole payload delta is `observed`
gaining `extra` and `extras`:

```
{"type":"ok","data":{"model":"carrier.toml","revision":"","observed":{"labels":"[\"security\"]","tier":"free"},"owned":{},"readers":[],"outcome":"decide","rule":"free","gates":[],"emit":{"plan":"basic"},"dispositions":{},"next":{},"writes":{},"clear":[],"escaped":false}}
exit=0
```

MVV 4 — the DECLARED scalar handed the same array literal MVV 2 carries.
Refused, now truthfully (fixture F3):

```
{"code":"flow-tag-invalid","message":"the tag `tier` is not set-valued; got the array literal [\"a\",\"b\"]","param":"tier"}
exit=2
```

MVV 5 — the undeclared key given an empty value, the one arm the carrier
does not escape (fixture F4, byte-identical to HEAD):

```
{"code":"flow-tag-invalid","message":"the tag `extra` was given an empty value","param":"extra"}
exit=2
```

All five steps demonstrably satisfied; `go test ./...` green at the same
commit.
