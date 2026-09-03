Model: claude-opus-5[1m]

# 3-amigos iter-2 — delta verification of the iter-1 rewrites (RDR 0026)

Scope: did the nine rewrites land coherently, did they create new contradictions, and
what elsewhere in the record is now stale. Not a fresh full review. Record read via the
projector only.

Verdict: **BLOCK** — 3 blocking findings, 6 non-blocking.

---

## Per-item disposition of the nine rewrites

| # | Rewrite | Landed | Notes |
| --- | --- | --- | --- |
| 1 | C1 `precedence:` — held = reported terminal condition | YES | agrees with S9 companion, `authority`, `trace` step 7, D-selection-predicate, C1 `whole:` |
| 2 | C1 `precedence:` — per-read re-arm | PARTLY | contradicted by the Illustrative Code (D-2) |
| 3 | held-pipe error type moved to `internal/accessor` | YES (claim true) | source confirms; A8's line cites drifted (D-4) |
| 4 | Phase 4 inserted, old Phase 4 → Phase 5 | YES | no stale four-phase reference; Phase 1 `readBounded` claim correct; one cite drift (D-6) |
| 5 | F5 pollability check declared NEW + Detail line | NO | duplicated clause + self-contradiction (D-1); `disposition` row stale (D-3) |
| 6 | +100 ms tolerance single-sourced in C1 `bound:` | YES | S8 and MVV row 3 both cite rather than restate; one soft spot (D-8) |
| 7 | S3 baseline + new scenario 3a | PARTLY | Detail form agrees with C1/MVV row 5; numbering is not element-addressable (D-5) |
| 8 | F6 → accepted trade-off | PARTLY | reasoning sound and consistent with JC1 and C1 `class:`; two stale echoes (D-7, D-9) |
| 9 | A9, A10 added | YES | well-formed, non-duplicative, cover items 2 and 5 |

---

## Findings, severity-ranked

### D-1 — BLOCKING — `0026:F5` contains a duplicated clause AND contradicts its own rewrite

F5 reads:

> `cmdbind` has no logger and gains none (its imports carry no `log`/`slog`), so "logged"
> is not a log line: the condition is carried on the invocation and surfaces as a `Detail`
> line on any refusal that invocation produces, which is what a test asserts on. On that
> host the join stays unbounded — on that host the join stays unbounded — the old
> hang, confined to a non-pollable host and **logged at the check** rather than discovered.

Two defects in one sentence:

1. The clause "on that host the join stays unbounded" appears **twice**, back to back,
   separated by an em dash. This is a mechanical artifact of the rewrite (the new sentence
   was spliced in front of the old one without deleting the overlap).
2. The trailing "logged at the check" is the *old* wording the rewrite exists to remove.
   F5 now explicitly says `"logged"` is not a log line and there is no logger — then, one
   clause later, says the condition is "logged at the check". The passage refutes itself
   inside a single sentence.

The rewrite is the fix for the iter-1 defect ("logged" was unassertable). It did not land.

### D-2 — BLOCKING — `0026:§illustrative-code` still shows a ONE-SHOT deadline, contradicting the new per-read re-arm in `0026:C1` `precedence:`

C1 `precedence:` now commits:

> The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain, so
> the grace bounds the IDLE GAP between reads and never the size of the tail […] A single
> absolute deadline would instead bound the whole remaining tail against the clock, which
> turns a large buffered payload into a false held-pipe report.

The Illustrative Code shape is exactly the mechanism C1 now names as the failure:

```go
if !waitBounded(&drains, WaitDelay) {
    deadlineIn(outR, errR, DrainGrace)  // a FUTURE deadline: the final read still runs,
                                        // so buffered bytes drain and only an empty pipe
                                        // with a live writer reports the deadline
    drains.Wait()                       // now returns; each drain reports EOF or deadline
}
```

`deadlineIn` is called ONCE, from the parent, from outside the drain goroutine, before
`drains.Wait()`. That is a single absolute deadline of `now + 50 ms` covering the entire
remaining tail — the shape MVV row 6 (1 MiB, drain started late) is now asserted to
survive, and which C1 says produces a false held report. Under this shape MVV row 6 fails.

Re-arming per read is only possible from **inside** the drain loop (`readBounded`), which
means the shape needs a second seam the code does not show — and Phase 1 already says
`readBounded` gains a reported terminal condition, so it is the natural home, but the
Illustrative Code does not reflect it.

