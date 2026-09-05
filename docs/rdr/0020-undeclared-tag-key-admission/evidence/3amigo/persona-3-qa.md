Model: claude-sonnet-5

# Persona 3 — QA / Tester

Widening note: I widened from the owned set (0020:S1-S4, 0020:MVV, 0020:C1) to 0020:§failure-modes and 0020:§risks-and-mitigations, because checking whether every documented behavior change has a corresponding pass/fail scenario requires comparing the owned scenario list against the full behavior surface the RDR itself claims — a silence in MVV/Testing Strategy has no line range to anchor to inside the owned set. I also grounded S1-S4 and C1 against the actual source (`internal/cli/flow_input.go`, `internal/cli/flow_next.go`) and the spike fixtures under `evidence/spikes/` to confirm every byte-exact assertion is either witnessed or explicitly marked unwitnessed.

## Findings

### LOW: array-carrier echo in 0020:S1 has no witnessed fixture for its exact wire form
- **Anchor:** 0020:S1
- **Finding:** S1's Expected clause asserts `observed` echoes the array-carrier value (`extras=["a","b"]`) "byte-for-byte," but unlike the scalar leg (normative fixture F1, read from `evidence/spikes/a2-undeclared-scalar-only.txt`), no spike fixture exists for the array leg — the RDR's own desk trace (0020:§desk-trace, MVV row 2) states this plainly: "the array leg is the only unwitnessed cell, correctly marked post-change." Grounding against `internal/cli/flow_next.go:45` (`Observed map[string]string`) shows the wire form is inferable — the array literal must serialize as the JSON-escaped string `"extras":"[\"a\",\"b\"]"`, matching the pattern already witnessed for the declared-set scalar-echo case (`labels` in F1). This makes the assertion writable by inference, not by a named fixture, so I can write the test but cannot cite a byte-exact precedent for the array leg specifically — only the scalar leg is pre-witnessed.
- **Blocks:** A test asserting the exact JSON serialization of an undeclared array-carrier value in `observed` cannot cite a pre-implementation fixture; it can only be justified by the `map[string]string` field type, one level removed from a direct run.

### LOW: no scenario exercises the documented "Silent" failure mode
- **Anchor:** 0020:§failure-modes
- **Finding:** §Failure Modes documents a "Silent" case: "a misspelled key (declared or not) passes as a carrier and the intended rule fails to match; resolution then refuses no-match or routes to an escape row." None of S1-S4 or MVV steps 1-5 exercises this path — there is no scenario that supplies a misspelled/undeclared key intended to guard a rule and asserts the resulting no-match refusal or escape-row routing. This is plausibly out of MVV scope (the failure mode is a *consequence* of the carrier design, not a new code path), but the RDR does not say so explicitly, and no Expected/pass-fail criteria exist anywhere in the owned scenario set for it.
- **Blocks:** A regression test confirming that a misspelled guard-intended key degrades to "no-match" or an escape row (rather than, say, panicking or silently matching a wrong rule) cannot be written from any stated Expected clause — the only guidance is descriptive prose in Failure Modes, not a scenario.

### INFO: 0020:S2 / 0020:MVV step 3 "identical...outcome" phrasing is testable only by cross-referencing the resolve JSON schema
- **Anchor:** 0020:S2
- **Finding:** MVV step 3 says "the selected rule and outcome are identical"; S2 expands this to the full field list (`rule`, `emit`, `gates`, `dispositions`, `next`, `writes`, `clear`, `escaped`). Neither element defines "outcome" in-line — it is only unambiguous once cross-checked against the resolve payload's actual JSON keys (`"outcome"` is the echoed `--outcome` flag value; `"rule"` is the matched rule id), confirmed via fixtures F1/F2. This resolved cleanly on inspection (not a gap), but the owned elements alone (without the fixture JSON) leave "outcome" one hop short of an unambiguous assertion, so I flag it as INFO rather than omit it.
- **Blocks:** Nothing currently — noted only because a tester reading S2/MVV in isolation, without opening the referenced fixture, could misassign "outcome" to a candidate-level field rather than the resolve payload's top-level echoed flag.

## Non-findings (grounded, resolved)
- 0020:S3 (declared scalar + array literal) and 0020:S4 (undeclared key, empty value) both cite pre-existing HEAD fixtures (F3, F4) with byte-exact message text; verified against `internal/cli/flow_input.go:714-737` and the spike fixtures — writable as stated, no gap.
- 0020:C1's ordering claim (empty-value arm hoisted ahead of the carrier branch, preceding any accessor) is consistent with the current code order in `flow_input.go` (duplicate check at :663-666, `canonicalValue` at :677, empty-value arm currently at :723) and the desk-trace's confirmed source-order row — testable via the same red-test scenario (S4) already specified.
- Refusal codes cited in C1 (`flow-tag-invalid`, `flow-tag-reserved`, `flow-tag-owned`, `flow-tag-duplicate`) all exist verbatim in `internal/cli/flow_input.go:36-39` — no invented or renamed code to chase.
