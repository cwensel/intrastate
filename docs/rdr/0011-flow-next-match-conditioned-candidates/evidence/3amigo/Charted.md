# Charted to successor — 3amigo, RDR cli/0011

Findings that were real but net-new scope. Each is dismissed from this loop and
recorded here rather than absorbed into the RDR (rdr-common: never edit the
current RDR to absorb net-new scope).

## 1. `unresolved` → `unknown` is an out-of-repo consumer migration with no migration note

**Origin**: PM F3 (ledger L10). **Why out of scope**: RDR 0011 decides the
payload shape; whether the project ships a documented migration note, a
deprecation window, or a version marker for `flow next --as=json` consumers is a
release-policy decision that spans every verb's payload, not this contract. The
RDR now states the cost honestly (Consequences: `--all` restores the candidate
set, not the field name/type), which is this record's obligation.

**Suggested successor**: an RDR (or kata) on CLI payload versioning policy —
whether `docs/cli-output-contract.md` should pin per-command payload fields as
well as the envelope, and what a consumer is owed when one is renamed. Note the
contract doc currently pins only the envelope, so a per-command field rename is
today an undocumented-surface change by construction.

## 2. `.rdr/resources.md` names a design doc that does not exist

**Origin**: dispatcher, during the grounding gate for the resolve half (not a
persona finding). **Why out of scope**: `.rdr/resources.md` lists `CLAUDE.md`
under *Design docs* as "Settled conventions (respond gateway, CLIError,
ValidateMode) — treat as contract", but no `CLAUDE.md` exists at the repo root.
The grounding gate directs every stage to check fixes against that file, so the
row silently resolves to nothing for every RDR in this consumer — a seam/index
problem, not a defect in 0011.

**Suggested successor**: a seam fix — point the row at the real location or drop
it. Worth raising with the driver.

*(Checked and NOT a finding: the JDR principle citations do resolve. `docs/jdr/0001-resolve-kernel-seam.md`
§Principles defines P1 "Missing artifact state is never masked", P2 "The kernel
refuses rather than guesses", P4 "A refusal names what was missing and who can
fix it", P7 "Pre-release, prefer the clean shape" — so the RDR's leans on P1 and
P7 are grounded. P1 and P2 independently corroborate this round's resolution of
the selection-site question.)*
