Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0008, iteration 2 (post-demotion re-entry)

Verdict: **RECONCILED**. Thirteen assumptions terminal, no BLOCKER, no unrun
spike. This pass follows the 08.1 cluster demotion and the scoped Stage 4
re-verify of {A6, A11}; iteration 1's dispositions are undisturbed except where
noted.

## Stage 5 preflight

- `Profile: foundational` → required lenses `cove 3amigo critique repeatability`.
  All four evidence dirs present.
- Determinacy / repeatability: `full` variant as `foundational` requires, three
  runs with three **distinct** stamps (`claude-opus-5[1m]`, `claude-fable-5`,
  `glm-5.2:cloud`), plus `diff.md` and `resolve.md`. No variant mismatch.
- **Caveat carried from iteration 1** (not a gate failure): `critique` ran
  single-model (Pass A + Pass B, both `claude-opus-5[1m]`) as a documented
  fallback recorded in `resources.md`. Repeatability supplied the genuine
  multi-model draw.
- The demotion did not invalidate any lens: the re-entry was declared
  STAGE-SCOPED to {A6, A11}, and no lens finding rested on the false
  negative-existential.

## Open set — four sources

1. **Pre-Lock lists**: fully dispositioned at iteration 1; nothing reopened.
2. **Still-Pending assumptions**: none. All thirteen carry a terminal Status.
3. **Named-but-unrun spikes**: none. A9's normalization spike is the only spike
   this RDR ever named; it ran at iteration 1 (`evidence/spikes/a9-normalization.md`).
4. **Exactness-word delta** (the post-mutation set — what the re-verify and the
   cluster gate introduced or disturbed): A6's `every constructor` / `every row`
   / `cannot read Input.Owned at all` / `every return pairs nil` / `no non-nil
   error path`; A11's disjoint-fields claim; and — surfaced by this pass — the
   Failure Modes claim `no user data is reclassified as a programmer mistake`.
   All verified below; the last one **did not survive as written**.

