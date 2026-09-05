# Post-Mortem: RDR-0020 An undeclared tag key is admitted as a pure carrier

Ledger-only. Opened at implementation-triage follow-up to record three drift
findings that converged on this record after it locked — a wording debt on
REQ-6, a refusal-envelope claim in Failure Modes, and a docs obligation naming
a file that does not exist. The remaining post-mortem sections are authored
later.

## Escaped-Defect Ledger

One row per finding, enriched with the two fields triage cannot assign. Rows
whose disposition is a record-level wording correction carry `n/a` odc columns
— there is no implementation defect to classify.

| Finding | odc-type | odc-trigger | Expected-catching stage | Precursor lens | Escape distance |
| --- | --- | --- | --- | --- | --- |
| REQ-6 / `0020:G-cross-cutting` claims a carried value is admitted "VERBATIM — byte-preserved, with no canonicalisation, folding, normalisation, or re-encoding — and echoes in the resolve payload's `observed` field as given". Not literally true for ill-formed UTF-8: the JSON-mode echo is serialized by `internal/cli/clierr/clierr.go::WriteJSONLine` (`json.NewEncoder` + `SetEscapeHTML(false)` + `Encode`), and `encoding/json` substitutes U+FFFD for each ill-formed byte, so the distinct carried values `a\xffb` and `a\xfeb` echo identically. Pre-existing and provenance-blind — `TestAdv1b_0020_...` proves the declared bare-scalar arm mangles identically, so this is NOT a 0020 regression — but `0020:C1` makes `observed` load-bearing for diagnosis and widens which bytes reach the seam. Recorded as ADV-1; `verification.md` already concludes it is "a bound on REQ-6's 'byte-preserved' wording that predates this change". | n/a | n/a | `4-resolve` | none — a literal "byte-preserved" claim is exactly what assumption verification should have checked against `encoding/json`'s documented U+FFFD substitution; grounding, 3amigo personas 1–3, repeatability and the spikes all assert it only for well-formed values | 2 stages (owed at `4-resolve`; caught post-lock at implementation-triage adversarial review as ADV-1) |
| The Failure Modes Diagnosis clause (REQ-38) rests the silent-typo acceptance on `observed` riding "the same envelope the refusal rides". It does not. `Observed` is a field of `resolvePayload` only (`internal/cli/flow_resolve.go`); the refusal envelope is `clierr.CLIError` = `{code, message, param, detail, hint, findings}` (`internal/cli/clierr/clierr.go`), which has no echo field. The Silent clause's outcome is a REFUSAL or an escape route — the two clauses describe different exits. On the refusing exit the operator sees the DECLARED key named absent, with no trace of the stray key they typed. Pinned by `TestAdv3_0020_TheStraySpellingDoesNotRideTheRefusalEnvelope` and `TestAdv3b_0020_OnTheEscapeExitTheEchoIsPresentAndCarriesTheTypo` — the echo is real, it is just on the wrong exit. Recorded as D3 (SPEC-DEFECT) / ADV-3. | n/a | n/a | `5-prelock` | none — `evidence/reconcile/reconcile.md` explicitly reviewed the Failure Modes carve-out and found "No drift against the Q1 ruling"; 3amigo's `consolidation.md` flagged only an adjacent gap on the same clause family (P3-L1, array-carrier echo wire form), not this one | 1 stage (owed at `5-prelock`; caught at Phase 3b adversarial review) |
| REQ-40 and REQ-41 obligate a line in `docs/model-schema.md`. That file does not exist and never has. The doc actually carrying the tag-authoring surface is `docs/model-authoring.md` §"Passing a set-valued tag", whose shipped text asserted the behaviour `0020:C1` retires ("An undeclared tag whose value is an array is refused."); corrected on the implementation branch by `d3d8eb6` to state that an undeclared tag is a pure carrier. Resolved as D1 (SPEC-UNDER) reading (b) — the obligation is satisfied against the doc that carries the surface; reading (a), creating `model-schema.md`, was rejected because it would leave a shipped doc asserting the exact behaviour the record retires. | n/a | n/a | `2-propose` | none — a path literal in an Implementation Plan should resolve against the tree before lock; grounding resolves citation anchors syntactically, not by filesystem existence, and no pre-lock lens checks path literals against the tree | 2 stages (owed at `2-propose`; caught at Stage 8 Phase-0 spec audit as D1) |

**Notes.**

- All three rows share one root cause worth naming: a locked record asserting a
  property of an artifact — an envelope's field set, a file's existence, an
  encoder's byte-fidelity — that was never mechanically checked against the tree
  or the standard library at any pre-lock stage. Each claim is checkable in
  seconds once someone thinks to check it; none was.
- A finalize-time check that Implementation Plan path literals resolve in the
  tree was considered and routed OUT of this repo. The RDR flow engine is
  external — `.claude/skills/rdr-finalize` is a symlink into the sibling `rdr`
  project — so the check belongs to that project, not to intrastate. Nothing is
  created for it here.
