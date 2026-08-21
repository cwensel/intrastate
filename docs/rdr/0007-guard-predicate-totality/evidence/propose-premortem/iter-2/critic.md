Model: claude-fable-5

# Premortem critic — RDR 0007 guard-predicate totality, propose iter-2

Scope: kernel-enforced presence check + kernel strong-Kleene combination over parsed atoms; pluggable evaluator narrowed to "compare one present value to one literal under one operator". Worked from the brief only; no repository files read.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|---|---|---|---|---|
| P-1 | Claim 1 "evaluator never called for an absent key" | The kernel's notion of *present* is "key in assembled view". The view assembly is upstream of the kernel and caller-supplied observed tags enter unconstrained; an observed tag with the same key as an owned tag but an empty/placeholder value (`""`, `"null"`, `"-"`) is *present*, so the kernel hands it to the evaluator, which compares `"" != "approved"` → FALSE → pruned → escapable `no_match`. Absence masked one hop earlier than the seam. | Escape row fires; no refusal; author never learns artifact state was missing | hindsight / negation C1 |
| P-2 | Claim 1, "provenance-blind" presence | Provenance-blindness is load-bearing for *presence* but the brief says "missing owned state among survivors is reported" in gate 2. An *observed* tag can satisfy presence for a key that the row's owned-state declaration expects *owned*. Kernel decides the guard on observed data, gate 2 sees the key present and says nothing. Caller-supplied tags shadow artifact state. | Row selected on a value the caller typed, not artifact state | negation C1 / refutation |
| P-3 | Claim 2 "`exists` is cardinality, not grammar" | The kernel now owns one token of the grammar. The grammar spec (separate, later, a moving peer) renames/aliases it (`present`, `has`, `exists:false` sugar, `!exists`, `missing`). Kernel matches neither → treats as non-`exists` → absent key → unevaluable. Every table that uses the new spelling refuses non-escapably. Conversely, a grammar that adds `exists` *with a path* or a *typed* presence (`exists_as:int`) is silently decided as raw presence. | Tables authored to the grammar spec refuse with `guard_unevaluable` naming a key the author deliberately asked `exists` about | negation C2 / refutation seed |
| P-4 | Claim 2, literal for `exists` | `presence == literal` requires the kernel to parse the literal as a bool. Literal is a string from the parser (`"true"`, `"yes"`, `"1"`, `"True"`). Kernel's parse disagrees with grammar's bool domain → `exists yes` is FALSE on a present key → pruned → escapable. | `exists` row silently never fires; escape row masks | negation C2 |
| P-5 | Claim 3 "`unless` empty → TRUE contribution" | Under strong-Kleene with `row = all ∧ ¬(unless)`: an `unless` block whose *only* atom is unevaluable yields `¬U = U`, so the whole row is U and vetoes the resolution. A table author writes `unless: legal_hold == true` intending "skip if on hold"; on artifacts with no `legal_hold` tag ever written (the common case), every resolution is refused. Correct by the stated rule, catastrophic in practice, and the only escape is rewriting every `unless` as `unless: exists legal_hold ∧ legal_hold == true` — which the kernel does *not* short-circuit: `F ∧ U = F`, fine, but authors will not know to write it. | All resolutions refuse on artifacts lacking a rarely-set tag | negation C3 |
| P-6 | Claim 3, "safe and complete" | Strong-Kleene ∧ is commutative but the brief's combination is across *blocks*; if the kernel evaluates `all` first and short-circuits on F (never consulting `unless`), vs evaluates both, verdicts agree — *but* the evaluator side-effects / which absent keys get *named* in the refusal differ. Two kernels agreeing on verdict but disagreeing on the refusal payload → snapshot tests flake, CLI envelope unstable. | Refusal names a different key set depending on atom order | negation C3 |
| P-7 | Claim 4 "seam loses nothing" | Cross-atom operators: `a == b` (tag-to-tag compare), `updated_at > now - 7d` (needs clock/context), `status in $approved_set` (needs a named set from the grammar spec), `count(reviewers) >= 2` (needs multi-valued tag). Each needs the view or context the seam no longer passes. The later grammar spec either forbids them (narrowing the language by accident, decided here by a kernel RDR) or the evaluator re-acquires the view via a side channel — the exact grammar-agnosticism violation the approach claims to prevent. | Grammar spec is forced to ship without relational/temporal guards, or evaluator gets a back-door view handle | negation C4 |
| P-8 | Claim 4, literal typing | The seam passes a *string* literal and a *string* value; the evaluator must know the declared kind to compare `"10" < "9"` correctly. Kinds live in tag declarations the kernel does not have. So the seam must also pass the key (so the evaluator can look up the declaration) — meaning the evaluator *does* hold per-key state and can return whatever it likes for a key it "knows" is absent by its own lookup. Claim 1's structural guarantee depends on the evaluator never being able to infer absence; with the key in hand it can. | An evaluator that consults its own declaration table folds "present but kind-mismatched" into false | negation C1+C4 |
| P-9 | Claim 5 "unevaluable for uncomparable present value is fine" | Inverse of masking: a present but malformed value (caller-supplied observed tag `"priority": "high!"` where kind is int) makes *every* row mentioning `priority` U → whole resolution refused, non-escapably. A caller can therefore **deny** any resolution by supplying one garbage observed tag with a key that any row's guard mentions. The refusal names `priority` as *absent*? No — it is present. The payload field is "absent keys"; this case has none, so the refusal is empty/misleading. | `guard_unevaluable` with an empty absent-key list; or a list that lies | negation C5 / C6 |
| P-10 | Claim 5 | Evaluator returns U for an operator token it does not recognise (grammar drift again). Every row using the new operator is U → veto. There is no distinct refusal kind for "grammar/evaluator mismatch"; it is folded into `guard_unevaluable` (closed five-kind taxonomy) and named as if artifact state were missing. | Author told a tag is missing when the real cause is evaluator version skew | negation C5 / refutation seed |
| P-11 | Claim 6 "naming absent keys is actionable, no leak" | The refusal is rendered in a CLI JSON envelope. Absent-key names reveal which tags the table's guards consult (the table is the policy; the keys are the policy surface) to whoever can trigger a resolution. With caller-supplied observed tags unconstrained, an attacker iterates: supply key K, see whether the refusal stops naming K — enumerates the full guard key set, and the literals by bisection once present (FALSE-prune vs TRUE-select is observable via exit code). Whether this matters depends on whether tables are secrets, which the RDR does not state. | Policy surface enumerable through refusal payloads | negation C6 |
| P-12 | Claim 6, "replacing the old guard-text field" | Removing the guard text from the refusal drops the *row* identity. "Absent: legal_hold" is actionable only if the author knows which row(s) mentioned it; in a 40-row table with 3 rows mentioning `legal_hold` under different blocks, they cannot tell whether the veto came from a row they expected to fire or a sibling. Actionability regressed relative to guard text + row id. | Author cannot locate the vetoing row | negation C6 |
| P-13 | Claim 7 "FALSE only from decided atoms" | True for a *single* atom. For a *row*, `F ∧ U = F` means a row with one present-and-false atom and one absent-key atom is pruned, and the absent key is **never reported** (prune-first precedes the unevaluable gate). If a sibling row with only the absent-key atom exists, the veto fires and names it; if no sibling mentions it, the absence is silently dropped. Whether missing state is reported depends on table shape, not on the artifact. Not masking per the stated rule — but the author's mental model ("you always tell me when needed state is missing") is false. | Same absent tag: reported on one table, silent on another | negation C7 |
| P-14 | Claim 8 "zero production importers" | This is the seed class exactly: a negative existential ("no importers") about a package whose importers a sibling spec (guard grammar, CLI wiring) has not written yet, stamped with no re-verification trigger. By the time 0007 implements, the grammar spec's evaluator and the CLI renderer will both depend on the row type and on `GuardResult`; "no migration cost" was true at propose time only. | Implementation conflicts with two sibling implementations that targeted the old row type | refutation seed (ii) |
| P-15 | Rejection R1 reason "conformance harness nothing forces the evaluator to run" | Same is true of the *kernel-side* rule: nothing forces the grammar spec's author to build against the kernel's operator token or bool-literal convention. The drift risk is not removed, it is moved to the kernel/grammar boundary (P-3, P-4, P-10) where it manifests as a wrong-kind refusal rather than masking. The RDR must say *why* that direction of drift is preferable (it is — fail-closed — but it must be stated and tested). | — | negation R1 |
| P-16 | R1 "kernel could not name absent keys" | A row-level seam returning `(GuardResult, absentKeys []string)` names them just as well. The reason is false as stated; the true reason is trust (who may decide absence), not capability. | — | negation R1 |
| P-17 | R2 reason "kernel has no tag declarations" | Yet the kernel now parses a bool literal for `exists` (P-4) and decides presence with a string-keyed map whose key normalisation (case, namespace prefix `ns:key`, trailing whitespace) is a declaration-level concern. The kernel already does a slice of R2; the RDR must draw the line explicitly: key identity = exact string equality after what normalisation? | `Exists Legal_Hold` vs `legal_hold` | negation R2 |
| P-18 | R3 "no reconstruction ⇒ no mapping failure" | The parse moved to the *loader*; an unparseable guard now fails at load, before Resolve. Fine — but the brief's refusal taxonomy is closed at five kinds, and load-time parse failure has no kind. If the CLI wires load+resolve in one command, the user sees a Go error, not an envelope. | Unstructured error for malformed guard where every other failure is an envelope | negation R3 |
| P-19 | R5 "precheck duplicates the atom key set" | True, but the approach *also* has a second key set: the row's owned-state declaration (gate 2). Atom keys and owned-state keys can drift from each other: a guard mentions `approver` but owned-state does not list it; an observed tag satisfies presence, gate 2 is silent. Which is P-2 in another guise. The RDR should state whether atom keys are required ⊆ owned-state keys, or not. | — | negation R5 |
| P-20 | Three-valued `GuardResult` reuse | Existing tests pin prune-first and escape gating over a *row-level* three-valued result. Replacing it with kernel-computed verdicts keeps the type but changes who produces it; the adversarial suite will keep passing because it injects row verdicts directly, **not** atoms. The suite no longer exercises the code path that matters. | Green suite, untested combinator | consumer artifact |

