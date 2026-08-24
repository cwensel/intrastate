Model: claude-fable-5

# Discharge check — RDR 0004 (accessor-execution-safety-model), iteration 4

Date: 2026-08-24. Read-only. Evidence only; the gate dispositions.

Inputs: `../iter-3/reconcile-report.md` §Dispositions ("SPEC-DEFECT — RDR 0004 → Draft
(RE-LOCK-ONLY)", report:137-145; findings row report:95); `../iter-3/answer-check-0004.md`
(rows 15.6, 15.7, 17.3, *Lands in* tables, Observations 1-3); home
`docs/jdr/0001-resolve-kernel-seam.md` at `567600f` — §D5 (JDR:222-255), §D7 (JDR:300-378),
JD-15 (JDR:544-553), JD-17 (JDR:560-575); current
`docs/rdr/0004-accessor-execution-safety-model.md` at `ffaf897` (1025 lines; re-entry commit
`1a5ddca` "land JDR 0001 §D5 clear-as-removal + §D7 definition shape; repair rebind/presence
prose; drop Refinement Context"), read in the regions the iter-3 clauses map to plus every
`clear`/`rebound`/`present` hit. Anchors are `0004:<n>` against the CURRENT file (iter-3 anchors
are quoted as "formerly").

## Iter-3-named contradicting clauses — successor in current text

### C1 — Approach prose, read-back "values are present" (formerly `0004:213-215`, unfenced)

Successor `0004:235-238` (Approach, unfenced):

> "A write accessor applies a planned owned-tag mutation to a caller-supplied artifact role, then
> re-reads the same role and verifies each planned owned tag for equality: an assigned value must
> be held exactly, and a key written as the reserved `<clear>` sentinel must be absent."

**LANDED.** "are present" is gone; the sentence now carries both arms (equality for assignment,
absence for clear). Consistent with §D5(a) "read-back asserts the key is **absent**" (JDR:238) and
§D7(iv) equality (JDR:353-354).

### C2 — Fidelity table row "present in the re-read" (formerly `0004:511`, unfenced)

Successor `0004:571` (Round-Trip / Inverse table, unfenced):

> "| `write -> read`, owned tags | For every planned owned key, the re-read holds exactly the plan's
> expected value; for a planned `<clear>`, the re-read does not hold the key at all | value
> equality (assignment) or absence (clear) over the planned key set | none |"

**LANDED.** The invariant now has the absence arm iter-3 said was missing; the strength column
names it explicitly. Value-for-value against §D5: clear → key not held. Consistent.

### C3 — Zero `<clear>` clauses (§D5 *Lands in 0004*: remove-key / read-back-absent /
idempotent-clear / unreadable-on-read + MVV scenario for a clearing rule)

New fence `0004:357-364` (**fenced**):

> "`<clear>` is a reserved tag value (JDR 0001 §D5). A planned write of `<clear>` MUST remove the
> key from the artifact, and its read-back MUST assert the key is absent — a re-read that still
> holds the key, including as the literal string `<clear>`, is `read_back_mismatch`. Clearing a
> key the artifact does not hold MUST succeed. A read that yields `<clear>` as a value MUST treat
> that key as unreadable and refuse `incomplete_read`."

Per landing item:

| §D5 item (JDR:253-255) | Status | Evidence |
| --- | --- | --- |
| `<clear>` write removes the key | LANDED (fenced) | `0004:358-359` "MUST remove the key" |
| read-back asserts absence | LANDED (fenced) | `0004:359-361`; also read-back fence `0004:366-374` wording "verify each planned owned tag for equality against the held value"; Technical Design `0004:276-277` "a `<clear>` write removes the key, so its read-back asserts absence" |
| idempotent clear succeeds | LANDED (fenced) | `0004:361-362` "Clearing a key the artifact does not hold MUST succeed" |
| unreadable-on-read | LANDED (fenced) | `0004:362-363`; classified as `incomplete_read` — the home says "unreadable" (JDR:239) without naming a class; 0004's own vocabulary already makes unreadable = `incomplete_read` refusal. Not a contradiction. |
| MVV scenario for a clearing rule | LANDED (unfenced) | Scenario 9 `0004:983-990`: holds-key / lacks-key / stores-literal fixtures plus a read of the literal; expected success/success/`read_back_mismatch`/`incomplete_read`; scenario-3 assignment fixture kept as control. Also disposition rows `0004:506-509`, MVV loudness row `0004:558`, Testing `0004:862-865`, Risks `0004:791-794`, Load-bearing "Clear semantics" `0004:419-423`, References `0004:1010`. |

Iter-3 Observation 2 (a literal `<clear>` stored and read back would pass green) is closed by
`0004:359-361` and disposition row `0004:507`. The premise at JDR:227 ("0004 contains the word
"clear" zero times") is no longer true of the body — expected, since the home asked for it.

New assumption **A11** `0004:209-223` (Pending, MVV Scenario 9) carries the clear semantics as a
property to prove; "A2 and A7 ... hold unchanged; this extends the expected-value predicate to
absence rather than re-opening it" (`0004:214-217`) — matches iter-3's "re-verify none".

### C4 — "cannot be rebound" (formerly `0004:370-372`, unfenced) vs §D7(ii)

Successor `0004:408-411` (Load-Bearing Decisions, unfenced):

> "**Identity** — an accessor is identified by `(flow id, accessor name, capability)`. Capability
> is part of the identity, so the same id may appear in both `[read.x]` and `[write.x]` — two
> identities, not a rebinding (JDR 0001 §D7(ii)); only a second binding of the *same* triple is
> multiply-bound."

**LANDED.** The sentence is retired and replaced by the home's rule (JDR:344-346 "The same id in
`[read.x]` and `[write.x]` is legal"). `grep -n rebound` in the body hits only the Status note
(`0004:14`) and this repair (`0004:411`, "not a rebinding"). Identity fence (formerly 17.2) is
unchanged at `0004:299-303`.

### C5 — §D7 *Lands in 0004* items iter-3 marked absent (JDR:366-368)

| Item | Status | Evidence |
| --- | --- | --- |
| definition is the capability-table entry | LANDED (unfenced) | Technical Design `0004:249-252` "A definition is the capability-table entry in RDR 0002's layout — `[read.<id>]`, `[write.<id>]`, or `[gate.<id>]` (JDR 0001 §D7(ii))"; Load-bearing "Definition shape" `0004:414-418` "by citation ... This RDR does not spell the layout; it consumes it." |
| read-back is equality | LANDED (fenced, strengthened) | read-back fence `0004:366-374` now reads "verify each planned owned tag for equality against the held value — a write replaces the whole value, so containment is not equality" (formerly "verify the expected owned-tag values"); cites §D7(iv) at `0004:276`. |
| writers carry `keys` | LANDED DIFFERENTLY (unfenced only) | `0004:253-255` "`keys` is the binding on readers and writers alike: a reader's `keys` is its requested key set, a writer's `keys` is the owned keys it may write or clear"; `0004:415-416`. The requested-key fence `0004:315-319` still binds **read** accessors only ("Every read accessor definition MUST declare the requested key set"). The writer-side `keys` obligation is prose plus a citation, not a fence. Does this contradict the home? No — it adopts JDR:340 "Readers and writers both declare `keys`". It does sit on one side of the home's own internal tension (JDR:336-337 "`keys` exists only on readers and `read_back` only on writers" vs JDR:340), iter-3 Observation 1; 0004 now follows the second sentence. Whoever reconciles D7(ii)'s two sentences owns that; 0004 cites rather than restates. |
| "cannot be rebound" repaired | LANDED | C4 above. |

Additional §D7(ii) restatement worth noting (not a contradiction): `0004:255-258` "RDR 0002 owns
the carrier and the binding validations (a key served by zero or two readers, a written or
cleared key not in exactly one writer); this RDR owns what it means to execute a referenced
accessor safely" — this repeats the home's loader rules verbatim as an ownership pointer to 0002.
It does **not** repeat "a writer key that is not owned" (that stays 0004's `write_non_owned_tag`,
`0004:589` Phase 1 row). JD-17's "Open at (ii)" provenance-scope question (owner 0002) is not
touched by 0004, which is the assigned ownership.

## Structural checks

| Check | Result | Evidence |
| --- | --- | --- |
| Refinement Context section absent | YES | Heading sweep (`grep '^## \|^### '`) lists no "Refinement Context"; commit `1a5ddca` subject "drop Refinement Context". |
| Status line | `Final [re-locked 2026-08-24 — Gate PASS. ... Three records remain Pending (A9, A10, A11) ...]` (`0004:9-16`) | Bracket is a landing note; it names no pending joint decision (iter-3's `[joint decision → §JD-15, §JD-17]` qualifier is gone). |
| README Index row | `| [0004](...) | Accessor execution safety model | Final | High |` (README:15) | Agrees: both `Final`, no qualifier on either side. |
| Gate artifact | `0004:1004` "Responses: `0004-accessor-execution-safety-model/artifacts/gate.md` (Gate PASS 2026-08-24)" | Not opened; cited by the Status line. |
| Prerequisites checklist (formerly `0004:771-780`) | **All four boxes still `- [ ]` (unchecked)** at `0004:841-849` | Text updated: "All Critical Assumptions verified (A1-A8 Verified; **A9, A10, and A11 Pending** — MVV-proven properties carried to lock with the MVV as the named implementation-time plan)"; the RDR 0001 / 0002 / 0003 boxes unchanged in wording. Iter-3 said "Prerequisites resolved at the gate" (report:226) and "named for the re-lock gate" (report:143); the re-lock left them unchecked and instead explains the Pending records in the Status line. Whether an unchecked Prerequisites block at Final is acceptable is the gate's call; the checklist's own first item cannot be ticked while A9-A11 are Pending, and the other three are peer-implementation facts (0001 Implemented is stateless per its own record; 0002/0003 unimplemented). |
| Re-verify set | none named by iter-3; none required | A11 is new (Pending, MVV), not a re-verify. A1-A8 stamps not re-examined (RE-LOCK-ONLY). |

## Summary of evidence

- All four iter-3-named contradiction sites have a successor that agrees with §D5 / §D7 in
  meaning; two of the landings are fenced (`0004:357-364`, `0004:366-374`), the rest unfenced.
- One item LANDED DIFFERENTLY: writer-side `keys` is prose-plus-citation, the fence at
  `0004:315-319` remains reader-only. Not a contradiction of the home; it inherits the home's
  D7(ii) internal tension (iter-3 Observation 1), which is unresolved in the home text at
  `567600f`.
- Prerequisites remain unchecked at Final; the Status line discloses it.
- No `<clear>` / rebind / presence wording survives in the body that contradicts the home.
