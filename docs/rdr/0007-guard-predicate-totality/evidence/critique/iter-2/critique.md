Model: claude-opus-5

# Hostile critique — RDR 0007 Guard predicate totality (iteration 2)

Independent pass. No prior critique evidence read.

---

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Normative Contracts, existence clause: "`exists = false` decides TRUE when it is absent" vs. A24 Evidence, which records the spike's atom as `{Key: "guard.missing", Op: resolve.OpExists}` with NO `Literal` and the spike probe asserting "an `exists` atom over an absent key is decided FALSE" | The only executed prototype of the existence rule implements the OPPOSITE polarity from the normative clause, and the only mutation-tested probe file the RDR hands Phase 1 as "the shape to start from" encodes the wrong rule. Phase 1 copies the probe, inherits `exists`-absent⇒FALSE, and the row prunes instead of selecting. | An author writes "this row applies when `legal_hold` was never set" (`legal_hold exists = false`), and the row silently never fires. The kernel reports `no_match` — the escapable class — so a modeled escape rescues it and the user gets a plan built from the wrong edge. This is the exact masking path the RDR exists to close, re-entered through the one operator the RDR calls total. | §1, premortem, AT-1 |
| C-2 | §Normative Contracts, payload clause: "The ordering is therefore the tuple `(RuleID, SourceLocator, key, block, operator token, literal)` … two entries equal on all six name the same atom" | The tuple is NOT total. RDR 0003's operator matrix gives `in` a "non-empty typed scalar set" and `contains` a "non-empty typed element set" literal; a set literal has no stated scalar spelling, and §Phase 3 concedes `resolve.Tag.Value` is "a bare `string`, so a set-valued tag reaches the seam as one opaque string with no stated element encoding." Two `in` atoms on one key with permuted literal sets compare on an unspecified encoding. Worse, the claimed totality is a claim about ATOM identity, but the kernel MUST NOT report a duplicate — and the RDR gives no dedup rule, only a prohibition. | `Refusal.Undecided` comes back in a different order for the same table on two runs (or two `intrastate` builds), and RDR 0001 REQ-1 replay determinism — the property ADV-3 freezes — fails intermittently. The operator sees a diff in a `--json` refusal payload with no input change, and cannot tell whether the artifact changed or the tool did. | §1, §3, premortem, AT-2 |
| C-3 | A21 Evidence: "It binds producers, NOT the kernel boundary — so this composition holds for rows a conforming producer built, and **the kernel enforces no part of it**" | Factually contradicted by RDR 0009 (`Final`), whose third normative block reads "The kernel enforces the same obligation as an entry precondition of Resolve … MUST cause Resolve to return a non-nil Go error … `Table.CheckValid() error` … `ErrEscapeShapeBreach` … `*EscapeShapeBreachError`." A21 read only 0009's producer clause and stopped. The kernel DOES enforce it, on the Go error path, over the whole table, before any evaluation. | Every existing escape fixture in the frozen suite (`fixtures_test.go::escapeRow` sets `Writes: []resolve.Tag{{Key:"status",Value:"Blocked"}}`) makes `Resolve` return a non-nil error once 0009 lands. Phase 1's "re-encode the frozen suite" produces a suite where ADV-1b, ADV-2, Fixup-1d and Fixup-1e all fail on `mustResolve`'s error check, and the implementer cannot tell whether 0007 or 0009 broke them. | §1, §2, premortem, AT-3 |
| C-4 | §Testing Strategy row 15: "the shipped `fixtures_test.go::escapeRow` is non-conforming on BOTH fields and must not be reused as-is" | The RDR names the defect and then does not sequence the repair. `escapeRow` is the constructor for every escape fixture in the frozen suite; "must not be reused as-is" for ONE test while every other escape test keeps calling it means two contradictory escape fixtures coexist, and A21's empty-`RequiresOwned` composition is asserted in exactly one place while being violated everywhere else in the same package. | The suite is green and the property is untested. Six weeks later an escape row with a populated `RequiresOwned` raises `owned_state_unavailable` in production, exactly the disposition A21 says composition makes impossible. | §1, §2, AT-3 |
| C-5 | A19 Status: "Pending — blocked on a peer-document change … **This is a REVERSAL of a closed joint decision**" and Prerequisites: "**Does not block lock**" | The RDR proposes to lock while carrying an unresolved reversal of JDR 0001 §D4 (recorded Closed at §JD-12), with two mutually exclusive downstream shapes (structured `omitempty` field vs. `Detail`-flattened text) and no decision. "Reversible in both directions" is asserted, not shown: under the fallback, machine consumers must re-parse a two-level sorted array out of a prose string, and any consumer written against the flattened form breaks when the field lands. | The operator's `jq '.detail'` script, written against the first shipped release, silently stops finding the absent key when the structured field lands — `Detail` becomes a summary and the payload moves. No error; the script's `// empty` branch fires and the runbook says "no missing keys." | §1, §2, §3, premortem, AT-4 |
| C-6 | A3 Status: "Pending — the BOUND is verified; the SHAPE the spike proved is not the shape this RDR specifies" | The RDR's single load-bearing spike proved a DIFFERENT design (retained `Refusal.Guard []GuardAtom` + `UndecidedAtoms` beside it) and the RDR then reasons about blast radius, `gate` byte-identity, `escapeOrRefuse` zero-diff, and the 154/154 count as if they transferred. They do not: the spike's `gate` hunk writes `Guard: copyAtoms(lowest.Guard)` and calls `slices.MinFunc`, both of which the specified design RETIRES. "154 PASS" is the count for a superset. | Phase 1's exit condition ("the re-spike's count with the field actually removed") is unknown at lock. The implementer discovers mid-phase that the frozen dispositions were never proven under the specified shape and must re-decide contracts under schedule pressure, at which point re-deciding a frozen test is the cheap move. | §1, §2, premortem, AT-5 |
| C-7 | §SEAM normative block: "Frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` keeps its verdict but changes its REASON … To keep testing the nil-seam rule the fixture needs a value atom over a PRESENT key (e.g. `reviews >= 3`)" | The RDR diagnoses that its own change destroys the only test of the nil-seam contract, prescribes the fix in prose, and then leaves the *shipped* nil-seam narrowing (Testing Strategy row 19: "a nil-seam row that now yields a PLAN") with no frozen guard at all. Row 19's discriminating leg — `exists` over a PRESENT key under a nil seam ⇒ TRUE ⇒ plan — is a behavior change from "any non-empty guard with no seam is unevaluable" to "the kernel plans without a seam," and it is asserted nowhere in the frozen suite. | A future flow verb ships with `Input.Guards` unset (a wiring bug — there is no production importer today, so nothing catches it). Every table whose guards are pure `exists` atoms resolves to plans and writes owned tags. Previously that configuration refused. | §1, §2, premortem, AT-6 |
| C-8 | A22 Status: "Pending — kernel half Verified; the producer half is an UNSTATED INFERENCE against RDR 0002's current text" + §PRESENCE IS PROVENANCE-BLIND: "Key identity is exact string equality … canonicalization … is the normalizer's (RDR 0002)" | The kernel's entire correctness rests on a canonicalization duty that RDR 0002 does not carry, is not in 0002's direction list, and — per REQ-37's frozen banned-name test — cannot be implemented kernel-side under any exported name containing `Canonicalize`. The fallback ("narrow A22 to an INPUT PRECONDITION") makes the kernel's guarantee conditional on an unchecked precondition, which is precisely the "held by discipline" failure the RDR rejects Alternative 1 for. | A guard authored `Legal_Hold.exists = true` and a write authored `legal_hold` reach the kernel as two keys. The user gets `guard_unevaluable` naming `Legal_Hold` as missing, goes to the artifact, finds `legal_hold` present, and concludes the tool is broken. The RDR's own "If wrong" admits it: "misdiagnoses a producer defect as missing state." | §1, §3, premortem, AT-7 |
| C-9 | §SURVIVOR MEMBERSHIP: "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists" + §Failure Modes, "`unless` over a rarely-set tag" | The veto is resolution-scoped and non-escapable, and the RDR's own Failure Modes entry concedes the flagship authoring idiom ("skip when `legal_hold` is true") vetoes every artifact that never had `legal_hold`. The mitigation is a hand-written conjoined existence atom that RDR 0003's grammar has not agreed to and that the RDR hands off in Phase 4 — after lock, after implementation. | The first real table an operator writes deadlocks the whole flow: one `unless` clause over an optional tag makes `intrastate flow` refuse non-escapably on every artifact in the repo, with a refusal naming a key the author deliberately did not set. Recovery is "fix the artifact state or the accessor" — neither of which is the problem. | §1, §3, premortem, AT-8 |
| C-10 | §GATE, THEN COUNT: "The honest cost is that a row failing both ways surfaces the write-dependency problem first" + §Disposition mini-check row "Survivor missing an owned `RequiresOwned` key ⇒ `owned_state_unavailable` … Exit group: **no code row yet**" | The RDR's own headline refusal loses precedence to a refusal that has NO CLI code, NO exit-group row in RDR 0005's stable-code table, and no named owner beyond a JDR entry. The precedence ordering was inherited from shipped code whose stated rationale the RDR simultaneously declares superseded. | The user hits the flagship case, gets `owned_state_unavailable` instead of `guard_unevaluable`, and RDR 0005's renderer has no `Code` for it — so it falls through to a generic envelope. The missing-key payload this whole RDR exists to produce is not rendered at all. | §1, §2, premortem, AT-9 |
| C-11 | §Normative Contracts, empty-block clause: "An empty or omitted `unless` block is ABSENT: the `¬(unless_conj)` TERM DROPS OUT" vs. RDR 0003 Technical Design: "a row's accepted assignments are the intersection of all positive `all` atom domains **minus the single conjunctive assignment set matched by the row's full `unless` block**" | Two different algebras for the same row. 0003's subtraction over an empty `unless` block subtracts the FULL scoped product (an empty conjunction over the product is the whole product), disabling the row; 0007's rule drops the term. A9's Evidence claims 0003 "requires the same reading" — it does not; 0003 states subtraction, not term-dropping, and 0003 is Draft with no clause fixing the empty case. | Lint (RDR 0006, consuming 0003's algebra) proves a row group exhaustive; the kernel (0007's algebra) selects a different row set. The user is told the table is exhaustively covered and then gets `no_match`. JDR 0001 P5 — "a green lint means resolution succeeds" — is violated in the one direction that makes lint worthless. | §1, §3, premortem, AT-10 |
| C-12 | §Load-Bearing Decisions, Naming: "The absent-key payload is a field on `Refusal`, not a new kind" + A23 Status "Pending" | The RDR states the exact new exported surface (`Undecided []UndecidedRow`, `UndecidedRow`, `UndecidedAtom`, `ReasonAbsent`, `ReasonUncomparable`, `OpExists`, `LiteralTrue`, `LiteralFalse`, `GuardAtom`, `Block` constants) while A23 — the assumption that this surface reopens no frozen contract — is Pending and its Evidence is literally the word "Pending." `boundary_test.go::exportedKernelSymbols` walks every exported top-level name and three frozen tests (REQ-25, REQ-36, REQ-37) exact-match against banned lists. Nine-plus new exported names were never swept against those lists. | Phase 1 compiles, then `TestReq37_KernelIntroducesNoHashOrCanonicalSerialization` or a sibling fails on a name nobody thought about, and the implementer's only options are renaming the normative surface this RDR fixed by name or reopening a frozen RDR 0001 boundary test. | §1, §2, AT-11 |
| C-13 | §Approach: "A value-comparing atom over an absent key is marked unevaluable by the kernel **and the evaluator is never called for it** — the strict-routine rule (research C3)" | The SQL:2003 analogy is inverted. `RETURNS NULL ON NULL INPUT` sets the result to NULL and continues; it does not halt the statement. The RDR's rule sets the atom to U and then VETOES the entire resolution (§SURVIVOR MEMBERSHIP). SQL's actual behavior for `WHERE x = NULL` is that the row does not qualify — SQL treats U as non-selection, the very "absence-as-false at resolution scope" this RDR rejects. The cited precedent supports the atom rule and contradicts the aggregation rule; the RDR presents it as supporting both. | Nothing at first — this is a reasoning defect, and it surfaces as C-9: the resolution-scoped veto is the design's most user-hostile property and it has no prior art at all, only a citation that covers a different level. | §1, §3, premortem, AT-12 |
| C-14 | §Prerequisites: five unchecked boxes, four of them peer-document changes on RDRs that are `Draft` (0002, 0003), `Final`-and-must-re-lock (0009, 0005), or a closed JDR decision requiring reversal | The RDR's own gate lists five prerequisites and then argues one of them "Does not block lock." Three of the five (0002's two duties, 0003's placement reading, 0009's re-lock) are hard blockers with no fallback stated for 0009. A16, A22 and A24's producer halves ALL ride on the same unlanded 0002 duty — one unwritten clause in a Draft peer carries three assumptions. | Implementation starts, the normalizer emits a token the kernel does not recognize (0002 was never told), the atom is treated as a value atom, and every `exists` guard goes unevaluable. The "fails closed" story is technically satisfied and the tool is unusable. | §1, §2, AT-13 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The existence operator ships with inverted polarity