## 1. Prospective hindsight — the failure narrative

0007 shipped in week 3. The kernel grew `atom` and `combineRow`; the evaluator interface shrank to `Compare(key, op, literal, value string) GuardResult`. All 61 existing `internal/resolve` tests passed unchanged because they construct `Input` with row verdicts already set — none of them built a row from atoms. Two new tests covered `exists` and one absent-key veto.

Week 5, the guard-grammar spec (0009) landed its evaluator. Its author read the grammar, not the kernel, and spelled presence `has`, with `exists` kept as an alias *in the evaluator*. The kernel did not know `has`; every `has X` atom over an absent key went to the unevaluable path and every `has X` over a present key went to the evaluator, which compared the value to `"true"` and returned FALSE. Tables written to the grammar doc had presence checks that were non-escapably refused on absence and pruned on presence — exactly backwards — and nobody noticed because the one integration fixture used `exists`.

Week 6, the CLI layer wired `resolve` and rendered `guard_unevaluable` with the new `absent_keys` field. The first customer table used `unless: legal_hold == true` on a corpus where 94% of artifacts had never had `legal_hold` written. Every resolution refused. The refusal named `legal_hold`, the customer added `legal_hold: false` to all artifacts by a bulk observed-tag injection — i.e., they supplied the state the guard asked for, from outside the artifact, with caller-controlled values. That is masking, done by the user, invited by the refusal message. The design had prevented the evaluator from folding absence into false and had instead taught the operator to do it.

