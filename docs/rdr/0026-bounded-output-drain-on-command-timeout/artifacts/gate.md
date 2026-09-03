# Finalization Gate — cli/0026

- **RDR**: 0026 — Deliver the command timeout even when a detached grandchild holds the output pipes
- **Slug**: `0026-bounded-output-drain-on-command-timeout`
- **Date**: 2026-09-03
- **Profile**: foundational (accretion floor; `floor-holds`)
- **Verdict**: **READY — PASS**. Locked to Final at this pass.

**Fence disposition (author judgement, recorded because the model has nowhere to
put it).** `--outcome fence` emits `stopped:overlap-uncited` on
`overlap_uncited=1+`. Both arms were FIRED, not synced, and both pairs are
INCIDENTAL — no cross-citation or JDR home is owed. Detail in
`evidence/tooling-pass/tooling-pass.md`; the short form:

- `0012 0026` on `internal/resolve` + `internal/table` — 0012 cites the two
  packages as an *import-direction* fact fixing where a test fixture type may live
  (`internal/table` imports `internal/resolve`, so `resolve` cannot import
  `table`), declaring no API or shape for either. 0026's use is a quoted
  `go list -deps ./internal/accessor` line proving `internal/accessor` has no
  `internal/cli` edge. Two records reading the same dependency graph from opposite
  corners for unrelated placement decisions.
- `0016 0026` on `read_back_incomplete` — `0016:C4` decides only *when the
  reader-resolution site emits* the class (fail closed on a multiply-bound role
  rather than first-match), reusing it as an already-defined surface, and
  `0016:F3` explicitly freezes its semantics as an unchanged fail-safe. 0026
  likewise declares the read-back leg unchanged and touches only the direct write
  leg's `execution_failure`. Both defer to RDR 0004 as the class's home; neither
  constrains the other.

**Why the fence fired at all, and why that is not a defect in it.** JC1 records
that at 0026's propose (2026-08-31) BOTH intersect arms returned `overlaps[]`
empty, and 0012/0016 were already open (both proposed 2026-08-28). So no peer
moved — 0026's own text grew these literals during Stages 4-6, as quoted tool
output and as a description of behavior this record explicitly declines to change.
The fence caught a real post-check accretion and correctly demanded the arms be
re-fired. They were, and the answer is that the new literals are evidence
quotations, not decisions.

