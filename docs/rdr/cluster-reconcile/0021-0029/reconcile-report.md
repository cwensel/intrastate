# Cluster Reconcile — `0021-0029`, iteration 1

Stage 7.1, 2026-09-12. Members resolved by
`rdr index --json --cluster-of <N> --closure --final-unimplemented`,
run from both seeds; each returns the other as `cross-cutting`, so the
membership is confirmed, not proposed.

| Member | Title | Status as examined | Revision |
| --- | --- | --- | --- |
| 0021 | Lint's normalized-graph export | Final (Priority Low) | `01412d2` |
| 0029 | What a released version number promises an agent about the machine-readable output | Final (Priority High) | `22ec5cb` |

Joint-decision homes examined: `cli/0029 §Normative Contracts` C3 (fired
→ 0022, out of cluster) and C4 (fired → 0021, in cluster). Both homes
live inside member 0029, so the member revisions above are also the home
revisions — nothing else to diff.

Neither record declares a `Cluster:` field and neither names the other as
a Predecessor; membership rests entirely on the reciprocal
`cross-cutting-owner` edges (0029 → `cli/0021:C2` at 0029:1868;
0021 → 0029 at 0021:1096).

**Build order** (`rdr index --topo=0021,0029 --edges predecessors,overrides`):
**0029, then 0021** — no explicit edge between them, so Priority decides
(High before Low). This is the order Stage 8's precheck enforces, and it
is load-bearing here: see Finding PW-4.

**Preconditions.** `rdr index --cycles --json` → `findings: []` (no
ownership, home, home-ahead-of-lock or open-at-lock cycle).
`rdr index --open-joint --json` → `open: []`. Neither member carries an
implementation capsule (`artifacts/status.md` absent on both), so both
were Final-and-unimplemented at entry.

**Cap.** `rdr-loop.toml --outcome cluster-cap --tag iter=1 --tag open=none`
→ `next=run`, `rule=cap-first`: "First pass: the set-scoped checks run in
full and its findings table is the origin ledger."

## Inputs run

| Check | Disposition | Detail |
| --- | --- | --- |
| Whole-set critique | RE-SCANNED | `critique-set.md`, 15 ledger rows, model `claude-opus-5[1m]` |
| Pairwise 0021 ↔ 0029 | RE-SCANNED | `pairwise-0021-0029.md`, 4 findings, model `claude-opus-5[1m]` |

**Pairs scanned / possible: 1 / 1.** Nothing trimmed. With two members
there is exactly one pair, and it is plausibly-interacting on the
strongest available grounds (reciprocal `cross-cutting-owner` edges plus
four `mentions` sites), so no trim decision arose.

## Findings

Critique rows `C-*`; pairwise rows `PW-*`. Rows naming the same defect are
merged and both ids kept.

| Pair / source | Finding | TYPE | SEVERITY | Disposition | Detail |
| --- | --- | --- | --- | --- | --- |
| PW-1, C-1, C-2, C-3 | 0021 never assigns the tiers `0029:C4` obliges for `--emit`, `graph-export-too-large`, and C2's field spellings | gap | blocks-impl | **SPEC-DEFECT → 0021** | Draft @refine, STAGE-SCOPED |
| PW-2 | The delegated surfaces have no enumeration seam, so any tier assigned is unassertable | gap | risks-impl | **SPEC-DEFECT → 0021** | same demotion; `0029:C4` obliges the seam |
| C-14 | `0021:C2`'s "closed 11-member list with `emit` appended last" states a cardinality and tail position `0029:C2` forbids asserting on an `append-only` vocabulary | contradiction | risks-impl | **SPEC-DEFECT → 0021** | same demotion; undecidable until the tier is assigned |
| PW-3, C-9 | `0021`'s joint-check reads "clear (12 peers) — open peers 0012–0020, 0022–0024"; 0029 is in neither range, yet 0029 fired a joint decision at 0021 | gap | risks-impl | **SPEC-DEFECT → 0021** | same demotion; re-run the check against 0029 |
| C-11 | `0021:A4` verified against the pre-0029 envelope; `0029:C1` adds a non-`omitempty` field to every terminal record | gap | risks-impl | **SPEC-DEFECT → 0021** | same demotion; A4 is the `re-verify` set |
| C-8 | `0021:S2`'s goldens captured under `--as=json` carry `schema_version`, which moves on 0029's schedule | gap | risks-impl | **SPEC-DEFECT → 0021** | same demotion; scope the goldens to the document |
| C-10 | `0029`'s Decision Rationale delegates to 0021 as "an early Draft with every assumption still Pending" (and "RDR 0021 (Draft, `large`)"); 0021 was Final 2026-08-28, two weeks before 0029's date | contradiction | blocks-impl | **SPEC-DEFECT → 0029** | Draft @finalize, RE-LOCK-ONLY, `re-verify none` |
| C-7 | Neither record sequences the two version markers: which does an agent read first, and what happens when one is supported and the other is not | gap | risks-impl | **JOINT-DECISION** | Home `cli/0029 §Normative Contracts` C4; see below |
| PW-4 | 0021's byte-identity goldens vs 0029's added envelope field | round-trip | — | **NO CONFLICT** | see below |
| C-15 | `0021:C2` uses the descriptive word "closed", which `0029:C2` retires repo-wide | duplication | cosmetic | **NO CONFLICT** | `0029:S9` exempts sites describing something C4 does not tier, naming `DumpColumns` — 0021's usage — explicitly |
| C-4 | `0029:A4`'s Verified status was earned against `main`, which lacks 0021's surfaces | gap | risks-impl | **NO CONFLICT** | a consequence of PW-1, not a separate repair; re-verifies once 0021 assigns tiers |
| C-5, C-6 | `A7`'s snapshot reads seams, so it cannot see an unassigned surface; the CLIError `code` vocabulary has no seam at all | gap | risks-impl | **NO CONFLICT** | both limits already disclosed inside 0029 (`C4` records `seam: none (prose-only)`; `A4`'s "If wrong" names the residual) |
| C-12 | The 1.0.0 census gate fails if 0021 ships untiered | gap | cosmetic | **NO CONFLICT** | `§Activation Step 3` records "take the 1.0.0 slip" as the designed remedy |
| C-13 | The set cannot be implemented in its stated order without a demotion | — | — | **NO CONFLICT** | correct, and this report is that demotion |