Compounding: the inline comment "only an **empty** pipe with a live writer reports the
deadline" is the precise predicate C1 `precedence:` was rewritten to DEMOTE ("'Empty, with
a writer the group signal could not reach' describes the common case; it is not the test").
The comment also contradicts the S9 trickle-writer companion, which is never empty and must
still refuse. Both the mechanism and the predicate in this block are the pre-rewrite ones.

### D-3 — BLOCKING — `0026:§conditional-mini-checks` `disposition` table still carries the pre-rewrite F5 wording

Row:

| Non-pollable pipe (`os.ErrNoDeadline`) | unbounded join — the old hang | **logged AT THE CHECK** | none | **loud at the check, by design** (F5, S5) |

Both the "Event / error" cell and the "Silent vs loud" cell restate "logged", which F5 and
S5 have replaced with "carried on the invocation, surfacing as a `Detail` line". This is
the exact stale passage the delta brief anticipated, and it was not updated. It is the
`authority`-class table — the canonical statement of the observable — so an implementer
reading the table alone builds a log line that no test can assert on, which is the iter-1
defect reproduced.

The "Artifact minted: none" cell is also now wrong: the rewrite mints a `Detail` line.

### D-4 — MAJOR — `0026:A8` line citations drifted from source and from `0026:§phase-4-...`

Verified against `internal/accessor/executor.go`, `internal/accessor/model.go`,
`internal/cli/flow_exec.go` on the working tree:

| A8 claim | Actual | Phase 4 says |
| --- | --- | --- |
| six `r.applied = true` sites at `:363, :375, :393, :401, :409, :472` | `363, 375, 400, 409, 423, 472` | `:363, :375, :400, :409, :423, :472` (correct) |
| `errors.As` match takes the `if err != nil` arm at `:368` | `if err != nil` is at `:367`; `:368` is the `refusalOf` return | `:368-371` (spans it, acceptable) |
| `ClassTimeout` + phase check at `:398-404` | `399-405` | `:399-405` (correct) |
| `ClassReadBackIncomplete` at `:410-414` | `419-424` — `:410-414` is the `ClassExecutionFailure` arm | `:419-424` (correct) |
| `Applied()` accessor at `:323` | `323` (correct) | `model.go:324` (off by one) |

Three of A8's five cites are wrong; Phase 4 (the newer text) is right on four of five. Two
passages now assert different line numbers for the same three symbols. The `ClassReadBackIncomplete`
cite is the damaging one: A8 uses `:410-414` to argue "the read-back arm already sets
`detailMayHaveApplied`", but `:410-414` is the arm that provably does NOT — A8's own
decisive evidence points at the wrong lines. The underlying claims are all TRUE (verified);
only the anchors are wrong.

Note: the substantive item-3 claim is **confirmed**. `internal/cli/cmdbind/cmdbind.go:44`
imports `github.com/cwensel/intrastate/internal/accessor`; `go list -deps ./internal/accessor`
returns only `internal/resolve`, `internal/table`, `internal/accessor` — no `cli` package.
`accessor.ExecError` is at `internal/accessor/model.go:423` and `executor.go:83` already
`errors.As`-es it. C1 `refusal:`, the `authority` "Held-pipe error type" row, the
`disposition` "Held pipe, deadline NOT elapsed" row and Phase 1 all agree the type lives in
`internal/accessor`. Item 3 landed; only A8's anchors are stale.

Minor rider: C1 `refusal:` cites the import at `cmdbind.go:42`; the import line is `:44`
(`:42` is `"time"`).

### D-5 — MAJOR — scenario `3a` is not element-addressable, but two passages cite it as `S3a`

The projector's element census for 0026 lists `S1 … S9` and no `S3a` — scenario 3a is
rendered as a free paragraph between list items 3 and 4, outside the ordered list, so the
projector never registers it. Meanwhile the `oracle` table has a row keyed **`S3a (stderr
held alone)`** and the Testing Strategy body refers to it. A citation that does not resolve
through the projector cannot be checked by the gate or cited by a downstream record, and
`rdr inspect --select 0026:S3a` will not resolve.

Separately, the numbering is a form break: every other scenario is a plain integer in one
continuous list. Renumbering 3a→4 and shifting 4-9 → 5-10 would fix addressability but
would renumber existing ids; appending it as `10` preserves them. This is a decision, not
a mechanical fix.

