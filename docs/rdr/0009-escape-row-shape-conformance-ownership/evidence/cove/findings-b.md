Model: claude-opus-5[1m]
Lens: cove (pass B — adversarial/implementer angle)

Persisted by the dispatcher: this pass ran with a read-only toolset and returned its
findings in its result packet rather than writing the file itself. Content below is the
pass's own report, transcribed verbatim in substance.

## Findings

### FB-1 (refuted-claim, blocking)
**Anchor**: "internal/resolve/fixtures_test.go::escapeRow populates **`NextTags` only — it
sets no `Writes` at all**" (A3 Scope correction; restated in Context, Existing
Infrastructure Audit "Builder sets `NextTags` only", and Phase 2 "The `escapeRow` builder
needs no change").

FALSE at HEAD: `internal/resolve/fixtures_test.go:249-260` — `escapeRow` returns a `Row`
with `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}},` on line 257. Verified
against `git show HEAD:internal/resolve/fixtures_test.go`. The RDR's own A3 spike ledger
contradicts the RDR: `evidence/spikes/a3-fixture-conformance.md:35-36` — "removed the line
`Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}},`" and `:194-195` "must be
stripped alongside the `escapeRow` builder." The RDR inverted its own evidence.

### FB-2 (refuted-claim, consequence of FB-1)
**Anchor**: "Conforming the two override sites is an implementation step of this RDR
(Phase 2)" / "Bring the write-bearing escape rows into conformance by dropping `Writes` at
the two call-site overrides".

Not two sites — `escapeRow` is called at 16 sites (`adversarial_test.go:128,171,193,227`;
`fixup_test.go:43,79,101,125,215,216,284`; `resolve_test.go:539,553,554,566,598`), and every
row it produces breaches `len(Escape)!=0 && len(Writes)!=0`. Phase 2 as written leaves those
breaching rows, so the whole-table precondition would error on most of the frozen escape
suite. A3's own "Edits 2 and 3 matter" note (`a3-fixture-conformance.md:42-43`) says the
opposite of the RDR's Phase 2.

### FB-3 (refuted-claim)
**Anchor**: "`errors.Join` costs nothing (nil when all-nil; degrades to the single error's
own message for one breach)".

Half true, and the half that's false is load-bearing. `errors/join.go` — with one non-nil
arg `Join` still returns `&joinError{}`, not the bare error. Only `Error()` degrades
(`if len(e.errs) == 1 { return e.errs[0].Error() }`). The returned *value* is a wrapper, so
a single-breach `Resolve` returns a `*joinError`, not the typed breach error. Scenario 1's
`errors.As`/`AsType` recovery still works (`Unwrap() []error`), but a caller doing a direct
type assertion or `err == ErrEscapeShapeBreach` fails. The RDR never says whether the
single-breach return is wrapped or bare — implementer-facing silence with a testable
behavioral fork.

### FB-4 (internal-contradiction)
**Anchor**: Existing Infrastructure Audit row "Breach is diagnosable from the error text
(A7), not a new code" vs. the Normative Contract "the verb MUST wrap it into a
*clierr.CLIError carrying a stable Code" and Failure Modes "the stable `CLIError.Code`
distinguishes this breach" / "not 'read the error text.'"