**Root cause in the RDR.** The RDR states two mutually exclusive rules about the same operator in two places, and the executed one is the wrong one.

The normative clause is unambiguous:

> Its verdict is `presence == literal`: `exists = true` decides TRUE when the key is present, `exists = false` decides TRUE when it is absent.

The only code that has ever run this rule says the opposite. A24's own Evidence records it:

> Note the A3 spike's existence atom is `{Key: "guard.missing", Op: resolve.OpExists}` with NO `Literal` set — it decides on a zero value, which the `presence == literal` rule forbids.

And the spike's `evaluateAtom`, per the spike record §4, is:

> absent key + `exists` op → decided FALSE; absent key + any other op → UNEVALUABLE; present key + `exists` → TRUE

That is `presence`, not `presence == literal`. The literal is not read at all. The spike's own probe case asserts it as a property:

> `t.Run("presence: exists atom over an absent key is FALSE, not undecidable", ...)` — "an `exists` atom over an absent key is DECIDED false, so the row is pruned"

**The specific passage that enabled it.** A3's Coverage-gap note hands that probe file to Phase 1 as the starting artifact:

> the 6-case probe at `evidence/spikes/a3-reshape/spike_kleene_test.go.txt` kills the mutant and **is the shape to start from**.

The RDR tells the implementer to start from a file that encodes the inverted rule, in a section whose whole purpose is establishing what the frozen suite fails to cover. The one place the RDR notices the discrepancy is A24 — status Pending, evidence "Pending" — which frames it as a missing-literal fail-closed question, not as a polarity inversion. The word "polarity" appears in the normative clause ("Polarity lives in the literal") and in A14; it does not appear in A24, in the spike record, or in Phase 1's instructions.