No `DEFER-TO-IMPLEMENTATION` disposition was taken: every surviving
finding is either fenced normative text or changes a clause's meaning,
and both guardrails forbid deferral. `deferred.txt` is empty.

### PW-4 — why NO CONFLICT

The scan reported that "neither record names the ordering" and typed the
finding `joint`. Two facts it did not have resolve it:

1. **Build order is 0029 → 0021** (above). 0021's goldens are therefore
   first captured over an envelope that already carries `schema_version`.
2. **0021 has no goldens on disk.** `artifacts/` holds only `gate.md` and
   `run-plan.md`; `0021:A5` states explicitly that the spike's sha256 is
   *not* a golden and that "the Phase 1 golden minted by the real exporter
   is the byte fixture for Testing Strategy scenario 2."

So nothing goes red. The re-capture obligation `0029` does carry
(Phase 1 Step 1) is scoped to `0023`'s already-checked-in golden, and is
already written there. The *wording* exposure — that a `--as=json` golden
would straddle the envelope boundary at all — is real and rides with
0021's demotion as C-8.

### Ownership — why PW-1/PW-2 are SPEC-DEFECT and not JOINT-DECISION

Both scans typed the tier gap `joint`. It is not. `0029:C4` decides
nothing jointly: it states a general rule over its own domain ("A surface
added later takes a tier assignment in the same document as part of the
change that adds it; an unassigned machine-readable surface is a
defect"), applies it to 0021 by name, and bounds its own census honestly
("complete as of this record's implementation, not complete for all
time"). 0021 simply never discharged it — the tier words appear nowhere
in its text, `0029:C1`/`0029:C4` are cited nowhere, and its own Phase 4
already edits `docs/cli-output-contract.md`, the very file `0029:C2`
requires the tier be recorded in. A defect internal to one record dressed
as a shared decision would have dodged a real demotion.

0021 is the less foundational of the pair (Priority Low against High; and
`0021:C2` already defers to 0029 — "an incompatible change bumps the
marker under the `0.x` promise RDR 0029 governs"), so 0021 yields.

0029's own demotion is independent and narrower: not the tier rule, which
stands, but a false factual premise about the delegee's status.

### C-7 — the joint decision

**Home**: `cli/0029 §Normative Contracts` C4 — already the normative home
for the envelope-versus-document versioning boundary, and already the
home 0029's lock-fence joint-check fired at 0021 (0029:1090).

**Open question**: *which version marker does a consumer check first —
the envelope's `schema_version` (`0029:C1`) or the document's `schema`
(`0021:C2`) — and what does it do when one is supported and the other is
not?*

Neither record answers it. `0029:C1` obliges a consumer to "reject an
envelope reporting an unsupported major" but is silent on a supported
envelope carrying an unrecognised document marker; `0021:C5` defines all
four `--as`×`--emit` cells and the `.data` embedding but never sequences
the two reads. Neither record solely owns it: 0029 owns the envelope
promise, 0021 owns the document format, and the ordering *between* them
is the boundary C4 already claims.

**Where the tolerance is recorded**: both records are demoted by this
report, so neither Status line can carry the
`Final [joint decision → …]` qualifier — that slot holds the demotion
qualifier. The tolerance is recorded in both re-entry notes instead, and
the qualifier applies at re-lock if C4 has not answered by then. 0029
answers it in C4 as one paragraph; 0021 cites the home rather than
restating it.

## Verdict

**NOT RECONCILED.**

Two SPEC-DEFECTs stand, and both records are demoted:

| Record | New status | Re-entry | Scope | Why that scope |
| --- | --- | --- | --- | --- |
| 0021 | `Draft [revised from Final 2026-09-12; re-verify A4 @refine]` | Stage 3 (refine) | **STAGE-SCOPED** | The approach, document shape and determinism proofs stand; what is owed is a stability declaration over surfaces already defined, plus two citations. Only A4 is disturbed (the envelope it verified against moves). |
| 0029 | `Draft [revised from Final 2026-09-12; re-verify none @finalize]` | fix + `/rdr-finalize` | **RE-LOCK-ONLY** | A factual correction to Decision Rationale prose. No assumption disturbed, no contract clause changes meaning; the census, tier vocabulary and enforcement design all stand. |

Both records are flipped, both README rows updated, and both carry a
`## Refinement Context (cluster re-entry — delete on re-lock)` note with
the defect quoted, the target stage, the scope and the resolution
direction — done, not recommended.

Do not implement either record until it re-locks at its named scope.
Re-run this gate after 0021's refine, since that rewrite touches the
cross-RDR seam directly.

## Ledger

`open.txt` carries `C-7` — the one finding that is neither repaired nor
closed by a demotion, and the entry the next iteration diffs against.
Every other row is dispositioned. `deferred.txt` is empty.
