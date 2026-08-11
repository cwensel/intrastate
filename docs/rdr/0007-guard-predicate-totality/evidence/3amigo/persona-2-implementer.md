Model: claude-opus-5[1m]

# 3amigo Persona 2 — Implementer

RDR: `docs/rdr/0007-guard-predicate-totality.md`
Lens question: if I started coding this Monday, what would I ask in the first hour?
Grounded against `internal/resolve/resolve.go` (@ `9ad06a3`), `internal/resolve/fixtures_test.go`,
`docs/rdr/0003-guard-predicate-exhaustiveness.md`, and
`0007-.../evidence/spikes/aggregation-probe.md`.

---

## IMP-1 — The domain rule is written over structured atoms; the seam it constrains carries an opaque `string`. Nothing says who parses.

**Severity: Blocker**

**Triggering passage** — *Technical Design*:

> "The contract constrains one seam — `GuardEvaluator.Evaluate`
> (`internal/resolve/resolve.go::GuardEvaluator`) — and is implemented by RDR 0003's
> evaluator when that RDR is built."

and *Normative Contracts*:

> "Value-comparing guard operators (equality, membership, bounded integer comparison,
> set containment) are PARTIAL over the assembled evaluation view: an atom whose
> referenced tag key is absent from the view MUST evaluate to unevaluable"

**The gap.** The shipped seam is:

```go
type GuardEvaluator interface {
    Evaluate(guard string, view TagSet) GuardResult
}
```

