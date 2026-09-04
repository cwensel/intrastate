# Triage — RDR 0028 declared-line-edit-writer

Single-pass roborev triage over the implementation branch.

- Window: `BASE` 790bff5 .. `HEAD` cea7294 (frozen at Phase 0)
- Batch label: `batch:rdr-0028`
- Spine: 7 in-window per-commit auto-review jobs (6774, 6779, 6783, 6786,
  6789, 6794, 6801), 18 findings. Plus a bounded one-round sweep of the
  fix-commit reviews (6806, 6807), 2 findings.
- Every finding was grounded against the RDR, `deviations.md` and current
  source by a sub-agent that PROBED rather than read the review text.

## Verdicts

| # | Job | Finding | Verdict | ODC | Outcome |
|---|---|---|---|---|---|
| F1 | 6779 | category fail-fast order after 0025's six | DROP rdr-adjudicated | design-conformance/design-conformance | C1.4 fixes REGISTRATION order; `edit_tag_argv0` and the edit set are disjoint by carrier (REQ-73) |
| F2 | 6779 | `anchor`/`replace` coerced to "" and accepted | **FIX-NOW** | checking/boundary | `80f7d5e` |
| F3 | 6779 | parsed template segments discarded | DROP over-engineering | n/a-not-a-defect/design-conformance | C1.2 forbids re-scanning runtime DATA, not re-parsing the authored template; both parses call one parser |
| F4 | 6779 | `EditAnchorProbe` empty substitution | DROP superseded-at-HEAD | n/a-not-a-defect/boundary | premise false: RE2 accepts a leading quantifier as literal text |
| F5 | 6774 | dump-carriage assertion tautological | KATA-BUG | test-oracle/test-coverage | `aj7z` |
| F6 | 6783 | `Artifact.Context` unwired in production | DROP superseded-at-HEAD | interface/design-conformance | `artifactMap` builds and attaches the context |
| F7 | 6783 | read-back reader argv not checked pre-`Apply` | **FIX-NOW** | checking/logic-flow | `0dc6219` |
| F8 | 6783 | unplanned sibling rules skipped | DROP superseded-at-HEAD (D10) | algorithm/logic-flow | Phase 3c |
| F9 | 6783 | apply-time refusals at exit 3 | DROP superseded-at-HEAD | assignment/design-conformance | both carriers wrap the sentinel |
| F10 | 6783 | MVV helper one-key shape vs its own fixture | KATA-BUG | test-oracle/design-conformance | `d8kp` |
| F11 | 6794 | reader refusal reports as a line-edit write refusal | KATA-BUG | interface/backward-compat | `j7jw` |
| F12 | 6789 | ADV-3 oracle accepts any error | KATA-BUG | test-oracle/test-coverage | `mqwx` |
| F13+F14 | 6786, 6801 | doc claims a uniform rule-scoped subject | KATA-BUG (consolidated) | documentation/doc-code-drift | `sr1e` ★ |
| F15 | 6774 | re-anchor identity test is cardinality-only | KATA-BUG | test-oracle/test-coverage | `mqwx` |
| F16 | 6774 | atomicity counts no renames; rename-failure never induced | KATA-BUG | test-oracle/test-coverage | `mqwx` |
| S1 | 6806 | contract clause names in user-facing errors | **FIX-NOW** | documentation/doc-code-drift | `bb290c6` |
| S2 | 6807 | single argv walk inverts precedence step (1) | **FIX-NOW** | checking/logic-flow | `bb290c6` |

★ carries an open question.

## The two that mattered

Both FIX-NOW findings from the main pass were data-safety defects that a
green suite did not catch.

**F2** — a rule omitting `anchor` or `replace` was accepted, and both
omissions destroy data silently. An empty anchor compiles to the
everything-matcher, so on a one-line file it selects that line with
cardinality 1 (the ambiguity check never fires) and rewrites it. An empty
template is the empty LINE, not "no write". Both returned a nil error.

**F7** — with the command gate ON, the read-back reader's argv placeholders
were never examined before `Apply`. The edit landed, the reader then refused
during read-back, and the caller was told the write may have been applied
and was not verified — for a pure defect of the request. C1.3's `order:`
forbids exactly that. Deviation D11 does not cover it: D11 reasoned over the
gate-OFF path, where no child is spawned and these refusals are unreachable.

Each fix is mutation-verified: neutralizing it makes the new test fail.

## Notes

- Every DROP rests on a probe of current source or a clause quote, not on
  the review text. Four of the six drops are stale findings from per-commit
  reviews of commits later fixed on the same branch.
- Three findings were confirmed by SURVIVING MUTANTS — a cardinality-only
  re-anchor, a per-rule rename loop, and a stripped dump surface each left
  the full suite green. Those are false coverage claims, not live defects,
  and are filed together as `mqwx`.
- F11 is a regression introduced by Phase 3c's own third-site fix, which
  widened the request marker's reach past the edit writer.
