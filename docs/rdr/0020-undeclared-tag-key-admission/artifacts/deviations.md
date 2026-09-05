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
- **Status**: needs author decision → RESOLVED (reading (b): the
  obligation is satisfied against `docs/model-authoring.md`, which is
  corrected on this branch). Grounds for resolving without the author:
  this is not a live design fork. `docs/model-schema.md` does not
  exist, so reading (a) cannot be satisfied by editing anything; and the
  `model-authoring.md` correction is owed under REQ-57 ("no shipped
  invocation changes meaning") and REQ-31 REGARDLESS of which file
  carries the new guidance — a shipped doc asserting the behaviour C1
  retires is a defect either way. What remains is a naming error in the
  locked record's REQ-40/41, which under the never-amend rule is a
  post-mortem note, not an implementation choice, and is tracked as
  `intrastate#g3ks` so the record's wording is corrected at its own stage
  rather than silently absorbed here.

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
- **Status**: needs author decision → RESOLVED (deferred to a
  successor RDR; no code change on this branch). Grounds: the two
  remedies — amend the Failure Modes claim, or specify the `observed`
  echo onto the refusal envelope — are BOTH record-level, and neither is
  an implementation choice this launch may make. Adding the echo would
  be a new public output surface no Normative Contract names
  (additive-is-not-exempt). Under the never-amend rule the locked text
  is not edited here. The finding is preserved as a behavioural pin
  (`TestAdv3_0020_...`, inverted to fail the moment the echo appears)
  and routed to RDR authoring as `intrastate#fk4n`
  (`kind:rdr-seed`), which is where a design fork of this shape belongs. Nothing about D3 is outstanding
  *for the code* on this branch.

## Spun-off work (filed to kata, `batch:rdr-0020`)

Findings that are real but NOT this branch's to fix. Filed so the branch can
close honestly rather than carrying them as open decisions.

| kata | route | what |
|---|---|---|
| `intrastate#r41t` | `type:bug` `lifecycle:queued` — drainable | ill-formed UTF-8 in a `--tag` value collapses to U+FFFD in the `observed` echo, so two distinct values echo identically. Pre-existing and provenance-blind (ADV-1b proves the declared-scalar arm mangles identically), so NOT a 0020 regression — but 0020 makes `observed` load-bearing and widens which bytes reach the seam. Any fix must move BOTH arms together or it breaches `0020:G-cross-cutting`. |
| `intrastate#fk4n` | `kind:rdr-seed` — RDR authoring | D3: the silent-typo acceptance rests on an echo the refusal envelope does not carry. Fork: amend the diagnosis claim, or spec `observed` onto the error envelope (a CLI output-contract change touching every refusing verb). |
| `intrastate#g3ks` | `kind:rdr-seed` — RDR authoring | D1's residue: REQ-40/41 names `docs/model-schema.md`, which does not exist. Post-mortem note for 0020, plus an optional finalize-time check that Implementation Plan path literals resolve. |

`kind:rdr-seed` children are refused by `kata-ship` gate 4 by design — their
deliverable is an RDR draft, not a red/green code change — so a
`--label batch:rdr-0020 --drain` flight ships `r41t` only.
