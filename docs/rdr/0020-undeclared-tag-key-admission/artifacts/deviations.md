# Deviations — 0020-undeclared-tag-key-admission

Unattended run (`/rdr-implement-triage`): a genuine author decision is
RECORDED here and the run continues to the completion gate as an open
item, per the unattended override. It is not a halt.

## D1 — REQ-40 docs obligation names a file that does not exist

- **Type**: SPEC-UNDER
- **Phase**: 0 (spec audit)
- **Gap**: REQ-40 (Implementation Plan Step 3) obligates a line in
  `docs/model-schema.md`. No such file exists in the tree. Meanwhile
  `docs/model-authoring.md:925` ships the assertion
  "**An undeclared tag whose value is an array is refused.**", which
  this record's C1 falsifies outright.
- **Evidence**: `docs/model-schema.md` absent (verified); the resources
  index names only `docs/cli-output-contract.md` + `CLAUDE.md` as
  authoritative contracts; predecessors 0005/0008 both treat docs
  clauses as obligations against the doc carrying the surface.
- **Resolution**: reading (b) — put the guidance in
  `docs/model-authoring.md` and correct its §"Passing a set-valued tag",
  whose HEAD text C1 falsifies. Reading (a) (create `model-schema.md`)
  would leave a shipped doc asserting the exact behaviour the RDR
  retires. The correction is owed either way under REQ-57 ("no shipped
  invocation changes meaning") and REQ-31.
- **Status**: needs author decision (proceeding under reading (b);
  evidence supports it but the filename literal in the locked record is
  wrong either way, so the author should confirm the record's intent)

## D2 — REQ-33 misspelled-key row carries no MVV obligation

- **Type**: IMPL-DECISION
- **Phase**: 0 (spec audit)
- **Gap**: Is REQ-33 (a misspelled key passes silently) a disposition row
  needing its own pinned test, or a recorded accepted cost?
- **Evidence**: the MVV's five steps and Testing Strategy's five
  scenarios never mention a misspelling; C1's normative text has no
  misspelling clause; the Premortem answers it with the `observed` echo
  (REQ-38), not a new check; its "Silent or loud" column is the only
  **silent** row.
- **Resolution**: treated as an accepted cost for the MVV bar — its
  behaviour is fully determined by REQ-1/REQ-31 (a misspelling IS an
  undeclared key). Pinned as OPTIONAL coverage, never counted as an MVV
  gate.
- **Status**: mechanical translation