**The symptom.** An `exists = false` atom is the RDR's *sanctioned* route from absence to a verdict — MVV scenario four, Testing Strategy row 7, A5's two-row pattern, A14's entire justification, and the mitigation in the `unless`-footgun Failure Mode all depend on it. Under the inverted rule, an `exists = false` atom over an absent key decides FALSE. The row prunes under D8. The resolution reports `no_match`. `no_match` is escapable per RDR 0002. So the user authoring "apply this row when the tag was never set" gets a plan through a *different* edge — the escape — and the output carries `Escaped: true` and nothing else. This is byte-for-byte the masking path in the RDR's own Background:

> A FALSE guard plus an absent `RequiresOwned` key plus a modeled `no_match` escape yields a plan via the escape, `Escaped:true`. This is the path where missing artifact state is masked.

The RDR closes that path for value operators and re-opens it for the one operator it declares total.

### 1.2 The escape-row fixtures break on RDR 0009's kernel precondition, and A21 says they cannot

**Root cause in the RDR.** A21 makes a factual claim about a `Final` peer that the peer's text refutes. A21 Evidence:

> escape-row shape conformance "is a PRODUCER obligation on every constructor of `resolve.Row` values: a Row with a non-empty Escape list MUST have an empty Writes slice." It binds producers, NOT the kernel boundary — so this composition holds for rows a conforming producer built, **and the kernel enforces no part of it.**

RDR 0009's third normative block, sixty lines below the one A21 quotes:

> The kernel enforces the same obligation as an entry precondition of Resolve. … A table containing a row that breaches the conformance predicate above MUST cause Resolve to return a non-nil Go error identifying the offending row by RuleID and SourceLocator, with no Result disposition. The precondition is evaluated over the WHOLE table before any evaluation step.

A21 read 0009's producer clause and stopped reading. The RDR then propagates the error into Testing Strategy row 15 ("The kernel enforces NO part of this — A21 is a producer composition"), into the Capability Dependencies table, and into the `authority` mini-check.

