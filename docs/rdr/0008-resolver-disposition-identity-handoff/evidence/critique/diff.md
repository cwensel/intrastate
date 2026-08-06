Model: gpt-5.6-sol

# Critique fallback diff — RDR 0008

Variant: fresh-context single-model fallback; no alternate-model roster is
configured.

## Agreement

- First-pass `CRIT-1` and fallback `FALLBACK-3` agree that the current public
  `resolver.Input.Table []Edge` cannot enforce atomic identity/action ownership.
  The fallback tightens the failure window from forged construction to mutation
  between validation, matching, and `plan`.
- First-pass `CRIT-2` and fallback `FALLBACK-3` agree that the typed locator and
  its copy/validation semantics need a named boundary before lock.
- Both passes require production-boundary tests rather than a CLI lookup or a
  hand-assembled resolver fixture.

## Disagreement / unique signal

- `FALLBACK-1` finds a source-schema contradiction absent from the first pass:
  RDR 0002 forbids action fields on escape rules while RDR 0008 expects a normal
  successful plan with next tags.
- `FALLBACK-2` treats omission of `TableRevision` as a carrier defect. The
  decided RDR cluster instead gives RDR 0007 revision derivation and charts
  table/revision enforcement to its successor. The valid residue is a scope
  clarification: `SelectedRule` is logical identity within a revision-bound
  input, not a standalone historical event id.
- `FALLBACK-3` shows A2's prior spike became stale when the locator changed from
  string to an RDR 0002-owned typed value and the design began promising a
  snapshot-owned table.

## Resolution summary

- **fixed** — `FALLBACK-1`, traced to origin `CRIT-1`: A6 now blocks lock until
  RDR 0001/RDR 0002 and the production parse-to-CLI path choose one modeled
  escape disposition/action meaning.
- **fixed** — `FALLBACK-3`, traced to origin `CRIT-1`/`CRIT-2`: A2 is Pending;
  the table must own caller storage, hide row replacement, return copies from
  inspection, and pass pre/post-resolution mutation plus race tests.
- **dismissed-with-cite** — the part of `FALLBACK-2` requiring `TableRevision`
  inside `SelectedRule`: RDR 0007 `Load-Bearing Decisions / Identity` owns the
  revision and its `Existing Infrastructure Audit / Replay revision claim`
  charts association enforcement to a successor. RDR 0008 now states the
  carrier's revision-scoped meaning explicitly.
- **fixed** — `FALLBACK-3` error-boundary residue: malformed authored identity is
  a load/config error mapped by RDR 0005; an impossible invalid value inside a
  validated table is an internal programmer error.

Charted: revision-bound historical audit identity and enforcement of the
normalized-table/revision association remain with RDR 0007's explicitly charted
enforcement successor; adding `TableRevision` to `SelectedRule` is outside RDR
0008's logical identity handoff.
