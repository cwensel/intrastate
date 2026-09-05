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

## D3 — the Diagnosis clause's echo is absent from the refusing exit

- **Type**: SPEC-DEFECT
- **Phase**: 3b (adversarial review)
- **Gap**: the record's Failure Modes section accepts the silent-typo cost
  on the strength of one compensating control: "the resolve payload's
  `observed` field echoes every carried key byte-for-byte — the stray
  spelling sits beside the declared keys **in the same envelope the
  refusal rides**." The Silent clause's outcome is a REFUSAL (no-match) or
  an escape route, but `observed` is a field of the SUCCESS payload only.
  On the refusing exit the operator gets an envelope naming the DECLARED
  key as absent and carrying no trace of the stray key they typed.
- **Evidence**: verified against the built CLI — `clierr.CLIError` is
  `{code, message, param, detail, hint, findings}`
  (`internal/cli/clierr/clierr.go:50`); `Observed` exists only on the
  resolve payload (`internal/cli/flow_resolve.go:53`). Reproduced by
  `TestAdv3_0020_TheStraySpellingDoesNotRideTheRefusalEnvelope`
  (`teir=free` vs declared `tier`): refusal envelope keys are
  `[code message findings]`. `TestAdv3b_...` pins the OTHER exit, where
  the echo IS present and carries the typo verbatim — the echo is real,
  it is just on the wrong exit.
- **NOT an implementation defect**: the admission seam does exactly what
  C1 says, and NO REQ in `req-list.md` obligates the refusal envelope to
  echo (REQ-38 asserts only that `observed` echoes every carried key).
  Adding the echo to the refusal envelope would be a NEW public output
  surface no Normative Contract names — additive-is-not-exempt, so it is
  not a Phase 3c fixup and was not added silently.
- **Resolution**: the ADV-3 test was inverted to PIN the shipped shape
  (it fails the moment the refusal envelope gains the echo, which is
  exactly when this deviation is resolved). The record-level remedy —
  amend the Failure Modes claim, or specify the echo on the refusal
  envelope in a follow-on RDR — is the author's call. Routed to triage as
  an rdr-seed candidate.
- **Status**: needs author decision
