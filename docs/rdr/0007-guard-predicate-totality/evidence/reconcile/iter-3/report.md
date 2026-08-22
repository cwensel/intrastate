Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0007 guard predicate totality (iteration 3)

Third Stage 6 on this RDR. Iteration 1 reconciled the pre-demotion draft
(`../report.md`); JDR 0001 then demoted it and the draft was re-proposed,
re-refined and re-resolved. Iteration 2 (`../iter-2/report.md`) returned
**NOT RECONCILED — return to Stage 5**: two of the four `foundational`
lenses had never reviewed the current design. Both have now run.

## Verdict: RECONCILED — ready for Finalize

Every item in the open set is terminal. No BLOCKER. No refutation. No
MVV-critical item deferred.

## Stage 5 completeness preflight — PASS

`Profile: foundational` → lens row `cove → 3amigo → critique → repeatability`
(rdr-common §lens-row), evaluated against the CURRENT draft:

| Lens | iter-2 evidence | Model stamps (§model-stamp) |
| --- | --- | --- |
| cove | `cove/iter-2/` | `claude-opus-5[1m]` |
| 3amigo | `3amigo/iter-2/` | `claude-opus-5[1m]` |
| critique | `critique/iter-2/` | `claude-opus-5` + `claude-fable-5`, diff by `claude-sonnet-5` — dual-model draw satisfied |
| repeatability | `repeatability/iter-2/` | `claude-opus-5` / `claude-fable-5` / `claude-sonnet-5`, all `variant: full (profile: foundational)`, plus `diff.md` |

No variant mismatch: `foundational` owes the full variant and the full
variant is what ran. Determinacy trigger: n/a as a separate obligation — at
`foundational` repeatability is a row member, not a trigger.

## Open set — built from all four sources

| Item | Source | Disposition | Evidence / plan |
| --- | --- | --- | --- |
| **A3** kernel change is bounded | 1 (cove F-2), 2 | **VERIFIED** | `spikes/a3-reshape-respike.md` §Q-A3 — 154/154 frozen, 180/180 with probes |
| **A23** payload surface is additive | 1 (3amigo H-1), 2 | **VERIFIED** | ibid. §Q-A23 — `TestReq7` passes unedited; zero `Refusal{}` literals in tests |
| **A24** existence literal comparison | 1 (3amigo), 2 | **VERIFIED** (kernel half) | ibid. §Q-A24 — both polarities + 6 fail-closed cases |
| **A26** exported-name sweep | 1 (critique C-12), 2 | **VERIFIED** | ibid. §Q-A26 — REQ-25/36/37 all PASS against compiled surface |
| **A27** five literal byte values | 1 (repeatability G-2/G-3), 2 | **VERIFIED** (kernel half) | RDR 0003 vocabulary + `0003/evidence/spikes/guard-fixture.toml` (`exists = true`, `domain = [true, false]`) |
| **A19** structured envelope field | 1, 2 | **DOWNGRADED** | Blocked on JDR §D4 reopening; `Detail`-flattening fallback; reversible either way |
| **A22** tag-key canonicalization | 1, 2 | **DOWNGRADED** | Blocked on RDR 0002 re-lock; kernel-side input-precondition fallback verified today |
| **A25** 0009 fixture collision | 1 (critique C-3/C-4), 2 | **DOWNGRADED** | Collision VERIFIED on source; only SEQUENCING open → cluster-reconcile |
| **A12** two-row absence pattern | 2 | **DOWNGRADED** | Home JDR 0001 §JD-4 (open, substance decided); rides RDR 0003's lock |
| **B-10** combined refusal reporting | 3 (absorption residue) | **ACCEPTED** (Design Decision) | Rejected alternative now named in GATE-THEN-COUNT clause |
| Exactness-word delta | 4 | **CLEAN** | Every new claim traces to a named assumption or quoted source |

**No assumption remains `Pending`** — verified by grep over the Critical
Assumptions section (0 hits).

## Absorption audit (delegated) — one residue, closed

Iteration 2's audit cleared cove + 3amigo (48 findings, no residue). This
pass audited the two lenses that ran after it.

- **critique/iter-2 — 28 findings, 27 absorbed.**
- **repeatability/iter-2 — 17 findings, all 17 absorbed** (13 contract
  "pins", 3 dismissed-with-cite, 1 left non-normative).
- Both `Charted.md` files verified consistent.

**Residue: B-10** (critique pass B, `claude-fable-5`). A first-class finding
with its own adversarial test AT-10, **dropped at the merge step**:
`diff-modelB.md` declares its queue the "Union of distinct defects" but B-10
appears in none of its three tables and none of its 17 queue items, so it
never reached `dispositions.md`. Confirmed by direct grep. Note the merge
dropped two IDs — C-2 was also absent from the queue but was caught
downstream and given an explicit ledger row; B-10 alone survived unnoticed.

*A cross-model merge that silently loses a finding from the alt-model pass
is the one defect the dual-model draw exists to prevent. Recorded here
because the mechanism, not this instance, is the lesson.*