**The specific passage that enabled it.** The Prerequisites box:

> **RDR 0009 re-locks against the loss of `Refusal.Guard`.** 0009 is `Final` and names `Refusal.Kind` / `Refusal.Guard` as the asserted properties carrying its non-vacuity evidence

This is the RDR's only engagement with 0009 as a *code* peer, and it is scoped to one deleted field. The RDR never asks what else 0009 puts in `Resolve`. It also never notices that 0009's own conformance evidence rests on the frozen suite passing at 154 — a count 0007 also claims, from a spike that ran the same suite. Both RDRs are counting the same 154 tests as evidence for incompatible reshapes of the same six fixtures.

**The symptom.** `fixtures_test.go::escapeRow` — the constructor behind every escape fixture in ADV-1b, ADV-2, ADV-2b, Fixup-1d and Fixup-1e — is:

```go
RequiresOwned: []string{"status"},
NextTags:      []resolve.Tag{{Key: "status", Value: "Blocked"}},
Writes:        []resolve.Tag{{Key: "status", Value: "Blocked"}},
Escape:        []resolve.RefusalKind{class},
```

Non-empty `Writes` on a non-empty `Escape`. Under 0009's precondition, `Resolve` returns `ErrEscapeShapeBreach` and no `Result` — so `mustResolve(t, in)` fails before any 0007 assertion runs. Whichever RDR implements second finds five tests failing on the other RDR's contract, with a stack trace pointing at a shared fixture constructor. 0007 explicitly declines to fix it: Testing Strategy row 15 says only that ONE test must not reuse `escapeRow` "as-is," leaving the other five on the non-conforming constructor.

The user-visible end of this is not a test failure; it is what the implementer does under pressure. The cheapest resolution is to strip `Writes` from `escapeRow` globally — which changes ADV-2b's premise ("an escape can emit owned-tag Writes derived from owned state the accessor snapshot never produced"), silently vacating an adversarial test that exists to prove escape rows cannot launder writes.

### 1.3 The payload sort is not total, and REQ-1 replay determinism breaks intermittently

**Root cause in the RDR.** The RDR reasons its way to a six-field tuple and declares it total by construction. It is not.

> The ordering is therefore the tuple `(RuleID, SourceLocator, key, block, operator token, literal)`, compared field by field in that order. Those are the row's identity plus the atom's four §D1 fields, so the tuple is total by construction: two entries equal on all six name the same atom, which the kernel MUST NOT report twice.

Three defects stack here.

First, `literal` is not a scalar for two of the five operators. RDR 0003's matrix gives `in` a "non-empty typed scalar set" and `contains` a "non-empty typed element set." A set has no canonical string order unless somebody specifies one, and Phase 3 of this very RDR concedes nobody has:

> `resolve.Tag.Value` is a bare `string`, so a set-valued tag reaches the seam as one opaque string with no stated element encoding. … This RDR deliberately does not invent the encoding.

The RDR declines to specify the encoding and then sorts on it.

Second, "two entries equal on all six name the same atom" is false as stated. RDR 0003's Identity decision keys an atom by "its source rule/context id **plus its position within `all` or `unless`**" — position is part of atom identity precisely because `(tag, operator, literal)` is not unique. Two textually identical atoms at different positions in one block are two atoms by 0003's rule and one entry by 0007's tuple.

Third, the MUST-NOT-report-twice clause is a prohibition with no mechanism. The RDR says the kernel must not emit duplicates; it does not say dedup on what, or before or after the sort, or whether a duplicate is a producer defect the kernel should surface. Row 17 of the Testing Strategy asserts permutation-equality, which a dedup-then-sort satisfies and a sort-with-duplicates does not.

**The specific passage that enabled it.** The RDR argues itself into confidence by ruling out the *smaller* tuple:

> Row identity then key is NOT total: this RDR's own sanctioned idioms put two atoms on one key in one row … Two entries would then tie on `(row, key)` and differ only in block, operator, or reason, leaving the order to an unstable tie-break

Having found one insufficient tuple and extended it, the RDR treats extension as proof and stops. It never asks whether the added fields are themselves totally ordered.

**The symptom.** For most tables, nothing — single-atom guards and scalar literals sort fine. Then a table with two `in` atoms over one key goes unevaluable, and the JSON refusal payload comes back with the two entries in map-iteration order or in whatever order the literal-set string happened to serialize. The same table, same artifact, same revision, run twice, produces two different `--json` payloads. RDR 0001 REQ-1 says the same tuple always replays the same disposition; ADV-3 freezes the row-order half of it. The operator diffing two refusal payloads in a CI log sees a change and starts looking for a state change that never happened.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`#### Normative Contracts` → the payload block** (the clause beginning "The payload is named surface, not shape-by-description" through "Any contract or test distinguishing WHICH row went unevaluable MUST assert on the row entries").

