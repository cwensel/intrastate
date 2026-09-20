Model: claude-opus-5

# Pre-lock critique — RDR 0012, declared-kind carrier at the guard value seam

Verdict: BLOCK. This record ships a silent-verdict-change engine wearing
the costume of a de-silencing bug fix, and it hands the loudest of its
new verdicts to a consumer — graphlint — that is wired to swallow
`GuardUnevaluable` as "yes, admissible." The record never notices,
because C4 books graphlint as a construction site and stops there.

---

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0012:C4` | Names `graphlint/reach.go::atomAdmitsValue` a construction site and asserts the site audit makes lint "the same decision the runtime evaluator makes" — but `atomAdmitsValue` collapses three verdicts to two via `verdict != GuardFalse`, so every new `GuardUnevaluable` C2 mints reads as ADMITTED. The RDR makes graphlint MORE permissive in the same change that makes the runtime stricter. | `lint` reports a model clean; `flow resolve` on the same model refuses `flow-guard-unevaluable`. The CI gate the user trusts stops catching the exact defect this RDR exists to surface. | §1, premortem, AT-1 |
| C-2 | `0012:C2` | Same `GuardUnevaluable` verdict flows into `graphlint/analysis.go::nodeMeetsAll` under a NEGATED test (`if !atomAdmitsValue(...) { return false }`), where over-approximation inverts: unevaluable now reads as "the terminal IS met," suppressing a real invariant-1 finding. Two consumers of one helper with opposite polarity; C4's one-line audit sees neither. | `lint` silently stops reporting a dangling/unreachable terminal it reported yesterday. No diagnostic, no exit-code change — a finding just disappears. | §1, premortem, AT-2 |
| C-3 | `0012:A6` | Pending, Method "Spike (corpus scan)", Evidence "none yet" — while `§Consequences` and `§MVV` already narrate the silent `GuardFalse → GuardTrue` flip as intended behavior, and `0012:S9` pins it as a normative scenario. Settled-fact prose depends on an unverified assumption; the Finalization Gate's own `§Assumption Verification` text forbids exactly this. | A persisted owned value spelled `07` against `iter eq 7` starts matching a row it never matched. The row fires. Nobody is told. | §1, §3, premortem, AT-3 |
| C-4 | `0012:A4` | Pending after C5 was narrowed at 3amigo, and the narrowing severed the assumption from its own evidence: A4's lint-agreement claim now rests on a check at `normalize.go::(*loader).atom` while the spike that measured the flip (`a4-lint-diff.md`) measured the pre-narrowing design. A Pending assumption gates the chosen design's central agreement claim. | Model whose `lint` verdict and `flow resolve` verdict disagree over an `int` guard — the divergence A4 was written to exclude. | §1, §2, AT-4 |
| C-5 | `0012:§failure-modes` | The "known, accepted residual" paragraph concedes that CLI `--tag iter=07` against `iter eq 7` moves `GuardFalse → GuardTrue` SILENTLY, then accepts it. An RDR whose Problem Statement is "silent verdict changes are worse than loud refusals" ships a new silent verdict change through its most-used ingress, and books it as scope discipline. | `intrastate flow resolve --tag iter=07` plans a different outcome than it did before upgrade. Exit 0 both times. No refusal, no note, no release signal. | §1, §3, premortem, AT-5 |
| C-6 | `0012:C2` | The `bool` arm restricts held values to the tokens `true`/`false` and makes EVERY other held value `GuardUnevaluable` — including `"1"`, `"0"`, `"True"`, `"yes"`, which today compare as raw strings and answer `GuardFalse` deterministically. No assumption, spike, or scenario bounds this population; A6 is scoped to "`int`/`bool` dimensions" in prose but its evidence and its worked example are int-only. | Every reader returning `"1"` for a bool dimension turns a working resolve into `flow-guard-unevaluable` on upgrade. Loud, but unbounded and unannounced. | §1, premortem, AT-6 |
| C-7 | `0012:C2` | The `in` poison rule — "ONE unparseable member poisons the whole list" — is stated for the held/literal parse, but `in` literals are load-conformed (`conformKind`) and C5 now additionally canonicalizes predicate literals, so the poison arm is unreachable from authored models and reachable ONLY from a hand-built atom. The rule's stated justification (lint/runtime divergence over a mixed list) is therefore about a state C5 makes unauthorable. A normative clause with no reachable input. | None directly — but the suite case pinning it (`0012:S1`'s mixed-list leg) passes vacuously, and the reviewer believes a risk is covered that the design already excluded by a different mechanism. | §1, AT-7 |
| C-8 | `0012:C3` | The migration is described as "re-pointed, not duplicated," but it simultaneously (a) widens `TestGuardEvaluatorContract`'s signature, (b) amends `0007:REQ-71`'s normative boundary pin, (c) makes the kernel's own `conformingContractSeam` kind-aware, (d) re-points a cross-package call site, and (e) migrates eighteen zero-value sites in `internal/guard` — ALL in phases that must land together or `main` is red. The record calls this "one package plus three call sites" in `§Consequences`. | Not user-facing. The implementer discovers mid-phase that Phase 1 and Phase 2 cannot be separated, `make check` is red between them, and the phase plan is abandoned. | §1, §2, premortem |
| C-9 | `0012:§consequences` | "kernel, atom shape, refusal taxonomy, and envelope are byte-untouched — the change is one package plus three call sites." Contradicted inside the same record by C3 (kernel test seam changes, `0007:REQ-71` amended), C5 (`internal/table/normalize.go` changes, new user-visible load refusal), and Phase 1 (eighteen `internal/guard` test sites). Five packages, not one. | Scope estimate the human approved is wrong by 5x; the change lands late or half-done. | §1, §2 |
| C-10 | `0012:S4` | Scenario 4's grep-style assertion ("no zero-value construction of `guard.Evaluator` outside `NewEvaluator` survives in non-test code") is unimplementable as stated against the `var` / `new` / embedded-field / struct-literal-field forms it claims to match. `guard.Evaluator{}` appears today as a FIELD VALUE inside `resolve.Input{... Guards: guard.Evaluator{}}` in nine test sites; a non-test analogue is a composite literal a text grep cannot distinguish from a constructor call without type resolution. The record names it "C4's mechanical backstop." | A new construction site added six months later is not caught. It answers `GuardUnevaluable` for every `eq`/`in` — or, if its guards use only `lt/gte/contains`, is not caught at all (the record concedes this in `0012:F6` and then still calls S4 the backstop). | §1, §2, AT-8 |
| C-11 | `0012:§technical-design` | The trailing non-normative paragraph between C4 and C5 asserts "No new refusal kind, error code, or envelope field for the SEAM," and C5 then adds a user-visible load refusal reusing `malformed_predicate_atom`. The disclaimer is technically scoped to "the SEAM," but the record's own `§Consequences` generalizes it to "refusal taxonomy and envelope are byte-untouched," which is what a reader carries away. | A downstream consumer parsing `malformed_predicate_atom` payloads receives a new message shape (`"00" is not the canonical spelling...`) it was told not to expect. | §1, §2 |
| C-12 | `0012:MVV` | The MVV requires an accessor reading `iter` as an OWNED value, returning `many`, then `07`, then `4` — three different reader behaviors — yet no step of the MVV, and no scenario `S1`-`S9`, says how the reader is made to return those values. Every other fixture (`S7`, `S8`) is a model file plus a CLI invocation; the MVV's load-bearing steps 2 and 3 need a programmable reader the record never specifies. | The MVV is not executed as written; the implementer substitutes a unit test at the seam, and the end-to-end claim ("the seed defect refuses loudly") is never actually demonstrated. | §1, §2, AT-9 |
| C-13 | `0012:A3` | Verified against three sites on `main`, and the verification is the whole guarantee for C4's "SAME loaded model" half — which `§Testing-Strategy`'s Known-coverage-gaps paragraph concedes has NO automated oracle. A Verified assumption about a snapshot of `main` is doing the work of an invariant. | A fourth site added later, or a refactor that moves `guardSeam` behind a cache, types comparisons against the wrong model's kinds. No panic, no refusal — wrong verdicts. | §1, premortem, AT-10 |
| C-14 | `0012:ALT3` | Rejected partly because "a new CLI error code is an unadjudicated envelope change" — but C5, the chosen design's Phase 4, adds a new user-visible load refusal message under an existing code, and `§Failure-Modes` concedes the CLI residual ALT3 would have closed. The rejection's cost argument is applied asymmetrically: ALT3 pays for a refusal path, O1 gets one free by reusing a code. | The CLI silent-flip (C-5) that ALT3 would have prevented ships, justified by a rejection that priced ALT3's refusal path as prohibitive. | §1, §3 |
| C-15 | `0012:G-proportionality` | The gate asks to confirm the RDR authors at most ONE independent load-bearing contract. This record authors C1 (carrier), C2 (comparison semantics), C3 (conformance suite shape + a `0007:REQ-71` amendment), C4 (producer obligation), and C5 (a load-time canonicalization rule in a DIFFERENT package, `internal/table`, with its own user-visible refusal and its own phase). C5 in particular is a separable contract with a separable blast radius. | The gate either rubber-stamps or blocks at lock, after four rounds of pre-lock review have already been spent. | §2 |
| C-16 | `0012:§problem-statement` | Frames the owned ingress as "the sole unconformed door," which is true of KIND conformance and false of the property C2 actually depends on: the model-literal and CLI doors conform kind but NOT canonical spelling (`conformKind` accepts `"00"`, `"+1"`, `"-0"`), which C5 then has to patch for one of the three doors and explicitly declines for the other. The Problem Statement's "two are already conformed" is load-bearing for the design's scope and is only half-true. | The user reads "the CLI path is safe" from the Problem Statement and is then surprised by the `--tag iter=07` flip buried in `§Failure-Modes`. | §1, §3 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 graphlint reads the new `GuardUnevaluable` as "admissible," and the lint gate goes quiet on the exact defect this RDR exists to surface

**Root cause in the RDR.** C4 enumerates
`internal/graphlint/reach.go::atomAdmitsValue` as one of three
construction sites and treats "construct over the right model" as the
whole of graphlint's obligation. It does not ask what graphlint DOES
with the verdict. It does, once, in `§Key Discoveries` — "An UNEVALUABLE
verdict is taken as satisfiable" is quoted nowhere, and the discovery
bullet that touches this path concerns match-atom byte-compare
divergence, not verdict collapse.

**The specific passage that enabled it.** `0012:C4`:

> On `main` that set is exactly three: `internal/cli/flow_resolve.go::guardSeam`,
> `internal/guard/product.go::valueSatisfies`,
> `internal/graphlint/reach.go::atomAdmitsValue` (A3).

and, in the same clause, the guarantee it claims to preserve:

> This is what keeps lint's stated property — "the same decision the
> runtime evaluator makes" (`valueSatisfies` doc) — true after the seam
> becomes declaration-aware.

That sentence attributes `valueSatisfies`' doc comment to the whole
three-site set. It is false for `atomAdmitsValue`. The source is
unambiguous:

```go
func atomAdmitsValue(a table.Atom, held string) bool {
	verdict := guard.Evaluator{}.Evaluate(...)
	return verdict != resolve.GuardFalse
}
```

`GuardUnevaluable` and `GuardTrue` are the same answer here. Before this
RDR, `eq`/`in` over an `int` key could only answer True or False — the
raw-string arm never returned Unevaluable. After C2, it returns
Unevaluable for every held value that does not parse and every
`set`-kind key and every key absent from the mapping. Each of those
newly becomes `true` — admitted, satisfiable, edge traversable — in
`ownedAtomSatisfiable`, which feeds `matchSatisfiable`, which feeds
`reach.go:210`, `groups.go:123`, and `export.go:121`.

**Symptom the user sees.** The user runs `intrastate lint --model m.toml`
and gets exit 0, clean. They then run `intrastate flow resolve` against
the same model with a real reader and get `flow-guard-unevaluable`. The
RDR's own `§Decision Rationale` names this orientation as the worst
possible one — "leaves lint more permissive than the kernel, the worst
orientation for the split" — and rejects an ALTERNATIVE for causing it,
while the chosen design causes it through a mechanism the record never
inspected.

The `§Key Discoveries` bullet that gets closest ("the seam decides MATCH
atoms too, on lint's side only") argues that C5 keeps the *match-atom
literal* divergence unobservable. That argument is about which BYTES the
two paths compare. It says nothing about what graphlint does with a
third verdict it has never seen before. C5 cannot help: C5 constrains
authored literals, and the new Unevaluables come from HELD values and
from the defensive `set`/absent-key arms.

### 1.2 The same verdict, in `analysis.go::nodeMeetsAll`, is read with the opposite polarity and silently deletes a finding

**Root cause in the RDR.** C2 specifies the verdict lattice and assumes
one consumer polarity. The record's only treatment of lint's reading of
`GuardUnevaluable` is `product.go::valueSatisfies`, which handles it
explicitly and blockingly. `atomAdmitsValue` has two callers and they
disagree.

**The specific passage that enabled it.** `0012:C2`:

> - kind `set`: the operator/kind matrix does not admit `eq`/`in` over
>   `set`; handed one anyway, the seam answers GuardUnevaluable
>   (defensive — the matrix rejection at load stays the primary guard).
> - key absent from the mapping: GuardUnevaluable — the seam cannot type
>   the comparison (defensive arm; unreachable when C4 holds, per A1 /
>   `0007:A7`).

The record calls these "defensive" and "unreachable," which is exactly
the framing that stops a reviewer from asking what happens when they
fire. But the mapping is `DeclaredKinds(m)`, which C1 specifies as
OMITTING any key whose `Kind` is the empty string. A key declared with
no explicit kind — a real authoring state `conformKind` tolerates,
since its switch has no default — is therefore absent from the mapping,
and the "unreachable" arm becomes the common case for that key.

Now follow it into `analysis.go:438`:

```go
for _, v := range held {
    if !atomAdmitsValue(atom, v) {
        return false
    }
}
```

`nodeMeetsAll` returns true when every held value is admitted.
`GuardUnevaluable` → `atomAdmitsValue` returns true → the negation is
false → the loop continues → `nodeMeetsAll` returns true → the terminal
is judged MET → the invariant-1 dangling finding does not fire.

`ownedAtomSatisfiable` wanted over-approximation ("pruning an edge lint
cannot decide is the false-green direction"). `nodeMeetsAll` wants
under-approximation and gets the same helper. One helper, two polarities,
one RDR clause that inspects neither.

**Symptom the user sees.** A finding that fired yesterday does not fire
today. No message says it stopped. The model still lints clean. The
dangling terminal reaches production.

### 1.3 The silent `GuardFalse → GuardTrue` leg ships unmeasured, and the record has already decided to ship it

**Root cause in the RDR.** A6 is Pending with Evidence "none yet," and
the rest of the record is written as if it were Verified.

**The specific passages that enabled it.** `0012:A6`:

> - **Status**: Pending
> - **Method**: Spike (corpus scan)
> - **Evidence**: none yet.

and, elsewhere in the same record, treated as settled: `0012:S9` is a
NORMATIVE scenario asserting the flip is the desired behavior; `0012:MVV`
step 3 makes it the MVV's "parsed-comparison witness"; and
`0012:§consequences` concludes:

> So the blast radius of the silent leg is not yet known, which is a
> rollout input, not a design defect.

That last clause is the failure. "Not yet known" and "not a defect" are
being asserted in the same breath about the same unmeasured population,
by a record whose Problem Statement is that silent verdict changes are
the harm. The Finalization Gate's own `§Assumption Verification` template
text forbids this construction verbatim: "no assumption marked `Pending`
or `Unverified` may have settled-fact prose elsewhere in the RDR
depending on it." S9, MVV step 3, and the Consequences bullet are all
settled-fact prose depending on A6.

And A6's scope is narrower than the exposure. A6 asks about
"non-canonical but parseable" int values. The `bool` arm's exposure
(C-6) is a different shape entirely: `"1"`, `"0"`, `"True"`, `"yes"`,
`"TRUE"` are all held values that today answer a clean `GuardFalse`
against `b eq true` and tomorrow answer `GuardUnevaluable`. That is not
in A6's worked example, not in S6's corpus scan, not in any spike.

**Symptom the user sees.** Two symptoms, opposite directions, same
upgrade. A persisted `iter = "07"` starts matching a row it never
matched — silently, a plan changes. A persisted `enabled = "1"` starts
refusing `flow-guard-unevaluable` — loudly, a working command breaks.
Neither is in a release note, because the record classified the first as
"intended" and never enumerated the second.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`§Implementation Plan` — specifically the Phase 1/2/3/4 decomposition
(`0012:§phase-1-seam-constructor-and-typed-arms` through
`0012:§phase-4-canonical-int-spelling-at-the-predicate-ingress`).**

It will be rewritten because it is already internally inconsistent, and
the inconsistency is visible in the record's own prose.

Phase 1 says: constructor, typed arms, zero-value retired, REQ-34 test
narrowed, plus eighteen `internal/guard` test-site migrations — and
gives the reason they cannot be deferred:

> It has to be this phase: once the typed arms land, every one of those
> evaluates `eq` against a nil mapping and each test asserting
> GuardTrue/GuardFalse on an `eq` goes red. Deferring them to Phase 3
> would leave `main` red between the two, which the project's
> run-`make check`-before-done convention does not allow.

That argument is correct and it generalizes, which the record does not
notice. `internal/guard/guard_evaluator_0003_test.go:67` drives
`TestGuardEvaluatorContract` over `var ev guard.Evaluator`. C3 says that
site is re-pointed "in the SAME phase as the signature widening" — which
is Phase 2. But the signature widening is what breaks `guard_mvv_test.go`
and the `0007:REQ-71` `var fn` pin, and those are in `internal/resolve`,
and the re-pointed site is in `internal/guard`, and the typed arms that
make the zero value answer Unevaluable are in Phase 1. Land Phase 1
alone: `guard_evaluator_0003_test.go:67` runs the OLD suite over a nil
mapping and its `eq`/`in` cases invert to Unevaluable — red. Land Phase 2
alone: the constructor does not exist. Land Phase 4 before Phase 1 and
`S7`'s regression has no typed comparison to guard.

Phase 4's own text concedes the ordering constraint ("Lands with S7's
regression and after Phase 1, so the literal side is closed before typed
comparison ships — the ordering A4 depends on") — but A4's ordering
requirement is that C5 lands BEFORE C2 ships, and Phase 1 IS C2. The
plan orders Phase 4 after Phase 1. The record states the dependency and
then violates it in the phase numbering.

Within six weeks the implementer will collapse Phases 1-3 into one
commit, because that is the only decomposition `make check` admits, and
Phase 4 will move ahead of it or be dropped. The section will be
rewritten to describe what was actually done. Adjacent casualties:
`§Consequences`' "one package plus three call sites" (C-9) and
`§Proportionality`'s one-contract claim (C-15) both have to be rewritten
with it.

---

## 3. The one assumption that will not survive first contact with a real user

**`0012:A6`** — and not because the scan will find something. Because
the user will find it first, and will not read it as a behavior fix.

A6's framing is the tell. It asks whether the installed base "contains no
value whose guard verdict this RDR changes SILENTLY." That is a question
about a corpus. The user does not have a corpus; the user has one
artifact, hand-edited over months, holding owned values written by a
reader that was never told `int` meant canonical decimal. The record
itself says so, twice:

> every persisted owned value and every hand edit to the artifact enters
> through it (`JDR 0004 §JD-2`) — `0012:§problem-statement`

> The held side is a different population and is unmeasured —
> `0012:A6`

A6 will not survive because its If-wrong already describes the outcome
and then declines to act on it:

> **If wrong**: the change needs a release note and possibly a scan or
> `--fix` affordance before rollout, rather than shipping as a behavior
> fix.

"Rather than shipping as a behavior fix" is the whole question, and it is
parked in an If-wrong of a Pending assumption. The record's `Type` is
`Bug Fix` and its `Priority` is `Medium`. A user upgrading a Medium bug
fix does not read release notes, does not run a scan, and does not expect
a plan to change. They expect `flow resolve` to do what it did last
week.

What first contact looks like concretely: a user with `iter = "07"`
persisted against a row guarded `iter eq 7`. Today the guarded row prunes
and the fallback plans; the user has built around the fallback for
months. After the upgrade the guarded row matches. Different outcome,
exit 0, no message. They will not file a bug against the guard seam —
they will file it against whatever downstream consumed the new outcome,
and it will take days to trace back.

And the CLI leg (C-5) makes it worse, because it is the ingress a user
actually types. `§Failure-Modes` documents the flip and accepts it:

> a CLI `--tag iter=07` against `iter eq 7` is admitted at input and
> then, under C2's parsed comparison, matches where today's raw-string
> seam prunes — a verdict moving GuardFalse → GuardTrue silently
> (measured, `evidence/spikes/a4-typed-compare.md`).

MEASURED. The record measured a silent verdict flip on the CLI ingress
and shipped it, in an RDR whose Problem Statement is that silent verdict
changes are worse than loud refusals. The justification — that closing it
at `conformKind` would retire a `0024` disposition — is a real constraint
and an irrelevant one: `flow_input.go::canonicalValue` is a CLI-specific
site that does not go through `conformKind`'s shared arm for the
canonicality check any more than `(*loader).atom` does. The record chose
to scope C5 to one of two available non-shared sites and called the
second out of scope.

---

## 4. Premortem — written from six months after the ship

RDR 0012 shipped in three commits over two weeks. `make check` was green.
The conformance suite was green, including all nine of the new
kind-discriminating cases. The MVV was executed — as a unit test at the
seam, not with a real accessor, because nobody ever specified how to make
a reader return `many`, then `07`, then `4` (C-12). Nobody noticed the
substitution because the seam-level assertions were identical.

**Week 3 — the lint gate goes quiet.** A platform team had been relying
on `intrastate lint --model` in CI to catch unreachable terminals before
merge. After the upgrade, a model with an `int`-declared `retries` tag and
a match atom `retries eq 3` over a node holding an opaque-abstracted
value stopped producing the `graph-coverage-gap` it had produced the week
before. Nobody investigated a finding that stopped firing. The path:
`atomAdmitsValue` returned `verdict != GuardFalse` for the newly-minted
`GuardUnevaluable`, `ownedAtomSatisfiable` took the existential witness,
`matchSatisfiable` returned true, `reach.go:210` marked the edge
traversable, and the terminal was reachable. Lint clean. Exit 0.

The RDR had asserted the opposite in C4: "This is what keeps lint's
stated property — 'the same decision the runtime evaluator makes' — true
after the seam becomes declaration-aware." That doc comment lives on
`valueSatisfies`, which handles `GuardUnevaluable` in a dedicated
blocking arm with a twelve-line comment explaining why. `atomAdmitsValue`
has no such arm. It was in C4's list of three sites and the audit checked
only that it constructed over the right model.

**Week 5 — a finding deletes itself in the other direction.** A second
team's model used `nodeMeetsAll` to certify terminal satisfaction. One
key was declared without an explicit `kind` — `conformKind`'s switch has
no default and had always tolerated it. `DeclaredKinds` OMITTED that key,
per C1's explicit "a key whose `Kind` is the empty string is OMITTED."
C2's key-absent arm — described in the RDR as "defensive; unreachable
when C4 holds" — fired on every evaluation. `atomAdmitsValue` returned
true. `nodeMeetsAll`'s `if !atomAdmitsValue(atom, v) { return false }`
never tripped. `nodeMeetsAll` returned true. The invariant-1 dangling
finding vanished. The team shipped a model with a dangling terminal,
certified by a lint pass that had certified nothing.

The RDR had reasoned about the key-absent arm exactly once, and only
about its reachability, never about its consumers: "defensive arm;
unreachable when C4 holds, per A1 / `0007:A7`." A1 is about guard atoms
referencing undeclared KEYS. It is not about declared keys with an empty
KIND. Different population; the assumption did not cover it and the arm
was not unreachable.

**Week 7 — the first silent misroute report.** A user upgraded and a
workflow started taking a different branch. Their persisted artifact held
`iter = "07"` against a row guarded `iter eq 7`. Before: raw-string
compare, `"07" != "7"`, `GuardFalse`, row prunes, fallback plans. After:
`strconv.Atoi` both sides, `7 == 7`, `GuardTrue`, guarded row plans.
Exit 0 both times. No refusal. No note.

This was in the RDR. `0012:S9` pins it as a normative expected scenario.
`0012:MVV` step 3 calls it "the parsed-comparison leg, on the one path
load cannot reach." `§Consequences` names it "a SILENT flip" and then
concludes "the blast radius of the silent leg is not yet known, which is
a rollout input, not a design defect." A6, the assumption that would have
measured it, was Pending with Evidence "none yet" at lock. It was never
run.

**Week 8 — the bool cliff.** A different user's reader returned `"1"` for
a `bool`-declared `enabled` tag. Before: raw-string compare against the
literal `true`, `GuardFalse`, deterministic. After: C2's `bool` arm
accepts only the tokens `true`/`false`, so `"1"` is
`GuardUnevaluable`, and `flow resolve` refuses. A working command started
failing on upgrade. Loud — which the RDR counts as success — but
unbounded and unannounced: A6 mentions "`int`/`bool` dimensions" in its
headline and every piece of its evidence, its worked example, and S6's
corpus scan is int-only. The bool population was never even estimated.

**Week 9 — the CLI flip, which had been measured.** `--tag iter=07`
against `iter eq 7`. Admitted at input by `canonicalValue`, then matched
by C2's parsed comparison where it used to prune. `§Failure-Modes`
documented it as an "accepted residual" and cited the spike that measured
it. Nobody who typed the flag had read the Failure Modes section of an
RDR.

**Week 11 — the phase plan, in retrospect.** The implementer had
collapsed Phases 1, 2, and 3 into one commit in week 1, because Phase 1's
own argument ("deferring them would leave `main` red between the two")
applied transitively to the suite signature widening in Phase 2 and the
re-pointed `guard_evaluator_0003_test.go:67` call site that C3 assigned
to "the SAME phase as the signature widening." Phase 4 landed second,
after C2 had already shipped — inverting the ordering A4 required and
Phase 4's own text acknowledged ("so the literal side is closed before
typed comparison ships — the ordering A4 depends on"). For two weeks the
literal side was open while typed comparison was live. `§Implementation
Plan` was rewritten in week 3 to match what happened.

**What the postmortem concluded.** The RDR was rigorous about its fences
and blind about its consumers. It verified A3 (three construction sites),
A5 (no import cycle), A1 (keys are declared), A2 (statelessness is a
proxy) — every assumption about whether the CHANGE could be made. It left
Pending both assumptions about what the change would DO to existing
behavior (A4, A6), and shipped anyway. The one consumer it enumerated and
did not read — `atomAdmitsValue` — cost two silent lint regressions in
opposite directions, from a single line the record never quoted:
`return verdict != resolve.GuardFalse`.

---

## 5. Acceptance tests, working backward

### AT-1 — graphlint does not admit a value the runtime refuses (catches week 3)

```gherkin
Scenario: an unevaluable held value is not treated as lint-admissible
  Given a model declaring tag "retries" with kind "int"
    And a match row whose atom is "retries eq 3"
    And a graph node holding the value "many" for "retries"
  When graphlint evaluates reachability for that row
  Then "atomAdmitsValue" MUST NOT report the value admissible
   And the lint verdict for the model MUST match the verdict
       "flow resolve" produces for the same held value