and `Row.Guard` is a bare `string` ("Guard is the predicate handed to the guard seam.
Empty means no guard."). Every normative clause in this RDR quantifies over things the
`string` does not carry: *atoms*, *operators* (`eq`/`in`/`lt..gte`/`exists`/`contains`),
a *referenced tag key* per atom, and the `all` / `unless` block split. RDR 0003's
structured shape is TOML-side (`[rule.guard.all.profile] in = [...]`,
`[rule.guard.unless.prelock_iterations] gte = 3` — its Illustrative Code and
`0003-.../evidence/spikes/guard-fixture.toml`). The kernel flattens that to a string
with no stated encoding, and the RDR never says whether the flattening is lossless.

Concretely, on Monday I cannot answer:

- Does RDR 0003's evaluator receive the *authored atom structure* through some other
  channel (a lookup keyed on `Row.Guard` as an opaque **name**, the way
  `fixtureGuards.decided map[string]bool` treats it), or is `Row.Guard` a *serialized
  predicate* the evaluator must parse at every `Evaluate` call?
- If it is a name/handle, then `Evaluate`'s `guard string` parameter is a key into a
  compiled-guard table the evaluator owns — and this RDR's `all`/`unless` normative
  clauses land in that compiler, not at the seam it says it constrains.
- If it is a serialized predicate, this RDR is silently requiring a guard-text grammar
  that neither RDR 0002 (normalization) nor RDR 0003 (grammar) defines, and every
  conformance vector in Phase 2 has to invent one.

The RDR's own diagnostics language leans on the ambiguity in both directions: *Failure
Modes* says "the refusal carries the guard text (`Refusal.Guard`) and the contributing
rows; the named tags point at the absent state" — which reads `Row.Guard` as
human-legible authored text — while the *Diagnostic gap* bullet says "'which tag was
absent' is inferable only from the guard text", which is only true if the text names
tags. Shipped fixture guards are `"reviews >= quorum"`, `"is-fast-lane"`,
`"unknown-predicate"` — a mix of both readings in one test tree.

**The question I'd ask.** Is `Row.Guard` (a) an opaque identifier resolved by the
evaluator against compiled RDR 0003 atoms, or (b) a serialized predicate the evaluator
parses? If (a), name the RDR that owns the string↔atom binding and change this RDR's
Technical Design to say the domain rule constrains the compiled-atom evaluator behind
the seam, not `Evaluate`'s parameter. If (b), name the RDR that owns the guard-text
grammar. Either way: does `Row.Guard` need a shape change (e.g. structured atoms on
`Row`), and if so, Phase 1 is no longer "doc comments only"?

**Decision blocked:** whether Phase 2's view-reading evaluator is a
`map[guardName]compiledAtoms` lookup or a text parser — i.e. the entire type signature
and construction of the artifact Phase 2 says must be built ("A view-reading evaluator
… This is the domain rule's first executable expression"), and whether Phase 1 stays
non-behavioral.

---

## IMP-2 — "present in the assembled view" is provenance-blind here, but the kernel's only presence predicate is provenance-gated. `exists` has two defensible implementations.

**Severity: Blocker**

**Triggering passage** — *Normative Contracts*:

> "The existence operator is the sole TOTAL operator: it MUST decide true or false from
> the presence or absence of its referenced tag key alone."

and A6a:

> "no tag key present in the input tuple is ever dropped from the assembled view …
> Owned-over-observed precedence overwrites a `taggedValue` in place, so the key
> survives with `ProvenanceOwned` — value-level shadowing only … Since absence is what
> `exists` decides on, shadowing cannot flip a negative existence atom."

**The gap.** `TagSet` exposes exactly two presence-answering readers, and they disagree:

```go
func (s TagSet) Lookup(key string) (value string, prov Provenance, ok bool)  // exported, provenance-blind
func (s TagSet) has(key string, prov Provenance) bool                         // UNEXPORTED, provenance-gated
```

`missingOwned` uses the provenance-gated one: `view.has(key, ProvenanceOwned)`. So the
kernel already has a precedent where "present" means "present **with owned
provenance**". This RDR's rule says "presence or absence of its referenced tag key
alone", which is `Lookup`'s `ok`. Those differ for a real, reachable case: a tag key
supplied only as `Observed` (or the `recognized` key injected by `assemble`) is
`Lookup`-present but `has(key, ProvenanceOwned)`-absent.

RDR 0003 declares provenance per tag (`[tags.status] provenance = "owned"`,
`[tags.cluster_eligible] provenance = "observed"`) and its spike fixture puts
`exists = true` on `cluster_eligible`, an **observed** tag. So the observed-provenance
`exists` case is not hypothetical — it is the only `exists` example anyone has written.

Left open, an evaluator author has to guess:

- Does `exists` on a tag *declared* owned but *supplied* observed decide TRUE (key
  present) or FALSE (owned snapshot lacks it)?
- Does a value-comparing atom over an observed-only tag evaluate the observed value, or
  is it unevaluable because the declared provenance says owned?
- Is the `recognized` key (`resolve.go` `const recognizedTagKey = "recognized"`) in the
  domain of `exists`? `assemble` inserts it only `if in.Recognized != ""`, so an empty
  recognized outcome makes `exists recognized` decide FALSE — a spec-silent behavior an
  author could reasonably read either way.

A6a's disclaimer sharpens rather than closes this: "shadowing changes a value, never a
key's presence" is only reassuring under the provenance-blind reading. Under the gated
reading, an owned-declared tag arriving observed *is* a presence flip.

**The question I'd ask.** Does the domain rule's "present in the view" mean
`Lookup(...) ok == true` regardless of provenance, and does `exists` decide on that same
predicate? If yes, does `TagSet` need an exported provenance-blind `Has(key) bool`
(Phase 1 becomes an API addition, not doc comments), and should the RDR state explicitly
that guard evaluation is provenance-blind while `RequiresOwned`/`missingOwned` stay
provenance-gated? If no, state the gating rule and how the evaluator learns each tag's
declared provenance from `TagSet` — which today it cannot, because it only gets
`Lookup`'s *actual* provenance, not the declared one.

**Decision blocked:** the body of every `exists` and value-comparing atom in Phase 2's
view-reading evaluator, plus Testing Strategy scenario 6 ("`exists` is total"), whose
expected values invert depending on which predicate is meant.

---

## IMP-3 — The aggregation contract's escape half and the shipped `Refusal.Guard` single-value payload under-specify which guard text a multi-row refusal reports; scenario 5's expected values are unreproducible as written.

**Severity: Major**

**Triggering passage** — *Normative Contracts* (aggregation block):

> "An unevaluable ESCAPE row yields guard_unevaluable in place of the candidate-set
> refusal: the kernel MUST NOT claim the escape failed to rescue when it could not
> decide the escape at all."

and *Testing Strategy* scenario 5:

> "(a) `guard_unevaluable` carrying the escape row's guard … (b) `guard_unevaluable`
> carrying the *candidate's* guard"

**The gap.** `Refusal.Guard` is a single `string` ("Guard names the predicate the seam
could not decide, on guard_unevaluable"), but `gate` can accumulate many undecidable
rows and picks one by a **row-identity** tiebreak, not by guard text:

```go
lowest := slices.MinFunc(undecidable, func(a, b Row) int {
    return compareRefs(refOf(a), refOf(b))
})
return nil, &Refusal{Kind: KindGuardUnevaluable, Guard: lowest.Guard, Rows: rowRefs(undecidable)}
```

`compareRefs` sorts on `(RuleID, SourceLocator)`. So "carrying the escape row's guard"
is only well-defined when exactly one escape row is undecidable. With two undecidable
escape rows the reported guard is whichever has the lexicographically smallest `RuleID`
— a fact the RDR never states, yet Phase 2 vectors have to assert on it, and the RDR
elsewhere calls this precedence load-bearing.

Worse for reproducing scenario 5: the spike output the RDR promotes to normative
(*Testing Strategy*: "its captured output is the normative expected value") records both
probes with the **same** guard string:

> `PROBE A: kind="guard_unevaluable" guard="unknown-predicate"`
> `PROBE B: kind="guard_unevaluable" guard="unknown-predicate"`

Probe A is the escape-row case, Probe B the candidate case. Because both fixtures used
the same guard text, the captured evidence does **not** discriminate "carrying the
escape row's guard" from "carrying the candidate's guard" — the exact distinction
scenario 5 asks a vector to assert. As an implementer I cannot derive the expected
`Refusal.Guard` for 5(a) vs 5(b) from the cited evidence.

Separately, the *Existing Infrastructure Audit* flags "Owned-before-unevaluable
precedence is unfrozen (no contending test) … Phase 2 pins the precedence" and the
Testing Strategy adds an "Owned-state ordering gap … held by implementation only" — but
neither says which way to pin it when a **single survivor row** is simultaneously
missing an owned key and carrying an unevaluable guard. `gate` reports
`owned_state_unavailable` (the `missingOwned` block returns before the undecidable
loop); the RDR asserts that is the desired precedence in `gate`'s doc comment's voice
but never states it in a ```normative``` block, so the Phase 2 vector has no normative
anchor to cite.

**The question I'd ask.** For a `guard_unevaluable` refusal over N > 1 undecidable rows,
is `Refusal.Guard` normatively "the guard of the row with the lowest `(RuleID,
SourceLocator)`"? And should the owned-before-unevaluable precedence be promoted into a
```normative``` block so Phase 2 can pin it as spec rather than as observed behavior?

**Decision blocked:** the assertion values in Phase 2 vectors for scenario 5 and for the
flagged owned-state ordering gap — i.e. what I literally type as `want.Guard` on Monday.

---

## IMP-4 — Phase 2's "exported conformance harness" has no stated package, signature, or vector encoding, and the two named obligations conflict on where the code lives.

**Severity: Major**

**Triggering passage** — *Phase 2*:

> "**An exported conformance harness**, not an internal test. The suite MUST be callable
> against any `GuardEvaluator` (an exported function taking the seam), so RDR 0003's
> implement stage can instantiate it against its own evaluator."

and *Risks and Mitigations*:

> "the vector suite MUST be exported as a reusable conformance harness (an exported
> function taking a `GuardEvaluator`)"

and *Existing Infrastructure Audit*:

> "Golden conformance vector harness | None found under `internal/resolve/` (no
> `testdata/` repo-wide) | — | Build"

**The gap.** Four unstated things, each of which changes the first file I create:

1. **Package.** `GuardEvaluator` lives in `internal/resolve` — an *internal* package. RDR
   0003's evaluator is described as "a separate future package". If it is under
   `internal/`, importing works; if it is a sibling module or a public surface, an
   `internal/resolve/...` harness is unimportable by construction, and the harness has to
   move (`pkg/guardconform`? a new `internal/guardconform`?). The RDR asserts the
   cross-package binding is the whole mitigation for P-10 but never checks that Go's
   `internal/` visibility permits it.
2. **Signature.** "an exported function taking the seam" is not a signature. Does it
   take `*testing.T` (making the harness a test-only dependency and forcing `testing`
   into a non-`_test.go` file), or return `[]VectorResult` / an `error` the caller
   asserts on? These give different files.
3. **Test-package placement.** The named reuse target `fixtureGuards` is in
   `fixtures_test.go`, package `resolve_test` — nothing in `_test.go` is importable by
   another package. So the harness cannot be built by extending `fixtureGuards`; it must
   be a new non-test file. The Audit row says "Reuse (rows 1–3 only)" for
   `fixtureGuards`, which is true for the MVV but is misleading for the harness: the MVV
   rows and the harness rows cannot share a stub across the test/non-test boundary.
4. **Vector encoding.** "named, versioned golden vector suite" — Go table literals, or
   `testdata/*.json`/`*.toml` golden files? The Audit explicitly notes "no `testdata/`
   repo-wide", so this is a new convention for the repo, and "versioned" is unexplained
   (a `Version` field? a directory per version? what happens when a vector changes?).
   And to encode a vector at all I need IMP-1 answered — a vector's input is a guard,
   and I don't know if that's a name or a serialized predicate.

**The question I'd ask.** Give me the harness's import path and exact exported signature
(e.g. `func RunConformance(t *testing.T, eval resolve.GuardEvaluator)` in
`internal/guardconform`, vs. `func Vectors() []Vector` + a caller-side runner), confirm
RDR 0003's evaluator will live somewhere that can import it under Go's `internal/` rule,
and state the vector encoding and what "versioned" obliges.

**Decision blocked:** the first file I create in Phase 2 — its path, package clause, and
whether `testing` is a production-tree import.

---

## IMP-5 — Phase 1 is scoped "doc comments only", but two normative clauses it introduces have no doc-comment home and one depends on an unverified assumption.

**Severity: Moderate**

**Triggering passage** — *Phase 1*:

> "Narrow `Row.RequiresOwned`'s doc contract to post-guard write-dependency keys and
> state the domain rule on `GuardEvaluator` — doc comments in
> `internal/resolve/resolve.go` only; no behavior change."

versus *Normative Contracts*:

> "Empty atom blocks take the conventional identities: an empty `all` block is TRUE
> (empty conjunction), and an empty `unless` block MUST be treated as absent — NOT as a
> vacuously-true conjunction"

and *Prerequisites*:

> "[ ] **A9 verified** (empty-`unless` identity against RDR 0002's normalization) …
> Must close before lock"

**The gap.** Three concrete Monday problems:

1. **The empty-block identity has no landing site in this repo.** `Row` has no `all` or
   `unless` field — it has `Guard string`. So the empty-`unless` identity cannot be
   stated as a doc comment on any existing kernel symbol; it belongs on RDR 0002's
   normalizer or RDR 0003's compiler, neither of which is built. Phase 1 as written
   cannot discharge this clause, and Phase 2's scenario 9(b) ("a row whose `unless`
   block is omitted … the row is *not* disabled") has nothing to test against at the
   kernel seam, because a row with no `unless` block and no `all` block flattens to
   `Guard == ""`, which `evaluateGuard` already answers `GuardTrue` before the seam is
   reached. Scenario 9(b) may be untestable at this seam for the same reason A6b's
   scenario is declared out of scope.
2. **A9 is `Pending`** ("Evidence: Needed") yet the ```normative``` block already states
   the identity as settled. If A9 resolves the other way — the identity belongs to RDR
   0002's normalization — that normative block has to be deleted, and any Phase 2 vector
   written against it discarded. I would not write scenario 9(b) before A9 closes.
3. **The `RequiresOwned` narrowing has a doc-comment/behavior tension the RDR
   acknowledges but does not resolve into words.** Today:

   ```go
   // RequiresOwned names owned tag keys the row's evaluation needs.
   // A key absent from the owned snapshot yields owned_state_unavailable.
   ```

   The RDR narrows the *first* sentence to "the owned tag keys the row's post-guard
   transition (writes and clears) depends on", but the *second* sentence describes
   `missingOwned`'s unchanged behavior, which fires for **any** listed key including
   guard-read ones — and the RDR explicitly preserves that ("Listing a guard-read owned
   key in RequiresOwned remains legal and yields the more precise
   owned_state_unavailable diagnosis among surviving rows"). So the narrowed doc comment
   must simultaneously say "this field means write-dependencies" and "listing a
   guard-read key here is still legal and changes the diagnosis". The RDR gives me the
   normative sentence but not the reconciled doc comment, and the field mentions
   "clears" — a concept `Row` does not model (it has `Writes []Tag` and `NextTags []Tag`,
   no clear list). Do I add "clears" to a doc comment for a capability that does not
   exist in the struct?

**The question I'd ask.** (a) Where does the empty-`unless` identity land given `Row`
has no `unless` field — does it move to RDR 0002/0003, and does scenario 9(b) come out
of this RDR's test matrix? (b) Should Phase 1 wait on A9? (c) Give me the exact
replacement doc comment for `RequiresOwned`, including whether "clears" is aspirational
or names a `Row` field that must be added (which would break "no behavior change").

**Decision blocked:** whether Phase 1 is a same-day doc-comment patch or a Phase-1
blocked on A9 plus a possible `Row`/`TagSet` API change — i.e. whether I can land
anything on Monday at all.