It will be rewritten because it is the only place in the RDR where a `Pending` assumption (A23, evidence: the literal word "Pending"), an unresolved reversal of a closed joint decision (A19), an unspecified value encoding (Phase 3's concession), and a frozen boundary test (`boundary_test.go::exportedKernelSymbols`, consumed by REQ-25/36/37 exact-match guards) all intersect on the same nine exported identifiers. Every one of those four is individually enough to force an edit:

- **A23 resolves against the frozen suite.** The RDR asserts the payload "is additive on `Refusal` and reopens no frozen kernel contract" and then marks the assumption Pending with no evidence. If any of `Undecided`, `UndecidedRow`, `UndecidedAtom`, `ReasonAbsent`, `ReasonUncomparable`, `GuardAtom`, `OpExists`, `LiteralTrue`, `LiteralFalse` trips a banned-name list, the clause changes.
- **A19 resolves either way.** Under the structured field the block stays; under the `Detail` fallback the block must state a flattening, which it currently forbids ("serializing a two-level sorted array into it puts a second, undocumented encoding inside a string").
- **The sort tuple is underspecified** (C-2). The first set-valued literal in a refusal payload forces either an element encoding (which the RDR says is 0003's) or a different sort key.
- **RDR 0005 must render it.** The clause says "RDR 0005's renderer type against it" — but 0005 is `Final`, its stable-code table has no row for `owned_state_unavailable`, and its `flow-guard-unevaluable` row is scoped to "guard cannot be evaluated from **supplied facts**," which the RDR's own `disposition` mini-check flags as mis-signalling this RDR's flagship missing-*artifact*-state case. The renderer cannot be written from this clause plus 0005 as they stand.

The RDR half-knows this. It says "the payload is named surface, not shape-by-description" — an unusually strong commitment — and in the very next paragraph concedes "The spike's `UndecidedAtoms` is NOT this shape." The section names surface no implementation has produced, no test has exercised, and no consumer has agreed to consume. That is the definition of a section that gets rewritten on contact.

---

## 3. The one assumption that will not survive first contact with a real user

**A5** — that authors who need "row applies when tag X is absent" will express it with an existence atom, and that the disjunction "absent OR equals v" is written as two rows with the value row carrying a conjoined `X exists = true`.

> Authors who need "row applies when tag X is absent" express it with an existence atom; the disjunction "absent OR equals v" needs two rows (RDR 0003 rejects intra-guard disjunction), and the value row MUST itself carry an `X exists = true` atom beside `X eq v`, or it survives unevaluable when X is absent and the aggregation veto refuses the absent leg.

Its status is "Verified (runtime); the load-time half rides on A12" — and A12 is `Deferred`, because RDR 0003 is silent on how an existence atom projects onto the declared-domain product. So the assumption is: authors will reliably write a two-row pattern whose *legality* nobody has established.

Here is what actually happens. The user does not write two rows. The user writes one row with `status.eq = "Draft"` and a second condition, forgets that `legal_hold` might not exist on older artifacts, and ships. Every artifact that predates the tag goes `guard_unevaluable`, non-escapably, across the whole flow. The RDR's own Failure Modes section already describes this as the expected outcome and calls it "correct by the rule and a footgun in practice":

> an `unless` block whose only atom is a value atom over an absent key is `¬U = U` and vetoes — correct by the rule and a footgun in practice ("skip when `legal_hold` is true" refuses on every artifact that never had `legal_hold`).

The mitigation is the conjoined idiom. Three things stop it from working on a real user:

1. **It is not in the grammar yet.** A5's own Evidence: "Whether RDR 0003's TOML grammar lets an author *write* it is 0003's (Draft) and is handed over in Phase 4." Phase 4 is after implementation. The user gets the kernel before the authoring rule.
2. **The `all`-placement question is open.** A14 records that the cluster gate read RDR 0003's "Positive guard atoms MUST live in `all`" as contradicted by an `exists = false` literal in `all`, and routes it to §JD-1 unresolved. So the sanctioned single-atom absence test may be ungrammatical in the block it needs to live in.
3. **It is boilerplate for a default.** Every optional tag every guard touches needs a paired existence atom. The moment an author writes the same two-line incantation four times, they will write it zero times on the fifth, and the difference between "I forgot the existence atom" and "I meant to require the tag" is invisible in review. The RDR's whole reviewability argument (RDR 0002: the table is reviewable data) depends on the table reading as intent; a table where half the atoms are load-bearing and half are ritual does not.

The RDR anticipates the *class* of failure and names its own worst case correctly ("authors cannot express a legitimate absence-conditional row and fall back to sentinel-stamping upstream — the anti-pattern whose only sanctioned alternative *is* this pattern"). What it does not do is notice that the sanctioned alternative is unavailable at ship time and boilerplate afterward. Sentinel-stamping is not a risk here; it is the predicted outcome.

---

## 4. Premortem

*Written from six months after ship.*

We shipped Phase 1 and Phase 2 in the same fortnight. The kernel reshape went exactly as the A3 spike predicted for the parts the spike covered, which is what made everything after it hard to see.

**Week 1 — `evaluateAtom` ships inverted.** The implementer opened `spike_kleene_test.go.txt`, because A3 said it "is the shape to start from," and lifted `evaluateAtom` out of the spike record's §4 summary. That code reads: absent key + `exists` → FALSE; present key + `exists` → TRUE. It never reads `atom.Literal`. The normative clause said `presence == literal`; nobody diffed prose against the probe. Both spike probe cases involving `exists` passed, because both were written against the same inverted rule. Testing Strategy row 7 ("`exists` total across BOTH polarity channels") was written in Phase 2, and the Phase 2 author wrote it against the *code*, asserting `exists = false` over an absent key prunes the row — because that is what the kernel did, and the kernel was the source of truth for a test being backfilled. The RDR's clause was never the oracle.

**Week 3 — `flow resolve` on the RDR flow.** The first real transition model had a rule "an RDR with no `blocked_on` tag advances to prelock," written as `blocked_on.exists = false` in `all`. Every RDR advanced correctly *except* the ones that had once been blocked and cleared — because a cleared tag normalizes to a `<clear>` write, not a deletion (RDR 0002: "Absence from both the write block and the clear list MUST NOT imply deletion"), so `blocked_on` was present-with-`<clear>`-value. The row that should have fired for genuinely-never-blocked RDRs did not fire, because `exists = false` over an absent key had been implemented as FALSE. `Resolve` returned `no_match`. The table had a modeled `no_match` escape — the operations team had added one months earlier so the flow would not hard-stop — and `escapeOrRefuse` planned through it with `Escaped: true`. Three RDRs were advanced to prelock by the escape edge instead of the intended edge, with the escape's `NextTags` (`status: Blocked`) rather than the intended ones.

Nobody noticed for eleven days, because `Escaped: true` is a field in the JSON plan and not an error, and the runbook check was on exit code.

**Week 4 — the first non-escapable deadlock.** A different flow had `[rule.guard.unless] legal_hold.eq = true`. Written exactly as RDR 0003's illustrative TOML shows. Every artifact that had never carried `legal_hold` produced `¬U = U` at the block, `U` at the row, and `gate`'s undecidable branch returned `guard_unevaluable`. Non-escapable by design. The whole flow stopped for every artifact in the repo simultaneously. The refusal named `legal_hold` with reason `absent` — accurate and useless, because the author's intent was "we do not have a legal hold," which is the absence itself. The remedy in the docs read "make the state readable, or rewrite the guard with an existence atom if absence was intended." RDR 0003 was still Draft; the conjoined-atom idiom had never been written down anywhere an author would look; and `[rule.guard.unless.legal_hold]` with both an `exists` and an `eq` sub-key was not something the normalizer had ever been asked to emit. We shipped a hotfix that stamped `legal_hold = "false"` on every artifact — the sentinel-stamping anti-pattern the RDR names by name in its own Risks section.

**Week 5 — `Refusal.Undecided` diffs on identical input.** The CI job that captures refusal payloads for regression started flapping. Two `in` atoms over one key, both unevaluable, and `compareAtoms` compared their literals — which for a set literal was whatever `strings.Join` over an unordered slice produced that run. The sort tuple `(RuleID, SourceLocator, key, block, operator token, literal)` was total on paper and not on values, because nobody had specified the encoding of a set literal (Phase 3 had said so explicitly and deferred it to RDR 0003, which had not shipped). `TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder` still passed — it permutes rows, not atoms. Testing Strategy row 17's tie case had been written with two `eq` atoms, which are scalar and sort fine.

**Week 7 — RDR 0009 lands and takes out five frozen tests.** `Table.CheckValid()` ran over `fixtures_test.go::escapeRow`'s non-empty `Writes` and returned `ErrEscapeShapeBreach`. `mustResolve` fataled in ADV-1b, ADV-2 (both sub-cases), ADV-2b, Fixup-1d and Fixup-1e. A21 had said, in writing, "the kernel enforces no part of it." The implementer, believing that, spent two days assuming their own `gate` change had broken the fixtures before finding 0009's third normative block. The fix — strip `Writes` from `escapeRow` — silently vacated ADV-2b, whose entire premise is that an escape row *can* carry writes derived from unavailable owned state. That test now constructs a row that cannot express the attack it defends against. It is green.

**Week 9 — the JSON consumer.** The skill wrapper had been parsing `.detail` for the missing key, because A19's reopening had been declined and the fallback shipped: the payload flattened into `Detail` as prose. Then §D4 was reopened anyway in a later cluster pass, the `omitempty` structured field landed, `Detail` shrank to a one-line summary, and the wrapper's `jq -r '.detail | capture("key=(?<k>\\S+)")'` returned empty. Its fallback branch reported "no missing keys" and the operator retried the same command four times.

**What we got right, and why it did not save us.** The kernel never folded absence into false for a value operator. `evaluateAtom` genuinely never called the seam on an absent key. Strong-Kleene combination was correct, and the multi-atom tests the spike's mutation analysis demanded were written and killed the mutant. The single defect in the domain rule was in the one operator the RDR called *total* and therefore did not model as a combination case — and the aggregation veto, the property with no prior art behind it (SQL:2003 sets the value to NULL and keeps going; it does not abort the statement), is what turned every remaining defect from a wrong answer into a stopped flow.

**Named functions:** `resolve.go::evaluateAtom` (inverted existence polarity), `resolve.go::evaluateGuard` (correct), `resolve.go::gate` (undecidable branch, unchanged partition, veto), `resolve.go::escapeOrRefuse` (planned the wrong edge in week 3), `resolve.go::compareAtoms` (non-total on set literals), `resolve.go::Table.CheckValid` (0009, week 7), `fixtures_test.go::escapeRow` (five broken tests, one vacated), `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` (re-decided under pressure), `adversarial_test.go::TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement` (vacated).

**Named user journeys:** advance an unblocked RDR to prelock (week 3, wrong edge, silent); run `flow` on any artifact under a `legal_hold` guard (week 4, total deadlock); diff two refusal payloads in CI (week 5, phantom change); read the missing key out of a `--json` refusal in the skill wrapper (week 9, silent empty).

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time tests: each is checkable against the RDR document and its evidence, before any code exists.

### AT-1 — Existence polarity agrees between clause, spike, and probe (catches C-1)

```gherkin
Scenario: the executed existence rule matches the normative existence rule
  Given the normative clause "an existence atom's verdict is presence == literal"
  When I read every artifact the RDR hands Phase 1 as a starting shape
    (the A3 spike record's evaluateAtom summary, and
     evidence/spikes/a3-reshape/spike_kleene_test.go.txt)
  Then every existence decision in those artifacts MUST read the atom's Literal
  And an exists atom with Literal = LiteralFalse over an ABSENT key
    MUST be recorded as deciding TRUE
  And no artifact may assert that an exists atom over an absent key
    decides FALSE or prunes the row
```

*Would have fired.* The spike record §4 states "absent key + `exists` op → decided FALSE" and the probe's fifth case asserts `KindNoMatch` with the comment "an `exists` atom over an absent key is DECIDED false, so the row is pruned." A24's Evidence half-notices — it flags the missing `Literal` — and stops at fail-closed rather than at polarity.

### AT-2 — The payload sort key is total over the operator vocabulary (catches C-2)

```gherkin
Scenario: every field in the sort tuple has a stated total order
  Given the sort tuple (RuleID, SourceLocator, key, block, operator token, literal)
  When I enumerate RDR 0003's five operator classes
  Then for each operator, the literal's ordering MUST be specified in this RDR
    or in a cited normative clause of a peer RDR
  And for `in` and `contains`, whose literals are SETS, the element encoding
    and element order MUST be cited, not deferred
  And the RDR MUST state a dedup rule (on which fields, at which step)
    rather than only a MUST-NOT-report-twice prohibition
```

*Would have fired.* Phase 3 says outright "This RDR deliberately does not invent the encoding — element representation is part of RDR 0003's typed operator semantics," and the payload clause sorts on it anyway. The dedup step is stated as a prohibition with no mechanism.

### AT-3 — Every kernel-boundary claim about a Final peer quotes that peer's kernel clauses (catches C-3, C-4)

```gherkin
Scenario: A21's "the kernel enforces no part of it" is checked against all of RDR 0009
  Given A21 asserts RDR 0009 "binds producers, NOT the kernel boundary"
  When I grep RDR 0009's normative blocks for "kernel" and "Resolve"
  Then I find "The kernel enforces the same obligation as an entry precondition
    of Resolve … MUST cause Resolve to return a non-nil Go error"
  And A21's claim is refuted
  And the RDR MUST state the sequencing of RDR 0009's Table.CheckValid against
    Phase 1's fixture migration, naming fixtures_test.go::escapeRow, whose
    Writes field breaches 0009's predicate in five frozen tests
```

*Would have fired.* A21 quotes 0009's producer clause verbatim and stops one clause short.

### AT-4 — No unresolved reversal of a closed joint decision rides into lock (catches C-5)

```gherkin
Scenario: A19's two shapes are decided before lock, not after
  Given A19 records "This is a REVERSAL of a closed joint decision"
    against JDR 0001 §D4, recorded Closed at §JD-12
  And Prerequisites assert the reversal "Does not block lock"
  When I ask what a downstream JSON consumer parses in the interim
  Then the RDR must name ONE shape, because Detail-as-prose and a structured
    omitempty field are not interchangeable to a consumer that has already
    shipped a parser against one of them
  And "reversible in both directions" must be shown for CONSUMERS,
    not only for the envelope's append-only property
```

*Would have fired.* The RDR argues reversibility from `CLIError`'s "append-only" doc comment, which is a claim about the envelope's schema, not about parsers written against the flattened form.

### AT-5 — The load-bearing spike proved the specified shape (catches C-6)

```gherkin
Scenario: A3's evidence is evidence for A3's claim
  Given A3 is Pending because "the SHAPE the spike proved is not the shape
    this RDR specifies"
  When I read the numbers the RDR carries forward from that spike
  Then "154 PASS / 0 FAIL", "escapeOrRefuse needed ZERO changes",
    "gate's partition is byte-identical", and "all six pass re-encoded"
    are all measurements of the SUPERSET shape
  And the RDR MUST NOT state a Phase 1 exit condition it cannot compute
    ("the re-spike's count with the field actually removed")
  And the spike MUST be re-run before lock, because gate's only hunk in the
    superset spike writes Guard: copyAtoms(lowest.Guard) and calls
    slices.MinFunc — both of which the specified design retires
```

*Would have fired.* The RDR says all of this about itself and then locks anyway.

### AT-6 — Every behavior the RDR narrows has a test that would fail if it were not narrowed (catches C-7)

```gherkin
Scenario: the nil-seam narrowing has a discriminating guard
  Given Testing Strategy row 19 states a nil-seam row "now yields a PLAN"
    where the shipped kernel refuses
  When I look for the frozen or added test that pins the DISCRIMINATING leg
    (an exists atom over a PRESENT key with Input.Guards == nil ⇒ plan)
  Then row 19 must name a test, not a narrowing
  And the RDR must state what protects a caller that leaves Input.Guards nil
    by mistake, given that shipped evaluateGuard refuses every guarded row
    in that case today and the new rule plans
```

*Would have fired.* The RDR diagnoses that Fixup-1d stops testing the nil-seam rule, prescribes a fixture change in prose, and adds no test for the direction that newly produces a plan.

### AT-7 — No kernel guarantee rests on an unwritten peer clause (catches C-8, and C-14)

```gherkin
Scenario: A16, A22, A24's producer halves have a landed home
  Given A16 (existence token emission), A22 (key canonicalization),
    and A24 (typed boolean literal) all state their producer half
    lands in RDR 0002
  When I read RDR 0002's Refinement Context Direction list
  Then it names only the RequiresOwned producer and the fixture rename
  And three assumptions of a foundational RDR ride on one unwritten clause
    in a Draft peer
  And the kernel-side fallback ("input precondition, unchecked") is the
    same "held by discipline" property this RDR rejects Alternative 1 for
  And REQ-37's frozen banned-name list contains "Canonicalize", so the
    kernel cannot implement the fallback under the obvious name
```

*Would have fired.* The RDR states the gap itself, in A16 ("The PRODUCER half is not yet a duty on 0002"), A22 ("an UNSTATED INFERENCE against RDR 0002's current text"), and A24 ("rides on the same unlanded 0002 duty as A16 and A22") — three times, without treating the concentration as a single point of failure.

### AT-8 — The aggregation veto is validated against a realistic authored table (catches C-9)

```gherkin
Scenario: the veto does not deadlock the reference fixtures
  Given the resolution-level veto refuses non-escapably when ANY surviving
    candidate is unevaluable
  When I classify every authored guard in
    0003-…/evidence/spikes/guard-fixture.toml and
    0002-…/evidence/spikes/rdr-fixture.toml as safe-or-migration
  Then that classification MUST exist before lock, not in Phase 2
  And for every guard classified "migration", the RDR MUST show the migrated
    form is authorable in RDR 0002's current TOML shape
    ([rule.guard.unless.<tag>] carrying both an exists and an eq sub-key)
  And rdr-fixture.toml's `unless.profile.eq = "small"` MUST be walked
    against an artifact with no profile tag
```

*Would have fired.* The RDR schedules the classification into Phase 2 ("classify each authored guard in the reference fixtures as safe-or-migration") — after lock — while its own Failure Modes section already predicts the deadlock for exactly this authoring shape, which is the shape RDR 0002's illustrative TOML uses.

### AT-9 — The winning refusal on the flagship case is renderable (catches C-10)

```gherkin
Scenario: the flagship user journey produces a refusal RDR 0005 can render
  Given GATE, THEN COUNT puts owned_state_unavailable ahead of guard_unevaluable
  And the RDR's disposition mini-check records owned_state_unavailable's
    exit group as "no code row yet"
  When a row fails BOTH ways — the RDR's own Testing Strategy row 13
  Then the user gets a refusal with no stable Code in RDR 0005's table
  And the per-atom payload this RDR exists to produce is not rendered
  And the RDR must either supply the Code row, cite where it lands, or
    justify the precedence against the renderability gap
```

*Would have fired.* The RDR records both facts — the precedence and the missing code row — in two different sections and never joins them.

### AT-10 — The kernel's row algebra agrees with the lint algebra (catches C-11)

```gherkin
Scenario: empty unless block has one meaning across the cluster
  Given this RDR: an omitted unless block means the ¬(unless_conj) term
    DROPS OUT, and the verdict reduces to all_result
  And RDR 0003: "a row's accepted assignments are the intersection of all
    positive all atom domains MINUS the single conjunctive assignment set
    matched by the row's full unless block"
  When the unless block is empty
  Then 0003's subtraction removes the empty conjunction's match set —
    which over a scoped product is the WHOLE product — disabling the row
  And this RDR's rule keeps the row
  And A9's claim that 0003 "requires the same reading" is unsupported:
    0003 states subtraction, never term-dropping, and fixes no empty-block
    identity anywhere
  And JDR 0001 P5 ("a green lint means resolution succeeds") is violated
    in the direction that makes lint's blocking authority worthless
```

*Would have fired.* A9's evidence argues from RDR 0002's silence and RDR 0003's subtractive algebra, and reads the subtractive algebra as agreeing when it is the thing that disagrees.

### AT-11 — New exported surface is swept against the frozen boundary tests (catches C-12)

```gherkin
Scenario: A23 is resolved against boundary_test.go before lock
  Given the RDR names Undecided, UndecidedRow, UndecidedAtom, ReasonAbsent,
    ReasonUncomparable, GuardAtom, OpExists, LiteralTrue, LiteralFalse
    as normative exported surface
  And boundary_test.go::exportedKernelSymbols enumerates every exported
    top-level name
  And REQ-25, REQ-36, REQ-37 exact-match against banned lists including
    Parse, Export, Canonicalize, Apply
  When I sweep the nine names against the three lists
  Then the sweep MUST be recorded, because A23's Evidence is currently
    the word "Pending"
  And any helper the design implies (a literal parser, a canonicalizer)
    must be checked for a banned exported spelling before it is specified
```

*Would have fired.* A23's own status admits the sweep was never done: "the names are new normative surface this RDR now states, so they must be checked against the frozen suite before lock rather than assumed."

### AT-12 — Every prior-art citation is checked at the level it is used (catches C-13)

```gherkin
Scenario: SQL:2003 strict routines support the atom rule, not the veto
  Given the RDR cites RETURNS NULL ON NULL INPUT as the precedent for
    "the kernel marks an absent-key atom unevaluable and never calls
    the evaluator for it"
  When I ask what SQL does with the resulting NULL
  Then SQL propagates it and CONTINUES: a WHERE clause evaluating to
    UNKNOWN does not select the row and does not abort the statement
  And SQL's actual treatment at resolution scope is "unknown ⇒ not selected"
    — the exact absence-as-false-at-resolution-scope this RDR rejects
  And the resolution-level veto (SURVIVOR MEMBERSHIP) therefore has
    NO prior art in this RDR: SCXML continues after its error event,
    SQL continues after its NULL
  And the RDR MUST either find a precedent for halting or state plainly
    that the veto is unprecedented and argue it on its own terms
```

*Would have fired.* The Decision Rationale presents C3 as the deciding row for the whole approach; the citation covers the argument-passing rule and is silent on what a host does after.

### AT-13 — Prerequisites are blockers or they are not prerequisites (catches C-14)

```gherkin
Scenario: five unchecked prerequisite boxes at lock
  Given Prerequisites lists JD-3 landing in 0002, 0003's re-entry reading,
    the §D4 reopening, two duties added to 0002's direction list, and
    RDR 0009's re-lock
  And one of them is annotated "Does not block lock"
  When I ask which of the remaining four have stated fallbacks
  Then A19 and A22 have fallbacks; A16's producer half and RDR 0009's
    re-lock do not
  And RDR 0009's re-lock is the one with a live code conflict (C-3),
    with no fallback and no sequencing
  And the RDR must not lock with an unsequenced hard dependency on a
    Final peer's kernel precondition
```

*Would have fired.* The Prerequisites box states all five and reasons about blocking only for A19.