Absorption audit over all six evidence dirs: every lens round is absorbed,
including both `critique/iter-2` survivors (A4's six-site migration inventory,
block 4's non-masking MUST demoted to an interim rule). Residue was staleness,
not defect: seven passages still routed forward to "Stage 7.1" as a future
stage, and that stage has since run and landed as JDR 0001.

## Dispositions

| Item | Source | Disposition | Evidence / plan |
|---|---|---|---|
| **A6** — co-residency of the two `Resolve`-entry preconditions | 4 | **VERIFIED** | All six quoted claims located verbatim in Final 0009 (Normative Contracts; A10). Source confirmed: `internal/resolve/resolve.go::Resolve` is `func Resolve(in Input) (Result, error)`; `::Input` carries `Table`, `Owned`, `Observed`, `Recognized`, `Guards`; `::Row` carries `RequiresOwned`/`Escape`/`Writes`; doc comment reserves the error path verbatim; all five returns pair `Result` with literal `nil`. Precedence correctly deferred to JD-5. |
| **A11** — independence of the two predicates | 4 | **VERIFIED** | Disjoint read-domains re-derived from each RDR's own contract. `Input` count in Final 0009 re-counted: **5**, not zero — the false negative survives only inside `Corrects:` bullets and the delete-on-re-lock defect record, correctly labeled withdrawn. |
| **JD-5 routing** (cited by A6, A11, block 4, scenario 6) | 4 | **VERIFIED** | `docs/jdr/0001-resolve-kernel-seam.md` §JD-5 exists, is **open** (no "Closed by §Dn" marker, unlike JD-1/2/6/7), and its text matches the RDR's quotes verbatim. RDR 0009 carries the matching Status latch. |
| **`Overrides` "Narrows nothing in either peer"** | 4 | **BLOCKER-avoided → ACCEPTED, repaired in place** | Cluster recorded this as `contradiction / blocks-impl` (`report.md:91`). The body always disclosed the collision honestly (A4: "three violations at HEAD"); only the Metadata summary denied it. JD-10 has since **ruled the substance** ("rename them; there are no users to migrate"), so correcting the summary takes no joint decision. Field now states the narrowing and cites JD-10. |
| **A1 Scope limit** routing the Overrides half to "a Stage 7.1 read must adjudicate" | 4 | **VERIFIED (re-anchored)** | Stage 7.1 ran; landed as JD-10, which accepts the narrowing. Passage now cites the landed decision. |
| **A7 risk line** "Reconciled at Stage 7.1" | 4 | **VERIFIED (re-anchored)** | The 0007×0008 pairwise **cleared** the encroachment outright ("No duplication finding"; name/meaning split respected both sides). It raised two *different* `blocks-impl` items on the absent-key evaluation domain, both routed to **JD-4**, open. Both now latched in Risks and in A7's residual line. Neither disturbs the name reservation. |
| **A12 handoff form** routed to "a Stage 7.1 call" (2 sites) | 4 | **DOWNGRADED (re-anchored)** | Stage 7.1 gave the choice a home without picking a form: JD-10 carries "No pointer from 0002 to 0008; handoff form undetermined," open. The *obligation* remains a hard Prerequisite gate item, unchanged. Survivable: the gate blocks claiming the Problem Statement outcome, so a wrong form cannot ship silently. |
| **Implementation-ordering** "belongs in the Stage 7.1 cluster read… if cluster-reconcile finds it does not" | 4 | **VERIFIED (condition resolved)** | The cluster read ran and answered it: verdict **NOT RECONCILED**, "Do not implement over these." The passage's own named disposition — hold at Final-unimplemented until 0002 starts and the depended-on JDs close — is now the operative one. Locking remains correct; the gate is on implementing. |
| **JD-10 totality sub-question** vs blocks 1–2 | 4 | **ACCEPTED (Design Decision — boundary stated)** | Not a collision. Block 1's "absent outcome yields a view with no `recognized` key" is a *source fact* about implemented RDR 0001 (`::assemble` injects only for non-empty `Input.Recognized`) — the input JD-10 reasons over, not its answer. Block 2's no-lower-bound rule quantifies over whether a model *declares* the tag; JD-10 asks what a *declared* one denotes. Distinct quantifiers. Boundary note added between blocks 2 and 3; naming rules unchanged either way. |
| **JD-8** (`reserved_tag_key` has no CLI code/exit mapping) | 4 | **ACCEPTED (Design Decision — latched)** | No contradiction: the RDR already owns the category token and disclaims the exit code to RDR 0005. Only the named home was missing. Latched at the Approach disclaimer; JD-8 notes 0005's A-block pre-authorizes additive `Code` values with no new envelope fields or exit groups. |
| **JD-9** — Failure Modes' "no user data is reclassified as a programmer mistake" | 4 | **DOWNGRADED (refuted as stated; repair landed)** | **The one genuine refutation this pass.** The claim rested on A10, whose sweep covers only the *accessor* half of the data channel. RDR 0005's planned `flow` verbs feed `Input.Observed` from `--tag name=value`, so `--tag recognized=x` is a reachable **user** invocation that this RDR classes a programmer mistake. A8's "zero non-test callers at HEAD" is why it escaped and is why it is unreachable today. Claim scoped to what A10 proves; A10's "Consequence for A8" corrected from "closes it" to "closes the accessor half"; JD-9 latched at both sites. |
| A1–A5, A8, A9, A10, A13 | 1,2 | **VERIFIED (undisturbed)** | Lock-time stamps preserved per the STAGE-SCOPED re-entry declaration; A10 gained a scope caveat (above) that narrows a consequence claim, not the assumption. |
| A12 | 1,2 | **DOWNGRADED (unchanged from iteration 1)** | Refuted as stated; repair is the explicit Prerequisite. |
| **Anchor doctrine** — 15 bare peer-RDR line citations | 4 | **VERIFIED (repaired)** | Bare peer line numbers are the exact mechanism that produced this demotion. All live ones converted to durable anchors (section heading / assumption ID): 0002 ×6, 0004 ×4, 0006 ×3, 0007 ×2. The two 0009 citations survive only inside the delete-on-re-lock defect record, where naming the retired ranges is the point. Two 0007 citations had cited *overlapping but different* ranges for one passage — a direct symptom of the fragility. |

## Hard rules

- **Refutation → BLOCKER?** No. One item was refuted as stated (JD-9, the
  "no user data" blast-radius claim). It is not a design premise: the naming
  rules, the `reserved_tag_key` category, the failure payload, and the `Input`
  predicate are all unchanged, and the key stays reserved on both halves of the
  data channel either way. What was wrong was a *scope* sentence claiming more
  than its cited assumption proved. The open remainder — which channel the
  breach is reported on — is jointly owned and already homed at JD-9. Repairing
  the over-reach in place is the correct disposition; no return stage is owed.
- **MVV-critical deferral?** None. The kernel half runs during this RDR's
  implementation. The normalizer half is scoped to the peer owning that code and
  is gated by a written Prerequisite, not deferred by choice. No JD this RDR
  latches is load-bearing for what the MVV proves: JD-4/5/8/9/10 all bear on
  peer surfaces, report channels, or lint semantics downstream of the naming
  rule the MVV exercises.

## Flow hygiene

The 08.1 demotion was a route-back that reopened a completed stage and owed a
§punt-ledger row; none existed. Opened
`docs/rdr/0008-recognized-tag-key-ownership-postmortem.md` with two rows: the
A6/A11 stale-peer defect, and the JD-9 scope over-reach found this pass. Both
share one root cause — an assumption verified against a *moving* peer or a
*future* caller, stamped with a point-in-time fact and no re-verification
trigger — and neither was caught by a pre-lock lens (all four ran clean).

## Completeness

- No `_Draft placeholder._`, no seed-skeleton header (grep: 0).
- `## References` filled with real citations; no template brackets.
- Six `normative` fences balanced; thirteen assumptions, zero Pending.
- Remaining bracketed text is the Finalization Gate's own response prompts —
  Stage 7's to answer, correctly unfilled here.
- The Status qualifier and `## Refinement Context` correctly remain: both are
  doctrinally cleared by the Stage 8 re-lock, not by Stage 4 or 6.

## RDR edits made by this stage

`Overrides` field corrected; A1 Scope limit, A7 residual + risk line, A12
handoff (×2), and the implementation-ordering passage re-anchored from the
now-run "Stage 7.1" to their landed JDR 0001 decisions; JD-10 totality boundary
note added between normative blocks 2 and 3; JD-8 latched at the Approach
disclaimer; Failure Modes' reclassification claim scoped and JD-9 latched;
A10's "Consequence for A8" corrected; 15 bare peer-RDR line citations converted
to durable anchors (source-file `file:line` refs paired with `::Symbol` anchors
are left as-is — the doctrine permits them; only peer-RDR body refs were the
fragile class).