The audit row is stale Propose-era text that survived the Resolve-stage revision (which the
Failure Modes bullet explicitly narrates: "An earlier draft recorded this as an accepted
diagnostic gap"). It also contradicts the typed-error contract "inspectable via errors.As /
errors.AsType without parsing message text." Three passages, one row wrong.

### FB-5 (internal-contradiction / unexpressible type)
**Anchor**: Technical Design "returns a nil `Result` and a non-nil error" vs. MVV scenario 1
"`Resolve` returns `Result{Plan: nil, Refusal: nil}`".

`resolve.go::Result` is a **struct**, not a pointer; `func Resolve(in Input) (Result, error)`.
"nil `Result`" is not expressible in Go. The MVV wording is the correct one (zero value).
Also unpinned: `Result{}.Refused()` returns `false` (`resolve.go::Refused`), so a caller that
checks `Refused()` before `err` sees a "success-shaped" zero value with `Plan == nil` — a
nil-deref waiting to happen. The RDR never states the error-first checking obligation.

### FB-6 (silence)
**Anchor**: "The conformance predicate MUST be exported by the kernel package as a
construction-time check … It MUST return error (nil when valid), and its name MUST follow
Go's error-returning convention — CheckValid or Validate" + "the exact type/sentinel
spelling is sharpened at Pre-Lock."

The receiver is never pinned: both cited precedents are **methods** (`pprof.Profile.CheckValid`,
`rsa.PrivateKey.Validate`), which would make this `func (t Table) CheckValid() error`. But the
predicate is stated as being over `Table.Rows`, while the per-row predicate is row-scoped, and
`Row`/`Table` have no existing *exported* methods (`rescues`, `models` are both unexported).
Scenario 7 ("The exported predicate called directly by a producer") is unwritable without this.
The sentinel's export status is likewise unstated — a lowercase package-level sentinel would
defeat `errors.Is` for external callers, which is the clause's stated purpose.

### FB-7 (refuted-claim, determinism / scenario 10 unsatisfiable as written)
**Anchor**: "The reported rows MUST be ordered by RowRef identity using the kernel's existing
compareRefs ordering, never by Table.Rows position" + scenario 10 "Expected: Identical
reported row sequence."

`resolve.go::compareRefs` compares only `RuleID` then `SourceLocator` — a **strict weak
ordering, not a total order over distinct rows**. Two distinct breaching `Row`s with equal
`(RuleID, SourceLocator)` compare 0. `slices.SortFunc` is explicitly **not stable**
(`SortStableFunc` is the stable variant), so for tied refs the output permutation is a
function of input position, not of the table value. Constructible from existing fixtures:
`fixup_test.go:79,101,125` all build `escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", …)`
— same identity, different `RequiresOwned`/`Writes`. Hand-built rows with zero-value
`RuleID`/`SourceLocator` (this RDR's own target population, per A5) collide trivially. Either
dedupe by `RowRef`, use `SortStableFunc`, or state that ties are collapsed.

### FB-8 (refuted-claim, minor)
**Anchor**: Consequences "`Resolve` already scans `Table.Rows` twice per call (the candidate
partition and `escapeOrRefuse`)" + Performance Expectations "`escapeOrRefuse` walks it again".

`resolve.go::escapeOrRefuse` runs **only** on the `len(selected)==0` and `>1` arms; on the
exact-one-match success path and on the `KindUnmodeledOutcome` early return (which returns
*before* the candidate loop) `Table.Rows` is scanned once or zero times. The "already twice"
framing overstates the existing baseline: for a conforming table that resolves successfully
the precondition is a 100% increase in `Table.Rows` passes, not a constant factor on an
existing double scan. The conclusion (still O(rows), no new cost class) survives; the measured
claim does not.

### FB-9 (silence / unverifiable premise)
**Anchor**: Load-Bearing Decisions, Nil-vs-empty: "an empty write set is an empty work set, so
`[]Tag{}` provably cannot produce the user-facing symptom (RDR 0004 applies only planned
writes)".

"Provably" is not checkable against source — RDR 0004 is `Final` but **unimplemented**, and
A4's own "Carried forward" text concedes this "rests entirely on RDR 0004's write-accessor
scoping, not on any kernel-local property." The kernel-local half is verifiable and should be
cited instead: `resolve.go::copyTags` returns `nil` for `nil` and a fresh zero-length slice for
`[]Tag{}`, so `len(Plan.Writes)==0` either way — scenario 3's assertion holds on kernel evidence
alone. Stating it as *proven* on an unimplemented peer overstates the grounding.

### FB-10 (grounding, confirmed-good — recorded so it is not re-chased)
The two cited line numbers resolve exactly at HEAD: `adversarial_test.go:229` and
`fixup_test.go:103` are both `escape.Writes = []resolve.Tag{{Key: "status", Value: "Escaped"}}`.
A1 confirmed (every `Resolve` return carries literal `nil`). A2's clierr claims confirmed
(`GroupUserEnv, GroupInternal: return 2`; non-`CLIError` → 1; `Unwrap` exposes `Cause`). A5
confirmed (`root.go` registers only `newVersionCmd()`). A7's field names confirmed (`Row`
`RuleID`/`SourceLocator`; `RowRef`). `errors.AsType` exists and `go.mod` is go 1.26.x, so the
`AsType` references are valid.

## Bottom line
FB-1/FB-2 are lock-blocking — Phase 2, the Existing Infrastructure Audit row, and A3's "Scope
correction" all rest on a claim about `escapeRow` that is false at HEAD and contradicted by the
RDR's own spike ledger. FB-7 is the second-most serious: the ordering contract is not
satisfiable as written.
