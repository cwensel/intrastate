Model: claude-opus-5[1m]

# Cove Findings — RDR 0002 (Transition Table As Reviewable Data), iteration 2

Chain-of-Verification lens, re-entry pass. Iteration 1 (loose files under
`evidence/cove/`) ran against the pre-§D7 draft. Since then the RDR was demoted
to Draft by cluster-reconcile iter-3 (SPEC-DEFECT, Stage 3, **STAGE-SCOPED**,
re-verify **A1, A9, A12**) and re-authored to the JDR 0001 §D7 closed layout
(refine `3c3cb3e`, resolve `929e931`), with a new spike at
`evidence/spikes/iter-2/`.

**Delta scope**: the §D7 re-authoring and the re-verify set A1/A9/A12. Claims
CONFIRMED unchanged at iteration 1 and untouched since are not re-litigated.

Findings are the origin ledger for this lens's resolve pass.

## Step 0 — Grounding sweep

Every codebase claim in the draft was read against source on `main` at
`9cd50b3`. **All kernel citations CONFIRM.** No claim was REFUTED. One claim is
NOT-FOUND but disclosed by the RDR itself, so it is not a finding.

| Claim | Site | Verdict |
| --- | --- | --- |
| `escapeOrRefuse` filters escape candidates on `row.Outcome != in.Recognized \|\| !view.matches(row.Match)` before gating, exactly as `::Resolve` does | escape clause; A11 | **CONFIRMED** — `internal/resolve/resolve.go::escapeOrRefuse`, the filter is verbatim as cited |
| `TagSet.matches` / `::Lookup` / `::has` index by raw map lookup with no folding | tag-key identity clause | **CONFIRMED** — all three read `s.tags[key]` directly |
| the repo's only `strings.ToLower` is on the `--as` flag value in `respond.go::ModeOf` | tag-key identity clause | **CONFIRMED** — `internal/cli/respond/respond.go:83,97`; no other `ToLower`/`ToUpper`/`TrimSpace` in non-test `internal/`+`cmd/` |
| `adversarial_test.go::TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement` freezes a populated `RequiresOwned` on an escape row yielding `owned_state_unavailable` | `RequiresOwned` clause | **CONFIRMED** — `internal/resolve/adversarial_test.go:226` |
| `resolve.go::Row` carries `Guard string` and `Match []Tag`, with no `Atom` type, no operator-token type, no per-atom block field | A9, A12, Prerequisites | **CONFIRMED** — `Row` is `{RuleID, SourceLocator, Outcome, Match []Tag, RequiresOwned []string, Guard string, NextTags, Writes []Tag, Escape []RefusalKind}` |
| `OpExists` / `LiteralTrue` / `LiteralFalse` do not exist in `internal/resolve` yet | existence-atom clause; Capability Dependencies | **NOT-FOUND — disclosed.** `grep` over `internal/` returns nothing. The RDR states this in three places and sequences Phases 2–3 behind RDR 0007's reshape. Not a finding. |
| `config.go::Load` has no parse-validation code; `config-invalid` is a doc comment and a `TODO` | A5 (corrected at iter-1) | **CONFIRMED** — `internal/cli/config/config.go:98` is `// TODO: parse `data` … once a TOML library is chosen`; no decoder is wired |
| `clierr.ErrorGroup` exists and exit mapping is RDR 0005's, not this RDR's | category clause; A5 | **CONFIRMED** — `internal/cli/clierr` declares `ErrorGroup` + five groups; the RDR correctly declines to assign them |

**Inverse sweep** (does a sibling path already make a decision this RDR
invents?). Searched `internal/`, `cmd/` for an existing transition table,
normalizer, strict TOML decode, or atom-identity rule: **none exists**. The one
adjacent decision — CLI error classification — is already `clierr.ErrorGroup`
and the RDR cites it as RDR 0005's rather than re-minting it. No
new-rule-with-an-existing-sibling finding.

## Findings

### CV2-001 — CONTRADICTION (blocking): the spike's merge key joins set-literal members on `,`, which the draft's own contract forbids, and it silently drops a distinct atom

