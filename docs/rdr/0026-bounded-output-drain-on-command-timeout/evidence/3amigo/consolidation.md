# 3amigo consolidation — RDR 0026

Three isolated persona passes (no cross-persona visibility; isolation verified —
no file references another's output). Model stamp on each file: `claude-opus-5[1m]`.

Hotspots below are a **mechanical count over the anchored ids** (`rdr anchors`,
ids cited by ≥2 personas). Overlap marks a hotspot *passage*, not a validated
finding; a finding raised by exactly one persona is not thereby weaker.

## Hotspot ids (≥2 personas, independently)

| Count | Id |
| --- | --- |
| 3 | `0026:C1` |
| 3 | `0026:S1` |
| 3 | `0026:§existing-infrastructure-audit` |
| 2 | `0026:A6` |
| 2 | `0026:MVV` |
| 2 | `0026:S5` |
| 2 | `0026:S8` |
| 2 | `0026:S9` |
| 2 | `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob` |
| 2 | `0026:§conditional-mini-checks` |
| 2 | `0026:§consequences` |
| 2 | `0026:§failure-modes` |
| 2 | `0026:§illustrative-code` |
| 2 | `0026:§implementation-plan` |

`0026:C1` is the record's one contract and every persona landed on it — expected,
and not itself signal. The load-bearing hotspots are `0026:S1` (all three),
`0026:§implementation-plan` × `0026:§existing-infrastructure-audit` (PM + implementer),
and the `DrainGrace` decision (PM + implementer, from opposite directions).

## Merged findings, severity-ranked

Persona key: **PM** = persona-1-pm, **IMP** = persona-2-implementer, **QA** = persona-3-qa.

### High

| # | Anchors | Persona(s) | Finding | Blocks |
| --- | --- | --- | --- | --- |
| H-a | `0026:C1` `refusal:`, `0026:A8` | IMP (S1) | `errors.As` on `cmdbind.HeldPipeError` is specified at a site inside `internal/accessor`, but `cmdbind` imports `accessor`, not the reverse — an import cycle that will not compile. Verified with `go list -deps`. | Where the held-pipe type is declared; the first file written. Everything downstream — the executor cannot recognise a held pipe until this resolves. |
| H-b | `0026:§implementation-plan`, `0026:C1` `refusal:`, `0026:A8`, `0026:S1` | IMP (S2) | No phase schedules the `executor.go::Write` and `flow_exec.go::accessorFailureOf` edits that C1 and A8 commit. Phase 1 says "no other line of `spawn` moves"; Phases 2–4 are comments, fixtures, docs. `0026:S1` asserts a CLI rendering with no phase landing the code it tests. | Branch file scope and commit boundary: one-package change vs three-package change touching CLI error rendering. |
| H-c | `0026:C1` `precedence:`, `0026:§illustrative-code`, `0026:S9` | QA (H1) | "Held" is defined twice and the definitions disagree on S9's trickle-writer arm: as a *condition* ("empty, with a writer the group signal could not reach") vs as a *mechanism* (whatever the drain reports as terminal condition = `os.ErrDeadlineExceeded`). A byte-every-10 ms writer is non-empty during the grace, so refuse and accept are each defensible from the same clause. | The S9 companion assertion — the arm that pins refuse-on-no-EOF rather than on elapsed time. |
| H-d | `0026:C1` `precedence:`, `0026:MVV` row 6, `0026:S9` | QA (H2) | The deadline is armed once at `now + DrainGrace`, but `readBounded` is a loop; C1 and MVV row 6 both say "final read" (singular). A 1 MiB tail through a ~64 KiB pipe buffer needs ≥16 transits. Whether that drains inside 50 ms is unmeasured — A2 measured 24 bytes, A3 measured under an unbounded grace. A false refusal (held reported on a pipe with no writer) is the failure shape. | MVV row 6 determinism; S9's prompt arm at 1 MiB. Re-arm rule, or a measured basis and a payload sized to it. |
| H-e | `0026:F6`, `0026:C1` `class:`, `0026:D-the-drain-grace-…`, `0026:§problem-statement` | PM (H1) | The record's own motivating scenario (author declares `timeout = "2s"`, deadline elapses) routes to `timeout`, which per 0025:C4 carries no `Err`/`Detail` — so the entire held-pipe diagnostic, and the `DrainGrace` constant justified in `D-the-drain-grace-…` by preserving exactly that text, is unreachable in the headline case. `F6` labels itself "under-explained"; the justification and the class rule are in tension and neither cites the other. | Whether this RDR — which owns amending 0025:C4 — admits a held-pipe sub-reason on `timeout`, or knowingly sacrifices the author diagnostic in the majority case. A product call, not deferrable to Stage 8. |
| H-f | `0026:§phase-4-surface`, `0026:§consequences`, `0026:MVV` | PM (H2) | Profile says `user-facing yes`, and the only user-facing artifact is one unlocated sentence in the last phase: no file, no surface, no owner, no acceptance check, no MVV/scenario row. Two `§consequences` Negative claims (refusal "visible, with the pipe named"; residue "explicitly admitted") hold only if this phase ships. | What "done" means for the user-facing half; Stage 8 verifiability. Name the artifact, or delete Phase 4 and record that the diagnostic lives only in the refusal string. |

### Medium

| # | Anchors | Persona(s) | Finding | Blocks |
| --- | --- | --- | --- | --- |
| M-a | `0026:F5`, `0026:S5`, `0026:§existing-infrastructure-audit`, `0026:§conditional-mini-checks` (`disposition`) | IMP (S3) + QA (L1) — **hotspot, independent** | Both clauses use the definite article for a "creation-time pollability check" that does not exist (`grep` finds no `NoDeadline`/`SetReadDeadline` outside spikes), is in no audit row, and is in no phase. IMP asks where it runs, what "logged" means (`spawn` has no logger), and whether the fallback refuses; QA asks what a test observes to conclude the log happened — the `disposition` table calls this row "loud at the check, by design", and loudness is the one property neither can assert. | Whether `spawn` gains a logging dependency; whether S5 needs an injectable seam ("simulated at the check") that C1 does not name; S5's positive assertion. |
| M-b | `0026:§illustrative-code`, `0026:§existing-infrastructure-audit` (`readBounded`), `0026:C1` `precedence:`/`whole:`, `0026:D-selection-predicate` | IMP (S4) | `readBounded` has three exits, not the two the record names: the overflow path's `io.Copy(io.Discard, r)` can itself hit the deadline, and today discards that error. So a drain can be simultaneously overflowed (`execution_failure` per 0025:C4) and held (`HeldPipeError`) — same class, different `Err`/`Detail`, and `D-selection-predicate` resolves EOF-vs-bound and timeout-vs-execution_failure but not overflow-vs-held. Also: a deadline inside the discard-Copy is on a demonstrably non-empty pipe, colliding with C1's "empty" definition (see H-c). | The shape of the value each drain reports out (bare `error` vs struct with bytes + terminal condition + overflow flag) and the ordering of the two post-join refusal arms. Named as the first data structure the implementer would write. |
| M-c | `0026:S8`, `0026:C1` `bound:`, `0026:MVV` row 3 | QA (M1) | `C1 bound:` commits `timeout + 2·WaitDelay` with no tolerance; `S8` expects `+ 100 ms` (premortem P-8). C1 names S8 as where "the total tightens", not where it may be exceeded. A run at `+60 ms` passes S8 and violates C1. | S8's timing assertion and any CI-hardening of MVV row 3 — whether an over-bound run is a bug. |
| M-d | `0026:S1`, `0026:§conditional-mini-checks` (`authority`) | QA (M2) | S1's direct-write leg — the proof the A8 gain reached a user — asserts a *rendered* string with no named observable: no exit code, stream, stable substring, or golden convention. The `Applied()` half is writable today; the rendering half is a second, unspecified surface, and S1 requires both. | The only row proving the A8 re-key reached a user. |
| M-e | `0026:C1` `refusal:`, `0026:S1`, `0026:MVV` row 5, `0026:D-the-drain-grace-…` | QA (M3) | Only the both-pipes-held case has an oracle. Single-pipe `Detail` forms are contract surface with no scenario and no stated rendering (separator, order, wording). Consequentially, **stderr-held-alone** is untested — and that is precisely the case `D-the-drain-grace-…` invokes to justify `DrainGrace` ("what the grace preserves is the stderr tail"). | A stdout-clean/stderr-held scenario; any single-pipe `Detail` assertion. Leaves the `DrainGrace` rationale unverified by any row — converging on H-e from the test side. |
| M-f | `0026:§consequences`, `0026:F2`, `0026:§risks-and-mitigations` | PM (M1) | The behaviour change for tools that *correctly* answer while permanently escaping the group is a one-line accepted negative with no survey. F2/A5's 160-trial measurement addresses the in-group helper that misses the bound, not this population. The mitigation ("the author sees exactly what to fix") presumes the author can change the helper — false for a third-party daemonizing binary. | Plain behaviour change vs migration note / release-note callout. Cheap to close: state whether any known reader exhibits the shape, or that it is unknown. |
| M-g | `0026:§problem-statement`, `0026:MVV`, `0026:§background` | PM (M2) | The promised outcome is that a `flow resolve` caller stops waiting forever with nothing to branch on; the defect was reproduced at `flow resolve` ("still blocked after 15 s"). Every MVV/S row asserts at the binding/executor layer. Validation moved a layer down from where the bug was observed. | Whether the tracked defect can close on binding-level tests alone. |

### Low

| # | Anchors | Persona(s) | Finding | Blocks |
| --- | --- | --- | --- | --- |
| L-a | `0026:MVV` row 6, `0026:S7`, `0026:S9` | QA (L3) | Three rows share a "stalled / artificially delayed drain" input whose construction is never named (test hook, build tag, injected sleep, seam in `readBounded`). S7 is an explicit long-lived regression trap, so the seam is production-adjacent, not scaffolding. Expectations are precise; the *input* is not. | Nothing outright — risks MVV row 6 and S7 building two different hooks. |
| L-b | `0026:S3`, `0026:§conditional-mini-checks` (`fidelity`) | QA (L2) | "Byte-for-byte unchanged from today" names no baseline artifact — no golden, no pre-change capture, no cited existing test. The `oracle` control discriminates the cap, not the byte-equality claim. | The byte-equality half of S3. |
| L-c | `0026:C1` `precedence:` (final sentence), `0026:§illustrative-code`, `0026:§existing-infrastructure-audit` | IMP (S5) | C1 says "after the join every read end is closed"; the defers close at *function return*, strictly later (the `werr` switch and overflow check run between). The audit resolves this correctly, but only in the audit table — against the two places an implementer reads first. Following the Illustrative Code's `closeReads` literally writes a double-close the audit then excuses. | Whether to write `closeReads` at all — dead code worth not writing. |
| L-d | `0026:§approach`, `0026:§decision-rationale` | PM (L1) | Both argue engineering precedence; neither states the user-visible promise in one line. The outcome is derivable (nearest is `§consequences` bullet 1) but requires reconstruction from mechanism. | Nothing hard; affects whether a later reader re-litigates a settled precedence. |

## Convergences worth naming

Two independent pairs, neither persona having seen the other:

1. **The `DrainGrace` justification is unbacked at both ends.** PM (H-e) found its
   stated benefit discarded by `C1 class:` in the headline case; QA (M-e) found the
   same benefit — the preserved stderr tail — exercised by no scenario. Same
   decision (`D-the-drain-grace-…`), reached from product and from test.
2. **The non-existent pollability check** (M-a): IMP found no such code and no
   phase; QA found no observable to assert. Both landed on `0026:F5`/`0026:S5` and
   the `disposition` table's "loud at the check, by design".

Also convergent on the record's honesty: PM, IMP and QA each independently
recorded `0026:A6`'s bound accounting as disclosed rather than hidden, and all
three declined to raise it.

## Checked and found sound (recorded so a re-run does not re-derive)

- `0026:A1`, `0026:A4`, `0026:A7` — verified against source by IMP (single `spawn`
  join site; deadline-first ordering at `invokeRead`/`Gate`/`Write`; `cmd.Stdin` a
  `bytes.Reader` with `WaitDelay` set).
- `0026:D-naming`, `0026:D-the-drains-stay-concurrent` — implementable as written;
  S7's 65536-vs-1 MiB oracle is discriminating.
- `0026:MVV` rows 1–5, `0026:S2`, `0026:S4`, `0026:S6`, `0026:S7` prompt arm — fully
  testable as specified.
- Read-end close: the audit settles the double-close question (swallowed `ErrClosed`,
  not a panic); no test owed. (L-c is about where a reader finds that answer.)
- Unfilled `0026:G-*` gate bodies are template text at Draft — Stage 7's, not findings.