**Disposition: ACCEPTED (Design Decision).** B-10's concern was raised one
round earlier as 3amigo PM-9 and dismissed-with-cite as a re-raise of an
adjudicated call — both verified. But B-10 went further than PM-9: it named
the specific alternative (one refusal carrying BOTH payloads, for which this
RDR's own per-row/per-atom surface is the vehicle) and asked why the RDR
refuses it. The RDR stated its side without recording the rejected one. Per
the ACCEPTED disposition's requirement that the rejected alternative be
named, the GATE-THEN-COUNT clause now carries it:

> Combined reporting … is REJECTED here: `Refusal.Kind` is RDR 0001's single
> stable discriminator that RDR 0005 maps to one CLI code, so a both-ways
> refusal would need either a new kind or a kind whose meaning depends on
> which payload fields are populated. Both reopen a frozen taxonomy (REQ-7
> pins the kind set at exactly five) …

Grounded before assertion: `TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`
does pin exactly five kinds, and `Refusal.Kind`'s doc comment confirms it is
"the stable discriminator RDR 0005 maps to a CLI code."

§amendment-sweep: the `disposition` mini-check row already recorded "owned
payload only … (the accepted cost)", so the new clause explains a call the
table had already made. No stale site.

## The re-spike — what one spike closed

The iteration-2 report scoped a re-spike correcting the first spike's three
divergences. It ran as specified (`Refusal.Guard` DELETED, seam view-free,
existence literals explicit) in a scratch copy; `internal/` was never edited
and `go vet ./...` exits 0.

**All four questions VERIFIED.** Three results are worth carrying forward:

1. **Fixup-1d is a RE-DECIDE, and for a deeper reason than the RDR stated.**
   The RDR argued the test's `Refusal.Guard` assertion loses its referent —
   true, and a compile error. But the re-spike found the fixture's guard
   `iterations >= 3` names a key the view does not carry, so the kernel
   decides that atom on ABSENCE (`Reason:absent`, captured verbatim) and the
   nil seam is never consulted. The test would go green on kind alone while
   testing **nothing about the nil seam**. The SEAM clause's prescribed fix
   (move to the present key `reviews`) restores the property. A3's own "If
   wrong" branch, demonstrated rather than inferred.

2. **The coverage gap is proved, and the RDR's mutant was the wrong one.**
   A raw integer `min` is caught by the frozen suite already (`T <= U` folds
   even a single-atom `U` to `T`). What survives all 154 is swapping
   FALSE-dominance for UNEVALUABLE-dominance, observable only across two
   atoms. Testing row 21 updated to assert that leg.

3. **A26's banned lists are EXACT-match, not substring.** All three tests
   compare `if n == b`. The verdict is unchanged (no collision under either
   reading), but the RDR's claim that REQ-37 forecloses A22's fallback
   "under any name containing `Canonicalize`" was overstated — only the bare
   name is banned, so `CanonicalizeKey` would pass. Corrected at A26; A22
   never repeated the constraint, so no stale site. Separately, `Undecided`
   is a struct FIELD and `exportedKernelSymbols` walks top-level decls only,
   so the frozen tests cannot see it — its clearance rests on the standalone
   sweep. Recorded at A26 as a scope fact.

## Hard rules

**Refutation → BLOCKER: not triggered.** Nothing refuted. The one property
change (Fixup-1d) is the RDR's own predicted branch with a fix already
written into the SEAM clause — a confirmed prediction, not a refutation.
A19's text records a refutation ("the original 'without a new envelope
field' reading is refuted"), but it was resolved in the honest direction
(the RDR corrected to match shipped `clierr.CLIError`, consequence routed to
a Prerequisite): a closed route-back, not an open blocker. No §punt-ledger
row is owed — no stage was reopened.

**No MVV-critical deferral: satisfied.** Iteration 2 left this question
explicitly open, because the MVV asserts "the absent key in the payload" and
the payload surface (A23) was Pending — "not obviously survivable, because
the MVV cannot be written without the field it asserts on." **A23 is now
VERIFIED**, proven by a spike that built the payload and asserted on it, so
the question is resolved rather than adjudicated. The four DOWNGRADED items
were then checked against the MVV text directly: it references none of A19,
A22, A25 or A12 — no CLI envelope, no canonicalization, no fixture
sequencing, no overlap check. The MVV runs against the real kernel with a
real atom, entirely inside the surface the re-spike proved.

## Completeness check

- No `_Draft placeholder._`, no seed-skeleton header, no surviving template
  bracket (grep: 0 hits).
- `## References` carries 13 real citations; no bracketed placeholders.
- Exactness-word delta swept over `git diff 757750b..HEAD` on the RDR body:
  every added all/every/only/exactly/never/total/canonical claim traces to a
  named assumption or a quoted source. Two spot-checked on source: Testing
  row 21's premise (`GuardFalse=0, GuardTrue=1, GuardUnevaluable=2` —
  confirmed, so a naive `min` really does return `T` for `T ∧ U`) and A25's
  fixture collision (confirmed verbatim on `fixtures_test.go::escapeRow`).

## RDR edits made this pass

All dispositions landed **in the RDR**, not only here (Stage 7 reads the RDR):

- A3 → Verified, with the re-decide finding and the corrected helper-name
  list (no `copyAtoms` — the "replace" shape cannot alias the caller's table).
- A23, A24, A26, A27 → Verified, each with its concrete evidence pointer.
- A19, A22, A25, A12 → DOWNGRADED, each with a named plan and a stated
  survivable fallback.
- GATE-THEN-COUNT clause → rejected combined-reporting alternative named
  (B-10).
- Testing Strategy preamble → rewritten onto the re-spike; row 21 →
  multi-atom dominance leg added.
- Prerequisites → assumption box checked and rewritten to the terminal state.

## Next

```
/rdr-finalize 0007
```

The Metadata `Status:` line still reads `Draft [revised from Final … —
re-verify A3, A16, A17, A18, A19, A21, A22]`. All seven are dispositioned;
the flip to `Final` is Stage 7's to make, not Stage 6's, so it is left
untouched deliberately.
