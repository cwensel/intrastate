Model: claude-opus-5[1m]

# Cove resolve — cli/0023, iteration 1

Origin ledger = `findings.md` F-1…F-7 (the lens's own first-pass findings)
plus `mini-checks.md` M-1/M-2 (the desk trace's). Every entry exits once.

| # | Anchor | Disposition | Section touched |
| --- | --- | --- | --- |
| F-1 | `0023:C2` | fixed | C2 (assignment ground rewritten); Overrides, Approach, census, Key Discoveries, Risks |
| F-2 | `0023:C2` | fixed | C2 always-keep clause |
| F-3 | `0023:A5` | fixed | A5 Evidence (retracted + re-grounded); Infrastructure Audit |
| F-4 | `0023:A5` → `0023:C1`/`0023:S4` | fixed | C1 whole-tree clause; S4 Expected; Infrastructure Audit |
| F-5 | `0023:A5` | fixed | A5 Evidence; C1 (the counterexample is now the clause's justification) |
| F-6 | `0023:C1`, `0023:S1` | fixed | C1 differential clause; S1 Expected; Risks P-12; oracle table |
| F-7 | `0023:A3` | fixed | A3 Evidence |
| M-1 | `0023:C2`, `0023:C1` | fixed | C1/C2 presence rules; S2, MVV 2, Decision Rationale |
| M-2 | `0023:S5` | fixed | Desk trace closing paragraph |

## The two that mattered

**F-1 — `revision` is constant-empty, so C2's rationale was false.**
C2 grounded `revision`'s PLAN assignment on it being "the plan's only
provenance" by which "a chained caller [detects] that the model changed."
`internal/cli/flow_exec.go::revision` returns `""` unconditionally, and its
doc comment explains why this is permanent-until-0002-changes: "RDR 0002's
`[model]` block admits `id`, `version`, `description`, and `metadata` and no
revision field, so NO model can declare one" (REQ-25's consequent, DEV-7).

Collapsed without escalation (tiebreaker-reduction gate). The ASSIGNMENT is
not this RDR's to reopen — JDR 0002 §D1's unclear-joins-PLAN rule and the
always-keep core both place it, and `revision` is structurally loader-produced
whatever its value. What was refuted is only the *justification*. So C2 now
assigns on the field's definition (a loader-produced identity token), records
the vacancy explicitly, and states what always-keep buys for it: a
forward-binding guarantee that a future projection mode cannot drop the
identity slot at the moment 0002 admits the key. The cost of being wrong is
14 bytes; the cost of the alternative is a silent provenance loss in a mode
nobody has designed. Recorded rather than assumed is the whole point — the
next reader would otherwise re-derive the same false claim.

Consequence for the Risks section: the "consumer loses provenance" risk was
mis-named. The field actually lost under projection is `model` (ECHO, and
populated), not `revision`. That risk now names `model` and explains why it
is nonetheless ECHO — the caller that opts in supplied it.

**F-3/F-4/F-5 — A5's implementability evidence was wrong in three ways.**
No shipped test walks `Commands()` recursively (every sweep is a fixed
two-level loop over the `flow` group), and the repo's one recursive walker
(`help_all.go::walkCommandTree`) skips `help`/`completion` by name with its
`docs.go` callers gating on `Hidden` — so the "the hidden `docs` command still
enumerates" evidence proved the opposite of what A5 claimed: `docs` enumerates
because that caller admits hidden commands, not because the walker is total.

The lens then found the counterexample that makes this load-bearing rather
than pedantic: `help_all.go::wireHelpSubcommandAll` registers a live second
`--all` on the auto-generated `help` command, invisible to the 0011 oracle
because that oracle descends the `flow` group only. A genuinely whole-tree
sweep for `all` would fail today.

`plan-only` is unaffected (different name), so the RDR's conclusion stands —
but A5's Verified stamp rested on retracted props, so it is flipped to
**Pending** with a spike plan, and C1/S4 now FIX the scope the walker must
cover (no name-skip, no `Hidden` gate, `help`/`completion` included) and name
`walkCommandTree` as explicitly not reusable. S4 gains a control asserting the
walk actually reaches `help`.

## Needs (re)verification — carried to Stage 6

- **A5 → Pending** (was Verified). Collision half still verified; the
  implementability half needs the spike named in its Verification plan: write
  the total walker, confirm it enumerates `help`/`completion`, assert the
  `plan-only` registrant set is exactly `{flow resolve}` under it.
- **New load-bearing claim, C1 width clause** — "the projected encoding is
  strictly shorter than the default." Grounded unconditionally by inspection
  (the five echo keys always render: `model`/`outcome` are non-`omitempty`
  strings, the three containers render `{}`/`[]`) and witnessed at 290 B → 152 B
  on the normative fixture. No new assumption booked: the claim is discharged
  by S1's own oracle at implementation, which is where it belongs.
- **No other assumption moved.** A1, A2, A3, A4, A6, A7 stay Verified — the
  lens confirmed each independently against source (A3's Evidence corrected
  for a count miscount only, 16 → 21 files, substantive claim unchanged).

## Not this RDR's to fix (recorded, not absorbed)

- `help --all` is a live second registration of the flag name `0011:C2`
  fences, unseen by the shipped 0011 oracle. Real, but 0011's to own — this
  RDR neither creates nor widens it. Cited in C1 as the justification for the
  total-walk scope; not otherwise acted on.
- The shipped 0005 oracles assert `escape_class` present on an escaped plan
  but exercise only single-class escape rows, leaving the producer's
  documented multi-class-unprobeable path untested (M-1). A 0005 gap.

Neither is charted to a successor: both are pre-existing conditions in locked
records that this projection leaves exactly as it found them.

## Convergence

No open ledger entries. Iteration 1 converged: F-1/F-3/F-4 were substantial
rewrites, but each replaced a refuted ground with a source-grounded one inside
the existing decision — no fix opened a new gap, and no finding traced outside
the origin ledger. `rdr lint` returns zero `resolution` findings.
