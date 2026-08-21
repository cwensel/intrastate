Model: claude-opus-5[1m]

# Critique iter-2 — grounding results

Code checks run against `main` before dispositioning. Delegated reads; verdicts
and durable anchors only. These are the authority for the grounding gate on the
rows below — a row whose cited behavior is refuted here does not edit the draft.

## G-1 — kernel collapses absent vs. unread (CONFIRMED, with a correction)

`internal/resolve/resolve.go::missingOwned` decides via
`internal/resolve/resolve.go::TagSet.has` (map presence + provenance equality),
producing `internal/resolve/resolve.go::KindOwnedStateUnavailable`. All symbols
the RDR cites exist as named.

**Correction the draft owes**: the collapse is *upstream* of the kernel. The
resolver's input is `internal/resolve/resolve.go::Input.Owned`, typed `[]Tag`,
and `internal/resolve/resolve.go::Tag` is a two-field Key/Value pair with no
third state. The distinction cannot survive the boundary type — the kernel does
not merely decline to recover it. A8's Evidence and the Load-Bearing "Absence
crosses the seam as omission" bullet both stop one layer short by saying the
kernel decides.

## G-2 — sentinel prohibition is not compiler-enforceable (CONFIRMED)

No accessor or executor package exists anywhere under `internal/` — the only Go
trees are `cmd/intrastate`, `internal/cli` (+ `clierr`, `config`, `respond`),
`internal/resolve`, `internal/version`. Every "accessor" occurrence in
`internal/` is doc prose. `TagSet`/`taggedValue` carry unexported fields with no
exported constructor, so an accessor's only contract with the kernel is the
`[]Tag` slice.

Consequence: `Tag{Key: "profile", Value: "<absent>"}` typechecks against
`Input.Owned` without complaint. The surface is greenfield, so the sentinel
prohibition must be carried by the RDR text and by tests — the type system will
not catch it. This is what makes the seam clause load-bearing rather than
decorative.

## G-3 — spike's derived key set is a bypassed fallback (CONFIRMED, narrows the row)

`…/evidence/spikes/main.go::expectedTagKeys` does derive the key set from
artifact contents, but it fires only when `len(requested) == 0`. All three read
scenarios pass explicit keys (`"status", "profile"`), so the circular path is
**not exercised** by them. A row claiming the derived path is "the only exercised
read path" is wrong on its stated basis.

## G-4 — absence masking via match-path + escape (CONFIRMED; the round's finding)

`missingOwned` operates over `internal/resolve/resolve.go::Row.RequiresOwned`, a
**separate hand-maintained declaration**, intersected with rows that survived
match filtering in `internal/resolve/resolve.go::gate`.

Masking chain: `Resolve` → `internal/resolve/resolve.go::TagSet.matches` (absent
owned key → row silently filtered, no refusal) → `gate` no longer sees the row
that would have declared the key → `missingOwned` returns empty →
`len(selected)==0` → `internal/resolve/resolve.go::escapeOrRefuse` → an escape
row whose `RequiresOwned` omits the absent key passes `gate` →
`Result{Plan, Escaped: true}`.

Two scoping facts that bound the finding:

- **This is specified behavior, not a defect.**
  `internal/resolve/adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`
  asserts the escape-plan outcome and fails if the kernel refuses. `gate`'s doc
  comment argues a pruned row's owned-state obligation is not part of the
  resolution. So "make the kernel refuse" is not an available fix — it breaks a
  locked peer's test.
- **`escapeOrRefuse` is itself gated.** An escape row that *does* declare the
  missing key in `RequiresOwned` raises `owned_state_unavailable` normally. The
  gap is not "escapes bypass the gate"; it is "the escape row need not declare
  the key at all."

**Residual sharp edge (unaddressed by this RDR or the kernel):**
`TagSet.matches` is provenance-blind and `RequiresOwned` is hand-maintained.
Nothing cross-checks that keys appearing in `Row.Match` or `Row.Guard` are also
listed in `Row.RequiresOwned`. An author expressing an owned-key dependency only
in `Match` gets no-match-then-escape instead of `owned_state_unavailable`. That
consistency obligation is pushed to RDR 0002 normalization and enforced nowhere.

Net for this RDR: the Disposition Table's summary "No input class exits
silently" is **false as written** — the guarantee holds only when the consuming
row declares the key in `RequiresOwned`. The table's own absent-key row already
carries the conditional ("a row requiring the key refuses…"); the summary
paragraph drops it.

## G-5 — "artifact unavailable" is unowned vocabulary (CONFIRMED)

`RefusalKind` is a closed set of exactly five, re-exported by
`internal/resolve/resolve.go::RefusalKinds`: `no_match`, `ambiguous_match`,
`owned_state_unavailable`, `guard_unevaluable`, `unmodeled_outcome`. There is no
artifact-unavailable kind; `KindOwnedStateUnavailable` is about a key absent from
an already-produced snapshot, not about the artifact being unreadable.

`internal/cli/clierr` has **no named code constants at all** — codes are bare
string literals (`command-error`, `config-not-found`, `config-read-error`,
`flag-invalid-value`). The only typed unavailability notion is the *group*
`internal/cli/clierr/clierr.go::GroupEnvUnavailable` (exit 3), scoped to external
facilities, not artifacts; `ExitCodeFor` collapses `GroupUserEnv` and
`GroupInternal` to exit 2.

Net: Failure Modes lists "artifact unavailable" as a typed refusal, but the RDR
introduces that vocabulary rather than reusing it, and the closed five-kind set
is RDR 0001's. The Disposition Table already routes the same input class to
`execution_failure`, so the two sections disagree.