The *content* of 3a is sound: the single-pipe `Detail` form ("the bare pipe name") agrees
with C1 `refusal:` ("the pipe(s) held ('stdout'/'stderr'/both)") and with MVV row 5's
both-pipes form `"stdout, stderr"` in that fixed order. No contradiction.

### D-6 — MINOR — two live citations for `readBounded`'s line range

`0026:§illustrative-code` and the `0026:§existing-infrastructure-audit` "Bounded,
overflow-detecting read" row both cite `cmdbind.go:352-356`. `0026:§phase-1-bounded-join`
and the `authority` "Held-vs-EOF-vs-overflow" row cite `cmdbind.go:352-364`, as does
`0026:D-selection-predicate`. Source: `readBounded` spans 352-364 (signature at 352, closing
brace at 364); the `:352-356` window stops mid-function, before the `io.Copy(io.Discard, r)`
that D-selection-predicate calls "the third exit". `:352-364` is right; `:352-356` is stale
and should be aligned.

The substantive Phase 1 claim is **correct**: `readBounded` today is
`func readBounded(r io.Reader, limit int) []byte` — no error return — so "gains a reported
terminal condition […] today it returns `[]byte` and NO error" is accurate, and consistent
with the `authority` table's "`readBounded` returns only `[]byte` today and reports nothing"
and with the Infra-Audit "SWALLOWS the error" row.

### D-7 — MINOR — the note closing `0026:§conditional-mini-checks` still calls F6 "under-explained"

> That is F6, recorded as **under-explained rather than contradictory**.

F6 was rewritten from "under-explained" to "**Accepted trade-off, not a gap**". The closing
note is the last surviving instance of the demoted framing, and it is doing load-bearing
work — it is the record's answer to why there is no CONTRADICTION row in the `trace` table.
It should now read "recorded as an accepted trade-off rather than a contradiction".

### D-8 — MINOR — the `oracle` table's MVV-row-3 and MVV-row-5 cells still carry bare bounds

Item 6's single-sourcing mostly held: C1 `bound:` states the +100 ms figure once, `0026:S8`
cites "the C1 `bound:` total plus its stated scheduling tolerance", and MVV row 3 says
"asserted with the C1 `bound:` tolerance". No passage restates "100 ms" a second time —
**the figure is single-sourced.** Two residues:

- `oracle` row for MVV row 3: "asserts elapsed ≤ 3 s AND class `timeout`" — no tolerance
  reference, so the table states a bare inequality the MVV itself qualifies.
- `oracle` row for MVV row 5 and MVV row 5 itself: "within `2·WaitDelay`" / "within
  `2·500 ms`" with no tolerance clause, while C1 `bound:` says the tolerance applies to the
  committed return bound generally.

Neither restates a rival figure, so this is a completeness gap in the tables rather than a
second source of truth.

### D-9 — MINOR — `0026:D-the-drain-grace-...` does not acknowledge the tension F6 now admits