```

RDR-review form: for each of the three C4 construction sites, quote the
line that consumes the `GuardResult` and state which of the three
verdicts it maps to which behavior. C4 names three sites and quotes the
consumption of ONE (`valueSatisfies`). A review that required all three
would have surfaced `return verdict != resolve.GuardFalse` immediately.
Raises: C-1.

### AT-2 — a helper with two callers is checked at both polarities (catches week 5)

```gherkin
Scenario: an unevaluable verdict does not suppress a terminal finding
  Given a model with a tag declared with no explicit kind
    And a terminal predicate over that tag
  When graphlint runs "nodeMeetsAll" over a node holding a value for it
  Then the invariant-1 dangling finding MUST fire exactly as it does
       before the seam becomes declaration-aware
```

RDR-review form: enumerate every CALLER of every function C4 names, not
just the functions. `atomAdmitsValue` has two (`reach.go:484`,
`analysis.go:438`) and they negate it oppositely. Raises: C-2.

### AT-3 — no settled-fact prose rests on a Pending assumption (catches week 7)

```
Step 1. List every RDR element marked Pending or Unverified. (A4, A6.)
Step 2. Grep the record for every passage that asserts, as fact, a
        behavior that assumption bounds.
Step 3. For A6 the hits are: S9 (a NORMATIVE scenario asserting the
        silent match is expected), MVV step 3 (the "parsed-comparison
        witness"), and the §Consequences bullet ("a rollout input, not a
        design defect").
Step 4. BLOCK: three settled-fact passages depend on an assumption whose
        Evidence field reads "none yet". The Finalization Gate's
        §Assumption Verification forbids exactly this.
```

Raises: C-3, and would have caught C-4 by the same sweep.

### AT-4 — a narrowed contract re-grounds its assumption's evidence

```
Step 1. C5 was narrowed at 3amigo from conformKind's shared int arm to
        normalize.go::(*loader).atom.
Step 2. A4's Evidence cites spikes (a4-render-path.md, a4-lint-diff.md)
        measured against the PRE-narrowing design.
Step 3. Require: for every contract narrowed after an assumption was
        Verified against it, either re-run the spike at the new home or
        mark the assumption Pending with a named re-verification step.
Step 4. A4 is correctly marked Pending. BLOCK stands until it is run,
        because A4 gates the chosen design's central lint/runtime
        agreement claim, which §Decision Rationale's "fit" row uses to
        beat O4.
```

Raises: C-4.

### AT-5 — the record introduces no new silent verdict change (catches week 9)

```gherkin
Scenario Outline: every verdict change is loud or bounded
  Given an ingress <ingress> and a held or supplied value <value>
    And a guard atom "iter eq 7" over an int-declared tag
  When the verdict before the change is <before>
   And the verdict after the change is <after>
  Then if <before> != <after>, the change MUST be observable to the user
       (a refusal, a warning, or a documented release-note entry)

  Examples:
    | ingress     | value | before      | after            | observable? |
    | owned read  | many  | GuardFalse  | GuardUnevaluable | yes         |
    | owned read  | 07    | GuardFalse  | GuardTrue        | NO  <-- BLOCK |
    | CLI --tag   | 07    | GuardFalse  | GuardTrue        | NO  <-- BLOCK |
    | owned read  | 4     | GuardFalse  | GuardFalse       | n/a         |
```

The RDR contains both BLOCK rows already — one in S9/MVV step 3, one in
`§Failure-Modes` — and classifies both as acceptable. A review that
demanded the table in this form could not have signed it. Raises: C-5,
C-16, and by extension C-14 (ALT3 was rejected for a refusal cost the
chosen design pays anyway).

### AT-6 — every kind arm's behavior-change population is enumerated (catches week 8)

```
For each kind token in {enum, bool, int, set, scalar}:
  Step 1. State the set of held values that answer GuardFalse today.
  Step 2. State the set that answers GuardFalse after C2.
  Step 3. State the difference, and cite the measurement that bounds it.
  Step 4. BLOCK any arm whose difference is non-empty and uncited.

  int:    difference = non-parsing values (-> Unevaluable, loud) AND
          non-canonical-parseable (-> True, SILENT). Cited: A6 (Pending).
  bool:   difference = every held value outside {"true","false"}:
          "1", "0", "True", "TRUE", "yes". -> Unevaluable, loud.
          Cited: NOTHING. <-- BLOCK
  set:    difference = eq/in over a set key -> Unevaluable.
          Cited as "defensive"; no population estimate. <-- BLOCK
  enum/scalar: difference = empty. OK.
  absent: difference = every eq/in over an unmapped key -> Unevaluable.
          C1 OMITS empty-kind keys, so this is not the empty set that
          A1 implies. <-- BLOCK
```

Raises: C-6, and the absent-key row raises C-2's precondition.

### AT-7 — every normative arm has a reachable input

```
Step 1. For each bullet in C2, name an authored model or a reader value
        that reaches it.
Step 2. The in-poison bullet ("ONE unparseable member poisons the whole
        list") requires an in-literal with a non-integer member on an
        int key. conformKind refuses that at load; C5 additionally
        canonicalizes it. Reachable only from a hand-built atom.
Step 3. Either mark the bullet DEFENSIVE (as the set and absent-key
        bullets are marked) or delete it. As written its stated
        justification -- preventing a lint/runtime divergence over a
        mixed list -- describes a state the design excludes elsewhere,
        so the suite case pinning it (S1's mixed-list leg) passes
        vacuously.
```

Raises: C-7.

### AT-8 — the mechanical backstop is mechanically expressible

```
Step 1. Write the actual check S4 specifies, as a command.
Step 2. It must reject: `Evaluator{}`, `var e Evaluator`, `new(Evaluator)`,
        `struct{ E Evaluator }{}`, and `resolve.Input{Guards:
        guard.Evaluator{}}` -- while accepting `NewEvaluator(...)`.
Step 3. A text grep cannot do this; it needs go/types or an AST pass.
Step 4. Either name the tool (a custom vet pass, a go/analysis checker)
        and put it in the Implementation Plan, or stop calling S4 "C4's
        mechanical backstop" and say the backstop is C1's LOUD
        disposition alone -- which F6 already concedes does not bind
        lt/lte/gt/gte or contains sites at all.
```

Raises: C-10.

### AT-9 — the MVV is executable as written

```
Step 1. For each MVV step, name the artifact (a model file, a command
        line, a reader implementation) and the expected output bytes.
Step 2. Steps 2 and 3 require a reader returning "many", then "07", then
        "4". No step, fixture, or scenario names that reader or says how
        it is configured. S7 and S8 give a model plus a full CLI
        invocation; the MVV's load-bearing steps do not.
Step 3. BLOCK until step 1 is a file path and a command, so the MVV
        cannot be silently downgraded to a seam-level unit test that
        proves the arm and not the journey.
```

Raises: C-12.

### AT-10 — C4's "same model" half has an oracle or an explicit waiver

```
Step 1. §Testing Strategy already concedes: "a site constructing over
        the WRONG model's kinds still compiles and types comparisons
        silently (F4). There is no observable."
Step 2. The stated guard is "the A3 site audit and ... there being
        exactly three sites."
Step 3. A3 is a snapshot of `main`. JDR 0004 §JD-3 already names a
        FOURTH site arriving with 0030 Phase 2.
Step 4. Require either a runtime assertion (the evaluator carries a
        model identity and Resolve rejects a mismatch) or an explicit,
        signed waiver that the guarantee is discipline -- which is the
        exact property §Decision Rationale used to reject O4.
```

Raises: C-13.

### AT-11 — the scope claim survives a file count

```
Step 1. §Consequences says "one package plus three call sites."
Step 2. Count packages the RDR's own text requires changing:
        internal/guard (constructor, arms, 18 test sites),
        internal/resolve (suite signature, fixture, conformingContractSeam,
                          the 0007:REQ-71 var fn pin),
        internal/cli (guardSeam, probeRow threading, 4 test sites),
        internal/graphlint (construction + verdict consumption),
        internal/table (C5 at (*loader).atom).
        Five.
Step 3. Correct §Consequences, and re-run §Proportionality's
        one-contract test against C1/C2/C3/C4/C5 -- C5 lives in a
        different package, carries its own user-visible refusal, and has
        its own phase. It is a separable contract.
```

Raises: C-9, C-11, C-15, C-8.