This is a **fired-and-cleared** check, not a skipped one. It is recorded here
rather than in the model because `tags.overlap_uncited` is `provenance =
"observed"` with domain `["0","1+","unchecked"]` and the fence group takes no
caller-supplied disposition — unlike `blocker_class` on `return` or
`user_facing`/`locks` on `profile`. A correctly-fired check therefore has nowhere
to land its verdict. Same shape as `joint_check_home=unhomed` in §5 below, where
the flow's own rule is that the gate's written response vouches for what the tool
cannot resolve. Locked on that precedent, by the author's explicit decision.

Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` (PASS; `lint` exit 0,
`blocking=0 resolution=0`). Item 4, Cross-Cutting Concerns, is authored once in the
record at `cli/0026:G-cross-cutting` and is not copied here.

## 1. Contradiction Check

No contradictions between Research Findings and the Proposed Solution. The three
findings that could have conflicted are each carried through into the design rather
than contradicted by it:

- Research established Go's own precedence for this exact conflict — `WaitDelay`
  "bounds the time spent waiting on … a child process that exits but leaves its I/O
  pipes unclosed"; `awaitGoroutines` "forcibly closes their pipes and returns
  `ErrWaitDelay`" — i.e. bound wins, the parent closes, success with an open pipe is
  an error. The Approach takes that disposition unchanged and quotes it rather than
  paraphrasing. The one place it goes further than Go is the *refusal* (Go's `Output`
  hands back the bytes beside `ErrWaitDelay`, leaving accept-or-refuse to the
  caller); that is stated as this RDR's choice, with both surveyed peer CLIs that met
  this failure having made it too. Declared, not smuggled.
- Research found `Cmd.WaitDelay` reaches only pipes `Cmd` creates, so `spawn`'s owned
  `*os.File` ends sit outside it. The design does not "fix" this by handing the pipes
  back — that is ALT3, rejected, because 76c121b's constraint (release the group
  BEFORE joining the drains) is what requires the manual ownership. The design instead
  applies Go's rule to the owned drains. Consistent with the finding, not against it.
- The measured finding that the drain-join leg alone runs 500.1-508.2 ms against a
  500 ms `WaitDelay` in every sample could have contradicted a per-leg bound. It does
  not, because C1 commits the TOTAL (`timeout + 2·WaitDelay`) and the RDR says so
  explicitly — "the total is the claim, the legs are not separately bounded", with the
  MVV asserting on the total. The RDR declines to state a figure its own evidence
  contradicts.

Planned features vs stated principles: the refusal narrows 0025:C4's "true EOF,
nothing truncated" to "whole within the bound". That is a real change to a peer's
locked contract, and it is handled as one — declared in Metadata **Overrides**
(0025:C4 and 0025:F4's residue class), not asserted as compatible. The stated
principle it serves (a partial answer must never be parsed as whole — the
silent-wrong-answer class 0025 exists to prevent) is upheld: no opt-out flag ships,
and the RDR names why one is not proposed (a per-binding "accept partial output" knob
is ALT1 by another name, rejected for ALT1's reason).

## 2. Assumption Verification

12 Critical Assumption records. **11 Verified, 1 Pending (A11, downgraded).** No
`Docs Only` record. No off-vocabulary Method label (`ca_off_vocabulary=0`). Every
record's Status / Method / Evidence agree and every "If wrong" is non-empty.

**Self-reference:** none. The five Source Search records (A1, A4, A7, A8, A12) all
anchor into the consumer source tree; none resolves to the record file or to this
RDR's artifact directory.

**Anchor resolution:** `anchors_total=58`, `anchors_unresolved=0`,
`anchors_unlooked=0`, `peer_evidence_unresolved=0`, checked with `--repo` — so every
cited `path::Symbol` was looked for and found. Nothing SKIPPED.

**Status consistency — passes.** A11 is the sole `Pending` record and no clause states
its claim as settled. C1 `precedence:` (b) states the happens-before edge as a
*requirement* and names `-race` as the check, rather than asserting the edge holds.
A11's downgrade is principled rather than convenient: it is unverifiable before
implementation **by construction** — its named check is S3 + S7 under `-race`, and S7
needs the drain-start stall seam that A12 confirms does not exist in `cmdbind` today
and that this RDR itself adds. There is no bounded join to race until Phase 1 writes
one, so no spike can stand in. It is not MVV-critical: the MVV provably cannot
exercise the edge (silent child, each drain blocked in its first read), which is why
the record names S3/S7 instead of an MVV row. The plan is binding, not aspirational —
**Phase 1 closes only when S3 and S7 pass under `-race -count=25`, and a race report
on the slice handoff is a Phase-1 blocker, not a follow-up.** The failure mode is loud
and bounded and the repair is local, so it is survivable into implementation.

The four assumptions that were `Pending` at Stage 5 are all terminal, and one of them
changed the contract rather than merely confirming it — worth recording, because it is
the reason this gate is not a rubber stamp:

- **A9 — Verified by spike** (`evidence/spikes/a9-a10-regrace/`, darwin + linux). The
  spike **corrected the assumption's own witness**: a 1 MiB *buffered* tail does not
  discriminate a re-armed per-read deadline from a single absolute deadline (both
  recover it whole — 1.6 ms darwin / 9.5 ms linux across 32 non-waiting reads). MVV
  row 6 was amended to carry a *paced* shape 6(b) that does discriminate: 0/12 vs
  12/12, the single-absolute implementation truncating at 15.6% of the payload
  (163840 of 1048576 bytes, terminal `os.ErrDeadlineExceeded`) on both OSes. C1 now
  reads "size or total duration" — duration, not size, is the operative variable.
- **A10 — Verified by spike** (same dir). The probe is detecting (`os.ErrNoDeadline`
  on a genuine dup'd pipe fd, so it detects poller registration rather than file kind)
  and a no-op on the hot path (parked-read 12/12 at 300 ms with a 10 ms probe window,
  round-trip identical in every cell). C1 `whole:`'s byte-for-byte claim and S3 are no
  longer contingent.
- **A8 — Verified by source search.** Every cited site holds with no drift. The two
  decisive facts: `flow_exec.go`'s `ClassExecutionFailure` arm carries no phase check
  and no `Detail`, so the direct write arm genuinely lacks the applied sense; and
  `Applied()` has zero non-test callers. The gain is real and one-sited.
- **A12 — Verified by source search.** The stall seam is new, the hook shape is viable
  at `cmdbind.go:273`/`:277`, and no third build tag is needed.

A6's Verified stamp is qualified in the record ("Verified, with the per-leg claim
narrowed") and the narrowing is carried into C1, so the stamp and the contract agree.

## 3. Scope Verification

**The Minimum Viable Validation is in scope and executed in Phase 1 — not deferred.**

The proof is fixture **FX-deadline-escape**, run through the real
`cmdbind.Reader.Read` under `--allow-commands` with a declared `timeout = "2s"`. A Go
helper child spawns a grandchild with `SysProcAttr{Setsid: true}` that inherits stdout
and stderr and sleeps 60 s, while the child itself writes nothing and sleeps past the
deadline — the exact escaped-grandchild shape the defect reported.

Six rows, of which these carry the contract:

- **Row 3 (the headline assertion)** — the read returns within `2 s + 2·500 ms`
  (Normative: elapsed ≤ 3 s, asserted with C1's `bound:` tolerance), the executor
  classifies `timeout`, and the refusal's `Detail` is empty as 0025:C4 states for
  `timeout`. Critically, this is asserted **at the `flow resolve` boundary, not only
  at the binding** — the defect was reported there ("still blocked after 15 s"), so
  the row also asserts that `flow resolve` itself terminates and emits the class in
  its envelope. A binding-level pass with a hung CLI would not close the original
  report, and the row is written to make that insufficiency impossible to mistake for
  success.
- **Row 4** — the grandchild is still alive, asserted. The admitted residue is
  *tested for*, not merely conceded in prose.
- **Row 5 (success shape)** — the child writes a whole envelope and exits 0 while the
  grandchild holds the pipes: `execution_failure` within `2·500 ms` of the exit, `Err`
  naming "stdout, stderr" and "exited 0", and the envelope **not parsed**. This is the
  row that proves the silent-wrong-answer class is closed.
- **Row 6, both shapes (Normative)** — the late-but-unheld control, one shape per
  rejected mechanism. 6(a) burst guards C1 `precedence:` against a deadline of
  literally `now` (which short-circuits before the syscall and returns n=0). 6(b)
  paced — 1 MiB in 32 chunks separated by 20 ms idle gaps — is the mechanism
  discriminator, and the RDR states in terms that **shape (a) alone is not sufficient
  and must never be the only witness**. That constraint is the direct product of the
  A9 spike correcting its own witness, so the MVV now tests the mechanism rather than
  the outcome.

Measured basis for the budget (A6, `evidence/spikes/a3-a6-drain-bound/`): worst
end-to-end 2.013 s (darwin) / 2.082 s (linux) against the 2.5 s budget — 487 ms and
418 ms of headroom.

## 5. Proportionality

**Right-sized. Nothing to trim before locking. No split owed.**

**Contract count — the split test.** This RDR is the sole author of exactly **one**
independent load-bearing contract: **C1**, the drain precedence at
`cmdbind::spawn` (`contracts=1`, `contracts_durable=1`, `contracts_transient=0`). C1's
four lines — `precedence:`, `bound:`, `whole:`, `refusal:` — are facets of that single
seam, not four seams: each answers "what does the binding promise when a bounded wait
and a whole read conflict?" for one part of the same question. The `refusal:` line
does not open a second seam either, because it explicitly **cites** JDR 0003 §D1 as
the carrier's owner rather than re-deciding it — this record is that policy's
instance. One seam, one contract, no split.

**Profile re-validation.** The Metadata field reads `foundational`, and it is correct
on both the matrix and the floor. `--outcome floor` returns `emit.floor: foundational`
(`floor-holds`) from `seam_lineage=2+` — three prior closed point-fixes at
`cmdbind::reapGroup` / `drains.Wait` (intrastate#1drn, #x5vq, #878e) plus 76c121b, the
`os.Pipe` ownership change that opened the gap, with the count confirmed at Resolve
against `git log -- internal/cli/cmdbind/`. The field also matches its own two
judgements: user-facing **yes** (a previously-accepted binding can now refuse), locks
**cross-rdr** (it amends 0025:C4). The lenses that actually ran agree with the value —
cove, 3amigo, critique (two base models, diffed; `critique_models=differ`) and
repeatability (variant `full`: run-1/2/3 plus diff) are all complete, which is the
foundational row, and the row resolver returns `lens-foundational-row-complete`. No
lens was routed past by a wrong Profile.

**Form.** The field is a value plus one clause naming the contract. No matrix or
provenance prose survives from the template or the Seed.

**Length.** ~1200 lines is large in absolute terms and proportionate here. The mass is
where the risk is and is load-bearing rather than decorative: 12 assumption records
with four spike-backed verdicts across two OSes; ten test scenarios (S1-S10); six
failure modes. The three sections that could invite a trim each earn their space —
§Consequences carries the unsurveyed-population admission and the Phase-5 survey that
converts it into a route-back trigger; F5 carries the non-pollable-host residual with
its scope honestly bounded ("unexercised for a real `os.Pipe`"); MVV row 6 carries two
shapes because the spike proved one is not a witness. Cutting any of them would remove
a stated limit, not prose.

**Not flagged for trim, deliberately:** the Assumption Verification section retains
the Stage-5 gate text beneath the Stage-6 disposition. That is a record of what was
open and how it closed, and it is in gate.md rather than the locked record, so it
costs the contract nothing.
