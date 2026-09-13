# Absorption audit — cli/0021, iteration 4

Scope: every pre-lock lens round with a dir on disk (grounding, 3amigo
base+iter-3, critique, repeatability run-1) plus the base/iter-2/iter-3
reconcile reports. Cove absent (expected). Verdict: PASS, no residue.

## Grounding
No findings raised (findings.md: "Refuted / Not-Found: None"). Nothing to
absorb.

## 3amigo (base, 12-row ledger; iter-3, 2-row delta)
All 12 base rows dispositioned within the lens itself: fixed (C2/A7, C4/A8,
S6, S5, C1, S1, S8), dismissed-with-cite (C5/JDR 0002 D1, problem-statement),
charted-to-successor (MVV/S6 diagram legibility, decision-rationale CI-diff).
iter-3 delta (S6 decoder mechanism, C4 checkNodeCeiling) both fixed. Verified
against current text: 0021:C1 (evidence/reconcile/iter-4 read), C2, C4 all
carry the cited fixes.

## Critique (dual-model, 18-row ledger)
4 fixed (A8, C2 <opaque>-reaches-wire, C5 `data.dot` key, RT1), 13
dismissed-with-cite, 1 charted-to-successor (taxonomy.go stale comment, kata
successor). All fixed items confirmed present in current C2/C5/RT1 text.

## Repeatability (run-1, 4 admitted findings F1-F4)
Per base absorption-audit.md lines 105-134, 227: F1 (--emit arm order), F2
(field-spelling closed list) and F3 (terminal-satisfaction predicate) were
ABSORBED at base — text fixes cited verbatim. F4 (model identity) was the
one base RESIDUE.

F4 closed at reconcile iter-2 item 4 ("ABSORBED... C2 pins `model` to
`table.Model.ID`") and reconfirmed at iter-3 item 10. Verified directly in
current 0021:C2 text: "`model` (`table.Model.ID`, the AUTHORED `[model] id`,
never the `--model <path>` argument or any path-derived string)".

F3's base fix (`ownedAtomSatisfiable` existential/opaque-admits semantics)
was later WITHDRAWN together with A9 at cluster-reconcile repair `d37a921`
(iter-3 item 1). Checked current C2 text for re-opened residue: C2 now reads
"NO per-node terminal marking: which merged nodes satisfy a `terminal`
predicate set is the dead-end quantifier RDR 0015 owns (JDR 0001 §JD-23),
and this record declares no evaluator of its own." This is a superseding,
equally-terminal absorption of F3 — disclaim-and-cite-the-owner rather than
define-the-predicate-here — consistent with the joint-decision "cite the
owner" rule (iter-3 items 2-3, MOOT-by-removal). Not residue: the record no
longer makes any claim F3 could contradict.

F1 and F2 verified present verbatim in current text: C1 states the --emit
arm position explicitly (read this pass); C2's field list is the closed
11-member `rows[]` enumeration "AS SPELLED HERE" citing `dumpColumns`.

## Reconcile base/iter-2/iter-3 cross-check
Base: NOT RECONCILED (2 blockers: A9 limb-a refuted, C2/JDR-0001-§JD-23
restatement) + 1 residue (F4). iter-2: RECONCILED — both blockers downgraded/
discharged, F4 absorbed, A9 stays Pending by design (not MVV-critical).
iter-3: RECONCILED — cluster re-entry (0021-0029) repair removed A9 and its
dependent clause entirely (d37a921); all cluster obligations VERIFIED
discharged except one non-blocking cross-record obligation owed at 0029:C4
(C-7 read-order question), which 0021 correctly cites without answering.

## Conclusion
No finding from any lens round, across any reconcile iteration, names
something the CURRENT RDR text still fails to address. The one
cross-record obligation still open (0029:C4 must answer C-7) is explicitly
non-blocking to 0021 per iter-3's own verdict and does not imply a re-opened
0021 assumption or a new 0021 spike.
