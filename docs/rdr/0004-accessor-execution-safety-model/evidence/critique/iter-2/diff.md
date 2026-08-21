Model: claude-sonnet-5 (diff pass; authored neither critique)

# Diff — Pass A (claude-sonnet-5) vs Pass B (claude-fable-5), RDR 0004 critique iter-2

Reconciled by passage anchor, not by ledger ID (ids are per-file and do not
correspond). Grounding beyond the supplied facts: confirmed live against
`internal/resolve/resolve.go` (`Resolve`, `gate`, `escapeOrRefuse`,
`missingOwned`, `Row.Match`/`RequiresOwned`/`Escape`/`rescues`) and
`docs/jdr/0001-resolve-kernel-seam.md` (JD-3), and against the spike's
`write()` (`main.go:224-278`, re-reads via `clone(art.tags)` at line 248, not
via `read()`).

---

## 1. AGREED

Both passes independently flag the same three RDR passages as under-specified
or unsound. In every case Pass B's framing is more precise: it names a
concrete mechanism (a Go type, a code path, a struct field) where Pass A
names the shape of the gap but stops at "the RDR doesn't say."

### AGREED-1 — Absence-as-omission has no pinned return type at the seam
- **Anchor**: Load-Bearing Decisions, "Absence crosses the seam as omission" ("how the accessor layer represents an absent value internally is free... what crosses the seam MUST leave the key absent from the owned snapshot") + Normative Contracts seam clause.
- **A's ID**: C-2 (also C-8, failure 1a, AT Scenario 1).
- **B's ID**: C-2 (also W2).
- **Shared claim**: The RDR states an outcome at the seam ("key absent from the owned snapshot") but never specifies the accessor's external return type, so nothing stops a sentinel-value implementation (`map[string]string{"key": "<absent>"}`) from typechecking as compliant.
- **Framing difference**: Pass A stops at "the return type is unspecified, a sentinel could slip through" and proposes a structural test (return type must have no code path for inserting a value). Pass B goes further: it observes the RDR states *two* contradictory shape contracts for the same success branch — the completeness clause ("exactly the keys it was asked for — no requested key missing") and the seam clause ("leave the key absent") — and that no value can satisfy both simultaneously; there must be an unnamed projection component between an accessor-internal value and the seam-crossing `Input.Owned` value. Pass B also catches that A9's own If-wrong clause names the likely fix as "moving the absence encoding into RDR 0001's `Input` shape," directly contradicting Proportionality's claim that "`Input` is unchanged" — a live self-contradiction within the RDR that Pass A does not cite. **Pass B is more precise**: it identifies the contradiction is present *in the document today*, not merely a risk of future underspecification.