F6's rewrite is sound and consistent with `0026:JC1` (the sub-reason constraint is homed in
JDR 0003 §D1, so declaring it here would force a 7.1 demotion — matching the memory-recorded
never-declare-in-both rule) and with C1 `class:` (which now itself carries the parenthetical
"per 0025:C4 a `timeout` refusal carries no Err/Detail, so the held-pipe reason survives only
on `execution_failure`"). The new `§consequences` Negative bullet matches. No contradiction.

The tension is acknowledged in ONE direction only. F6 names it explicitly — "the stderr tail
`D-the-drain-grace` cites to justify `DrainGrace` [is] unreachable in exactly that case".
`D-the-drain-grace` still argues the grace's whole justification is "the stderr tail that
`wrap` carries into `Detail` — the text that tells a user WHICH helper held the pipe",
with no note that the record's headline scenario never renders that text. A reader arriving
at the decision first gets an unqualified justification.

Mitigating: the justification is not *entirely* unreachable — the `execution_failure` leg
(MVV row 5, S3a: child exits 0, grandchild holds) does render the tail, and S3a was added
precisely so "the grace's load-bearing justification has [a] test". So the decision survives;
it is the one-line cross-reference back to F6 that is missing. Non-blocking.

---

## Items that landed clean (no finding)

- **Item 1.** C1 `precedence:`'s new predicate is consistent everywhere it matters: `0026:S9`'s
  trickle companion (a byte every 10 ms, never empty, never EOF ⇒ REFUSED) is exactly the case
  the new predicate admits and the old one excluded; the `authority` table's canonical column
  reads "the drain's reported condition; NOT the clock, NOT the byte count […] overflow
  outranks held"; `trace` step 7 reads "held = the drain's reported terminal condition […]
  deadline error ⇒ held, whatever bytes the grace delivered"; `0026:D-selection-predicate`
  independently states overflow-wins with the same reason ("the bytes were read, which makes
  the pipe's terminal state not what the invocation turns on") — C1 uses near-identical
  wording, which is citation-consistent rather than drift since D- is the owner of the
  tie-break. C1 `whole:` (EOF ⇒ whole) is the complement and does not collide.
- **Item 2 (partial).** Setting aside D-2, the re-arm claim is consistent with MVV row 6
  (1 MiB late-but-unheld, now explicitly the failure a single absolute deadline would cause),
  the `fidelity` table's "Final drain after the timer, writer gone | whole remaining tail,
  however many reads it takes — the grace bounds the idle gap between reads, not the tail",
  `0026:S7`, `0026:S9`, and `0026:A9`. No passage other than the Illustrative Code implies a
  single absolute deadline. C1 itself no longer says "final read" as a singular event in the
  mechanism sentence — it says "each read of the final drain", and the surviving "final read"
  phrases (in the EOF/held sentence and in `trace` step 7) name the drain's *terminating*
  read, which is compatible.
- **Item 4 (structure).** No passage anywhere refers to a four-phase plan or uses "Phase 4"
  to mean the surface phase. `§consequences` correctly forward-references "which is what
  Phase 5 documents", and Phase 5 correctly claims to be "the record's only user-facing
  artifact". `docs/cli-output-contract.md` exists on the working tree (16 KB), so the named
  artifact is a real file to amend rather than one to create. Phase 4's behavioral claims all
  verify: `executor.go`'s `if err != nil` write arm returns `refusalOf(def, timeout,
  ClassExecutionFailure, err)` and sets no applied sense; `flow_exec.go`'s
  `ClassExecutionFailure` arm (410-412) is a bare `envErr` with no phase check and no
  `Detail`; `ClassTimeout` (399-405) does carry `at == phaseWrite` plus `detailMayHaveApplied`;
  `ClassReadBackIncomplete` (419-424) does set `detailMayHaveApplied`. `Applied()` has zero
  non-test callers. S1 correctly names the direct-write leg as the proof and the read-back
  leg as regression-only, matching C1 `refusal:`'s "one site changes, not two".
- **Item 7 (content).** S3's baseline is now falsifiable ("the existing `cmdbind` suite passing
  unmodified — no new golden is minted […] a diff in them is the failure") and agrees with the
  `oracle` row ("absence-of-change oracle; discriminated by the cap boundaries") and the
  `fidelity` row. The single-pipe Detail form is consistent with C1 `refusal:` and MVV row 5.
- **Item 9.** `0026:A9` and `0026:A10` are both well-formed: Status `Pending`, Method `Spike`,
  non-empty Evidence naming what is unmeasured and how, non-empty "If wrong". Neither
  duplicates A1-A8 — A9 is distinct from A2 (24 bytes) and A3 (delayed drain, unbounded
  grace), and A9's Evidence says so explicitly; A10 is distinct from F5's error-value
  confirmation because it covers the *ordinary*-path no-op, which nothing else asserts. Each
  covers the load-bearing claim the delta introduced: A9 covers the item-2 re-arm, A10 covers
  the item-5 probe. The `oracle` MVV-row-6 cell correctly flags "its 1 MiB tail under a
  re-armed grace is A9, unmeasured until that spike runs", so the Pending status is not
  contradicted by settled-fact prose (the `§assumption-verification` gate's status-consistency
  rule).

---

## Blocking set

D-1, D-2, D-3. All three are rewrite-mechanics failures rather than design defects: the
decisions are sound, but three passages still carry the pre-rewrite text and one of them
(D-2) is the code shape an implementer would build from.

D-4 and D-5 are major but do not block on their own — D-4's claims are all true against
source and only the anchors drifted; D-5 is an addressability/form decision.
