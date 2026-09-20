Model: claude-sonnet-5

# Barrier diff — RDR 0012 pre-lock critique, Pass A vs Pass B

---

## Convergence

Defects both passes raise on the same element id. Where the anchor
matches but the stated defect differs, that is flagged as TWO defects
sharing an anchor, not one.

- **`0012:C2`** — A: C-2, C-6, C-7 / B: C-2, C-9, C-15. Both passes
  independently derive that C2's typed comparison creates readings that
  diverge from raw-string behavior in ways the record does not bound.
  A's C-2 and B's C-15 converge exactly: the key-absent "defensive,
  unreachable" arm is not unreachable (A1 covers undeclared keys, not
  declared keys with empty `Kind`) and both trace it into
  `atomAdmitsValue`/`nodeMeetsAll` consumers. A's C-6 (bool-token
  population unbounded) and B's C-9 (GOARCH-dependent int overflow) are
  different populations under the same "C2 introduces an unbounded,
  uncited behavior-change population" defect shape — genuinely
  convergent framing, different instances.

- **`0012:C4`** — A: C-1 / B: C-6, C-7. Both passes independently derive,
  from the same source line (`return verdict != resolve.GuardFalse` at
  `internal/graphlint/reach.go:502-513`), that C4's "LOUD" guarantee and
  its "same decision the runtime evaluator makes" claim are false at the
  `atomAdmitsValue` site: an unevaluable/zero-value verdict reads as
  SATISFIABLE there, not as a refusal. This is the single highest-
  confidence finding in both passes — same root cause, same source
  citation, arrived at independently. B additionally raises a distinct
  second defect under the same anchor (C-7: no threading decision is
  documented for wiring the constructed evaluator into graphlint, unlike
  the CLI's `probeRow` threading, which C4 does address) — that is a
  second, non-overlapping defect on the same anchor, not a restatement.

- **`0012:A4`** — A: C-4 / B: C-1 (partial), C-9 sourced via C2/C3 but
  bearing on A4. Both mark A4 Pending and both find its re-verification
  scope insufficient — A frames this as "narrowed after evidence was
  gathered against the pre-narrowing design," B frames it as "the
  re-verification note targets literals, when literals were never the
  leak — held/`[initial]`/`[rule.write]` cells are." These are
  compatible, not identical: B's version is more specific and, per the
  source check below, correct about which population is actually
  uncovered.

- **`0012:A6`** — A: C-3 / B: C-4. Both flag A6 (Pending, "evidence: none
  yet") as gating settled-fact prose or as excluded from the lock gate.
  A's C-3 finds S9/MVV-step-3/§Consequences narrating the silent flip as
  accepted fact while A6 is Pending; B's C-4 finds the literal mechanism
  — Prerequisites lists "A1–A5," so A6 is excluded from the checklist by
  construction, not merely under-verified. Same underlying defect (A6 is
  Pending and load-bearing, and the record proceeds as if it were
  settled), two different textual proofs. Verified: `docs/rdr/0012-...md:1161`
  reads exactly `- [ ] All Critical Assumptions verified (A1–A5)`.

- **`0012:§consequences`** — A: C-9, C-15 / B: C-14. Both passes attack
  the same sentence, "the change is one package plus three call sites"
  (`docs/rdr/0012-...md:1050`), with overlapping but not identical file
  counts (A: five packages named; B: 18 test sites + 4 CLI fixtures +
  graphlint fixtures + 3 call sites + the REQ-71 pin + the REQ-34 test).
  Independent derivations of the same scope-understatement defect.

- **`0012:S4`** — A: C-10 / B: C-13. Both find Scenario 4's "no
  zero-value construction survives in non-test code" claim unimplementable
  as a mechanical check (grep cannot distinguish the composite-literal and
  embedded-field forms from a constructor call) and both note no
  `make check` target is named to run it. Same defect, same reasoning.

- **`0012:C5`** — A raises C5 only tangentially (via C-7's `in`-poison
  argument and C-11/§technical-design); B makes C5 the center of the
  critique (C-1, C-3, C-8, C-12) with source-verified specificity (the
  narrowing from the shared `conformKind` arm to `(*loader).atom` only).
  Not scored as full convergence — see Pass-B-only below — but flagged
  here because both passes independently distrust C5's completeness.

## Pass-A only

- **A C-8** (`0012:C3`) — the migration's "re-pointed, not duplicated"
  framing hides that Phase 1 and Phase 2 cannot land separately without
  redding `main` (widened `TestGuardEvaluatorContract`, the `0007:REQ-71`
  pin, and eighteen `internal/guard` zero-value test sites all move
  together). B's C-14 counts the same sites but does not make A's
  ordering-dependency argument.

- **A C-11** (`0012:§technical-design`) — the "no new refusal kind,
  error code, or envelope field for the SEAM" disclaimer is scoped to
  "the SEAM" but generalized by `§Consequences` to the whole record, and
  C5 (Phase 4) then adds a user-visible load refusal under an existing
  code, contradicting the generalized reading.

- **A C-12** (`0012:MVV`) — MVV steps 2 and 3 require a reader returning
  `many`, then `07`, then `4`, and no scenario or fixture specifies how
  that reader is configured, unlike S7/S8 which pair a model file with a
  full CLI invocation. Risk: the MVV is silently downgraded to a
  seam-level unit test.

- **A C-13** (`0012:A3`) — A3's "exactly three construction sites" is
  verified against a snapshot of `main`; `§Testing-Strategy`'s
  Known-coverage-gaps paragraph concedes there is no automated oracle
  for a fourth site arriving later (and JDR 0004 §JD-3 already names one
  incoming with RDR 0030 Phase 2).

- **A C-14** (`0012:ALT3`) — ALT3 was rejected partly because "a new CLI
  error code is an unadjudicated envelope change," but C5's Phase 4 adds
  a comparable envelope-adjacent change (a new refusal message under an
  existing code) and is not held to the same cost standard — an
  asymmetric rejection argument.

- **A C-16** (`0012:§problem-statement`) — "the owned ingress is the sole
  unconformed door" is true only for KIND conformance, not for canonical
  spelling; the model-literal and CLI doors both admit non-canonical
  spellings under `conformKind`'s permissive int arm, which the Problem
  Statement's framing obscures.

## Pass-B only

- **B C-1** (`0012:§key-discoveries`, `C5`, `A4`) — SOURCE-CONFIRMED, see
  Contradictions below. C5's narrowing to `(*loader).atom` leaves
  `heldValues`/`writesOf`-fed `[initial]` and `[rule.write]` cells
  uncanonicalized (`canonicalValues` at `internal/graphlint/reach.go:343`
  only sorts/dedupes, no `Atoi` round-trip), reopening the `matchlit`
  counterexample the record's own ruling used to reject an earlier
  design. A's coverage of the graphlint site (C-1/C-2) stays at the
  "defensive arm" and consumer-polarity level and does not reach this
  ingress-specific gap.

- **B C-2** (`0012:C2`, `D-canonical-spelling-at-load`) — the record
  rejects an alternative for "institutionalizing two comparison
  semantics for one declared kind" and then C2 does exactly that between
  `[match]` (byte-compared, `internal/resolve/resolve.go:175-182`,
  confirmed) and `[guard]` (parse-compared) atoms inside the kernel
  itself — an argue-against-what-you-built defect A does not name in
  this form.

- **B C-3** (`0012:F1`, `C5`) — names the specific contradiction between
  A6's If-wrong ("an RDR premised on removing silent verdict changes
  cannot introduce them") and `§Failure-Modes`' acceptance of the
  `--tag iter=07` flip as a residual, citing both passages against each
  other directly. A's C-5 raises the same failure mode but does not pair
  it against A6's own If-wrong text.

- **B C-5** (`0012:§problem-statement`, `MVV`) — the named beneficiary
  (the "rdr navigator" consumer used to motivate the record) executes
  the CLI, passes every fact as `--tag k=v`, declares zero `int` tags,
  and its `bool` tags are refused upstream by `canonicalValue` before
  reaching the seam — so the motivating consumer cannot experience
  either the defect or the fix. Not raised by A in any form.

- **B C-8** (`0012:C5`) — C5's match-literal canonicalization is itself a
  regression pathway: a model authored consistently non-canonical
  (`[match] n = "007"` with callers passing `--tag n=007`) works today
  under byte-compare; C5 refuses it at load, and once rewritten to the
  canonical form, the same callers' `--tag n=007` stops matching. Not
  raised by A.

- **B C-9** (`0012:C2`, `C3`) — `strconv.Atoi` is platform-width; the
  conformance suite's overflow case has a GOARCH-dependent expected
  verdict with no width fixed in the contract. Not raised by A.

- **B C-10** (`0012:A2`, `S3`) — the narrowed reflection assertion ("no
  view-typed or runtime-valued state") is not expressible by reflection
  as stated: a `map[string]string` of kinds and a `map[string]string` of
  held values are the same type, so the proxy collapses to "no field of
  type `resolve.TagSet`," which a cached view trivially evades. Not
  raised by A.

- **B C-11** (`0012:§approach`) — the record narrows a peer RDR's
  Verified assumption (`0007:A17`) in its own prose without that reading
  a shared declare-in-both check for anything but kind tokens, risking a
  7.1 cluster demotion. Not raised by A.

- **B C-12** (`0012:C5`) — C5's refusal reuses the `malformed_predicate_atom`
  category for a value `conformKind` itself calls well-formed, so the
  error a user sees misnames the defect class (spelling vs syntax). Not
  raised by A (related to but distinct from A's C-11, which is about the
  envelope-disclaimer scope, not the category's semantic accuracy).

## Contradictions

**1. Is graphlint's site coverage in C4 adequate, or does it miss the
actual leak population?**

- A's framing (C-1, C-2): the three C4 construction sites are the right
  set; the defect is that the record audits only construction ("built
  over the right model") and never inspects consumption — the two
  callers of `atomAdmitsValue` read the verdict with opposite polarity,
  and neither is checked.
- B's framing (C-1): the three sites are also mis-scoped at the INPUT
  side — C5 (the canonicalization contract meant to make lint/runtime
  agreement possible) was narrowed at 3amigo to predicate literals only,
  but graphlint's `reach.go` feeds held values from `heldValues`/
  `writesOf`, which are sourced from `[initial]`/`[rule.write]` cells
  conformed only by `conformKind`'s permissive int arm (admits `"007"`,
  `"+1"`, `"-0"`), never by C5's round-trip.
- **Settled by source.** Both are correct and compose rather than
  conflict, but B's claim is the one grounded in the actual data flow:
  `internal/graphlint/reach.go:436-440` (`heldValues`) calls
  `canonicalValues` (`reach.go:343`, sort/dedupe only, confirmed — no
  numeric round-trip), not any C5 gate. The RDR's own scope table at
  `docs/rdr/0012-...md:625` states this itself: "YES for predicate atoms
  — deliberately NOT total over the other ingresses (S5 measured only
  this one)," and line 760 confirms `conformKind`'s int arm (which
  admits `"00"`/`"+1"`/`"-0"`) is "SHARED by every authoring site ...
  `[initial]`, `[rule.write]`, `[emit]`" and stays unchanged. So B's more
  specific claim is the one the record's own text and the source both
  confirm: C5 does not close the `[initial]`/`[rule.write]` leg feeding
  reach, and the "no model can author a value the two read differently"
  sentence at `docs/rdr/0012-...md:1033-1035` is false under the
  narrowed C5. Cite: `internal/graphlint/reach.go::heldValues`,
  `internal/graphlint/reach.go::canonicalValues`,
  `internal/table/normalize.go::(*loader).atom`. A's polarity argument
  (C-2, the `nodeMeetsAll` negation) is independently true and stands
  alongside this, not against it — they are two distinct defects on the
  same function, both real.

**2. Does C1's "LOUD" disposition hold as a design-wide guarantee, or is
it a property of only some consumer sites?**

- A treats C4's coverage of `atomAdmitsValue` as adequate ON THE MERITS
  of A3's site count, with the defect being that the audit never reads
  what the site DOES with the verdict (C-1).
- B goes further and claims the "LOUD" property C1 asserts and F6/S4/S5
  lean on is flatly false at the `atomAdmitsValue` site specifically —
  it is true at the two refusal-producing sites (`flow_resolve.go`,
  `product.go::valueSatisfies`) and false at the third, and the record
  generalizes from two sites to three without checking (B's C-6).
- **These are not in tension** — B's claim is the source-confirmed
  generalization of A's claim (verified: `atomAdmitsValue` returns
  `verdict != resolve.GuardFalse`, confirmed at
  `internal/graphlint/reach.go:502-513`, called positively at
  `reach.go:484` and negated at `analysis.go:438`). Recorded here only
  because the two passes state it at different scope and a reader
  merging the ledger could mistake them for duplicates when B's is
  strictly the sharper claim of the two.

**3. Is A4's re-verification scope wrong because it targets the wrong
population (B), or wrong because its evidence predates a narrowing
event (A)?**

- A (C-4): A4 was Pending after C5 was narrowed at 3amigo, severing the
  assumption from its supporting spike (`a4-lint-diff.md`), which
  measured the pre-narrowing design.
- B (C-1's closing argument): A4's Pending note specifically asks to
  "re-verify that no `eq`/`in` int comparison reads a LITERAL that
  escapes the narrowed check" — but per the source finding above,
  literals were never the leak; held cells are. So the re-verification
  as scoped will PASS (no literal escapes) while the actual leak (held
  `[initial]`/`[rule.write]` cells) ships uncaught.
- **Settled by source, and this is the sharper reading.** Given the
  `heldValues`/`canonicalValues` trace above, B's claim is correct: the
  re-verification instruction in A4's Pending note (targeting literals)
  cannot detect the actual gap. A's claim (evidence predates the
  narrowing) is also true but less actionable — it would be satisfied by
  re-running the same spike, which would still not test the held-cell
  path. B's version should govern how A4 gets re-verified.

## Merged ledger

| ID | RDR passage | Failure mode | Symptom | Raised by | Confidence |
|----|-------------|---------------|---------|-----------|------------|
| D-1 | `0012:C4` | `atomAdmitsValue` reads `GuardUnevaluable` as satisfiable (`verdict != resolve.GuardFalse`); C4's "same decision as the runtime evaluator" / "LOUD" claims are false at this site. | `lint` exit 0 on a model the runtime refuses or that has a dead-end/unreachable node. | both (A-1, B-6) | high |
| D-2 | `0012:C2` (via `analysis.go::nodeMeetsAll`) | Same `atomAdmitsValue` verdict is consumed with negated polarity at a second call site; over-approximation inverts to under-approximation and suppresses a real finding. | A dangling/unreachable-terminal finding that fired yesterday silently stops firing. | A-2 | high |
| D-3 | `0012:C4` / graphlint threading | No documented decision for how the constructed evaluator reaches `atomAdmitsValue` (2 frames below the only function holding the model), unlike the CLI's `probeRow` threading which C4 does address. | Implementer either rebuilds the kinds map per atom evaluation (allocation storm) or widens three signatures with no contract backing the choice. | B-7 | med |
| D-4 | `0012:C5`, `§key-discoveries`, `A4` | C5 was narrowed to predicate literals only (`normalize.go::(*loader).atom`); reach's held values (`heldValues`/`writesOf`, fed by `[initial]`/`[rule.write]`) are conformed only by `conformKind`'s permissive int arm and are never canonicalized, reopening the `matchlit` counterexample the record's own ruling rejected. "No model can author a value the two read differently" (line 1033) is false. | `lint` clean on a model whose root is a dead end; `flow resolve` refuses `flow-no-match` at that root. | B-1 (source-confirmed) | high |
| D-5 | `0012:C2`, `D-canonical-spelling-at-load` | The record rejects an alternative for "institutionalizing two comparison semantics for one declared kind," then C2 does exactly that between byte-compared `[match]` atoms (`resolve.go::TagSet.matches`) and parse-compared `[guard]` atoms — confirmed in source. | Same view, same key, same value: a match row is excluded while a sibling guard row is selected, with no reason token explaining the divergence. | B-2 (source-confirmed) | high |
| D-6 | `0012:A6`, `F1`/`§failure-modes` | A6 (held-side census) is Pending, "evidence: none yet," excluded from the lock gate (Prerequisites lists only A1–A5, confirmed at line 1161), while S9/MVV-step-3/§Consequences narrate the silent `GuardFalse→GuardTrue` flip as settled/accepted — contradicting A6's own If-wrong, which says an RDR premised on removing silent verdict changes cannot introduce an unbounded set of them. | A persisted owned value like `iter="07"` silently starts matching `iter eq 7`; exit 0, no note, no release signal. | both (A-3, B-4) | high |
| D-7 | `0012:C2` | The `bool` arm accepts only `true`/`false`; every other held value (`"1"`, `"0"`, `"True"`, `"yes"`) becomes `GuardUnevaluable`, an unbounded and uncited population (A6's evidence and worked examples are int-only). | A working `flow resolve` over a bool-declared tag starts refusing `flow-guard-unevaluable` on upgrade, unannounced. | A-6 | med |
| D-8 | `0012:C2`, `C3` | `strconv.Atoi`/parsing is platform-width; the conformance suite's overflow case has a GOARCH-dependent expected verdict with no width fixed in the contract. | Conformance suite is green on amd64 and red on 32-bit builds; verdict differs by build target. | B-9 | med |
| D-9 | `0012:C2` (key-absent arm) | `DeclaredKinds` omits keys whose `Kind` is empty string; `conformKind`'s switch has no default, so an unkinded declared key is a real authoring state. The "defensive, unreachable under A1" framing is wrong — A1 covers undeclared keys, not declared-but-unkinded ones — so this arm is the common case for that key, not unreachable. | A guard over a previously-deciding key on an unkinded tag flips to `GuardUnevaluable`/is treated as admitted, depending on the consumer. | both (A-2 secondary, B-15) | high |
| D-10 | `0012:C2` | The `in` poison rule ("one unparseable member poisons the list") is stated for a state C5's literal canonicalization (plus load-conformance) makes unauthorable from any authored model; reachable only from a hand-built atom. The suite case pinning it passes vacuously. | None user-facing directly; a reviewer believes a risk is covered by a rule with no reachable input. | A-7 | low |
| D-11 | `0012:C3` | The migration's "re-pointed, not duplicated" framing hides that Phase 1 and Phase 2 cannot land separately without redding `main` (widened contract test, the `0007:REQ-71` pin, and 18 `internal/guard` zero-value test sites move together by the record's own Phase-1 argument, which generalizes to Phase 2 without the record noticing). | Not user-facing; implementer discovers mid-phase the plan can't be split and `make check` is red between phases. | A-8 | med |
| D-12 | `0012:§consequences` | "One package plus three call sites" understates scope by the record's own text: 5 packages named (A) / 18 test sites + 4 CLI fixtures + graphlint fixtures + 3 call sites + REQ-71 + REQ-34 (B). | Scope estimate approved by the human is wrong by roughly an order of magnitude; change lands late or half-done. | both (A-9/15, B-14) | high |
| D-13 | `0012:S4` | Scenario 4's "no zero-value construction survives in non-test code" is a grep-style claim that cannot mechanically distinguish `var`/`new`/embedded-field/struct-literal-field forms from a constructor call, and no `make check` target implements it. | A construction site added later is not caught by the named "mechanical backstop." | both (A-10, B-13) | high |
| D-14 | `0012:§technical-design` | "No new refusal kind, error code, or envelope field for the SEAM" is scoped to the seam but generalized by `§Consequences` to the whole record; C5 (Phase 4) adds a new user-visible load refusal message under an existing code, which the generalized reading told the reader not to expect. | A downstream consumer parsing `malformed_predicate_atom` payloads receives a new message shape it wasn't told to expect. | A-11 | med |
| D-15 | `0012:MVV` | MVV steps 2–3 need a reader returning `many`, then `07`, then `4`; no scenario or fixture specifies how that reader is configured, unlike S7/S8 which pair a model file with a full CLI invocation. | The MVV is silently downgraded to a seam-level unit test; the end-to-end claim is never actually demonstrated. | A-12 | med |
| D-16 | `0012:A3` | A3 ("exactly three construction sites") is a verified snapshot of `main` with no automated oracle for a fourth site (JDR 0004 §JD-3 already names one incoming via RDR 0030 Phase 2); `§Testing-Strategy` concedes a wrong-model construction "compiles and types comparisons silently." | A future construction site over the wrong model types comparisons with no observable, wrong verdicts with no refusal. | A-13 | med |
| D-17 | `0012:ALT3` | ALT3 was rejected partly on "a new CLI error code is an unadjudicated envelope change," but the chosen design's C5/Phase 4 pays a comparable envelope-adjacent cost (new refusal message, existing code) without the same scrutiny — an asymmetric rejection argument. | The CLI silent flip (D-6) that ALT3 would have prevented ships anyway, justified by a rejection that priced ALT3's cost as prohibitive while the chosen design pays a similar cost elsewhere. | A-14 | med |
| D-18 | `0012:§problem-statement` | "The owned ingress is the sole unconformed door" is true for KIND conformance only; the model-literal and CLI doors both admit non-canonical spelling under `conformKind`'s shared permissive int arm, which C5 patches for one door and explicitly declines for the CLI. | Reader concludes the CLI path is safe from the Problem Statement, then hits the `--tag iter=07` flip buried in `§Failure-Modes`. | A-16 | med |
| D-19 | `0012:C5` | C5's match-literal canonicalization is a regression pathway on its own: a model authored consistently non-canonical (`[match] n="007"`, callers using `--tag n=007`) works today under byte-compare; C5 refuses it at load, and after the diagnostic-suggested rewrite the same callers stop matching. | A model that loaded yesterday refuses at load today; after the fix, previously-working caller invocations silently stop matching. | B-8 | med |
| D-20 | `0012:C5` | C5's refusal reuses the `malformed_predicate_atom` category for a value `conformKind` itself treats as well-formed, misnaming the defect class (spelling vs syntax) to the user and to downstream envelope routing. | User sees "malformed" and checks TOML syntax instead of the actual issue (non-canonical spelling). | B-12 | med |
| D-21 | `0012:A2`, `S3` | The narrowed reflection assertion ("no view-typed or runtime-valued state") collapses to "no field of type `resolve.TagSet`," which a cached `map[string]string` view trivially evades; S3's stated fail condition cannot be implemented as written. | The structural test stays green for an evaluator that caches the view differently; the guarantee REQ-34 was meant to enforce is gone with nothing replacing it. | B-10 | med |
| D-22 | `0012:§approach` | The record narrows a peer RDR's Verified assumption (`0007:A17`) in its own prose; the declare-in-both check fired only for kind tokens, not for this narrowing, risking two records stating A17 differently. | 7.1 cluster gate demotes 0012 or 0007 for declare-in-both; implementation stalls on a re-lock. | B-11 | low |
| D-23 | `0012:§problem-statement`, `MVV` | The named motivating beneficiary (the "rdr navigator" consumer) executes the CLI exclusively, passes every fact as `--tag k=v`, declares zero `int` tags, and its bool tags are refused upstream by `canonicalValue` before reaching the seam — the consumer cannot experience either the defect or the fix. | The motivating user journey in the Problem Statement and the MVV's acceptance paragraph describe a program that does not exist; a reviewer cannot verify any real consumer reaches the owned door with a typed tag. | B-5 | med |

## Health

**Pass A (claude-opus-5): PASS.** Every ledger row anchors to a named
element id and, in the body sections, to a specific function and line
(`atomAdmitsValue`'s `return verdict != resolve.GuardFalse`,
`analysis.go:438`'s negated loop, `guard_evaluator_0003_test.go:67`).
The acceptance tests are written as runnable Gherkin/step sequences
tied to specific source symbols, not generic advice. Anchors span
findings, alternatives, scenarios, consequences, and the proportionality
gate — several distinct origins, not clustered on §1. The premortem
traces concrete source-level causal chains (`ownedAtomSatisfiable` →
`matchSatisfiable` → `reach.go:210`) rather than narrating generically.

**Pass B (claude-fable-5): PASS.** Ledger rows are similarly anchored to
named functions and exact call sites (`internal/resolve/resolve.go::
TagSet.matches`, `internal/graphlint/reach.go::heldValues`,
`internal/table/load.go::conformKind`'s `strconv.Atoi`), and several
claims verified independently against source in this diff pass turned
out to be correct and, in two cases (C-1's `[initial]`/`[rule.write]`
gap, A4's mis-scoped re-verification), more precisely grounded than
Pass A's treatment of the same anchors. The named-beneficiary defect
(C-5) is a concrete user-journey check, not generic advice. Anchors span
findings, key discoveries, prerequisites, approach, and scenarios.

Both passes meet the lens's Expected signal: concrete named passages,
concrete functions, concrete (or falsified-concrete, in B-5's case)
user journeys, and ledger rows with real anchors spanning several
origins. Neither shows generic advice, unanchored rows, or all-§1
clustering.