**This is the defect the draft names, reproduced against the normative spike.**

The literals clause (Normative Contracts, "Literals carry the same byte-exact
identity as keys") states:

> A set-valued literal MUST normalize to an ordered sequence of its members …
> it MUST NOT be rendered into a single delimiter-joined string as its
> normalized value. Joining on any delimiter makes membership ambiguous whenever
> a member contains that delimiter — space-joining collapses `["needs work"]`
> and `["needs", "work"]` into one literal, so two different predicates become
> one atom and the dump's atom sort certifies the collision as canonical.

`evidence/spikes/iter-2/main.go` does exactly this, on `,`:

```go
func (a Atom) literalString() string { return strings.Join(a.Literal, ",") }

func (a Atom) identity() string {
	return a.Key + "\x00" + a.Block + "\x00" + a.Operator + "\x00" + a.literalString()
}
```

`identity()` is the **merge map key** — `mergeAtoms` writes
`out[atom.identity()] = atom` — so two atoms differing only in where the comma
falls collapse to one. `compareAtoms` is correct (it uses
`slices.Compare(a.Literal, b.Literal)` over the sequence), which is precisely
why the loss is invisible: it happens in the merge, before the sort the draft
says "certifies the collision as canonical" can observe it.

Witness, run against the committed spike binary:

```toml
[tags.reviewers]
provenance = "observed"
kind = "string"

[context.revs-ab.match.reviewers]
in = ["a,b", "c"]

[[rule]]
id = "collide-merge"
use = ["draft", "revs-ab"]
[rule.match.recognized]
eq = "round-clean"
[rule.match.reviewers]
in = ["a", "b,c"]
[rule.write]
stage = "prelock"
```

```
$ ./rdr0002spike2 collide2.toml
rdr.collide-merge#a     … atoms=[reviewers.eq=a@match; status.eq=Draft@match]   …
rdr.collide-merge#b,c   … atoms=[reviewers.eq=b,c@match; status.eq=Draft@match] …
```

The context contributed `["a,b", "c"]` and the rule contributed `["a", "b,c"]` —
two **distinct** atoms under the draft's identity tuple `(key, block, operator
token, literal)`. Only the rule's survived. Under the contexts clause the two
should both survive and the rule should be a **dead rule** (two non-identical
constraints on one key accumulated through inheritance) reported by RDR 0006's
unreachable-rule check. Instead it normalizes to a live, matching, expanding
rule. That is the failure mode the contexts clause explicitly rules out:

> Keying the merge on `(block, key, operator)` alone silently drops one
> constraint and makes the survivor depend on map iteration order over the `use`
> list … the loss happens before the sort can see it, so ordering determinism
> cannot detect it.

The draft anticipated a *prefix* key and forbade it; the spike instead lost the
same information through a *lossy rendering of the last field*. Same outcome,
route the clause does not name.

Reachability is unrestricted: only `#` is banned, and only from rule ids,
alphabet members, and match-block `in` members (`main.go:301,484,567`). A comma
is authorable in any `string`-kind tag value and in any guard-block set literal.
`neg-in-member-hash.toml` and `neg-alphabet-hash.toml` guard `#` only.

The rendering is lossy at the dump too — both atoms print as
`reviewers.in=a,b,c`, so the review surface cannot distinguish them either.

**Why this is a draft finding and not just a spike bug.** The spike is the
draft's normative evidence: `output.txt` is cited as the witness for A1, A11,
A13, the desk trace, and the Round-Trip fidelity table. The contract clause is
correct; what is missing is that the clause bans *one* delimiter-joining site
(the normalized value) and is silent on the **atom identity/merge key** being
derived from a joined rendering. A conforming implementation could read the
clause, keep `Literal []string` as the normalized value exactly as required, and
still key its merge on a joined string — as this one does.