### AGREED-2 — A8's "binding's call" delegation has no enforcement mechanism
- **Anchor**: Load-Bearing Decisions, "Absent vs unreadable is the binding's call, defaulting to unreadable."
- **A's ID**: C-5 (also "the one section that will be rewritten," §2).
- **B's ID**: C-6 (also §3, "the one assumption that will not survive first contact").
- **Shared claim**: A8 is stamped Verified on evidence (the spike's `unreadable map[string]bool`) that proves branch machinery, not a binding's ability to actually classify absent-vs-unreadable at a real transport boundary (TOML parse errors, HTTP responses, `gh`/`git` exit codes are all named by both). The RDR gives no shared type or validation-time check to enforce the classification.
- **Framing difference**: Pass A frames this as a governance/enforcement gap ("delegated to N independently-written binding implementations... no shared type, no validation-time check") and proposes a structural fix (a shared helper enforcing the default). Pass B goes further and shows the conservative default is not merely unenforced but *self-defeating*: forcing every ambiguous case to `incomplete_read` makes every optional key read through a CLI-backed binding refuse constantly, which creates a rational incentive for the binding author to guess absence on exit 1 — "the exact guess the contract forbids, invisible below the seam." **Pass B is more precise**: it identifies the incentive mechanism, not just the missing enforcement, and names why the conservative default cannot survive first contact rather than merely "might get it wrong."

### AGREED-3 — A9 improperly bundles four independent rules into one Pending status
- **Anchor**: Critical Assumptions A9 ("the read-completeness contract's seam and boundary rules are implementable as stated") — Status Pending, Method MVV Test.
- **A's ID**: C-1 (also §3, "the one assumption that will not survive first contact").
- **B's ID**: C-8 (also §2 runner-up mention, §4 premortem "Week 2").
- **Shared claim**: A9 packages four genuinely independent claims (seam omission / omission encoding, validated requested-key set, timeout-outranks-incomplete_read, read_back_incomplete) under one Pending/Verified unit, with wildly different verification difficulty and no per-rule fallback.
- **Framing difference**: Pass A's version is an implementability critique — it argues the timeout-precedence rule specifically presumes key-granular observability that atomic I/O primitives (a single `os.ReadFile`) structurally cannot provide, so that sub-rule may be vacuous or force artificial per-key polling. Pass B's version is a lock-governance critique — bundling means "if any one fails, all four are formally refuted together," and it ties this directly to AGREED-1 (the omission encoding is not just Pending but arguably already refuted on paper by the `RequiresOwned`/escape mechanism Pass B found, see DISAGREEMENT/B-ONLY below). **Pass B is more precise for the lock decision** (which sub-rule is actually load-bearing and already falsified vs. genuinely open), but **Pass A is more precise on the mechanics** of why timeout-precedence specifically is a bad fit for the codebase's only real I/O shape (TOML file read). Treat these as complementary, not overlapping duplicates — see RANKED.

---

## 2. A-ONLY

### A-ONLY-1 — Whole-read refusal poisons multi-key accessors on one flaky field
- **Anchor**: Normative Contracts, "One unreadable requested key MUST refuse the whole read..." + Load-Bearing "Read refusal granularity."
- **A's ID**: C-4 (failure 1c, AT Scenario "multi-key accessor with one permanently-flaky key").
- **Claim**: A single field that is consistently unreadable poisons every read of a multi-key accessor, even for keys unrelated to it and even for transitions that don't consume the flaky field, with no runtime remedy short of re-authoring the transition model into narrower accessors.
- **B's silence assessment**: Likely **oversight**, not a considered non-finding. Pass B's own C-7/AT-7 (requested-set-to-table-consumption drift) and its general focus on the escape/no_match mechanism absorbed most of its attention on the *absence* side of read completeness; it never engages with the *reliability* implication of whole-read refusal for accessors backed by real external systems (RDR explicitly scopes in "External API accessors" via A5). This is a genuine, distinct, code-adjacent concern (the RDR's own text half-concedes it: "a decision, not an artifact of the spike's early return" protests exactly the point Pass A is making) that Pass B had room to raise given its granular, mechanism-first style, and didn't.

### A-ONLY-2 — Replay-stability evidence rides the disowned fallback path, not the normative one
- **Anchor**: A4 Evidence + Load-Bearing "Requested key set" (`expectedTagKeys` cited as "circular shape — fixture convenience, not the contract").
- **A's ID**: C-6 (AT Scenario "Replay stability is proven on the normative requested-key-set path").
- **Claim**: 10 of 13 transcript lines exercise the spike's fallback/derived-key-set path; the explicit requested-key-set path (the one the normative clause actually requires) is exercised by only 3 lines, and A4's "Verified" replay-stability claim is asserted on the disowned path, not the contract path.
- **B's silence assessment**: **Oversight** with moderate confidence. This is a specific, checkable claim about transcript-line coverage that materially undercuts A4's "Verified" stamp — exactly the kind of self-reference/evidence-mismatch issue the RDR's own Method vocabulary rules exist to catch. Pass B's C-3 gets close (it also worries about A4's determinism claim, but via the *timeout* nondeterminism angle, not the *which code path was actually witnessed* angle). The two are adjacent but not the same finding; Pass B does not independently arrive at the coverage-gap observation, and nothing in its methodology should have prevented it from noticing (it does close-read the spike elsewhere, e.g. C-6/AT-3). Treat as a real gap in Pass B's ledger, not a considered pass.

### A-ONLY-3 — The three-item silent-failure enumeration is reactive, not derived
- **Anchor**: Failure Modes, "There are three silent-failure shapes, each with a mandatory guard."
- **A's ID**: C-7 (AT Scenario "silent-failure enumeration is derived, not just listed").
- **Claim**: The enumeration lists exactly the three shapes the last two review rounds happened to catch, with no argument it is the complete set; `read_back_incomplete` is itself a fourth near-silent shape that had to be discovered by review, undermining confidence there isn't a fifth.
- **B's silence assessment**: **Considered non-finding**, likely. Pass B's own premortem (§4, "Week 5") independently surfaces a *new* near-silent shape not on the RDR's list — the stuck-Final / no-modeled-recovery-edge problem (its C-4/W3) — which is functionally a demonstration of the same underlying worry (the enumeration is incomplete) but Pass B chose to present it as a standalone substantive defect (C-4) rather than as a meta-critique of the enumeration's methodology. This suggests Pass B saw the same gap and made a reasonable stylistic choice to cash it out as a concrete failure mode rather than a process critique — not that it missed the pattern.

---

## 3. B-ONLY

### B-ONLY-1 — The omission encoding is masked by `no_match` + escape, not just "unenforced" (grounded, confirmed live against code)
- **Anchor**: Normative Contracts seam clause + Disposition Table row "Requested key absent from artifact... loud one layer down" + the table's closing claim "No input class exits silently."
- **B's ID**: C-1 (§1 W1, §4 premortem Week 4, AT-1).
- **Claim**: `owned_state_unavailable` fires in exactly one place — `missingOwned`, which quantifies only over `Row.RequiresOwned`. A row that consumes the same missing key via `Match` (not `RequiresOwned`) is silently excluded from the candidate set in `Resolve()`'s first loop (`view.matches(row.Match)` returns false, confirmed at `resolve.go` around the candidate-filtering loop) before `gate`/`missingOwned` ever run. Zero candidates yields `KindNoMatch`, which is escapable via `escapeOrRefuse` (confirmed: re-filters `Escape` rows by `row.rescues(r.Kind)` and `view.matches(row.Match)`, and on exactly one viable escape row returns `Result{Plan: ...}` with `Escaped: true`). So a genuinely-missing owned key can produce a *committed write* through an escape edge, never touching `owned_state_unavailable` at all. JD-3 in JDR 0001 independently confirms `RequiresOwned`'s producer is unresolved — "The field appears **zero** times in 0002... so no layer is obliged to populate it" — matching Pass B's claim exactly.
- **A's silence assessment**: **Oversight**, and a significant one. This is not a matter of interpretation or style; it is a verifiable code-path claim (now independently confirmed against `resolve.go` and `docs/jdr/0001-resolve-kernel-seam.md` JD-3) that directly falsifies the RDR's Disposition Table claim "No input class exits silently." Pass A's methodology (heavy on RDR-internal close reading, spike-code tracing, and Gherkin scenario construction) never traced the *consumer side* of the seam clause — it stayed within RDR 0004's own text and the resolver's `missingOwned`/`TagSet.has` surface (per the supplied grounding facts) but did not walk `Resolve()`'s candidate-selection and escape logic. Given the grounding facts handed to both passes explicitly named `missingOwned` and `TagSet.has` but not `Row.Match`, `Row.Escape`, or `escapeOrRefuse`, Pass A may have been anchored by the grounding brief's scope; Pass B did the extra legwork. This is the single highest-value asymmetry between the two passes.

### B-ONLY-2 — Refusal-class enumeration disagrees with itself across three RDR sections
- **Anchor**: Failure Modes list vs. Disposition Table vs. Normative Contracts (artifact-unavailable is both its own class and folded into `execution_failure`; gate-denied is both listed as a refusal and explicitly "typed gate result, not a refusal"; non-owned-write is both a runtime refusal-shaped item and a pre-runtime validation code).
- **B's ID**: C-5 (§2 runner-up, AT-6).
- **Claim**: A mechanical set-equality check across the three enumerations fails today on three named items; this is exactly the kind of drift the Finalization Gate's "mechanical pre-sweep" is supposed to catch and didn't.
- **A's silence assessment**: **Oversight**, moderate confidence. This is a five-minute mechanical cross-reference (AT-6 describes it as such) that Pass A's methodology — heavy on prose-level scenario construction — did not perform. Nothing about Pass A's approach should have prevented catching it; it appears simply not to have compared the three lists side by side.

### B-ONLY-3 — Requested-key-set has no validation tie to table consumption
- **Anchor**: Normative Contracts requested-key-set clause + the eight validation arms (none of which checks "requested set covers what the table's rows consume").
- **B's ID**: C-7 (AT-7).
- **Claim**: Completeness is relative to the declared requested-key set, but nothing ties that set to the keys the transition table's rows actually consume (via `Match` or `RequiresOwned`). A table edit that adds a consumed key no read definition requests passes all eight validation arms and fails only at runtime, surfacing as a config-drift defect misdiagnosed as an artifact problem.
- **A's silence assessment**: **Oversight**, moderate-to-low confidence. This is a natural extension of Pass A's own C-6 concern (requested-key-set fidelity / circularity) but pointed at a different failure mode (validation coverage vs. table drift rather than replay-path coverage). Pass A was clearly attentive to the requested-key-set clause elsewhere in the document, which makes the miss slightly more notable, but the specific "validate coverage against table consumption" angle requires cross-referencing RDR 0002's table-row shape (`Row.Match`/`RequiresOwned`) against RDR 0004's validation arms — the same cross-file traversal that produced B-ONLY-1. Consistent with a genuine gap in Pass A's search scope rather than a considered rejection.

### B-ONLY-4 — Post-mutation `read_back_incomplete`/timeout leaves the artifact advanced with no modeled recovery
- **Anchor**: Normative Contracts `read_back_incomplete` clause ("the verification did not run") + Round-Trip Invariants ("undo is not claimed") + absence of any Disposition Table row for "write succeeded, then post-mutation timeout/incomplete."
- **B's ID**: C-4 (§1 W3, §4 premortem Week 5, AT-5).
- **Claim**: By the time `read_back_incomplete` (or a post-mutation timeout) is reported, the authoritative artifact is already mutated. The RDR requires no retry, no idempotency, and models no transition-table edge for "owned tag advanced, resolution reported failure" — so a user who re-runs the flow after a transient re-read blip either gets `no_match` forever (no edge out of the advanced state) or, per B-ONLY-1's mechanism, silently escapes.
- **A's silence assessment**: **Oversight**, moderate confidence. Pass A's C-1 flags that `read_back_incomplete` is unwitnessed evidence-wise (failure 1b), but stops at "the refusal class might get miscoded as `read_back_mismatch` under time pressure" — a classification-error framing. It does not follow through to the *consequence* framing (artifact already mutated, no recovery edge modeled), which is a distinct and arguably more load-bearing defect. Given Pass A's general thoroughness on this exact normative clause, this reads as a genuine blind spot on consequence-tracing rather than a considered decision that the point doesn't matter — it's the kind of finding that follows naturally from the clause Pass A was already staring at.

---

## 4. DISAGREEMENT

**No direct logical contradiction between the two passes' claims was found.** Both passes converge on the read-completeness/absence machinery being the weakest part of the RDR, and neither asserts something the other's claims make false. The closest thing to tension is a difference in *severity framing* on the omission-encoding question, worth flagging explicitly:

- **Pass A treats the sentinel-vs-omission ambiguity as a future implementation risk** — "nothing stops a Phase 1 implementer from choosing a sentinel," a *possible* wrong turn that acceptance tests (its Gherkin scenarios) can prevent before it happens.
- **Pass B treats the same general area as already broken on paper, independent of which return type gets chosen** — even with the omission encoding implemented *exactly as specified*, the `no_match`/`Row.Match`/escape mechanism (B-ONLY-1) means a genuinely-absent key still produces a silent write whenever the missing key is consumed via `Match` rather than `RequiresOwned`. This is not the sentinel bug Pass A describes; it is a distinct failure that occurs even in the *correct* implementation of the RDR's stated rule.

These are not contradictory — they are two different defects in the same neighborhood, one about *which value crosses the seam* (A) and one about *what the resolver does downstream even given the correct value* (B) — but a reader could mistake them for restatements of each other. They are not: fixing Pass A's concern (pin the return type, forbid sentinels) does **nothing** to fix Pass B's concern (the `RequiresOwned`-vs-`Match` gap), because Pass B's failure mode is triggered by the omission encoding working exactly as designed. Any implementation response to this RDR must treat both as live, independent defects rather than assuming a type-level fix subsumes the resolver-level one.

One near-miss worth flagging for the record: A9's If-wrong clause proposes "moving the absence encoding into RDR 0001's `Input` shape (JDR 0001 §D3 option (c))" as the fallback if the current approach doesn't hold up. Pass B's B-ONLY-1 finding suggests that fallback would not even fix its defect on its own — option (c) changes how absence is *carried*, not whether a `Match`-consuming row is exempted from the owned-state check. Neither pass states this explicitly, but it follows from combining AGREED-1 (A's/B's shared concern about the `Input`-shape fallback) with B-ONLY-1 (B's escape-mechanism finding); flagging it here since it affects how much weight the A9 fallback plan should be given at lock.

---

## 5. RANKED

Ordered by how load-bearing each distinct defect is for an RDR about to lock — weighted on (a) does it falsify or require rewriting a normative clause, (b) does it affect the accessor/resolver consumer seam contract, (c) is it grounded in code that exists today (vs. a spike-only or hypothetical binding).

1. **B-ONLY-1 — `no_match`/escape masks missing owned state consumed via `Match`.** Falsifies the Disposition Table's explicit claim ("No input class exits silently") against code that exists on `main` today (`Resolve`, `gate`, `escapeOrRefuse`, confirmed live). Directly affects the resolver consumer seam. Not a future risk — reproducible today with a hand-built table, no accessor implementation required. Highest-confidence, highest-severity finding in either pass.

2. **AGREED-1 — No pinned return type for the seam; two normative clauses bind the success branch incompatibly.** Blocks Phase 1 from starting cleanly; the RDR's own text (A9 If-wrong vs. Proportionality) already contradicts itself on whether RDR 0001's `Input` is really out of scope. Load-bearing for the consumer contract and for whether A9 can even be verified as stated.

3. **AGREED-2 / B-ONLY-4 (combined) — Absent-vs-unreadable delegation is unenforced, and post-mutation `read_back_incomplete` has no recovery path.** Both concern normative clauses whose failure mode is operational (not merely a spec gap): one creates an incentive to violate the safety rule; the other leaves a mutated artifact with no modeled way forward. Neither requires a hypothetical binding shape — both apply to any real accessor.

4. **A-ONLY-2 — A4's replay-stability evidence rides the non-normative fallback path.** Directly undercuts a Verified status stamp with a specific, checkable transcript-coverage claim; affects confidence in A4 and, by extension, A9's neighbors.

5. **B-ONLY-3 — No validation ties requested-key set to table consumption.** Real gap in the eight-arm validation story, but its failure mode (config drift after a table edit) is less acute than 1-4 — it degrades gracefully into `owned_state_unavailable` or `incomplete_read` rather than a silent write.

6. **A-ONLY-1 — Whole-read refusal poisons multi-key accessors on one flaky field.** Real operational concern for external-API-backed accessors, but it is a reliability/usability defect, not a safety-contract violation — the system fails loud (`incomplete_read`), just too often.

7. **B-ONLY-2 — Refusal-class enumeration disagrees with itself across sections.** Mechanical, real, but cheap to fix (a five-minute cross-reference) and does not change any normative clause's substance, just its editorial consistency.

8. **AGREED-3 — A9 bundles four independent rules under one Pending status.** Important for how the lock decision is framed, but somewhat subsumed by items 1-2 above (which already show at least one of the four bundled rules is either contradictory or actively falsifiable today) — ranked here mainly for the governance point that Pending-bundling is itself a lock-process risk, not a new technical defect.

9. **A-ONLY-3 — Silent-failure enumeration is reactive, not derived.** Meta-level critique of methodology; useful for process hygiene but does not itself identify a new defect beyond what items 1-8 already surface (per the assessment in B-ONLY discussion, Pass B's own premortem effectively demonstrates the same point via C-4).