Week 8, an integration partner discovered that sending one observed tag `priority: "n/a"` made every resolution on every artifact refuse (P-9), and that the refusal's `absent_keys` was `[]`. Triage spent two days looking for a missing tag.

Week 9, post-incident: the "zero importers" assumption in 0007 was re-read. Two importers existed by then. The assumption had no re-verification trigger. Same class as ledger (i) and (ii).

## 2. Obstacle negations

### C1 — "structural enforcement; evaluator never sees absence"
Negated by P-1, P-2, P-8. The structure guarantees the evaluator is not called for *kernel-absent* keys. It does not guarantee that kernel-presence means artifact-state-presence: observed tags are unconstrained and provenance-blind presence lets them stand in for owned state. And the seam must pass the key (for kind lookup), so an evaluator can still carry per-key knowledge. Enforcement is structural for one hop; masking has two other hops (view assembly, evaluator's own declaration lookup).
Required answer: either (a) presence for a key the row lists as owned-state must be satisfied by an *owned* tag, or the RDR says explicitly that observed tags may satisfy guards and that this is the intended trust model; and (b) the seam passes the key or it does not — say which and accept the consequence.

### C2 — "`exists` is cardinality, not grammar"
Negated by P-3, P-4, P-17. The kernel owns one token's spelling, one literal's parse, and one key-normalisation rule. Those are three grammar facts. They are small, and fail-closed when wrong, but they are a coupling to a moving peer with no re-verification trigger — the seed class.
Required answer: a single shared constant/package for the `exists` token and bool-literal form that both the kernel and the grammar spec import, plus a test in the kernel that fails if the grammar spec's token list does not contain the kernel's.

### C3 — "strong-Kleene over blocks is safe and complete"
Negated by P-5, P-6. Semantically complete; pragmatically, `unless` over commonly-absent tags vetoes everything, and authors will reach for the masking idiom. Also, the set of absent keys *named* must be defined independent of evaluation order.
Required answer: the refusal's `absent_keys` is the set over *all* atoms of *all* surviving rows that were U due to absence, computed without short-circuit; and the RDR acknowledges the `unless`-over-absent behaviour as intended and documents the `exists ∧` idiom, or introduces a grammar-level sugar.

### C4 — "seam loses nothing"
Negated by P-7, P-8. Relational, temporal, set-membership, multi-valued guards are now impossible without a side channel. This kernel RDR is deciding the grammar's expressive ceiling.
Required answer: state the ceiling explicitly as a decision ("guards are unary atoms over one tag and one literal; no cross-tag, no context") and record it as a constraint on 0009, or widen the seam to pass an opaque context handle (not the view) so the grammar can grow.

### C5 — "U for uncomparable present values is fine"
Negated by P-9, P-10. It is fine for masking; it opens denial, and it makes the refusal payload lie (present-but-uncomparable is not absent).
Required answer: refusal payload carries two lists — `absent_keys` and `uncomparable_keys` (or an `atoms: [{key, reason}]` list) — and the RDR decides whether caller-supplied observed tags may shadow a present owned tag of the same key at all.

### C6 — "naming absent keys is actionable, no leak"
Negated by P-11, P-12. Regression in row locatability; enumeration of the guard surface. Both are fixable by including `row_id` and by deciding whether tables are confidential.
Required answer: payload includes `{row, key, block}` per unevaluable atom.

### C7 — "prune-first remains sound"
Negated by P-13 only in the sense of expectation: it is sound per the rule but a pruned row's absent keys vanish. The RDR should state this as a consequence so the author's mental model is correct: "missing state is reported only if it blocks a row that would otherwise survive."

### C8 — "zero importers, no migration cost"
Negated by P-14. This is the seed class. Stamp it with a trigger: "re-verify at implement start; if importers > 0 the row-type change needs a migration note for each."

### R1 — row-level seam
P-15, P-16. The stated reasons are weak; the real reason (kernel must be the trust boundary for *absence*) is good. Rewrite the rejection to say that, and note that drift moves to the kernel/grammar boundary where it fails closed.

### R2 — kernel evaluates everything
P-17. Partially adopted de facto. State the line.

### R3 — opaque string + error channel
P-18. Load-time parse failure has no envelope kind. Name where it surfaces.

### R4 — closed world
Stands.

### R5 — declared read list
P-19. The approach has an analogous second key set (owned-state). State the relation.

## 3. Consumer artifacts (what would have caught each)

| Finding | Artifact |
|---|---|
| P-1, P-2 | `TestResolve_ObservedTagShadowsAbsentOwnedKey`: row owned-state lists `approver`; view has only an observed `approver: ""`. Expect `guard_unevaluable` (or owned-state refusal), not prune. Currently would select/prune. |
| P-3 | `TestExistsToken_MatchesGrammarSpec`: kernel's token constant must be in the grammar package's operator table (cross-package compile-time import, or a golden list checked in both repos). Also user journey: "author writes the presence operator exactly as the grammar doc spells it". |
| P-4 | Table test over literals `true/True/yes/1` for `exists`: define which are accepted; the rest must be a load-time parse error, never FALSE. |
| P-5 | User journey "unless over a rarely-set tag": table `unless: legal_hold == true`, corpus with no `legal_hold`. Expect documented behaviour (refuse) and documented idiom; test `TestUnlessOverAbsent_Vetoes` pins it so it is a decision, not a surprise. |
| P-6 | `TestAbsentKeys_OrderIndependent`: same row with atoms permuted; refusal payload equal. |
| P-7 | Negative spec test: a fixture guard `a == b` (two keys) must fail at parse with a message naming the limitation; plus a recorded constraint on 0009. |
| P-8 | Interface review: does `Compare` receive `key`? If yes, `TestEvaluator_CannotObserveAbsence` cannot be written structurally — record that and rely on P-1's test. |
| P-9 | `TestGarbageObservedTagDoesNotDenyAllRows` or, if denial is accepted, `TestUncomparableValue_RefusalNamesKeyAsUncomparable` — payload must not be empty. |
| P-10 | `TestUnknownOperator_IsLoadError_NotUnevaluable`: evaluator returning U for an unknown op must be distinguishable from absence (or operators validated at load). |
| P-11 | Threat-model line in RDR: "tables are / are not confidential"; if confidential, refusal names count not keys to unauthenticated callers. |
| P-12 | `TestRefusalPayload_IncludesRowID` and CLI envelope golden file. |
| P-13 | `TestPrunedRowAbsentKeysNotReported` — pins the rule so the docs can state it. |
| P-14 | Assumption with re-verification trigger; `rg -l 'internal/resolve'` at implement start recorded in evidence. |
| P-18 | CLI journey: malformed guard in table → structured envelope with a kind; decide whether taxonomy grows to six. |
| P-20 | `TestResolve_FromAtoms_*`: every adversarial case in the existing suite re-expressed with atoms + a stub evaluator, so the combinator is on the tested path. Coverage diff on `combineRow` must be non-zero from the adversarial suite. |

## 4. Refutation targets (seed class in this approach)

Seed root cause: a point-in-time fact about a moving peer or future caller, stamped with no re-verification trigger. Instances here:

1. "Kernel package has zero production importers" (Claim 8) — future callers from 0009 and the CLI. Needs a trigger.
2. "The `exists` operator token is `exists`" (Claim 2) — fact about a peer spec not yet written. Needs a shared constant or a cross-spec test.
3. "Evaluator returns true/false for present values" — assumes the grammar spec's kind system can always decide; the peer may add kinds (durations, sets) where comparison is partial. Needs the uncomparable channel (P-9).
4. "Tag values arrive as strings" — environment fact; if the view later carries typed values, `presence == literal` and the seam's string signature both move. Trigger: re-verify at implement.
5. "Closed five-kind refusal taxonomy" — a load-time parse failure (R3) and a present-but-uncomparable case (P-9) both want a kind. The taxonomy's closure is a point-in-time fact about the CLI contract.

None of these forces a switch; all are foldable as stated mitigations plus assumption triggers.