**Owed**: the clause should state that the merge key and any identity rendering
compare the member **sequence**, not a joined string; and the Round-Trip lossy-site
list, which currently names one site ("set-valued and multi-entry fields are
rendered with unescaped separators"), should say that this makes the *rendered*
identity non-recoverable for set literals — the `#`-banning argument that makes
row identity recoverable does not extend to atom literals.

### CV2-002 — SILENCE: three enumerated load categories have no negative fixture, and the Testing Strategy asserts one mutated fixture per category

The category block enumerates the load categories; scenario 3 of the Testing
Strategy requires "one mutated fixture per category, asserted **by category, not
message text**, tripping one category and no other." Three enumerated categories
have no fixture in `evidence/spikes/iter-2/neg/` (37 cases) or
`negative-cases.txt` (92 lines):

| Category | Implemented? | Fixture? |
| --- | --- | --- |
| `cyclic context inheritance` | **yes** — `main.go::mergeContext` cycle guard | **no** |
| `malformed tag declaration` | **yes** — `main.go::validateTags` (bad provenance; missing kind) | **no** |
| `duplicate model id` | n/a — cross-document by construction | **no** (disclosed: A10 states the pair-of-documents fixture is owed) |

I confirmed the first two fire, by mutating the committed fixture:

```
$ ./rdr0002spike2 neg-cyclic.toml
refused: cyclic context inheritance at "prelock"
$ ./rdr0002spike2 neg-no-kind.toml
refused: malformed tag declaration: "cluster_ready" has no kind
```

So this is a **fixture gap, not a contract gap** — the contract and the
implementation agree, and nothing is refuted. But it is exactly the gap the
"every oracle must have a failing control" rail exists to close, and iteration 1's
CV-007 recorded the same thing for `cyclic context inheritance` on the iter-1
spike. It survived the §D7 re-authoring unfixed, which is why it is worth
re-raising rather than treating as closed.

`duplicate model id` is already disclosed on A10 and needs no new finding.

### CV2-003 — SILENCE: `[model.metadata]` is exempted from strictness, but nothing bounds what the exemption may carry

The layout clause says:

> **`[model.metadata]` is the one sanctioned extension namespace.** The loader
> MUST decode it as a free-form table, carry it through to the normalized model
> untouched, and MUST NOT interpret any key in it; strictness applies everywhere
> else.

The draft is silent on three consequences of an uninterpreted free-form table:

1. **Does it participate in the Round-Trip invariant?** The invariant compares
   "the full field list the dump contract carries," and the dump field list
   (row identity, locator, outcome, atoms, gate ids, next-state tags, writes,
   required-owned keys, escape classes) does **not** include model metadata —
   it is a model-level, not row-level, field. So two documents differing only in
   `[model.metadata]` satisfy the row-set invariant while carrying different
   normalized models. Whether that is intended is unstated.
2. **Does it reach the dump?** The dump MUST carry "every field" of the
   normalized *candidate row*; metadata is not one. The desk trace witnesses
   `metadata_keys=owner,review_cadence` in the load step, so the spike surfaces
   it — but no clause requires or forbids that.
3. **Is nesting bounded?** "Free-form table" admits arbitrary depth. The spike's
   `ModelMeta` handling and the strict decoder interact here in a way no clause
   fixes.

None of this refutes a claim. It is a genuine silence at a spot the strict-decoding
contract deliberately opens a hole, and the hole's edges are undefined.

### CV2-004 — SILENCE: the fail-fast ordering freedom and the one-defect-per-fixture rule are stated as jointly sufficient, but the merge happens before validation

The version-gate clause says load is fail-fast, that the order of independent
checks is deliberately unspecified, and that:

> the Testing Strategy's one-defect-per-fixture rule makes order unobservable by
> construction

That holds for *validation* categories. It does not obviously hold across the
**validate/normalize boundary**, because several categories are decided during
normalization (the clause itself says so: outcome binding, lifted-literal
alphabet membership, the write-free escape check). CV2-001 is a concrete case
where a defect is consumed by the merge *before* any validation can classify it:
the dropped atom never becomes a category at all, so a fixture carrying that one
defect trips **zero** categories rather than one. The "unobservable by
construction" argument assumes every defect eventually reaches a classifier.

This is a silence about the argument's scope, not a contradiction in it. Worth a
sentence bounding the claim to defects that reach a check.

## Step 1/2 — verification questions (12; 8 source-answerable)

Answered independently, each source answer citing what was actually read.

| # | Question | Kind | Answer |
| --- | --- | --- | --- |
| 1 | Does `escapeOrRefuse` filter on outcome before gating, as the escape clause claims? | source | **Yes.** `resolve.go::escapeOrRefuse` — `if row.Outcome != in.Recognized \|\| !view.matches(row.Match) { continue }`, then `gate(...)`. |
| 2 | Does any `internal/` symbol named `OpExists`/`LiteralTrue`/`LiteralFalse` exist? | source | **No.** `grep -rn` over `internal/` returns nothing. Disclosed by the RDR. |
| 3 | Does `resolve.Row` carry a per-atom block field or an `Atom` type? | source | **No.** `Row` has `Match []Tag` + `Guard string`; no `Atom` type in the package. Confirms A9/A12's stated block. |
| 4 | Does the kernel canonicalize tag keys anywhere? | source | **No.** `TagSet.Lookup`/`has`/`matches` all raw-index `s.tags[key]`. Only `ModeOf` folds, on a flag value. |
| 5 | Is the merge in the spike a set over the full atom identity tuple? | source | **No — this is CV2-001.** `mergeAtoms` keys on `Atom.identity()`, whose last field is `strings.Join(a.Literal, ",")`. Distinct member sequences collide. |
| 6 | Does any negative fixture exercise `cyclic context inheritance` or `malformed tag declaration`? | source | **No — CV2-002.** 37 files in `neg/`, none for either; both checks exist and fire when I mutate a fixture. |
| 7 | Is `,` rejected anywhere in an authored literal? | source | **No.** Only `#`, and only at `main.go:301` (alphabet), `:484` (rule id), `:567` (match-block `in` member). |
| 8 | Does `config.go::Load` contain parse-validation code? | source | **No.** `internal/cli/config/config.go:98` is a `TODO`; no TOML library is chosen. Confirms A5 as corrected. |
| 9 | Does the draft define what `[model.metadata]` may contain or whether it round-trips? | RDR-internal | **RDR is silent — CV2-003.** The clause fixes only "free-form, carried untouched, never interpreted." |
| 10 | Do "combined predicate set" 's two extents ever collide in the outcome-binding path? | RDR-internal | **No.** The narrow extent is fixed twice (Normative Contracts + Load-Bearing Decisions) with an explicit MUST NOT-read-here. `findOutcomeAtom` enforces block `match` and refuses otherwise. Coherent. |
| 11 | Is row kind stored as a field anywhere in the normalized value? | source | **No.** `spike main.go::Row` has no kind field; `func (r Row) kind()` derives it from `len(r.EscapeClasses) > 0`, matching the derived-view-property clause and RDR 0009. |
| 12 | Does the one-defect-per-fixture rule cover defects consumed before classification? | RDR-internal | **RDR is silent — CV2-004.** CV2-001 is an instance: the defect is absorbed by the merge and trips no category. |

## Summary

- **1 blocking contradiction** — CV2-001, source-grounded with a reproduction
  against the committed spike. The draft's own literals clause names the failure;
  the spike commits it through a route the clause does not close.
- **3 silences** — CV2-002 (fixture gap on two implemented categories, one a
  re-raise of iter-1's CV-007), CV2-003 (`[model.metadata]` edges), CV2-004
  (scope of the fail-fast argument).
- **0 refuted codebase claims.** Every kernel citation the draft makes CONFIRMS
  verbatim, including all three re-verify targets' evidence (A1's §D7 layout,
  A9's and A12's stated structural block).
- **0 new-rule-with-existing-sibling.** The inverse sweep found no adjacent path
  making any decision this RDR mints.

The re-verify set A1/A9/A12 is grounded: A1's closed-layout claims hold against
the re-authored fixture and the strict decoder; A9's and A12's `Pending` status
is correct and their blocking reason (`resolve.Row` has no atom shape) is
CONFIRMED at source.
